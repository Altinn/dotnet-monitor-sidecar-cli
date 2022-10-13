package utils

import (
	"path/filepath"
	"strings"

	dmskube "github.com/altinn/dotnet-monitor-sidecar-cli/pkg/kubernetes"
	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// AutoCompleteDaemonSets implements autocompletion for the daemonset commands
func AutoCompleteDaemonSets(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	kubeconfig, namespace, err := getFlags(cmd)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	h, namespace, err := getKubernetesHelper(kubeconfig, namespace)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	daemonsets, err := h.ListDaemonsetsInNamespace(cmd.Context(), namespace)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	names := getFilteredDaemonSetNames(daemonsets.Items, toComplete)
	return names, cobra.ShellCompDirectiveDefault
}

// AutoCompleteDeployments implements autocompletion for the deployment commands
func AutoCompleteDeployments(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	kubeconfig, namespace, err := getFlags(cmd)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	h, namespace, err := getKubernetesHelper(kubeconfig, namespace)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	deployments, err := h.ListDeploymentsInNamespace(cmd.Context(), namespace)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	names := getFilteredDeploymentNames(deployments.Items, toComplete)
	return names, cobra.ShellCompDirectiveDefault
}

// AutoCompletePodsWithDebugContainer implements autocompletion for the pod commands where debug contianer is present
func AutoCompletePodsWithDebugContainer(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	kubeconfig, namespace, err := getFlags(cmd)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	h, namespace, err := getKubernetesHelper(kubeconfig, namespace)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	pods, err := h.ListPodsInNamespace(cmd.Context(), namespace)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	names := getFilteredPodNamesWithDebugContainer(pods.Items, toComplete)
	return names, cobra.ShellCompDirectiveDefault
}

// getKubernetesHelper returns a kubernetes helper
func getKubernetesHelper(kubeconfig, namespace string) (h dmskube.Helper, ns string, err error) {
	if home := homedir.HomeDir(); home != "" && kubeconfig == "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return dmskube.Helper{}, "", err
	}
	if namespace == "" {
		namespace, err = GetNamespaceFromCurrentContext()
		if err != nil {
			namespace = "default"
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return dmskube.Helper{
		Client: clientset,
	}, namespace, nil
}

// getFlags returns kubeconfig and namespace flags
func getFlags(cmd *cobra.Command) (namespace, kubeconfig string, err error) {
	kubeconfig, err = cmd.Flags().GetString("kubeconfig")
	if err != nil {
		return
	}
	namespace, err = cmd.Flags().GetString("namespace")
	return
}

func getFilteredDeploymentNames(deployments []appsv1.Deployment, filter string) (names []string) {
	for _, d := range deployments {
		if strings.HasPrefix(d.Name, filter) {
			names = append(names, d.Name)
		}
	}
	return
}

func getFilteredDaemonSetNames(daemonsets []appsv1.DaemonSet, filter string) (names []string) {
	for _, d := range daemonsets {
		if strings.HasPrefix(d.Name, filter) {
			names = append(names, d.Name)
		}
	}
	return
}

func getFilteredPodNamesWithDebugContainer(pods []corev1.Pod, filter string) (names []string) {
	for _, p := range pods {
		if strings.HasPrefix(p.Name, filter) && p.Annotations["dev.local/dd-added"] == "true" {
			names = append(names, p.Name)
		}
	}
	return
}
