package main

import (
	// "fmt"
	"github.com/djedjethai/goframework"
	"log"
	"myapp/data"
	"myapp/handlers"
	"myapp/middleware"
	"os"
)

func initApplication() *application {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// init goframework
	gofram := &goframework.Goframework{}
	err = gofram.New(path)
	if err != nil {
		log.Fatal(err)
	}

	gofram.AppName = "myapp"

	myMiddleware := &middleware.Middleware{
		App: gofram,
	}

	myHandlers := &handlers.Handlers{
		App: gofram,
	}

	app := &application{
		App:        gofram,
		Handlers:   myHandlers,
		Middleware: myMiddleware,
	}

	app.App.Routes = app.routes()

	app.Models = data.New(app.App.DB.Pool)
	myHandlers.Models = app.Models
	app.Middleware.Models = app.Models

	return app
}
