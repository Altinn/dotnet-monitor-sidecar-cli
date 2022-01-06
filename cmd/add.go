package cmd

import (
	dmscmd "github.com/altinn/dotnet-monitor-sidecar-cli/pkg/cmd"
	"github.com/altinn/dotnet-monitor-sidecar-cli/pkg/utils"
	"github.com/spf13/cobra"
)

// addCmd represents the dmsctl add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a debug sidecar to your pods",
	Long: `To debug your pods, you can add the debug sidecar to your pods.
Example:
	# Add the debug sidecar to a Deployments pods
	dmsctl add deployment my-deployment
	# Add the debug sidecar to a DaemonSets pods
	dmsctl add daemonset my-daemonset`,
}

// addDeploymentCmd represents the dmsctl add deployment command
var addDeploymentCmd = &cobra.Command{
	Use:   "deployment [name]",
	Short: "Add the debug sidecar to a deployments pods",
	Long: `To debug your Deployment, you can add the debug sidecar to your pods.
After you have added the sidecar, you can port forward to one of the pods with dmsctl port-forward [podname].
Example:
	# Add the debug sidecar to a Deployments pods
	dmsctl add deployment my-deployment`,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: utils.AutoCompleteDeployments,
	Run: func(cmd *cobra.Command, args []string) {
		dmscmd.AddToDeployment(cmd.Context(), kubeconfig, namespace, args[0], containername, debugimage)
	},
}

// addDaemonSetCmd represents the dmsctl add deployment command
var addDaemonSetCmd = &cobra.Command{
	Use:   "daemonset [name]",
	Short: "Add the debug sidecar to a daemonsets pods",
	Long: `To debug your Daemonset, you can add the debug sidecar to your pods.
After you have added the sidecar, you can port forward to one of the pods with dmsctl port-forward [podname].
Example:
	# Add the debug sidecar to a Deployments pods
	dmsctl add daemonset my-daemonset`,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: utils.AutoCompleteDaemonSets,
	Run: func(cmd *cobra.Command, args []string) {
		dmscmd.AddToDaemonset(cmd.Context(), kubeconfig, namespace, args[0], containername, debugimage)
	},
}

var (
	containername     string
	debugimage        string
	defaultDebugImage string = "mcr.microsoft.com/dotnet/monitor:6.0"
)

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.AddCommand(addDeploymentCmd)
	addDeploymentCmd.Flags().StringVarP(&containername, "container", "c", "", "Supply container name if deployment contains multiple pods")
	addDeploymentCmd.Flags().StringVar(&debugimage, "debugimage", defaultDebugImage, "image to add as a debug sidecar")

	addCmd.AddCommand(addDaemonSetCmd)
	addDaemonSetCmd.Flags().StringVarP(&containername, "container", "c", "", "Supply container name if deployment contains multiple pods")
	addDaemonSetCmd.Flags().StringVar(&debugimage, "debugimage", defaultDebugImage, "image to add as a debug sidecar")
}
