package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/manifoldco/promptui"
	"github.com/otiai10/copy"
	log "github.com/sirupsen/logrus"
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
	Run: func(cmd *cobra.Command, args []string) {
		checkGenDirExistsOrCreateIt()

		// Go to interactive mode or not
		if len(args) == 0 {
			interactiveCliInput()
		} else {
			createProjectFromTemplate(args[0])
		}

	},
}

/**
 * Interactive mode to create a project from a templates
 */
func interactiveCliInput() {
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

	// Open with VSCode
	promptOpenVsCode := promptui.Prompt{
		Label:     "Open with VSCode",
		IsConfirm: true,
	}

	_, err = promptOpenVsCode.Run()
	if err != nil {
		openWithCode = false
	} else {
		openWithCode = true
	}

	createProjectFromTemplate(result)
}

/**
 *	Create the project from the specified template
 */
func createProjectFromTemplate(template string) {
	pathTemplate := getPathTemplateFromName(template)

	if targetPath == "" {
		current, err := os.Executable()
		checkError(err)
		targetPath = filepath.Dir(current)
	}

	if nameDir == "" {
		nameDir = template
	}

	projectDir := filepath.Join(targetPath, nameDir)
	log.Debug("Path of the generated project : " + projectDir)

	os.Mkdir(projectDir, os.ModePerm)
	copy.Copy(pathTemplate, projectDir)
	log.Info("Project successfully generated at " + projectDir)

	// Open VSCode if necessary
	if openWithCode {
		_, err := exec.Command("code", projectDir).Output()
		checkError(err)
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
		log.Error("The template %s is not available", templateName)
		panic(1)
	}

	log.Debug("Template specified : " + templateName)
	log.Debug("Path of the template : " + path)

	return path
}
