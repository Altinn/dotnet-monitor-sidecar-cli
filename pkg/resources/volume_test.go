package resources

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"

	corev1 "k8s.io/api/core/v1"
	utilrand "k8s.io/apimachinery/pkg/util/rand"
)

func Test_getTmpVolume(t *testing.T) {
	type args struct {
		podSpec          corev1.PodSpec
		containerToDebug string
	}
	tests := []struct {
		name        string
		args        args
		expExisting bool
		expVolume   corev1.Volume
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			name: "Returns volume with random name if no tmp volume is found",
			args: args{
				podSpec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: "test",
						},
					},
				},
				containerToDebug: "",
			},
			expExisting: false,
			expVolume: corev1.Volume{
				Name: fmt.Sprintf("tmpfolder-%s", utilrand.String(5)),
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				},
			},
			wantErr: false,
		},
		{
			name: "Returns volume with random name if tmp volume is mounted to other container",
			args: args{
				podSpec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: "test",
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "tmpfolder",
									MountPath: "/tmp",
								},
							},
						},
						{
							Name: "test2",
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "tmpfolder",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
				},
				containerToDebug: "test2",
			},
			expExisting: false,
			expVolume: corev1.Volume{
				Name: fmt.Sprintf("tmpfolder"),
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				},
			},
			wantErr: false,
		},
		{
			name: "Returns volume tmp volume if tmp volume is found for container to debug",
			args: args{
				podSpec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: "test",
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "tmpfolder",
									MountPath: "/tmp",
								},
							},
						},
						{
							Name: "test2",
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "tmpfolder",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
				},
				containerToDebug: "test",
			},
			expExisting: true,
			expVolume: corev1.Volume{
				Name: fmt.Sprintf("tmpfolder"),
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				},
			},
			wantErr: false,
		},
		{
			name: "Returns error if multiple containers and container to debug is empty",
			args: args{
				podSpec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: "test",
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "tmpfolder",
									MountPath: "/tmp",
								},
							},
						},
						{
							Name: "test2",
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "tmpfolder",
									MountPath: "/tmp",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "tmpfolder",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
				},
				containerToDebug: "",
			},
			expExisting: true,
			expVolume: corev1.Volume{
				Name: fmt.Sprintf("tmpfolder"),
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				},
			},
			wantErr: true,
		},
		{
			name: "Returns error if container not found",
			args: args{
				podSpec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: "test",
						},
						{
							Name: "test2",
						},
					},
				},
				containerToDebug: "test3",
			},
			expExisting: false,
			expVolume: corev1.Volume{
				Name: fmt.Sprintf("tmpfolder-%s", utilrand.String(5)),
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			existing, volume, err := getTmpVolume(tt.args.podSpec, tt.args.containerToDebug)
			if (err != nil) != tt.wantErr {
				t.Error("Expected error but non was returned")
				return
			}
			if !tt.wantErr {
				if existing != tt.expExisting {
					t.Errorf("getTmpVolume() expected = %v\n got %v", tt.expExisting, existing)
				}
				if !tt.expExisting {
					nameRegex := regexp.MustCompile("tmpfolder-[a-zA-Z0-9_.-]{5}")
					if !nameRegex.MatchString(volume.Name) {
						t.Errorf("tmp volume name does not match regex %s", nameRegex)
					}
					if !reflect.DeepEqual(volume.VolumeSource, tt.expVolume.VolumeSource) {
						t.Errorf("tmp volume source does not match expected %v\n got %v", tt.expVolume.VolumeSource, volume.VolumeSource)
					}
				}
				if tt.expExisting && !reflect.DeepEqual(volume, tt.expVolume) {
					t.Errorf("getTmpVolume() expected = %v\n got %v", tt.expVolume, volume)
				}
			}
		})
	}
}

func Test_removeTmpVolumeMount(t *testing.T) {
	type args struct {
		containers       []corev1.Container
		containerToDebug string
		tmpVolumeName    string
	}
	tests := []struct {
		name string
		args args
		want []corev1.Container
	}{
		{
			name: "Removes volume mount from container",
			args: args{
				containers: []corev1.Container{
					{
						Name: "test",
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "tmpfolder",
								MountPath: "/tmp",
							},
							{
								Name:      "othervolume",
								MountPath: "/other",
							},
						},
					},
				},
				containerToDebug: "test",
				tmpVolumeName:    "tmpfolder",
			},
			want: []corev1.Container{
				{
					Name: "test",
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "othervolume",
							MountPath: "/other",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeTmpVolumeMount(tt.args.containers, tt.args.containerToDebug, tt.args.tmpVolumeName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeTmpVolumeMount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_removeVolume(t *testing.T) {
	type args struct {
		volumes       []corev1.Volume
		tmpVolumeName string
	}
	tests := []struct {
		name string
		args args
		want []corev1.Volume
	}{
		{
			name: "Removes volume from volume slice",
			args: args{
				volumes: []corev1.Volume{
					{
						Name: "tmpfolder",
						VolumeSource: corev1.VolumeSource{
							EmptyDir: &corev1.EmptyDirVolumeSource{},
						},
					},
					{
						Name: "othervolume",
						VolumeSource: corev1.VolumeSource{
							EmptyDir: &corev1.EmptyDirVolumeSource{},
						},
					},
				},
				tmpVolumeName: "tmpfolder",
			},
			want: []corev1.Volume{
				{
					Name: "othervolume",
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeVolume(tt.args.volumes, tt.args.tmpVolumeName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeVolume() = %v, want %v", got, tt.want)
			}
		})
	}
}
