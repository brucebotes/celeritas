package celeritas

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

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

func (c *Celeritas) CheckForMaintenanceMode(next http.Handler) http.Handler {
	var allowedUrls = strings.Split(os.Getenv("ALLOWED_URLS"), ",")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maintenanceMode {
			//if !strings.Contains(r.URL.Path, "/public/maintenance.html") {
			if !stringInSlice(r.URL.Path, allowedUrls) {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Header().Set("Retry-After:", "300")
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
				http.ServeFile(w, r, fmt.Sprintf("%s/public/maintenance.html", c.RootPath))
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// stringInSlice returns true if needle exists in haystack
func stringInSlice(needle string, haystack []string) bool {
	for _, straw := range haystack {
		if strings.ToLower(straw) == strings.ToLower(needle) {
			return true
		}
	}
	return false
}
