package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
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
		// Init the command
		checkGenDirExistsOrCreateIt()

		templates := getAllAvailableTemplates()

		// Display the list of available templates
		if len(templates) > 0 {
			for i := 0; i < len(templates); i++ {
				fmt.Println("======================================")
				fmt.Println()
				fmt.Println("Name : " + templates[i].Name)
				fmt.Println("Description : " + templates[i].Description)
				fmt.Println()
			}
			fmt.Println("======================================")
		} else {
			fmt.Println("No templates are available")
		}

	},
}

/**
 *	Get the list of all available templates
 */
func getAllAvailableTemplates() []Template {
	// Get configuration of templates
	templatesConf, err := os.Open(getTemplatesConfigurationFile())
	checkError(err)
	defer templatesConf.Close()

	// Parsing the JSON file
	var templates Templates
	byteValue, _ := ioutil.ReadAll(templatesConf)
	json.Unmarshal(byteValue, &templates)
	log.Debug("Templates are recovered")

	return templates.Templates
}
