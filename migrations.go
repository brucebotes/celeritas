package celeritas

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gobuffalo/genny"
	pop "github.com/gobuffalo/pop/v6"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (c *Celeritas) PopConnect() (*pop.Connection, error) {
	//log.Println("Info - Supported migration dialetcs configured are:")
	//for i, d := range pop.AvailableDialects {
	//	log.Printf("\tdialect %d = %s", i, d)
	//}

	tx, err := pop.Connect("development")
	if err != nil {
		return nil, err
	}
	return tx, err
}

func (c *Celeritas) CreatePopMigration(up, down []byte, migrationName, migrationType string) error {
	var migrationPath = c.RootPath + "/migrations"
	err := popMigrationCreate(migrationPath, migrationName, migrationType, up, down)
	if err != nil {
		return err
	}
	return nil
}

func (c *Celeritas) RunPopMigrations(tx *pop.Connection) error {
	var migrationPath = c.RootPath + "/migrations"

	fm, err := pop.NewFileMigrator(migrationPath, tx)
	if err != nil {
		return err
	}

	err = fm.Up()
	if err != nil {
		return err
	}

	return nil
}

func (c *Celeritas) PopMigrateDown(tx *pop.Connection, steps ...int) error {
	var migrationPath = c.RootPath + "/migrations"

	step := 1
	if len(steps) > 0 {
		step = steps[0]
	}

	fm, err := pop.NewFileMigrator(migrationPath, tx)
	if err != nil {
		return err
	}

	err = fm.Down(step)
	if err != nil {
		return err
	}

	return nil
}

func (c *Celeritas) PopMigrateReset(tx *pop.Connection) error {
	var migrationPath = c.RootPath + "/migrations"

	fm, err := pop.NewFileMigrator(migrationPath, tx)
	if err != nil {
		return err
	}

	err = fm.Reset()
	if err != nil {
		return err
	}

	return nil
}

func (c *Celeritas) MigrateUp(dsn string) error {
	m, err := migrate.New("file://"+c.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		log.Println("Error running migration: ", err)
		return err
	}
	return nil
}

func (c *Celeritas) MigrateDownAll(dsn string) error {
	m, err := migrate.New("file://"+c.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Down(); err != nil {
		log.Println("Error running migration: ", err)
		return err
	}
	return nil
}

func (c *Celeritas) Steps(n int, dsn string) error {
	m, err := migrate.New("file://"+c.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Steps(n); err != nil {
		return err
	}
	return nil
}

func (c *Celeritas) MigrateForce(dsn string) error {
	m, err := migrate.New("file://"+c.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Force(-1); err != nil {
		return err
	}
	return err
}

// Copied from the github.com/gobuffalo/pop (version 4) module
// - it was deleted in version 6
// - orignal function in version 4 was named as MigrationCreate()
func popMigrationCreate(path, name, ext string, up, down []byte) error {
	run := genny.WetRunner(context.Background())
	g := genny.New()

	n := time.Now().UTC()
	s := n.Format("20060102150405")

	upf := filepath.Join(path, fmt.Sprintf("%s_%s.up.%s", s, name, ext))
	g.File(genny.NewFileB(upf, up))

	downf := filepath.Join(path, fmt.Sprintf("%s_%s.down.%s", s, name, ext))
	g.File(genny.NewFileB(downf, down))

	run.With(g)

	return run.Run()
}
