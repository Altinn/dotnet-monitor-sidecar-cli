package utils

import (
	"fmt"

	"k8s.io/client-go/tools/clientcmd"
)

// GetNamespaceFromCurrentContext returns the namespace from the current kubeconfig context
func GetNamespaceFromCurrentContext() (namespace string, err error) {
	kubeconfigs := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
	namespace, _, err = kubeconfigs.Namespace()
	if err != nil {
		err = fmt.Errorf("failed to get namespace from current context: %v", err)
	}
	return
}
