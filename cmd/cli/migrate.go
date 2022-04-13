package main

// Note we replaced (commented out) the SQL based migrations with the Pop/fizz  migrations

func doMigrate(arg2, arg3 string) error {
	// dsn := getDSN() // is used by github.com/golang-migrate/migrate/v4 - we have replaced this with fizz

	checkForDB()
	tx, err := cel.PopConnect()
	if err != nil {
		exitGracefully(err)
	}
	defer tx.Close()

	// run the migration command
	switch arg2 {
	case "up":
		//err := cel.MigrateUp(dsn)
		err := cel.RunPopMigrations(tx)
		if err != nil {
			return err
		}
	case "down":
		/*
		 *if arg3 == "all" {
		 *  err := cel.MigrateDownAll(dsn)
		 *  if err != nil {
		 *    return err
		 *  }
		 *} else {
		 *  err := cel.Steps(-1, dsn)
		 *  if err != nil {
		 *    return err
		 *  }
		 *}
		 */
		if arg3 == "all" {
			err := cel.PopMigrateDown(tx, -1)
			if err != nil {
				return err
			}
		} else {
			err := cel.PopMigrateDown(tx, 1)
			if err != nil {
				return err
			}

		}
	case "reset":
		/*
		 *err := cel.MigrateDownAll(dsn)
		 *if err != nil {
		 *  return err
		 *}
		 *err = cel.MigrateUp(dsn)
		 *if err != nil {
		 *  return err
		 *}
		 */
		err := cel.PopMigrateReset(tx)
		if err != nil {
			return err
		}
	default:
		showHelp()
	}
	return nil
}
