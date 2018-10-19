package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available project templates",
	Long:  "List all templates with you can regenerate a project",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("List all templates")
		checkGenDirExistsOrCreateIt()
	},
}
