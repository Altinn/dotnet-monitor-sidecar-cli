package kubernetes

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	"github.com/altinn/dotnet-monitor-sidecar-cli/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// PortForward runs port-forward to the given pod and waits for Interrupt
func (h *Helper) PortForward(ctx context.Context, namespace string, podname string) error {
	pod, err := h.Client.CoreV1().Pods(namespace).Get(ctx, podname, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if pod.Status.Phase != corev1.PodRunning {
		return fmt.Errorf("unable to forward port because pod is not running. Current status=%v", pod.Status.Phase)
	}
	_, err = h.GetDDPodApplyInfo(ctx, namespace, podname)
	if err != nil {
		if errors.IsNotPresent(err) {
			return fmt.Errorf("debug sidecar not attached to pod %s", podname)
		}
	}

	f, config, err := getRestSetup()
	if err != nil {
		return err
	}

	restClient, err := f.RESTClient()
	if err != nil {
		return err
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	defer signal.Stop(signals)

	stopCh := make(chan struct{}, 1)
	readyCh := make(chan struct{})

	go func() {
		<-signals
		if stopCh != nil {
			close(stopCh)
		}
	}()

	req := restClient.Post().
		Resource("pods").
		Namespace(pod.Namespace).
		Name(pod.Name).
		SubResource("portforward")

	return forwardPorts("POST", req.URL(), config, stopCh, readyCh)

}

func getRestSetup() (cmdutil.Factory, *rest.Config, error) {
	kubeConfigFlags := genericclioptions.NewConfigFlags(true)
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)
	f := cmdutil.NewFactory(matchVersionKubeConfigFlags)
	config, err := f.ToRESTConfig()
	config.GroupVersion = &schema.GroupVersion{Group: "", Version: "v1"}
	return f, config, err
}

func forwardPorts(method string, url *url.URL, config *rest.Config, stop, ready chan struct{}) error {
	transport, upgrader, err := spdy.RoundTripperFor(config)
	if err != nil {
		return err
	}
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, method, url)
	fw, err := portforward.NewOnAddresses(dialer, []string{"localhost"}, []string{"52323:52323"}, stop, ready, os.Stdout, os.Stderr)
	if err != nil {
		return err
	}
	return fw.ForwardPorts()
}
