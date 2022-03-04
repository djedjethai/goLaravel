package goframework

// up migration means: load some files
// which here we will put in migration folder in myapp project
// and it will exec the sql from these files

// migrate down means: reverse the process of an up migration

// we can do it by step but if something goes wrong we can force a migration

// to do that we gonna use code from golang-migrate/migrate
// go get github.com/golang-migrate/migrate/v4
// go get github.com/golang-migrate/migrate/v4/database/mysql
// go get github.com/golang-migrate/migrate/v4/database/postgres
// go get github.com/golang-migrate/migrate/v4/source/file

import (
	// "database/sql"
	//_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

// dsn = Data Source Name
// golang-migrate use a different postgres driver than we do
// so we need to rebuild the dsn (that has been done in the /cmd/cli/helpers.go)
func (g *Goframework) MigrateUp(dsn string) error {
	m, err := migrate.New("file://"+g.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err = m.Up(); err != nil {
		log.Println("Error running migration", err)
		return err
	}

	return nil
}

func (g *Goframework) MigrateDownAll(dsn string) error {
	m, err := migrate.New("file://"+g.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err = m.Down(); err != nil {
		log.Println("Error running migration", err)
		return err
	}

	return nil
}

func (g *Goframework) Steps(n int, dsn string) error {
	m, err := migrate.New("file://"+g.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err = m.Steps(n); err != nil {
		log.Println("Error running migration", err)
		return err
	}

	return nil
}

// if an err occur during migration
// golang/migrate package save a dirty message in db
// to avoid that we use the following function
func (g *Goframework) MigrateForce(dsn string) error {
	m, err := migrate.New("file://"+g.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err = m.Force(-1); err != nil {
		return err
	}

	return nil
}
