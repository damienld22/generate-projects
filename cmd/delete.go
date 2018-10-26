package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a template",
	Long:  "Delete a template from his name",
	Run: func(cmd *cobra.Command, args []string) {
		checkGenDirExistsOrCreateIt()

		// Go to interactive mode or not
		if len(args) == 0 {
			interactiveCliDelete()
		} else {
			deleteTemplate(args[0])
		}

	},
}

/**
 * Interactive mode to delete a template
 */
func interactiveCliDelete() {
	// Select the template
	availableTemplates := getAllAvailableTemplates()
	itemsTemplates := make([]string, 0)
	for _, elt := range availableTemplates {
		itemsTemplates = append(itemsTemplates, elt.Name)
	}

	promptSelectTemplate := promptui.Select{
		Label: "Select the template",
		Items: itemsTemplates,
	}

	_, result, err := promptSelectTemplate.Run()
	checkError(err)
	deleteTemplate(result)
}

/**
 *	Delete the template from his name
 */
func deleteTemplate(template string) {
	pathTemplate := getPathTemplateFromName(template)
	// Delete the directory
	err := os.RemoveAll(pathTemplate)
	checkError(err)
	
	// Delete the configuration
	templatesConf, err := os.Open(getTemplatesConfigurationFile())
	checkError(err)
	defer templatesConf.Close()

	// Parsing the JSON file
	var templates Templates
	byteValue, _ := ioutil.ReadAll(templatesConf)
	json.Unmarshal(byteValue, &templates)

	newSlice := []Template{}
	for _, elt := range templates.Templates {
		if elt.Name != template {
			newSlice = append(newSlice, elt)
			break
		}
	}

	// Rewrite to JSON configuration file
	templates.Templates = newSlice
	jsonConf, err := json.Marshal(templates)
	checkError(err)
	err = ioutil.WriteFile(getTemplatesConfigurationFile(), jsonConf, os.ModePerm)
	checkError(err)

	log.Info("Template deleted")
}
