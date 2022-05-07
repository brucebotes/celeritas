package main

import (
	"os"
	"path"

	"github.com/fatih/color"
)

func doAuth() error {
	checkForDB()
	tx, err := cel.PopConnect()
	if err != nil {
		exitGracefully(err)
	}
	defer tx.Close()

	// migrations
	dbType := cel.DB.DataType

	upBytes, err := templateFS.ReadFile("templates/migrations/auth_tables." + dbType + ".sql")
	if err != nil {
		exitGracefully(err)
	}

	downBytes := []byte("drop table if exists users cascade; drop table if exists tokens cascade; drop table if exists remember_tokens;")

	err = cel.CreatePopMigration(upBytes, downBytes, "auth", "sql")
	if err != nil {
		exitGracefully(err)
	}

	// run migrations
	err = cel.RunPopMigrations(tx)
	if err != nil {
		exitGracefully(err)
	}

	//copy files over
	err = copyFileFromTemplate("templates/data/user.go.txt", cel.RootPath+"/data/user.go")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/data/token.go.txt", cel.RootPath+"/data/token.go")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/data/remember_token.go.txt", cel.RootPath+"/data/remember_token.go")
	if err != nil {
		exitGracefully(err)
	}

	// copy over middleware
	err = copyFileFromTemplate("templates/middleware/auth.go.txt", cel.RootPath+"/middleware/auth.go")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/middleware/auth-token.go.txt", cel.RootPath+"/middleware/auth-token.go")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/middleware/remember.go.txt", cel.RootPath+"/middleware/remember.go")
	if err != nil {
		exitGracefully(err)
	}
	// reset "myapp" package references in templates
	color.Yellow("\tUpdating middleware source files...")
	appURL = path.Base(cel.RootPath) // get appName from path
	os.Chdir(cel.RootPath + "/middleware")
	updateSource()

	err = copyFileFromTemplate("templates/handlers/auth-handlers.go.txt", cel.RootPath+"/handlers/auth-handlers.go")
	if err != nil {
		exitGracefully(err)
	}
	// reset "myapp" package references in templates
	color.Yellow("\tUpdating handlers source files...")
	os.Chdir(cel.RootPath + "/handlers")
	updateSource()

	// copy templates
	err = copyFileFromTemplate("templates/mailer/password-reset.html.tmpl", cel.RootPath+"/mail/password-reset.html.tmpl")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/mailer/password-reset.plain.tmpl", cel.RootPath+"/mail/password-reset.plain.tmpl")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/views/login.jet", cel.RootPath+"/views/login.jet")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/views/forgot.jet", cel.RootPath+"/views/forgot.jet")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/views/reset-password.jet", cel.RootPath+"/views/reset-password.jet")
	if err != nil {
		exitGracefully(err)
	}

	// copy routes
	err = copyFileFromTemplate("templates/routes/users-routes.go.txt", cel.RootPath+"/routes-users.go")
	if err != nil {
		exitGracefully(err)
	}
	err = copyFileFromTemplate("templates/routes/util-routes.go.txt", cel.RootPath+"/routes-util.go")
	if err != nil {
		exitGracefully(err)
	}
	appURL = path.Base(cel.RootPath) // get appName from path
	os.Chdir(cel.RootPath)
	updateSource()

	// reset path
	//os.Chdir(cel.RootPath)

	color.Yellow("\t  - users, tokens and remember_tokens migrations created and executed")
	color.Yellow("\t  - user and token models created")
	color.Yellow("\t  - auth and middleware created")
	color.Yellow("")
	color.Yellow("Don't forget to add user and token models in data/models.go, and to add appropiate middleware to your routes!")
	color.Yellow("")
	color.Yellow("Ps. Run http://localhost:port/appname/util/create-user to create the first admin user")

	return nil
}
