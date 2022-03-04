package main

import (
	"net/http"
)

// helper func to just call a.get(s, h)
// instead of a.App.Routes.Get(s, h)
func (a *application) get(s string, h http.HandlerFunc) {
	a.App.Routes.Get(s, h)
}

func (a *application) post(s string, h http.HandlerFunc) {
	a.App.Routes.Post(s, h)
}

// same signature as the chi package(for the argument)
func (a *application) use(m ...func(http.Handler) http.Handler) {
	a.App.Routes.Use(m...)
}
