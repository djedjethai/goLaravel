package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"os"
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

		gof.RootPath = path
		gof.DB.DataType = os.Getenv("DATABASE_TYPE")
	}
}

func getDSN() string {
	dbType := gof.DB.DataType

	// "pgx" is the name we gave to postgres in the DSN(Data Source Name)
	// in the driver(package jackc)
	// but now goMigration package use a different driver, so need to reset the dsn
	if dbType == "pgx" {
		dbType = "postgres"
	}

	if dbType == "postgres" {
		var dsn string
		// case we are working with Docker images the DB password is set
		// so we set the dsn this way
		if os.Getenv("DATABASE_PASS") != "" {
			dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_PASS"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"))
		} else {
			// case we work with Postgres on the computer
			// password is already set
			dsn = fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"))

		}
		return dsn
	}
	// if not "postgres", mariadb or mysql got their dsn like this
	// just prepend the "mysql://" to what the Jackc driver normally use
	return "mysql://" + gof.BuildDSN()

}

func showHelp() {
	color.Yellow(`Available commands:
	help                     - show the help commands
	version                  - print application version
	migrate                  - runs all up migrations that have not been run previously
	migrate down             - reverses the most recent migrations
	migrate reset            - run all down migrations in reverse order, and then all up migrations
	make migration <name>    - creates two new up and down migrations in the migrations folder
	make auth                - creates and runs migrations for authentication tables, and creates models and middleware
	make handler <name>      - creates a stub handler in the handlers directory
	make model <name>        - creates a new model in the data directory
	make session             - creates a table in the database as a session store
	make mail <name>	- creates two starter mail templates in the mail directory
	`)
}
