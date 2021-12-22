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

// AddToDaemonset setup debug sidecar to a Daemonset and configures it
func AddToDaemonset(ctx context.Context, kubeconfig string, namespace string, deploymentname, containername, debugimage string) {
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
	d, token, err := h.AddDebugSidecarDaemonSet(ctx, namespace, deploymentname, containername, debugimage)

	if err != nil {
		if errors.IsAlreadyPresent(err) {
			fmt.Printf("Debug sidecar already attached to deployment %s\n", d.Name)
			return
		}
		panic(err.Error())
	}
	fmt.Printf("Added sidecar to daemonset %s with uid %s\n", d.Name, d.UID)
	fmt.Printf("Portforward to one of the pods with ddcli port-forward [podname].\nQuery the API with this auth header:\nAuthorization: Bearer %s\n", token)
}

// RemoveFromDaemonset removes the debug sidecar and configuration from a daemonset
func RemoveFromDaemonset(ctx context.Context, kubeconfig string, namespace string, daemonsetname string) {
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
	d, err := h.RemoveDebugSidecarDaemonSet(ctx, namespace, daemonsetname)

	if err != nil {
		if errors.IsNotPresent(err) {
			fmt.Printf("Debug sidecar not attached to daemonset %s\n", d.Name)
			return
		}
		panic(err.Error())
	}
	fmt.Printf("Removed sidecar from daemonset %s with uid %s\n", d.Name, d.UID)
}
