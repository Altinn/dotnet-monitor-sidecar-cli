package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/altinn/dotnet-monitor-sidecar-cli/pkg/resources"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Helper struct for kubernetes helper methods for managing debug sidecars
type Helper struct {
	Client kubernetes.Interface
}

// AddDebugSidecarDeployment adds debug sidecar to a Deployment
func (h *Helper) AddDebugSidecarDeployment(ctx context.Context, namespace, deploymentname, containerToDebug, debugimage string) (*appsv1.Deployment, string, error) {
	d, err := h.Client.AppsV1().Deployments(namespace).Get(ctx, deploymentname, metav1.GetOptions{})
	if err != nil {
		return nil, "", err
	}
	sn, token, err := h.CreateJWKSecret(ctx, namespace, deploymentname)
	if err != nil {
		return nil, "", err
	}
	template, err := resources.AddDebugContainerPodTemplate(d.Spec.Template, namespace, containerToDebug, debugimage, sn)
	if err != nil {
		h.RemoveJWKSecret(ctx, namespace, sn)
		return nil, "", err
	}
	d.Spec.Template = template
	dep, err := h.Client.AppsV1().Deployments(namespace).Update(ctx, d, metav1.UpdateOptions{})
	if err != nil {
		h.RemoveJWKSecret(ctx, namespace, sn)
		return nil, "", err
	}
	return dep, token, nil
}

// RemoveDebugSidecarDeployment removes debug sidecar from a Deployment
func (h *Helper) RemoveDebugSidecarDeployment(ctx context.Context, namespace, deploymentname string) (*appsv1.Deployment, error) {
	d, err := h.Client.AppsV1().Deployments(namespace).Get(ctx, deploymentname, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	ddConfig, err := resources.DDConfigFromPodTemplate(d.Spec.Template)
	if err != nil {
		return nil, err
	}
	d.Spec.Template, err = resources.RemoveDebugContainerPodTemplate(d.Spec.Template, namespace, ddConfig.ContainerToDebug)
	if err != nil {
		return nil, err
	}
	err = h.RemoveJWKSecret(ctx, namespace, ddConfig.SecretName)
	if err != nil {
		return nil, err
	}
	return h.Client.AppsV1().Deployments(namespace).Update(ctx, d, metav1.UpdateOptions{})
}

// GetDDApplyInfo returns the debug sidecar apply info for a deployment from kubernetes
func (h *Helper) GetDDApplyInfo(ctx context.Context, namespace, deploymentname string) (resources.DDConfig, error) {
	d, err := h.Client.AppsV1().Deployments(namespace).Get(ctx, deploymentname, metav1.GetOptions{})
	if err != nil {
		return resources.DDConfig{}, err
	}
	return resources.DDConfigFromPodTemplate(d.Spec.Template)
}

// GetDDPodApplyInfo returns the debug sidecar apply info for a pod from kubernetes
func (h *Helper) GetDDPodApplyInfo(ctx context.Context, namespace, podname string) (resources.DDConfig, error) {
	p, err := h.Client.CoreV1().Pods(namespace).Get(ctx, podname, metav1.GetOptions{})
	if err != nil {
		return resources.DDConfig{}, err
	}
	if p.ObjectMeta.Annotations["dev.local/dd-added"] != "true" {
		return resources.DDConfig{}, fmt.Errorf("debug sidecar not present")
	}
	appliedConfig := resources.DDConfig{}
	err = json.Unmarshal([]byte(p.ObjectMeta.Annotations["dev.local/dd-apply"]), &appliedConfig)
	if err != nil {
		return resources.DDConfig{}, err
	}
	return appliedConfig, nil
}

//ListPodsInNamespace returns list of pods in a namespace
func (h *Helper) ListPodsInNamespace(ctx context.Context, namespace string) (*corev1.PodList, error) {
	return h.Client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
}

//ListDeploymentsInNamespace returns list of deployments in a namespace
func (h *Helper) ListDeploymentsInNamespace(ctx context.Context, namespace string) (*appsv1.DeploymentList, error) {
	return h.Client.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
}

//ListNamespaces returns list of namespaces
func (h *Helper) ListNamespaces(ctx context.Context) (*corev1.NamespaceList, error) {
	return h.Client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
}
