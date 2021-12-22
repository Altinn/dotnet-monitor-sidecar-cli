package resources

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	utilrand "k8s.io/apimachinery/pkg/util/rand"
)

const secretMountPath = "/etc/dotnet-monitor"

func generateSidecarContainerSpec(containername, mountname, debugimage, secretname string) corev1.Container {
	return corev1.Container{
		Name:            containername,
		Image:           debugimage,
		ImagePullPolicy: corev1.PullIfNotPresent,
		Ports:           []corev1.ContainerPort{{ContainerPort: 52323}},
		Args: []string{
			"--urls",
			"http://*:52323",
		},
		TerminationMessagePath:   "/dev/termination-log",
		TerminationMessagePolicy: "File",
		SecurityContext: &corev1.SecurityContext{
			Capabilities: &corev1.Capabilities{
				Add: []corev1.Capability{
					corev1.Capability("SYS_PTRACE"),
				},
			},
		},
		Resources: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("50m"),
				corev1.ResourceMemory: resource.MustParse("32Mi"),
			},
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("250m"),
				corev1.ResourceMemory: resource.MustParse("256Mi"),
			},
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      mountname,
				MountPath: "/tmp",
			},
			{
				Name:      secretname,
				MountPath: secretMountPath,
			},
		},
	}
}

func getDebugContainerName(existingContainers []corev1.Container) string {
	if !debugContainerNameTaken(existingContainers) {
		return "debug"
	}
	return fmt.Sprintf("%s-%s", "debug", utilrand.String(5))
}

func debugContainerNameTaken(containers []corev1.Container) bool {
	for _, c := range containers {
		if c.Name == "debug" {
			return true
		}
	}
	return false
}

func removeContainer(containers []corev1.Container, containerName string) []corev1.Container {
	for i, c := range containers {
		if c.Name == containerName {
			return append(containers[:i], containers[i+1:]...)
		}
	}
	return containers
}
