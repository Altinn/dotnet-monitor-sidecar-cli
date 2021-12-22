package kubernetes

import (
	"context"
	"reflect"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestHelper_ListDaemonsetsInNamespace(t *testing.T) {
	type args struct {
		ctx       context.Context
		namespace string
	}
	tests := []struct {
		name       string
		args       args
		initfile   string
		goldenfile string
		wantErr    bool
	}{
		{
			name: "Success",
			args: args{
				namespace: "test",
			},
			initfile:   "testdata/daemonset/list.yaml",
			goldenfile: "testdata/daemonset/list.golden",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d appsv1.DaemonSet
			err := getObjectFromFile(tt.initfile, &d)
			if err != nil {
				t.Errorf("getDeploymentFromFile() error = %v", err)
				return
			}
			c := testclient.NewSimpleClientset(&d)
			h := &Helper{
				Client: c,
			}
			got, err := h.ListDaemonsetsInNamespace(tt.args.ctx, tt.args.namespace)
			if (err != nil) != tt.wantErr {
				t.Errorf("Helper.ListDaemonsetsInNamespace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var expected *appsv1.DaemonSetList
			err = readUpdateGoldenFile(tt.goldenfile, *update, &expected, got)
			if err != nil {
				t.Errorf("readUpdateGoldenFile() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, expected) {
				t.Errorf("Helper.ListDaemonsetsInNamespace() = %v, want %v", got, expected)
			}
		})
	}
}
