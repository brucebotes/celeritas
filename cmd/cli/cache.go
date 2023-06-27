package main

import (
	"github.com/fatih/color"
)

func doCacheApi() error {
	err := copyFileFromTemplate("templates/routes/api-routes.go.txt", cel.RootPath+"/routes-api.go")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/handlers/cache-handlers.go.txt", cel.RootPath+"/handlers/cache-handlers.go")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/views/templates/cache.jet", cel.RootPath+"/views/templates/cache.jet")
	if err != nil {
		exitGracefully(err)
	}

	color.Yellow("\tCache handlers, routes/api and cache.jet template created")
	color.Yellow("\tConfigure the CACHE (and REDIS connection if CACHE=redis)  in .env to use!")
	return nil
}
