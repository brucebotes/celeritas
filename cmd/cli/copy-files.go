package main

import (
	"embed"
	"errors"
	"io/ioutil"
	"os"
)

//go:embed templates
var templateFS embed.FS

func copyFileFromTemplate(templatePath, targetFile string) error {
	if fileExists(targetFile) {
		return errors.New(targetFile + " already exists!")
	}

	data, err := templateFS.ReadFile(templatePath)
	if err != nil {
		exitGracefully(err)
	}

	err = copyDataToFile(data, targetFile)
	if err != nil {
		exitGracefully(err)
	}

	return nil
}

func copyDataToFile(data []byte, to string) error {
	err := ioutil.WriteFile(to, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func fileExists(fileToCheck string) bool {
	if _, err := os.Stat(fileToCheck); os.IsNotExist(err) {
		return false
	}
	return true
}

func copyFromTemplateFolderToDestinationFolder(templatePath, destPath string) error {
	templates, err := templateFS.ReadDir(templatePath)
	if err != nil {
		return err
	}
	for _, f := range templates {
		if !f.IsDir() {
			err := copyFileFromTemplate(templatePath+"/"+f.Name(), destPath+"/"+f.Name())
			if err != nil {
				return err
			}
		}
	}
	return nil
}
