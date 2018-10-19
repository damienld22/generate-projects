package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(saveCmd)
}

var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save your template project",
	Long:  "Give the path of your template and we save it to regenerate it later",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args[0])
		checkGenDirExistsOrCreateIt()
	},
}
