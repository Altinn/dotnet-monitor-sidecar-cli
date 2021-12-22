package cmd

import (
	dmscmd "github.com/altinn/dotnet-monitor-sidecar-cli/pkg/cmd"
	"github.com/spf13/cobra"
)

// removeCmd represents the dmsctl remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove debug sidecar from your pods",
	Long: `After you are done debugging, you can remove the debug sidecar from your pods.
Example:
	# Remove the debug sidecar from a Deployments pods
	dmsctl remove deployment my-deployment
	# Add the debug sidecar to a DaemonSets pods
	dmsctl remove daemonset my-daemonset`,
}

// removeDeploymentCmd represents the dmsctl remove deployment command
var removeDeploymentCmd = &cobra.Command{
	Use:   "deployment [name]",
	Short: "Remove a debug sidecar to a Deployment.apps",
	Long: `After you are done debugging, you can remove the debug sidecar from your pods.
Example:
	# Remove the debug sidecar from a Deployments pods
	dmsctl remove deployment my-deployment`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dmscmd.RemoveFromDeployment(cmd.Context(), kubeconfig, namespace, args[0])
	},
}

// removeDaemonSetCmd represents the dmsctl remove deployment command
var removeDaemonSetCmd = &cobra.Command{
	Use:   "daemonset [name]",
	Short: "Remove a debug sidecar to a Daemonset.apps",
	Long: `After you are done debugging, you can remove the debug sidecar from your pods.
Example:
	# Remove the debug sidecar from a DaemonSets pods
	dmsctl remove daemonset my-daemonset`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dmscmd.RemoveFromDaemonset(cmd.Context(), kubeconfig, namespace, args[0])
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.AddCommand(removeDeploymentCmd)

	removeCmd.AddCommand(removeDaemonSetCmd)
}
