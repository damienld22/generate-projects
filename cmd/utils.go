package cmd

import (
	"os"
	"os/user"
	"path/filepath"
)

// CONSTANTS
const nameTemplatesDirectory = ".gen"
const nameTemplatesConfigFile = "templates.json"

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
