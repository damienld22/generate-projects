package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
)

var targetPath string
var nameDir string
var openWithCode bool

func init() {
	createCmd.Flags().StringVarP(&targetPath, "path", "p", "", "Optional path of the generated project (Default is the current directory")
	createCmd.Flags().StringVarP(&nameDir, "directory", "d", "", "Optional name of the generated directory project")
	createCmd.Flags().BoolVarP(&openWithCode, "open", "o", false, "Open VSCode inside the generated project")
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a project",
	Long:  "Create a project from a saved template",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		checkGenDirExistsOrCreateIt()
		createProjectFromTemplate(args[0])
	},
}

/**
 *	Create the project from the specified template
 */
func createProjectFromTemplate(template string) {
	pathTemplate := getPathTemplateFromName(template)

	if pathTemplate != "" {
		if targetPath == "" {
			current, err := os.Executable()
			checkError(err)
			targetPath = filepath.Dir(current)
		}

		if nameDir == "" {
			nameDir = template
		}

		projectDir := filepath.Join(targetPath, nameDir)
		os.Mkdir(projectDir, os.ModePerm)
		copy.Copy(pathTemplate, projectDir)

		// Open VSCode if necessary
		if openWithCode {
			_, err := exec.Command("code", projectDir).Output()
			checkError(err)
		}
	}
}

/**
 * Get the path of the template from his name
 */
func getPathTemplateFromName(templateName string) string {
	// Get configuration of templates
	templatesConf, err := os.Open(getTemplatesConfigurationFile())
	checkError(err)
	defer templatesConf.Close()

	// Parsing the JSON file
	var templates Templates
	byteValue, _ := ioutil.ReadAll(templatesConf)
	json.Unmarshal(byteValue, &templates)

	var path string
	for _, elt := range templates.Templates {
		if elt.Name == templateName {
			path = elt.Path
			break
		}
	}

	if path == "" {
		fmt.Println("The template %s doesn't exist", templateName)
	}

	return path
}
