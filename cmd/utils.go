package cmd

import (
	"os"
	"os/user"
	"path/filepath"
)

const nameTemplatesDirectory = ".gen"
const nameTemplatesConfigFile = "templates.json"

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func checkGenDirExistsOrCreateIt() {
	currentUser, _ := user.Current()
	pathGenDir := filepath.Join(currentUser.HomeDir, nameTemplatesDirectory)

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
