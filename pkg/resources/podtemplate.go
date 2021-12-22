package resources

import (
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

// AddDebugContainerPodTemplate adds debug sidecar to a PodTemplateSpec object
func AddDebugContainerPodTemplate(template corev1.PodTemplateSpec, namespace, containerToDebug, debugimage, secretname string) (corev1.PodTemplateSpec, error) {
	if template.Annotations["dev.local/dd-added"] == "true" {
		return corev1.PodTemplateSpec{}, fmt.Errorf("debug sidecar already present")
	}
	debugSidecarName := getDebugContainerName(template.Spec.Containers)
	existingVolume, tmpVolume, err := getTmpVolume(template.Spec, containerToDebug)
	if err != nil {
		return corev1.PodTemplateSpec{}, err
	}
	if err != nil {
		return corev1.PodTemplateSpec{}, err
	}

	debugSidecar := generateSidecarContainerSpec(debugSidecarName, tmpVolume.Name, debugimage, secretname)
	appliedConfig := DDConfig{
		ContainerToDebug:   containerToDebug,
		DebugContainerName: debugSidecarName,
		TmpdirAdded:        !existingVolume,
		SecretName:         secretname,
	}

	if len(template.Spec.Containers) > 1 {
		for _, c := range template.Spec.Containers {
			if c.Name == containerToDebug {
				if !existingVolume {
					c.VolumeMounts = append(c.VolumeMounts, corev1.VolumeMount{
						Name:      tmpVolume.Name,
						MountPath: "/tmp",
					})
				}
			}
		}
	} else {
		if !existingVolume {
			template.Spec.Containers[0].VolumeMounts = append(template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
				Name:      tmpVolume.Name,
				MountPath: "/tmp",
			})
		}
		appliedConfig.ContainerToDebug = template.Spec.Containers[0].Name
	}
	if !existingVolume {
		template.Spec.Volumes = append(template.Spec.Volumes, tmpVolume)
	}
	template.Spec.Volumes = append(template.Spec.Volumes, corev1.Volume{
		Name: secretname,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: secretname,
			},
		},
	})
	template.Spec.Containers = append(template.Spec.Containers, debugSidecar)
	if template.Annotations == nil {
		template.Annotations = make(map[string]string)
	}
	template.Annotations["dev.local/dd-added"] = "true"
	b, err := json.Marshal(appliedConfig)
	if err != nil {
		return corev1.PodTemplateSpec{}, err
	}
	template.Annotations["dev.local/dd-apply"] = string(b)
	return template, nil
}

// RemoveDebugContainerPodTemplate removes debug sidecar from a PodTemplateSpec object
func RemoveDebugContainerPodTemplate(template corev1.PodTemplateSpec, namespace, containerToDebug string) (corev1.PodTemplateSpec, error) {
	if template.Annotations["dev.local/dd-added"] == "" {
		return corev1.PodTemplateSpec{}, fmt.Errorf("debug sidecar not present")
	}
	appliedConfig := DDConfig{}
	err := json.Unmarshal([]byte(template.Annotations["dev.local/dd-apply"]), &appliedConfig)
	if err != nil {
		return corev1.PodTemplateSpec{}, err
	}
	template.Spec.Containers = removeContainer(template.Spec.Containers, appliedConfig.DebugContainerName)

	_, tmpVolume, err := getTmpVolume(template.Spec, containerToDebug)
	if err != nil {
		return corev1.PodTemplateSpec{}, err
	}
	if appliedConfig.TmpdirAdded {
		template.Spec.Volumes = removeVolume(template.Spec.Volumes, tmpVolume.Name)
		template.Spec.Containers = removeTmpVolumeMount(template.Spec.Containers, containerToDebug, tmpVolume.Name)
	}
	template.Spec.Volumes = removeVolume(template.Spec.Volumes, appliedConfig.SecretName)
	delete(template.Annotations, "dev.local/dd-added")
	delete(template.Annotations, "dev.local/dd-apply")
	return template, nil
}

// DDConfigFromPodTemplate returns the DDConfig defined in PodTemplateSpec object
func DDConfigFromPodTemplate(template corev1.PodTemplateSpec) (DDConfig, error) {
	if template.Annotations["dev.local/dd-added"] != "true" {
		return DDConfig{}, fmt.Errorf("debug sidecar not present")
	}
	appliedConfig := DDConfig{}
	err := json.Unmarshal([]byte(template.Annotations["dev.local/dd-apply"]), &appliedConfig)
	if err != nil {
		return DDConfig{}, err
	}
	return appliedConfig, nil
}
