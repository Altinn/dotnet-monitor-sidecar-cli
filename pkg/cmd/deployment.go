package cmd

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/altinn/dotnet-monitor-sidecar-cli/pkg/errors"
	dmskube "github.com/altinn/dotnet-monitor-sidecar-cli/pkg/kubernetes"
	"github.com/altinn/dotnet-monitor-sidecar-cli/pkg/utils"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// AddToDeployment adds a debug sidecar to a deployment and configures it
func AddToDeployment(ctx context.Context, kubeconfig string, namespace string, deploymentname, containername, debugimage string) {
	if home := homedir.HomeDir(); home != "" && kubeconfig == "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	if namespace == "" {
		namespace, err = utils.GetNamespaceFromCurrentContext()
		if err != nil {
			fmt.Printf("Error getting namespace: %v", err)
			return
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	h := dmskube.Helper{
		Client: clientset,
	}
	d, token, err := h.AddDebugSidecarDeployment(ctx, namespace, deploymentname, containername, debugimage)

	if err != nil {
		if errors.IsAlreadyPresent(err) {
			fmt.Printf("Debug sidecar already attached to deployment %s\n", d.Name)
			return
		}
		fmt.Printf("Failed to attach sidecar to deployment %s: %v\n", deploymentname, err)
		return
	}
	fmt.Printf("Added sidecar to deployment %s with uid %s\n", d.Name, d.UID)
	fmt.Printf("Portforward to one of the pods with ddcli port-forward [podname].\nQuery the API with this auth header:\nAuthorization: Bearer %s\n", token)
}

// RemoveFromDeployment removes the debug sidecar and configuration from a deployment
func RemoveFromDeployment(ctx context.Context, kubeconfig string, namespace string, deploymentname string) {
	if home := homedir.HomeDir(); home != "" && kubeconfig == "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	if namespace == "" {
		namespace, err = utils.GetNamespaceFromCurrentContext()
		if err != nil {
			fmt.Printf("Error getting namespace: %v", err)
			return
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	h := dmskube.Helper{
		Client: clientset,
	}
	d, err := h.RemoveDebugSidecarDeployment(ctx, namespace, deploymentname)

	if err != nil {
		if errors.IsNotPresent(err) {
			fmt.Printf("Debug sidecar not attached to deployment %s\n", d.Name)
			return
		}
		panic(err.Error())
	}
	fmt.Printf("Removed sidecar from deployment %s with uid %s\n", d.Name, d.UID)
}
