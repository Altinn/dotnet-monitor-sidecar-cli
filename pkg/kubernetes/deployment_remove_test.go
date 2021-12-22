package kubernetes

import (
	"context"
	"testing"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
)

func TestHelper_RemoveDebugSidecarDeployment(t *testing.T) {
	type args struct {
		namespace      string
		deploymentname string
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
				namespace:      "test",
				deploymentname: "test",
			},
			initfile:   "testdata/deployment/remove.yaml",
			goldenfile: "testdata/deployment/remove.golden",
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
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "dd-monitor-apikey-629pf",
						Namespace: "test",
					},
					Data: map[string][]byte{
						"test": []byte("test"),
					},
				},
			)
			h := &Helper{
				Client: c,
			}
			var expected appsv1.Deployment
			actual, err := h.RemoveDebugSidecarDeployment(ctx, tt.args.namespace, tt.args.deploymentname)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveDebugSidecarDeployment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			err = readUpdateGoldenFile(tt.goldenfile, *update, &expected, actual)
			if err != nil {
				t.Errorf("readUpdateGoldenFile() error = %v", err)
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Helper.AddDebugSidecarDeployment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if objToString(expected) != objToString(*actual) {
				t.Errorf("Unexpected deployment returned. Diff: %s", getDiffs(objToString(expected), objToString(*actual)))
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
