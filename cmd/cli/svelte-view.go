package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
)

func doSvelteView(modName string) error {
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
	data, err := templateFS.ReadFile("templates/views/svelte.index.jet")
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
	data, err = templateFS.ReadFile("templates/views/npm-packages/svelte.package.json")
	if err != nil {
		exitGracefully(err)
	}
	temp = string(data)
	temp = strings.ReplaceAll(temp, "${MOD_NAME}", modName)
	err = copyDataToFile([]byte(temp), modPath+"/package.json")
	if err != nil {
		exitGracefully(err)
	}
	err = copyFileFromTemplate("templates/views/npm-packages/svelte.tsconfig.json", modPath+"/tsconfig.json")
	if err != nil {
		exitGracefully(err)
	}

	/* temporary disabled
	err = os.Chdir(modPath)
	if err != nil {
		exitGracefully(err)
	}
	if d, err := os.Getwd(); err == nil {
		color.Yellow("\tCurrent working dir: " + d)
	} else {
		exitGracefully(err)
	}
	color.Yellow("\tInstalling node modules. This will take some time...")
	cmd := exec.Command("npm", "install")
	//err = cmd.Start()
	err = cmd.Run()
	if err != nil {
		color.Red("Error executing: npm install")
		exitGracefully(err)
	}
	*/
	color.Yellow("\tPlease chdir to views/" + modName + " and run: npm install")

	// install svelte builder
	data, err = templateFS.ReadFile("templates/views/node-scripts/svelte.esbuild.js")
	if err != nil {
		exitGracefully(err)
	}
	temp = string(data)
	temp = strings.ReplaceAll(temp, "${MOD_NAME}", modName)
	err = copyDataToFile([]byte(temp), modPath+"/esbuild.js")
	if err != nil {
		exitGracefully(err)
	}

	// copy svelte template stubs
	err = os.Mkdir(modPath+"/src", 0755)
	if err != nil {
		exitGracefully(err)
	}
	templates, err := templateFS.ReadDir("templates/views/svelte-scripts/src")
	if err != nil {
		exitGracefully(err)
	}
	for _, f := range templates {
		err = copyFileFromTemplate("templates/views/svelte-scripts/src/"+f.Name(), modPath+"/src/"+f.Name())
		if err != nil {
			exitGracefully(err)
		}
	}

	// create the public folder and contents for the dev environment
	err = os.Mkdir(modPath+"/public", 0755)
	if err != nil {
		exitGracefully(err)
	}
	templates, err = templateFS.ReadDir("templates/views/svelte-scripts/public")
	if err != nil {
		exitGracefully(err)
	}
	for _, f := range templates {
		err = copyFileFromTemplate("templates/views/svelte-scripts/public/"+f.Name(), modPath+"/public/"+f.Name())
		if err != nil {
			exitGracefully(err)
		}
	}

	makeSvelteViewsHandler()

	color.Yellow("\tYou may execute 'npm run dev' after the 'npm install' to test the svelte config")
	color.Yellow("\tPlease remember to add the following to your routes.go file")
	color.Yellow("\ta.get(\"/svh/{module}\", a.Handlers.SvelteViews)")

	return nil
}

func makeSvelteViewsHandler() {
	fileName := cel.RootPath + "/handlers/svelte-view-handler.go"
	if fileExists(fileName) {
		return
	}

	err := copyFileFromTemplate("templates/handlers/svelte.handler.go.txt", cel.RootPath+"/handlers/svelte-view-handler.go")
	if err != nil {
		exitGracefully(err)
	}

	/* incase we want to find and replace
	data, err := templateFS.ReadFile("templates/handlers/svelte.handler.go.txt")
	if err != nil {
		exitGracefully(err)
	}

	handler := string(data)
	handler = strings.ReplaceAll(handler, "$HANDLERNAME$", strcase.ToCamel(arg3))

	err = ioutil.WriteFile(fileName, []byte(handler), 0644)
	if err != nil {
		exitGracefully(err)
	}
	*/

}
