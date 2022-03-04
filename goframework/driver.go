package goframework

// go get github.com/jackc/pgconn
// go get github.com/jackc/pgx/v4
// go get github.com/jackc/pgx/v4/stdlib

import (
	"database/sql"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// dsn means Data-Source, the stream(or the string ..??) connected to the db
// we read the dbtype from the .env file
func (g *Goframework) OpenDB(dbtype, dsn string) (*sql.DB, error) {
	if dbtype == "postgres" || dbtype == "postgresql" {
		dbtype = "pgx"
	}

	// open a standart Go connection to the db
	db, err := sql.Open(dbtype, dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// db is the pool of connection
	return db, nil
}
