package main

// ss -lptn 'sport = :4000'
// sudo ss -tulwn | grep LISTEN

// get ENV var into a specifique process
// cat /proc/42979/environ | tr '\0' '\n'

import (
	"github.com/djedjethai/goframework"
	"myapp/data"
	"myapp/handlers"
	"myapp/middleware"
)

type application struct {
	App        *goframework.Goframework
	Handlers   *handlers.Handlers
	Models     data.Models
	Middleware *middleware.Middleware
}

func main() {
	g := initApplication()
	g.App.ListenAndServe()
}
