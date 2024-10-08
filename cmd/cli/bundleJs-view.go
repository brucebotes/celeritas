package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
)

func doBundleJSView(modName string) error {
	modName = strings.ToLower(modName)
	modPath := cel.RootPath + "/views/" + modName

	// sanitize the application name ( convert url to single word )
	if strings.Contains(modName, "/") {
		exploded := strings.SplitAfter(modName, "/")
		modName = exploded[len(exploded)-1]
	}
	appURL = modName

	log.Println("Module name : ", modName)

	// create folder in views sub directory
	if fileExists(modPath) {
		exitGracefully(errors.New(fmt.Sprintf("the %s template folder already exists in ./views", modName)))
	}
	err := os.Mkdir(modPath, 0755)
	if err != nil {
		exitGracefully(err)
	}

	// copy jet template
	color.Yellow("\tCreating index.jet file...")
	data, err := templateFS.ReadFile("templates/views/bundleJS/index.jet")
	if err != nil {
		exitGracefully(err)
	}
	temp := string(data)
	temp = strings.ReplaceAll(temp, "${MOD_NAME}", modName)
	err = copyDataToFile([]byte(temp), modPath+"/index.jet")
	if err != nil {
		exitGracefully(err)
	}

	// create package.json and run npm install
	data, err = templateFS.ReadFile("templates/views/bundleJS/npm-packages/package.json")
	if err != nil {
		exitGracefully(err)
	}
	temp = string(data)
	temp = strings.ReplaceAll(temp, "${MOD_NAME}", modName)
	err = copyDataToFile([]byte(temp), modPath+"/package.json")
	if err != nil {
		exitGracefully(err)
	}
	err = copyFileFromTemplate("templates/views/bundleJS/npm-packages/tsconfig.json", modPath+"/tsconfig.json")
	if err != nil {
		exitGracefully(err)
	}

	color.Yellow("\tPlease chdir to views/" + modName + " and run: npm install")

	// copy javascript template stubs
	err = os.Mkdir(modPath+"/src", 0755)
	if err != nil {
		exitGracefully(err)
	}
	templates, err := templateFS.ReadDir("templates/views/bundleJS/source/src")
	if err != nil {
		exitGracefully(err)
	}
	for _, f := range templates {
		if !f.IsDir() {
			data, err = templateFS.ReadFile("templates/views/bundleJS/source/src/" + f.Name())
			if err != nil {
				exitGracefully(err)
			}
			temp = string(data)
			temp = strings.ReplaceAll(temp, "${MOD_NAME}", modName)
			err = copyDataToFile([]byte(temp), modPath+"/src/"+f.Name())
			if err != nil {
				exitGracefully(err)
			}
		}
	}
	// create the src/components folder and contents
	err = os.Mkdir(modPath+"/src/components", 0755)
	if err != nil {
		exitGracefully(err)
	}
	err = copyFromTemplateFolderToDestinationFolder("templates/views/bundleJS/source/src/components", modPath+"/src/components")
	if err != nil {
		exitGracefully(err)
	}
	// create the src/components/timer folder and contents
	err = os.Mkdir(modPath+"/src/components/timer", 0755)
	if err != nil {
		exitGracefully(err)
	}
	err = copyFromTemplateFolderToDestinationFolder("templates/views/bundleJS/source/src/components/timer", modPath+"/src/components/timer")
	if err != nil {
		exitGracefully(err)
	}
	// create the src/components/controller folder and contents
	err = os.Mkdir(modPath+"/src/components/controller", 0755)
	if err != nil {
		exitGracefully(err)
	}
	err = copyFromTemplateFolderToDestinationFolder("templates/views/bundleJS/source/src/components/controller", modPath+"/src/components/controller")
	if err != nil {
		exitGracefully(err)
	}
	// create the src/pages folder and contents
	err = os.Mkdir(modPath+"/src/pages", 0755)
	if err != nil {
		exitGracefully(err)
	}
	err = copyFromTemplateFolderToDestinationFolder("templates/views/bundleJS/source/src/pages", modPath+"/src/pages")
	if err != nil {
		exitGracefully(err)
	}

	makeJSBundleViewsHandler()

	color.Yellow("\tPlease remember to add the following to your routes.go file")
	color.Yellow("\ta.get(\"/jsb/{module}\", a.Handlers.JSBundleViews)")
	color.Yellow("\ta.get(\"/jsb/{module}/*\", a.Handlers.JSBundleViews)")

	return nil
}

func makeJSBundleViewsHandler() {
	fileName := cel.RootPath + "/handlers/js-bundle-view-handler.go"
	if fileExists(fileName) {
		return
	}

	err := copyFileFromTemplate("templates/handlers/bundleJS.handler.go.txt", fileName)
	if err != nil {
		exitGracefully(err)
	}
}
