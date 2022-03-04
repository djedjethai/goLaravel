package main

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"io/ioutil"
	"strings"
	"time"
)

func doMake(arg2, arg3 string) error {
	switch arg2 {
	case "key":
		rnd := gof.RandomString(32)
		color.Yellow("32 characters encryptions key: %s", rnd)

	case "migration":
		dbType := gof.DB.DataType
		if arg3 == "" {
			exitGracefully(errors.New("You must give the migration a name"))
		}

		fileName := fmt.Sprintf("%d_%s", time.Now().UnixMicro(), arg3)

		upFile := gof.RootPath + "/migrations/" + fileName + "." + dbType + ".up.sql"
		downFile := gof.RootPath + "/migrations/" + fileName + "." + dbType + ".down.sql"

		err := copyFilefromTemplate("templates/migrations/migration."+dbType+".up.sql", upFile)
		if err != nil {
			exitGracefully(err)
		}
		err = copyFilefromTemplate("templates/migrations/migration."+dbType+".down.sql", downFile)
		if err != nil {
			exitGracefully(err)
		}

	case "auth":
		err := doAuth()
		if err != nil {
			exitGracefully(err)
		}

	case "handler":
		if arg3 == "" {
			exitGracefully(errors.New("You must give the handler a name"))
		}

		fileName := gof.RootPath + "/handlers/" + strings.ToLower(arg3) + ".go"
		if fileExist(fileName) {
			exitGracefully(errors.New(fileName + "already exists."))
		}

		data, err := templateFS.ReadFile("templates/handlers/handler.go.txt")
		if err != nil {
			exitGracefully(err)
		}

		// that will change what the name for their handler and convert it to camelCase
		// which is the norm foe naming an handler,
		// and replace any occurance of $HANDLERNAME$
		// with the correct name for this particular handler
		// .ReplaceAll() will replace all $HANDLERNAME$ in the handler string with arg3
		handler := string(data)
		handler = strings.ReplaceAll(handler, "$HANDLERNAME$", strcase.ToCamel(arg3))

		err = ioutil.WriteFile(fileName, []byte(handler), 0644)
		if err != nil {
			exitGracefully(err)
		}

	case "model":
		if arg3 == "" {
			exitGracefully(errors.New("You must give the model a name"))
		}

		data, err := templateFS.ReadFile("templates/data/model.go.txt")
		if err != nil {
			exitGracefully(err)
		}

		// i need to pluralize(mettre au pluriel) the models name
		// as we name them with plurial from the beginning
		model := string(data)
		plur := pluralize.NewClient()

		var modelName = arg3
		var tableName = arg3

		// set the names following our convention
		// modelName are wrote like: User or Token
		// tableName are wrote like: users or tokens
		if plur.IsPlural(arg3) {
			modelName = plur.Singular(arg3)
			tableName = strings.ToLower(tableName)
		} else {
			tableName = strings.ToLower(plur.Plural(arg3))
		}

		fileName := gof.RootPath + "/data/" + strings.ToLower(arg3) + ".go"
		if fileExist(fileName) {
			exitGracefully(errors.New(fileName + "already exists."))
		}

		// replace $MODELNAME$ with modelName
		// replace $TABLENAME$ with tableName
		model = strings.ReplaceAll(model, "$MODELNAME$", strcase.ToCamel(modelName))
		model = strings.ReplaceAll(model, "$TABLENAME$", tableName)

		err = copyDataToFile([]byte(model), fileName)
		if err != nil {
			exitGracefully(err)
		}

	case "mail":
		if arg3 == "" {
			exitGracefully(errors.New("you must give the mail template a name"))
		}
		htmlMail := gof.RootPath + "/mail/" + strings.ToLower(arg3) + ".html.tmpl"
		plainMail := gof.RootPath + "/mail/" + strings.ToLower(arg3) + ".plain.tmpl"

		err := copyFilefromTemplate("templates/mailer/mail.html.tmpl", htmlMail)
		if err != nil {
			exitGracefully(err)
		}

		err = copyFilefromTemplate("templates/mailer/mail.plain.tmpl", plainMail)
		if err != nil {
			exitGracefully(err)
		}

	case "session":
		err := doSessionTable()
		if err != nil {
			exitGracefully(err)
		}

	}

	return nil
}
