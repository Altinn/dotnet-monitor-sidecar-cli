package cmd

import (
	dmscmd "github.com/altinn/dotnet-monitor-sidecar-cli/pkg/cmd"
	"github.com/spf13/cobra"
)

// portforwardCmd represents the dmsctl port-forward command
var portforwardCmd = &cobra.Command{
	Use:   "port-forward [podname]",
	Short: "Forward port 52323 from your local machine to port 52323 in a pod",
	Long: `The debug image does not expose its endpoint out of the pod.
This command will forward port 52323 from your local machine to port 52323 in a pod.
Example:
	# Forward port 52323 from your local machine to port 52323 in the pod my-pod
	dmsctl port-forward my-pod`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dmscmd.ForwardPort(cmd.Context(), kubeconfig, namespace, args[0])
	},
}

func init() {
	rootCmd.AddCommand(portforwardCmd)
}
