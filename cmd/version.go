package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionString represents the version of the application
var versionString = "unset"

// versionCmd represents the dmsctl version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the cli version",
	Long:  `Print the cli version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("CLI Version: %s\n", versionString)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
