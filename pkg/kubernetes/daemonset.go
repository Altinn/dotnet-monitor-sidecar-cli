package kubernetes

import (
	"context"

	"github.com/altinn/dotnet-monitor-sidecar-cli/pkg/resources"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AddDebugSidecarDaemonSet adds a debug sidecar to a daemonset
func (h *Helper) AddDebugSidecarDaemonSet(ctx context.Context, namespace, daemonsetname, containerToDebug, debugimage string) (*appsv1.DaemonSet, string, error) {
	d, err := h.Client.AppsV1().DaemonSets(namespace).Get(ctx, daemonsetname, metav1.GetOptions{})
	if err != nil {
		return nil, "", err
	}
	sn, token, err := h.CreateJWKSecret(ctx, namespace, daemonsetname)
	if err != nil {
		return nil, "", err
	}
	template, err := resources.AddDebugContainerPodTemplate(d.Spec.Template, namespace, containerToDebug, debugimage, sn)
	if err != nil {
		h.RemoveJWKSecret(ctx, namespace, sn)
		return nil, "", err
	}
	d.Spec.Template = template
	ds, err := h.Client.AppsV1().DaemonSets(namespace).Update(ctx, d, metav1.UpdateOptions{})
	if err != nil {
		h.RemoveJWKSecret(ctx, namespace, sn)
		return nil, "", err
	}
	return ds, token, nil
}

// RemoveDebugSidecarDaemonSet removes the debug sidecar from a daemonset
func (h *Helper) RemoveDebugSidecarDaemonSet(ctx context.Context, namespace, deploymentname string) (*appsv1.DaemonSet, error) {
	d, err := h.Client.AppsV1().DaemonSets(namespace).Get(ctx, deploymentname, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	ddConfig, err := resources.DDConfigFromPodTemplate(d.Spec.Template)
	if err != nil {
		return nil, err
	}
	err = h.RemoveJWKSecret(ctx, namespace, ddConfig.SecretName)
	if err != nil {
		return nil, err
	}
	d.Spec.Template, err = resources.RemoveDebugContainerPodTemplate(d.Spec.Template, namespace, ddConfig.ContainerToDebug)
	if err != nil {
		return nil, err
	}
	return h.Client.AppsV1().DaemonSets(namespace).Update(ctx, d, metav1.UpdateOptions{})
}

//ListDaemonsetsInNamespace returns list of deployments in a namespace
func (h *Helper) ListDaemonsetsInNamespace(ctx context.Context, namespace string) (*appsv1.DaemonSetList, error) {
	return h.Client.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
}
