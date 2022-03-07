package main

// this package will make the color uniform for any os terminal
// go get github.com/fatih/color

import (
	// "fmt"
	"errors"
	"github.com/djedjethai/goframework"
	"github.com/fatih/color"
	// "log"
	"os"
)

const version = "1.0.0"

var gof goframework.Goframework

func main() {
	var message string
	arg1, arg2, arg3, err := validateInput()
	if err != nil {
		exitGracefully(err)
	}

	// populate the gof (the app) var
	// func in folder helper.go
	setup()

	switch arg1 {
	case "help":
		showHelp()
	case "new":
		if arg2 == "" {
			exitGracefully(errors.New("new require an application name"))
		}
		doNew(arg2)
	case "version":
		color.Yellow("Application version: " + version)
	case "migrate":
		if arg2 == "" {
			arg2 = "up"
		}
		err = doMigrate(arg2, arg3)
		if err != nil {
			exitGracefully(err)
		}
		message = "Migration complete"
	case "make":
		if arg2 == "" {
			exitGracefully(errors.New("make requires a subcommand: (migration|model|handler)"))
		}
		err = doMake(arg2, arg3)
		if err != nil {
			exitGracefully(err)
		}
	default:
		showHelp()
	}
	exitGracefully(nil, message)
}

func validateInput() (string, string, string, error) {
	var arg1, arg2, arg3 string

	// means we got at leat 1 arg (we have 2 but the 1st one will be goframework)
	if len(os.Args) > 1 {
		arg1 = os.Args[1]

		// 2 or 3 args
		if len(os.Args) >= 3 {
			arg2 = os.Args[2]
		}

		// if 3 args(which will be position 4)
		if len(os.Args) >= 4 {
			arg3 = os.Args[3]
		}
	} else {
		// use the imported color package
		color.Red("Error: command required")
		showHelp()
		return "", "", "", errors.New("command required")
	}

	return arg1, arg2, arg3, nil
}

// one or more strings, variatic parameters
// means 0 or more strings
func exitGracefully(err error, msg ...string) {
	message := ""
	if len(msg) > 0 {
		message = msg[0]
	}

	if err != nil {
		color.Red("Error: %v\n", err)
	}

	if len(message) > 0 {
		color.Yellow(message)
	} else {
		color.Green("Finished")
	}

	os.Exit(0)
}
