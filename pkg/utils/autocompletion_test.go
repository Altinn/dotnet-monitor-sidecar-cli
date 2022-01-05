package utils

import (
	"reflect"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_getFilteredDeploymentNames(t *testing.T) {
	type args struct {
		deployments []appsv1.Deployment
		filter      string
	}
	tests := []struct {
		name      string
		args      args
		wantNames []string
	}{
		{
			name: "Filter is blank",
			args: args{
				deployments: []appsv1.Deployment{
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "test-deployment",
						},
					},
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "another-deployment",
						},
					},
				},
				filter: "",
			},
			wantNames: []string{"test-deployment", "another-deployment"},
		},
		{
			name: "Filter matches one of the deployments",
			args: args{
				deployments: []appsv1.Deployment{
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "test-deployment",
						},
					},
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "another-deployment",
						},
					},
				},
				filter: "test",
			},
			wantNames: []string{"test-deployment"},
		},
		{
			name: "Filter matches non of the deployments",
			args: args{
				deployments: []appsv1.Deployment{
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "test-deployment",
						},
					},
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "another-deployment",
						},
					},
				},
				filter: "not-found",
			},
			wantNames: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNames := getFilteredDeploymentNames(tt.args.deployments, tt.args.filter); !reflect.DeepEqual(gotNames, tt.wantNames) {
				t.Errorf("getFilteredDeploymentNames() = %v, want %v", gotNames, tt.wantNames)
			}
		})
	}
}

func Test_getFilteredDaemonSetNames(t *testing.T) {
	type args struct {
		daemonsets []appsv1.DaemonSet
		filter     string
	}
	tests := []struct {
		name      string
		args      args
		wantNames []string
	}{
		{
			name: "Filter is blank",
			args: args{
				daemonsets: []appsv1.DaemonSet{
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "test-deployment",
						},
					},
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "another-deployment",
						},
					},
				},
				filter: "",
			},
			wantNames: []string{"test-deployment", "another-deployment"},
		},
		{
			name: "Filter matches one of the daemonsets",
			args: args{
				daemonsets: []appsv1.DaemonSet{
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "test-deployment",
						},
					},
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "another-deployment",
						},
					},
				},
				filter: "test",
			},
			wantNames: []string{"test-deployment"},
		},
		{
			name: "Filter matches non of the daemonsets",
			args: args{
				daemonsets: []appsv1.DaemonSet{
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "test-deployment",
						},
					},
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "another-deployment",
						},
					},
				},
				filter: "not-found",
			},
			wantNames: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNames := getFilteredDaemonSetNames(tt.args.daemonsets, tt.args.filter); !reflect.DeepEqual(gotNames, tt.wantNames) {
				t.Errorf("getFilteredDaemonSetNames() = %v, want %v", gotNames, tt.wantNames)
			}
		})
	}
}

func Test_getFilteredPodNamesWithDebugContainer(t *testing.T) {
	type args struct {
		daemonsets []corev1.Pod
		filter     string
	}
	tests := []struct {
		name      string
		args      args
		wantNames []string
	}{
		{
			name: "Filter is blank all pods have debug container",
			args: args{
				daemonsets: []corev1.Pod{
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "test-deployment",
							Annotations: map[string]string{
								"dev.local/dd-added": "true",
							},
						},
					},
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "another-deployment",
							Annotations: map[string]string{
								"dev.local/dd-added": "true",
							},
						},
					},
				},
				filter: "",
			},
			wantNames: []string{"test-deployment", "another-deployment"},
		},
		{
			name: "Filter matches one of the daemonsets with debug container",
			args: args{
				daemonsets: []corev1.Pod{
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "test-deployment",
							Annotations: map[string]string{
								"dev.local/dd-added": "true",
							},
						},
					},
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "another-deployment",
							Annotations: map[string]string{
								"dev.local/dd-added": "true",
							},
						},
					},
				},
				filter: "test",
			},
			wantNames: []string{"test-deployment"},
		},
		{
			name: "Filter matches non of the daemonsets with debug container",
			args: args{
				daemonsets: []corev1.Pod{
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "test-deployment",
							Annotations: map[string]string{
								"dev.local/dd-added": "true",
							},
						},
					},
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "another-deployment",
							Annotations: map[string]string{
								"dev.local/dd-added": "true",
							},
						},
					},
				},
				filter: "not-found",
			},
			wantNames: nil,
		},
		{
			name: "Filter matches daemonsets without debug container",
			args: args{
				daemonsets: []corev1.Pod{
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "test-deployment",
						},
					},
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "another-deployment",
							Annotations: map[string]string{
								"dev.local/dd-added": "true",
							},
						},
					},
				},
				filter: "test",
			},
			wantNames: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNames := getFilteredPodNamesWithDebugContainer(tt.args.daemonsets, tt.args.filter); !reflect.DeepEqual(gotNames, tt.wantNames) {
				t.Errorf("getFilteredDaemonSetNames() = %v, want %v", gotNames, tt.wantNames)
			}
		})
	}
}
