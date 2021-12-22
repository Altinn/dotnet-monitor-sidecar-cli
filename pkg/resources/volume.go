package resources

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	utilrand "k8s.io/apimachinery/pkg/util/rand"
)

func getTmpVolume(podSpec corev1.PodSpec, containerToDebug string) (bool, corev1.Volume, error) {
	containerCount := len(podSpec.Containers)
	if containerCount > 1 && containerToDebug == "" {
		return false, corev1.Volume{}, fmt.Errorf("multiple containers present, pleas supply the one you want to debug")
	}
	var volumeMounts []corev1.VolumeMount
	if containerCount > 1 {
		found := false
		for _, c := range podSpec.Containers {
			if c.Name == containerToDebug {
				found = true
				volumeMounts = c.VolumeMounts
			}
		}
		if !found {
			return false, corev1.Volume{}, fmt.Errorf("could not find container with name %s", containerToDebug)
		}
	} else {
		volumeMounts = podSpec.Containers[0].VolumeMounts
	}
	for _, vm := range volumeMounts {
		if vm.MountPath == "/tmp" {
			for _, v := range podSpec.Volumes {
				if v.Name == vm.Name {
					return true, v, nil
				}
			}
		}
	}
	return false, corev1.Volume{
		Name: fmt.Sprintf("tmpfolder-%s", utilrand.String(5)),
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}, nil
}

func removeTmpVolumeMount(containers []corev1.Container, containerToDebug, tmpVolumeName string) []corev1.Container {
	for i, c := range containers {
		if containerToDebug == "" || c.Name == containerToDebug {
			for j, v := range c.VolumeMounts {
				if v.Name == tmpVolumeName {
					c.VolumeMounts = append(c.VolumeMounts[:j], c.VolumeMounts[j+1:]...)
					break
				}
			}
			containers[i].VolumeMounts = c.VolumeMounts
		}
	}
	return containers
}

func removeVolume(volumes []corev1.Volume, tmpVolumeName string) []corev1.Volume {
	for i, v := range volumes {
		if v.Name == tmpVolumeName {
			volumes = append(volumes[:i], volumes[i+1:]...)
			break
		}
	}
	return volumes
}
