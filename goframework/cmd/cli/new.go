package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"io"
	"os/exec"
	"runtime"

	"os"
	"strings"
)

var appURL string

func doNew(appName string) {
	appName = strings.ToLower(appName)
	appURL = appName

	// sanitize the app name, convert any potentiel url to single word
	if strings.Contains(appName, "/") {
		exploded := strings.SplitAfter(appName, "/")
		// like if we clone a repo from github url, we will get only the project name
		appName = exploded[len(exploded)-1]
	}
	// log.Println("app name is: ", appName)

	// git clone the skeleton application
	color.Green("\tCloning Repository...")
	_, err := git.PlainClone("./"+appName, false, &git.CloneOptions{
		URL:      "https://github.com/djedjethai/goframework-app.git",
		Progress: os.Stdout,
		Depth:    1,
	})
	if err != nil {
		exitGracefully(err)
	}

	// remove the .git repository
	err = os.RemoveAll(fmt.Sprintf("./%s/.git", appName))
	if err != nil {
		exitGracefully(err)
	}

	// create a ready.go .env file(which should not be in the git repo for secu)
	color.Yellow("\tCreating .env file...")
	data, err := templateFS.ReadFile("templates/env.txt")
	if err != nil {
		exitGracefully(err)
	}
	env := string(data)
	env = strings.ReplaceAll(env, "${APP_NAME}", appName)
	env = strings.ReplaceAll(env, "${KEY}", gof.RandomString(32))

	err = copyDataToFile([]byte(env), fmt.Sprintf("./%s/.env", appName))
	if err != nil {
		exitGracefully(err)
	}

	// create a makefile
	if runtime.GOOS == "windows" {
		source, err := os.Open(fmt.Sprintf("./%s/Makefile.windows", appName))
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

	} else {
		source, err := os.Open(fmt.Sprintf("./%s/Makefile.mac", appName))
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

	}
	//  after setted the Makefile for mac or windows(depend of the user os)
	// delete the Makefile templates
	_ = os.Remove("./" + appName + "/Makefile.windows")
	_ = os.Remove("./" + appName + "/Makefile.mac")

	// update the go.mod file
	color.Yellow("\tCreating go.mod")
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

	// update the existing .go files with correct name/imports
	color.Yellow("\tUpdating source files...")
	// move a step above the dir where the app is
	os.Chdir("./" + appName)
	updateSource()

	// run go mod tidy in the project directory(to clean up)
	color.Yellow("\tRunning go mod tidy...")
	cmd := exec.Command("go", "mod", "tidy")
	err = cmd.Start()
	if err != nil {
		exitGracefully(err)
	}

	color.Green("Done Building " + appURL)
	color.Green("Go build something awesome")

}
