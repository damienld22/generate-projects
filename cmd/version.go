package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.0.1"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of 'gen'",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}
