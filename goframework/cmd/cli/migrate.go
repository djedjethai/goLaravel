package main

import (
// "fmt"
)

func doMigrate(arg2, arg3 string) error {
	dsn := getDSN()

	// run migration command
	switch arg2 {
	case "up":
		err := gof.MigrateUp(dsn)
		if err != nil {
			return err
		}
	case "down":
		if arg3 == "all" {
			err := gof.MigrateDownAll(dsn)
			if err != nil {
				return err
			}
		} else {
			err := gof.Steps(-1, dsn)
			if err != nil {
				return err
			}
		}
	case "reset":
		err := gof.MigrateDownAll(dsn)
		if err != nil {
			return err
		}

		err = gof.MigrateUp(dsn)
		if err != nil {
			return err
		}
	case "default":
		showHelp()
	}

	return nil
}
