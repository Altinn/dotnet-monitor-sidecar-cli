package kubernetes

import (
	"context"
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/jwt"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/informers"
	testclient "k8s.io/client-go/kubernetes/fake"
	clienttesting "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/cache"
)

func TestHelper_AddDebugSidecarDaemonSet(t *testing.T) {
	type args struct {
		ctx              context.Context
		namespace        string
		daemonsetname    string
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
			name: "Success",
			args: args{
				namespace:        "test",
				daemonsetname:    "test",
				containerToDebug: "test",
				debugimage:       "test",
			},
			initfile:   "testdata/daemonset/add.yaml",
			goldenfile: "testdata/daemonset/add.golden",
			wantErr:    false,
		},
		{
			name: "Daemonset notfound",
			args: args{
				namespace:        "test",
				daemonsetname:    "not-found",
				containerToDebug: "test",
				debugimage:       "test",
			},
			initfile: "testdata/daemonset/add.yaml",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			var existing *appsv1.DaemonSet
			if tt.initfile != "" {
				var daemonset appsv1.DaemonSet
				err := getObjectFromFile(tt.initfile, &daemonset)
				if err != nil {
					t.Errorf("getDeploymentFromFile() error = %v", err)
					return
				}
				existing = &daemonset
			}
			updates := make(chan *appsv1.DaemonSet, 1)
			c := createFakeClientWithDaemonsetWatcher(
				ctx,
				&cache.ResourceEventHandlerFuncs{
					UpdateFunc: func(old, new interface{}) {
						updates <- new.(*appsv1.DaemonSet)
					},
				},
				existing,
			)
			h := &Helper{
				Client: c,
			}
			var expected *appsv1.DaemonSet
			actual, gotToken, err := h.AddDebugSidecarDaemonSet(tt.args.ctx, tt.args.namespace, tt.args.daemonsetname, tt.args.containerToDebug, tt.args.debugimage)
			if (err != nil) != tt.wantErr {
				t.Errorf("Helper.AddDebugSidecarDaemonSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				err = readUpdateGoldenFile(tt.goldenfile, *update, &expected, actual)
				if err != nil {
					t.Errorf("failed to read golden file %v", err)
				}
				if len(actual.Spec.Template.Spec.Containers) != len(expected.Spec.Template.Spec.Containers) {
					t.Errorf("Unexpected number of containers in the deployment. Expected %v, got %v", len(expected.Spec.Template.Spec.Containers), len(actual.Spec.Template.Spec.Containers))
				}
				if len(actual.Spec.Template.Spec.Volumes) != len(expected.Spec.Template.Spec.Volumes) {
					t.Errorf("unexpected number of volumes, got %v, want %v", len(actual.Spec.Template.Spec.Volumes), len(expected.Spec.Template.Spec.Volumes))
				}
				// Simple token validation
				jt, err := jwt.Parse([]byte(gotToken))
				if err != nil {
					t.Errorf("Fail to parse token: %v", err)
					return
				}
				err = jwt.Validate(jt)
				if err != nil {
					t.Errorf("Fail to validate token: %v", err)
				}
				select {
				case ds := <-updates:
					if objToString(actual) != objToString(*ds) {
						t.Errorf("Daemonset returned not same as applied:\n%s", getDiffs(objToString(actual), objToString(*ds)))
					}
				case <-time.After(wait.ForeverTestTimeout):
					t.Errorf("Timed out waiting for daemonset to be updated")
				}
			}
		})
	}
}

func createFakeClientWithDaemonsetWatcher(ctx context.Context, handlers *cache.ResourceEventHandlerFuncs, objs ...runtime.Object) *testclient.Clientset {
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
	dsInformers := i.Apps().V1().DaemonSets().Informer()
	dsInformers.AddEventHandler(handlers)
	i.Start(ctx.Done())
	<-watcherStarted
	return c
}
