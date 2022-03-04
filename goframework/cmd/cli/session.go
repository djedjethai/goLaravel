package main

import (
	"fmt"
	"time"
)

func doSessionTable() error {

	dbType := gof.DB.DataType

	// in case user use another name
	if dbType == "mariadb" {
		dbType = "mysql"
	}

	if dbType == "postgresql" {
		dbType = "postgres"
	}

	filename := fmt.Sprintf("%d_create_session_tables", time.Now().UnixMicro())
	upFile := gof.RootPath + "/migrations/" + filename + "." + dbType + ".up.sql"
	downFile := gof.RootPath + "/migrations/" + filename + "." + dbType + ".down.sql"

	err := copyFilefromTemplate("templates/migrations/"+dbType+"_session.sql", upFile)
	if err != nil {
		exitGracefully(err)
	}

	err = copyDataToFile([]byte("drop table if exists sessions"), downFile)
	if err != nil {
		exitGracefully(err)
	}

	// run migration
	err = doMigrate("up", "")
	if err != nil {
		exitGracefully(err)
	}

	return nil
}
