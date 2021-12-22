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

func TestHelper_RemoveDebugSidecarDaemonSet(t *testing.T) {
	type args struct {
		ctx            context.Context
		namespace      string
		deploymentname string
	}
	tests := []struct {
		name       string
		initfile   string
		goldenfile string
		args       args
		want       *appsv1.DaemonSet
		wantErr    bool
	}{
		{
			name: "Remove success",
			args: args{
				namespace:      "test",
				deploymentname: "test",
			},
			initfile:   "testdata/daemonset/remove.yaml",
			goldenfile: "testdata/daemonset/remove.golden",
			wantErr:    false,
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
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "dd-monitor-apikey-w4hhs",
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
			var expected *appsv1.DaemonSet
			actual, err := h.RemoveDebugSidecarDaemonSet(tt.args.ctx, tt.args.namespace, tt.args.deploymentname)
			if (err != nil) != tt.wantErr {
				t.Errorf("Helper.RemoveDebugSidecarDaemonSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			err = readUpdateGoldenFile(tt.goldenfile, *update, &expected, actual)
			if err != nil {
				t.Errorf("readUpdateGoldenFile() error = %v", err)
				return
			}
			if objToString(expected) != objToString(*actual) {
				t.Errorf("Unexpected deployment returned. Diff: %s", getDiffs(objToString(expected), objToString(*actual)))
			}
			select {
			case dep := <-updates:
				if objToString(actual) != objToString(*dep) {
					t.Errorf("Daemonset returned not same as applied:\n%s", getDiffs(objToString(actual), objToString(*dep)))
				}
			case <-time.After(wait.ForeverTestTimeout):
				t.Errorf("Timed out waiting for deployment to be updated")
			}
		})
	}
}
