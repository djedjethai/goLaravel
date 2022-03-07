package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"

	// "log"
	"os"
	"strings"
)

func doNew(appName string) {
	appName = strings.ToLower(appName)

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

	// update the go.mod file

	// update the existing .go files with correct name/imports

	// run go mod tidy in the project directory(to clean up)
}
