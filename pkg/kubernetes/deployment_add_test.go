package kubernetes

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/ghodss/yaml"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/sergi/go-diff/diffmatchpatch"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/informers"
	testclient "k8s.io/client-go/kubernetes/fake"
	clienttesting "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/cache"
)

var (
	update = flag.Bool("update", false, "update the golden files of this test")
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestHelper_AddDebugSidecarDeployment(t *testing.T) {
	type args struct {
		namespace        string
		deploymentname   string
		containerToDebug string
		debugimage       string
	}
	tests := []struct {
		name       string
		args       args
		goldenfile string
		initfile   string
		wantErr    bool
	}{
		{
			name: "Returns deployment",
			args: args{
				namespace:        "test",
				deploymentname:   "test",
				containerToDebug: "test",
				debugimage:       "test",
			},
			initfile:   "testdata/deployment/add.yaml",
			goldenfile: "testdata/deployment/add.golden",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			var existing *appsv1.Deployment
			if tt.initfile != "" {
				var deployment appsv1.Deployment
				err := getObjectFromFile(tt.initfile, &deployment)
				if err != nil {
					t.Errorf("getDeploymentFromFile() error = %v", err)
					return
				}
				existing = &deployment
			}
			updateChan := make(chan *appsv1.Deployment, 1)
			c := createFakeClientWithDeploymentWatcher(
				ctx,
				&cache.ResourceEventHandlerFuncs{
					UpdateFunc: func(old, new interface{}) {
						updateChan <- new.(*appsv1.Deployment)
					},
				},
				existing,
			)
			h := &Helper{
				Client: c,
			}
			var expected appsv1.Deployment
			actual, token, err := h.AddDebugSidecarDeployment(ctx, tt.args.namespace, tt.args.deploymentname, tt.args.containerToDebug, tt.args.debugimage)
			err = readUpdateGoldenFile(tt.goldenfile, *update, &expected, actual)
			if err != nil {
				t.Errorf("readUpdateGoldenFile() error = %v", err)
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Helper.AddDebugSidecarDeployment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(actual.Spec.Template.Spec.Containers) != len(expected.Spec.Template.Spec.Containers) {
				t.Errorf("Unexpected number of containers in the deployment. Expected %v, got %v", len(expected.Spec.Template.Spec.Containers), len(actual.Spec.Template.Spec.Containers))
			}
			if len(actual.Spec.Template.Spec.Volumes) != len(expected.Spec.Template.Spec.Volumes) {
				t.Errorf("unexpected number of volumes, got %v, want %v", len(actual.Spec.Template.Spec.Volumes), len(expected.Spec.Template.Spec.Volumes))
			}

			// Simple token validation
			jt, err := jwt.Parse([]byte(token))
			if err != nil {
				t.Errorf("Fail to parse token: %v", err)
				return
			}
			err = jwt.Validate(jt)
			if err != nil {
				t.Errorf("Fail to validate token: %v", err)
			}
			select {
			case dep := <-updateChan:
				if objToString(actual) != objToString(*dep) {
					t.Errorf("Deployment returned not same as applied:\n%s", getDiffs(objToString(actual), objToString(*dep)))
				}
			case <-time.After(wait.ForeverTestTimeout):
				t.Errorf("Timed out waiting for deployment to be updated")
			}
		})
	}
}

func getObjectFromFile(filename string, obj interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open inputfile %v", err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read inputfile %v", err)
	}
	return yaml.Unmarshal(b, obj)
}

func objToString(input interface{}) string {
	data, err := yaml.Marshal(input)
	if err != nil {
		return ""
	}
	return string(data)
}

func getDiffs(expected, actual string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(objToString(expected), objToString(actual), false)
	return dmp.DiffPrettyText(diffs)
}

func readUpdateGoldenFile(goldenfile string, update bool, expected interface{}, actual interface{}) error {
	if update {
		data, err := yaml.Marshal(actual)
		if err != nil {
			return fmt.Errorf("failed to marhal input %v", err)
		}
		err = ioutil.WriteFile(goldenfile, data, 0644)
		if err != nil {
			return fmt.Errorf("failed to write goldefile %v", err)
		}
	}
	f, err := os.Open(goldenfile)
	if err != nil {
		return fmt.Errorf("failed to open goldenfile %v", err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read goldenfile %v", err)
	}
	err = yaml.Unmarshal(b, expected)
	return err
}

func createFakeClientWithDeploymentWatcher(ctx context.Context, handlers *cache.ResourceEventHandlerFuncs, objs ...runtime.Object) *testclient.Clientset {
	watcherStarted := make(chan struct{})
	c := testclient.NewSimpleClientset(objs...)
	c.PrependWatchReactor("*", func(action clienttesting.Action) (handled bool, ret watch.Interface, err error) {
		gvr := action.GetResource()
		ns := action.GetNamespace()
		watch, err := c.Tracker().Watch(gvr, ns)
		if err != nil {
			return false, nil, err
		}
		close(watcherStarted)
		return true, watch, nil
	})

	i := informers.NewSharedInformerFactory(c, 0)
	depInformers := i.Apps().V1().Deployments().Informer()
	depInformers.AddEventHandler(handlers)
	i.Start(ctx.Done())
	<-watcherStarted
	return c
}
