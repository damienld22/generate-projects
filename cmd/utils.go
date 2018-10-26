package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	log "github.com/sirupsen/logrus"
)

// CONSTANTS
const nameTemplatesDirectory = ".gen"
const nameTemplatesConfigFile = "templates.json"
const CustomWorkspace = "Custom"

// TYPES
type Templates struct {
	Templates []Template `json:"templates"`
}

type Template struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
}

// HANDLE ERRORS
func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

// CONFIGURATION TEMPLATES
func getTemplatesDirectory() string {
	currentUser, _ := user.Current()
	return filepath.Join(currentUser.HomeDir, nameTemplatesDirectory)
}

func getTemplatesConfigurationFile() string {
	return filepath.Join(getTemplatesDirectory(), nameTemplatesConfigFile)
}

func checkGenDirExistsOrCreateIt() {
	pathGenDir := getTemplatesDirectory()

	if _, err := os.Stat(pathGenDir); os.IsNotExist(err) {
		// Create the $HOME/.gen directory
		err := os.Mkdir(pathGenDir, os.ModePerm)
		checkError(err)

		file, err := os.Create(filepath.Join(pathGenDir, nameTemplatesConfigFile))
		checkError(err)
		defer file.Close()

		// Write an empty JSON array
		file.WriteString("[]")
		file.Sync()
	}
}

func getListOfWorkspaces() []string {
	currentUser, _ := user.Current()
	homeDir := currentUser.HomeDir

	listWorkspaces := make([]string, 0)

	listWorkspaces = append(listWorkspaces, filepath.Join(homeDir, "Bureau"))
	listWorkspaces = append(listWorkspaces, filepath.Join(homeDir, "workspace/sandbox"))
	listWorkspaces = append(listWorkspaces, CustomWorkspace)

	return listWorkspaces
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