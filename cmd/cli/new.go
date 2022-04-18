package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
)

var appURL string

func doNew(appName string) {
	appName = strings.ToLower(appName)
	appURL = appName

	// sanitize the application name ( convert url to single word )
	if strings.Contains(appName, "/") {
		exploded := strings.SplitAfter(appName, "/")
		appName = exploded[len(exploded)-1]
	}

	log.Println("App name : ", appName)

	// git clone the skeleton application
	color.Green("\tCloning repository...")
	_, err := git.PlainClone("./"+appName, false, &git.CloneOptions{
		//URL:      "git@github.com:brucebotes/celeritas-app.git",
		URL:      "https://github.com/brucebotes/celeritas-app.git",
		Progress: os.Stdout,
		Depth:    1,
	})
	if err != nil {
		exitGracefully(err)
	}

	// remove .git directory
	err = os.RemoveAll(fmt.Sprintf("./%s/.git", appName))
	if err != nil {
		exitGracefully(err)
	}

	// create a ready to go .env file
	color.Yellow("\tCreating .env file...")
	data, err := templateFS.ReadFile("templates/env.txt")
	if err != nil {
		exitGracefully(err)
	}

	env := string(data)
	env = strings.ReplaceAll(env, "${APP_NAME}", appName)
	env = strings.ReplaceAll(env, "${KEY}", cel.RandomString(32))

	err = copyDataToFile([]byte(env), fmt.Sprintf("./%s/.env", appName))
	if err != nil {
		exitGracefully(err)
	}

	// create a makefile
	var makeFileName = "Makefile.mac"
	if runtime.GOOS == "windows" {
		makeFileName = "Makefile.windows"
	}
	source, err := os.Open(fmt.Sprintf("./%s/%s", appName, makeFileName))
	if err != nil {
		exitGracefully(err)
	}
	defer source.Close()

	destination, err := os.Create(fmt.Sprintf("./%s/Makefile", appName))
	if err != nil {
		exitGracefully(err)
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		exitGracefully(err)
	}
	_ = os.Remove("./" + appName + "/Makefile.mac")
	_ = os.Remove("./" + appName + "/Makefile.windows")

	// update the go.mod file
	color.Yellow("\tCreating go.mod file...")
	_ = os.Remove("./" + appName + "/go.mod")

	data, err = templateFS.ReadFile("templates/go.mod.txt")
	if err != nil {
		exitGracefully(err)
	}

	mod := string(data)
	mod = strings.ReplaceAll(mod, "${APP_NAME}", appURL)

	err = copyDataToFile([]byte(mod), "./"+appName+"/go.mod")
	if err != nil {
		exitGracefully(err)
	}

	// update existing .go file with the correct name/imports
	color.Yellow("\tUpdating source files...")
	os.Chdir("./" + appName)
	updateSource()

	// run go mod tidy in the project directory
	color.Yellow("\tGet lastest version of celeritas and running go mod tidy...")
	cmd := exec.Command("go", "get", "github.com/brucebotes/celeritas")
	//err = cmd.Start()
	err = cmd.Run()
	if err != nil {
		color.Red("Error doing go get github.com/brucebotes/celeritas")
		exitGracefully(err)
	}

	cmd = exec.Command("go", "mod", "tidy")
	//err = cmd.Start()
	err = cmd.Run()
	if err != nil {
		color.Red("Error doing go mod tidy")
		exitGracefully(err)
	}

	color.Green("Done building " + appURL)
	color.Green("Go build something awesome!")
}
