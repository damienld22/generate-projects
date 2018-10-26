package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/otiai10/copy"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/manifoldco/promptui"

)

var name string
var description string

func init() {
	saveCmd.Flags().StringVarP(&name, "name", "n", "", "Optional name of the template")
	saveCmd.Flags().StringVarP(&description, "description", "d", "", "Optional description of the template")
	rootCmd.AddCommand(saveCmd)
}

var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save your template project",
	Long:  "Give the path of your template and we save it to regenerate it later",
	Run: func(cmd *cobra.Command, args []string) {
		checkGenDirExistsOrCreateIt()

		// Go to interactive mode or not
		if len(args) == 0 {
			interactiveCliSave()
		} else {
			pathTemplate := saveTemplate(args[0])
			saveTemplateConfig(pathTemplate)
			log.Info("The template has been saved at ", pathTemplate)
		}

	},
}

/**
 *	Save a template using the interactive mode
 */
func interactiveCliSave() {
	// Select the path of the template
	promptTemplatePath := promptui.Prompt{
		Label: "Path of the template",
	}
	result, err := promptTemplatePath.Run()
	checkError(err)

	template := result

	// Select the name of the template
	promptTemplateName := promptui.Prompt{
		Label: "Name of the template",
	}
	result, err = promptTemplateName.Run()
	checkError(err)

	name = result

	// Select the description of the template
	promptTemplateDescription := promptui.Prompt{
		Label: "Description of the template",
	}
	result, err = promptTemplateDescription.Run()
	checkError(err)
	description = result

	pathTemplate := saveTemplate(template)
	saveTemplateConfig(pathTemplate)

	log.Info("The template has been saved at ", pathTemplate)
}

/**
 *	Save the template configuration
 *	in the templates configuration file
 */
func saveTemplateConfig(path string) {

	path, err := filepath.Abs(path)
	checkError(err)

	if name == "" {
		name = filepath.Base(path)
	}

	if description == "" {
		description = filepath.Base(path)
	}

	template := Template{
		Path:        path,
		Name:        name,
		Description: description,
	}

	// Get configuration of templates
	templatesConf, err := os.Open(getTemplatesConfigurationFile())
	checkError(err)

	// Parsing the existing configuration file
	var templates Templates
	byteValue, err := ioutil.ReadAll(templatesConf)
	checkError(err)
	json.Unmarshal(byteValue, &templates)
	templatesConf.Close()

	templates.Templates = append(templates.Templates, template)

	// Rewrite to JSON configuration file
	jsonConf, err := json.Marshal(templates)
	checkError(err)
	err = ioutil.WriteFile(getTemplatesConfigurationFile(), jsonConf, os.ModePerm)
	checkError(err)
}

/**
 *	Save the template in the
 *	proper directory
 */
func saveTemplate(path string) string {

	// Get the absolute path
	path, err := filepath.Abs(path)
	checkError(err)

	// Copy the path to de templates folder
	targetPath := filepath.Join(getTemplatesDirectory(), filepath.Base(path))
	copy.Copy(path, targetPath)

	return targetPath
}
