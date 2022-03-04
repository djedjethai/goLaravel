package goframework

import (
	// "fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (g *Goframework) routes() http.Handler {
	mux := chi.NewRouter()
	// middleware that injects a req id into the context of each req
	mux.Use(middleware.RequestID)
	// will give us the ip address of our visitors
	mux.Use(middleware.RealIP)
	if g.Debug {
		mux.Use(middleware.Logger)
	}
	// recover from Panic
	mux.Use(middleware.Recoverer)

	// use the session we setted up in the session package/and gofram.Init()
	mux.Use(g.SessionLoad)

	// use the csrf middleware
	mux.Use(g.NoSurf)

	return mux
}
