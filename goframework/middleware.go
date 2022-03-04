package goframework

import (
	"net/http"
	"strconv"

	"github.com/justinas/nosurf"
)

func (g *Goframework) SessionLoad(next http.Handler) http.Handler {
	g.InfoLog.Println("sessionLoad called")
	// load and save session on every requests
	return g.Session.LoadAndSave(next)
}

func (g *Goframework) NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	// that will set a cookie on every request
	// so we want it only for production
	// inton c.config.cookie.secure we have a string to true or false indiquant if prod
	// let convert it to bool
	secure, _ := strconv.ParseBool(g.config.cookie.secure)

	// that will set a csrf token on any request
	// so if some api call, the way to desable the cookie is
	csrfHandler.ExemptGlob("/api/*")

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
		Domain:   g.config.cookie.domain,
	})

	return csrfHandler
}
