package main

import (
)

func doSessionTable() error {
	dbType := cel.DB.DataType
	tx, err := cel.PopConnect()
	if err != nil {
		exitGracefully(err)
	}
	defer tx.Close()

	if dbType == "mariadb" {
		dbType = "mysql"
	}

	if dbType == "postgresql" {
		dbType = "postgres"
	}

	if dbType == "sqlite3" {
		dbType = "sqlite"
	}

	/*
	fileName := fmt.Sprintf("%d_create_sessions_table", time.Now().UnixMicro())

	upFile := cel.RootPath + "/migrations/" + fileName + "." + dbType + ".up.sql"
	downFile := cel.RootPath + "/migrations/" + fileName + "." + dbType + ".down.sql"

	err := copyFileFromTemplate("templates/migrations/"+dbType+"_session.sql", upFile)
	if err != nil {
		exitGracefully(err)
	}

	err = copyDataToFile([]byte("drop table sessions"), downFile)
	if err != nil {
		exitGracefully(err)
	}

	err = doMigrate("up", "")
	if err != nil {
		exitGracefully(err)
	}
	*/

	upBytes, err := templateFS.ReadFile("templates/migrations/"+dbType+"_session.sql")
	if err != nil {
		exitGracefully(err)
	}
	downBytes := []byte("drop table if exists sessions;")

	err = cel.CreatePopMigration(upBytes, downBytes, "sessions", "sql")
	if err != nil {
		exitGracefully(err)
	}

	// run migrations
	err = cel.RunPopMigrations(tx)
	if err != nil {
		exitGracefully(err)
	}

	return nil
}
