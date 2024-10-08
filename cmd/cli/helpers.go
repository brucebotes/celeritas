package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

func setup(arg1, arg2 string) {
	if arg1 != "new" && arg1 != "version" && arg1 != "help" {
		err := godotenv.Load()
		if err != nil {
			exitGracefully(err)
		}

		path, err := os.Getwd()
		if err != nil {
			exitGracefully(err)
		}

		cel.RootPath = path
		cel.DB.DataType = os.Getenv("DATABASE_TYPE")
	}
}

func getDSN() string {
	dbType := cel.DB.DataType

	// the migrations tool uses a different driver
	// for migrations
	if dbType == "pgx" {
		dbType = "postgres"
	}

	if dbType == "postgres" {
		var dsn string
		if os.Getenv("DATABASE_PASS") != "" {
			dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_PASS"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"))
		} else {
			dsn = fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"))
		}
		return dsn
	}

	if dbType == "sqlite" {
		dbType = "sqlite3"
	}

	if dbType == "sqlite3" {
		var dsn string
		dsn = fmt.Sprintf("%s", os.Getenv("DATABASE_NAME"))
		return dsn
	}

	return "mysql://" + cel.BuildDSN()
}

func checkForDB() {
	dbType := cel.DB.DataType

	if dbType == "" {
		exitGracefully(errors.New("no database connection provided in .env"))
	}

	if !fileExists(cel.RootPath + "/config/database.yml") {
		exitGracefully(errors.New("config/database.yml does not exist"))
	}
}

func showHelp() {
	color.Yellow(`Avialable commands:

	help                             - show the help commands
	new <name>                       - create a new project with <name>
	down                             - put the server into maintenance mode
	up                               - take the server out of maintenance mode
	version                          - print application version
	migrate                          - runs all up migrations that have not been run previously
	migrate down                     - reverses the most recent migration
	migrate reset                    - runs all down migrations in reverse order, and then all up migrations
	make migration <name> <format>   - creates two new up and down migrations in the migrations folder containing <name>; format=sql/fizz (default fizz)
	make auth                        - creates and runs migrations for authentication tables, and creates models and middleware
	make handler <name>              - creates a stub handler in the handlers directory
	make model <name>                - creates a new model in the data directory
	make session                     - creates a table in the database as a session store
	make mail <mail>                 - create two starter mail templates in the mail directory
	make bundleJS <name>             - creates a new jet template with bundled javascript ESM integration in the views folder with <name>
	make cacheapi                    - creates the handlers and api routes for the cache (Redis or Badger) with a cache.jet test template 
	`)
}

func updateSourceFiles(path string, fi os.FileInfo, err error) error {
	// check for error before doing anything else
	if err != nil {
		return err
	}

	// check if current file is directory
	if fi.IsDir() {
		return nil
	}

	// only check go file
	matched, err := filepath.Match("*.go", fi.Name())
	if err != nil {
		return err
	}

	// we have a matching file
	if matched {
		// read file contents
		read, err := os.ReadFile(path)
		if err != nil {
			exitGracefully(err)
		}

		newContents := strings.Replace(string(read), "myapp", appURL, -1)

		// replace the changed file
		err = os.WriteFile(path, []byte(newContents), 0)
		if err != nil {
			exitGracefully(err)
		}
	}

	return nil
}

func updateSource() {
	// walk entire project folder, including subfolders
	err := filepath.Walk(".", updateSourceFiles)
	if err != nil {
		exitGracefully(err)
	}
}

/* No need to define these - the equivalent exist
   with the (cel *Celeritas) reciever
func createDirIfNotExist(path string) error {
	const mode = 0755
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, mode)
		if err != nil {
			return err
		}
	}
	return nil
}

func createFileIfNotExists(path string) error {
	var _, err = os.Stat(path)
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return err
		}

		defer func(file *os.File) {
			_ = file.Close()
		}(file)
	}

	return nil
}
*/
