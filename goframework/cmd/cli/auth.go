package main

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"time"
)

func doAuth() error {
	// migrations
	dbType := gof.DB.DataType
	fileName := fmt.Sprintf("%d_create_auth_tables", time.Now().UnixMicro())
	upFile := gof.RootPath + "/migrations/" + fileName + ".up.sql"
	downFile := gof.RootPath + "/migrations/" + fileName + ".down.sql"

	// temporary to make sure we do not have any err
	log.Println(dbType, upFile, downFile)

	err := copyFilefromTemplate("templates/migrations/auth_tables."+dbType+".sql", upFile)
	if err != nil {
		exitGracefully(err)
	}

	err = copyDataToFile([]byte("drop table if exists users cascade; drop table if exists tokens cascade; drop table if exists remember_tokens"), downFile)
	if err != nil {
		exitGracefully(err)
	}

	// run migrations
	err = doMigrate("up", "")
	if err != nil {
		exitGracefully(err)
	}

	// copy files over
	err = copyFilefromTemplate("templates/data/user.go.txt", gof.RootPath+"/data/user.go")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFilefromTemplate("templates/data/token.go.txt", gof.RootPath+"/data/token.go")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFilefromTemplate("templates/data/remember_token.go.txt", gof.RootPath+"/data/remember_token.go")
	if err != nil {
		exitGracefully(err)
	}

	// copy over middleware
	err = copyFilefromTemplate("templates/middleware/auth.go.txt", gof.RootPath+"/middleware/auth.go")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFilefromTemplate("templates/middleware/auth-token.go.txt", gof.RootPath+"/middleware/auth-token.go")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFilefromTemplate("templates/middleware/remember.go.txt", gof.RootPath+"/middleware/remember.go")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFilefromTemplate("templates/handlers/auth-handlers.go.txt", gof.RootPath+"/handlers/auth-handlers.go")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFilefromTemplate("templates/mailer/password-reset.html.tmpl", gof.RootPath+"/mail/password-reset.html.tmpl")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFilefromTemplate("templates/mailer/password-reset.plain.tmpl", gof.RootPath+"/mail/password-reset.plain.tmpl")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFilefromTemplate("templates/views/login.jet", gof.RootPath+"/views/login.jet")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFilefromTemplate("templates/views/forgot.jet", gof.RootPath+"/views/forgot.jet")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFilefromTemplate("templates/views/reset-password.jet", gof.RootPath+"/views/reset-password.jet")
	if err != nil {
		exitGracefully(err)
	}

	color.Yellow(" - users, tokens, and remember_tokens migrations create and executed")
	color.Yellow(" - user and token models created")
	color.Yellow(" - auth middleware created")
	color.Yellow("")
	color.Yellow("")
	color.Yellow("Don't forget to add user and token models in data/models.go, and to add appropriate middleware to your routes!")

	return nil
}
