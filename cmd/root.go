package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gen",
	Short: "Gen is a basic generator of custom projects",
	Long:  "Save your project templates and Gen can recreate them.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gen generator")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
