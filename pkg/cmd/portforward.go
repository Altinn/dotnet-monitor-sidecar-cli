package cmd

import (
	"context"
	"fmt"
	"path/filepath"

	dmskube "github.com/altinn/dotnet-monitor-sidecar-cli/pkg/kubernetes"
	"github.com/altinn/dotnet-monitor-sidecar-cli/pkg/utils"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// ForwardPort forwards port 52323 from host to a pod
func ForwardPort(ctx context.Context, kubeconfig string, namespace string, podname string) {
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
	h.PortForward(ctx, namespace, podname)
}
