package celeritas

import (
	"net/http"
	"strconv"

	"github.com/justinas/nosurf"
)

func (c *Celeritas) SessionLoad(next http.Handler) http.Handler {
	c.InfoLog.Println("SessionLoad called")
	return c.Session.LoadAndSave(next) // load and save our session on every request
}

func (c *Celeritas) NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	secure, _ := strconv.ParseBool(c.config.cookie.secure)

	// If you would like to exclude an endpoint from
	// being cheched for a cookie (in a form post
	// before hitting the route endpoint) like
	// in our api routes
	// - use ExemptGlob
	csrfHandler.ExemptGlob("/api/*")

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
		Domain:   c.config.cookie.domain,
	})

	return csrfHandler
}
