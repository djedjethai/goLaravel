package main

import (
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"log"
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
		URL:      "git@github.com/djedjethai/goframework-app.git",
		Progress: os.Stdout,
		Depth:    1,
	})

	// remove the .git repository

	// create a ready.go .env file(which should not be in the git repo for secu)

	// create a makefile

	// update the go.mod file

	// update the existing .go files with correct name/imports

	// run go mod tidy in the project directory(to clean up)
}
