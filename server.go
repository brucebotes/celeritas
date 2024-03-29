package celeritas

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

// ListenAndServe starts the web server
func (c *Celeritas) ListenAndServe() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		ErrorLog:     c.ErrorLog,
		Handler:      c.Routes,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
		TLSConfig:    c.Server.TLSConfig,
	}

	if c.DB.Pool != nil {
		defer c.DB.Pool.Close()
	}

	if redisPool != nil {
		defer redisPool.Close()
	}

	if badgerConn != nil {
		defer badgerConn.Close()
	}

	go c.listenRPC()

	c.InfoLog.Printf("Listening on port %s", os.Getenv("PORT"))
	if c.Server.Secure != true || c.Server.TLSConfig == nil {
		return srv.ListenAndServe()
	} else {
		return srv.ListenAndServeTLS("","")
	}
}
