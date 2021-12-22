package kubernetes

import (
	"context"
	"reflect"
	"testing"

	"github.com/altinn/dotnet-monitor-sidecar-cli/pkg/resources"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestHelper_GetDDApplyInfo(t *testing.T) {
	type args struct {
		namespace      string
		deploymentname string
	}
	tests := []struct {
		name     string
		args     args
		initfile string
		want     resources.DDConfig
		wantErr  bool
	}{
		{
			name: "Returns DDConfig",
			args: args{
				namespace:      "test",
				deploymentname: "test",
			},
			initfile: `testdata/deployment/get-ddinfo.yaml`,
			want: resources.DDConfig{
				ContainerToDebug:   "dotnet-container",
				TmpdirAdded:        true,
				DebugContainerName: "debug",
				SecretName:         "dd-monitor-apikey-629pf",
			},
			wantErr: false,
		},
		{
			name: "Deployment not found",
			args: args{
				namespace:      "test",
				deploymentname: "missing",
			},
			initfile: `testdata/deployment/get-ddinfo.yaml`,
			want:     resources.DDConfig{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			var d appsv1.Deployment
			err := getObjectFromFile(tt.initfile, &d)
			if err != nil {
				t.Errorf("failed to get deployment from file: %v", err)
				return
			}
			c := testclient.NewSimpleClientset(&d)
			h := &Helper{
				Client: c,
			}
			got, err := h.GetDDApplyInfo(ctx, tt.args.namespace, tt.args.deploymentname)
			if (err != nil) != tt.wantErr {
				t.Errorf("Helper.GetDDApplyInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Helper.GetDDApplyInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHelper_GetDDPodApplyInfo(t *testing.T) {
	type args struct {
		namespace string
		podname   string
	}
	tests := []struct {
		name     string
		args     args
		initfile string
		want     resources.DDConfig
		wantErr  bool
	}{
		{
			name: "Returns DDConfig",
			args: args{
				namespace: "default",
				podname:   "dotnet-app-75c6b8c7cf-4h6zw",
			},
			initfile: `testdata/deployment/get-ddinfo-pod.yaml`,
			want: resources.DDConfig{
				ContainerToDebug:   "dotnet-container",
				TmpdirAdded:        true,
				DebugContainerName: "debug",
				SecretName:         "dd-monitor-apikey-n6kps",
			},
			wantErr: false,
		},
		{
			name: "Pod not found",
			args: args{
				namespace: "default",
				podname:   "dotnet-app-75c6b8c7cf-aaaa",
			},
			initfile: `testdata/deployment/get-ddinfo-pod.yaml`,
			want:     resources.DDConfig{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			var pod corev1.Pod
			err := getObjectFromFile(tt.initfile, &pod)
			if err != nil {
				t.Errorf("failed to get deployment from file: %v", err)
				return
			}
			c := testclient.NewSimpleClientset(&pod)
			h := &Helper{
				Client: c,
			}
			got, err := h.GetDDPodApplyInfo(ctx, tt.args.namespace, tt.args.podname)
			if (err != nil) != tt.wantErr {
				t.Errorf("Helper.GetDDApplyInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Helper.GetDDApplyInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
