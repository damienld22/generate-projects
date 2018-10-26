package cmd

import (
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
	templateName := result

	// Open with VSCode
	promptOpenVsCode := promptui.Prompt{
		Label:     "Open with VSCode : ",
		IsConfirm: true,
	}

	_, err = promptOpenVsCode.Run()
	if err != nil {
		openWithCode = false
	} else {
		openWithCode = true
	}

	// Select the target workspace
	promptSelectWorkspace := promptui.Select{
		Label: "Select the workspace",
		Items: getListOfWorkspaces(),
	}

	_, result, err = promptSelectWorkspace.Run()
	checkError(err)

	if result == CustomWorkspace {
		promptCustomWorkspace := promptui.Prompt{
			Label: "Custom workspace",
		}

		result, err = promptCustomWorkspace.Run()
		checkError(err)
	}
	targetPath = result

	// Select the name of the directory
	promptNameDirectory := promptui.Prompt{
		Label:   "Name of the project",
		Default: templateName,
	}

	result, err = promptNameDirectory.Run()
	checkError(err)
	nameDir = result

	// Create the project
	createProjectFromTemplate(templateName)
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
