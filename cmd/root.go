package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	namespace  string
	kubeconfig string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dmsctl",
	Short: "CLI to add, remove and connect to dotnet-moniter sidecar in kubernetes",
	Long: `Small cli tool to ease the process of adding, removing and connecting to dotnet-monitor container to debug a dotnet application running in kubernetes.
Examples:
	# Add sidecars to the pods assosiated with a deployment in kubernetes
	dmsctl add deployment my-deployment

	# Add sidecars to the pods assosiated with a daemonset in kubernetes
	dmsctl add daemonset my-daemonset

	# Setup a port forward from your machine to the dotnet-monitor sidecar container in the pod
	dmsctl port-forward my-pod-13fa7

	# Remove sidecars from the pods assosiated with a deployment in kubernetes
	dmsctl remove deployment my-deployment

	# Remove sidecars from the pods assosiated with a daemonset in kubernetes
	dmsctl remove daemonset my-daemonset 
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dmsconfig.yaml)")
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "If present, the namespace scope for this CLI request. Otherwise, the current namespace is used.")
	rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "Override path to the kubeconfig file to use for CLI requests.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".dmsconfig" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".dmsconfig")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// NewDmsctlCommand returns rootCmd. Used to generate docs
func NewDmsctlCommand(ctx context.Context) *cobra.Command {
	return rootCmd
}
