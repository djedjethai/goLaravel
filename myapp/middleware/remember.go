package middleware

import (
	"fmt"
	"myapp/data"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (m *Middleware) CheckRemember(next http.Handler) http.Handler {
	// first check if user is already checked in via the session(userID)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.App.Session.Exists(r.Context(), "userID") {
			// not logged in
			// check for existance of a cookie(remember cookie)
			// which will determine if the user is a valid user
			cookie, err := r.Cookie(fmt.Sprintf("_%s_remember", m.App.AppName))
			if err != nil {
				// there is no cookie, so go to the next middleware
				next.ServeHTTP(w, r)
			} else {
				// we found a cookie so check it
				key := cookie.Value
				var u data.User
				if len(key) > 0 {
					// cookie have some datas
					split := strings.Split(key, "|")
					uid, hash := split[0], split[1]
					id, _ := strconv.Atoi(uid)
					validHash := u.CheckForRememberToken(id, hash)
					if !validHash {
						// there is a cookie, but not the good one
						m.deleteRememberCookie(w, r)
						m.App.Session.Put(r.Context(), "error", "You've been logged out from another device")
						next.ServeHTTP(w, r)
					} else {
						// valid hash, so log the user in,
						// that's the automatic logging
						user, _ := u.Get(id)
						m.App.Session.Put(r.Context(), "userID", user.ID)
						m.App.Session.Put(r.Context(), "remember_token", hash)
						next.ServeHTTP(w, r)
					}
				} else {
					// key length is zero, so it's probably a leftover cookie
					// often happend when user has not closed browser
					m.deleteRememberCookie(w, r)
					next.ServeHTTP(w, r)
				}
			}

		} else {
			// logged in

			next.ServeHTTP(w, r)
		}
	})
}

func (m *Middleware) deleteRememberCookie(w http.ResponseWriter, r *http.Request) {
	// first renew the token we have in the session (good practice)
	_ = m.App.Session.RenewToken(r.Context())

	// delete the cookie, and the way to do it is to set a new "expired" one
	newCookie := http.Cookie{
		Name:     fmt.Sprintf("_%s_remember", m.App.AppName),
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-100 * time.Hour),
		HttpOnly: true,
		Domain:   m.App.Session.Cookie.Domain,
		MaxAge:   -1,
		Secure:   m.App.Session.Cookie.Secure,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, &newCookie)

	// log the user out
	m.App.Session.Remove(r.Context(), "userID")
	m.App.Session.Destroy(r.Context())

	// renew the token one more time, not necessary but better for security
	_ = m.App.Session.RenewToken(r.Context())
}
