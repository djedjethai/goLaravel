package data

// to download the ORM
// go get -u github.com/upper/db/v4/adapter/postgresql
// go get -u github.com/upper/db/v4/adapter/mysql

import (
	"database/sql"
	"fmt"
	db2 "github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mysql"
	"github.com/upper/db/v4/adapter/postgresql"
	"os"
)

var db *sql.DB
var upper db2.Session

type Models struct {
	// any models inserted here (and in the New func)
	// are easily accessible throughout the entire application

	// when the framework is started for a new project
	// we do not want this 2 models to be load pragmaticaly
	// so the user will have to add them himself
	Users  User
	Tokens Token
}

func New(databasePool *sql.DB) Models {
	db = databasePool

	switch os.Getenv("DATABASE_TYPE") {

	case "mysql", "mariadb":
		// set upper/io ORM for mysql
		upper, _ = mysql.New(databasePool)
	case "postgres", "postgresql":
		// set upper/io ORM for postgres
		upper, _ = postgresql.New(databasePool)
	default:
		// do nothing
	}

	return Models{

		// when the framework is started for a new project
		// we do not want this 2 models to be load pragmaticaly
		// so the user will have to add them himself
		Users:  User{},
		Tokens: Token{},
	}
}

func GetInsertID(i db2.ID) int {
	// get the id type from i
	// this .Sprintf("%T") return the type
	idType := fmt.Sprintf("%T", i)
	// that the id type postgres return
	if idType == "int64" {
		// cast to int
		return int(i.(int64))
	}

	// as here we just support sql and postgres, for sql we return
	return i.(int)
}
