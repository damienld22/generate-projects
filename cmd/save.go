package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
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
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		checkGenDirExistsOrCreateIt()
		pathTemplate := saveTemplate(args[0])
		saveTemplateConfig(pathTemplate)
		fmt.Println("The template has been saved at " + pathTemplate)
	},
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
