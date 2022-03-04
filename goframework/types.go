package goframework

import (
	"database/sql"
)

type initPaths struct {
	rootPath    string
	folderNames []string
}

type cookieConfig struct {
	name     string
	lifetime string
	persist  string // if browser is closed and re-open
	secure   string // ecrypted cookie
	domain   string
}

type databaseConfig struct {
	dsn      string
	database string
}

// this type will be exported
type Database struct {
	DataType string
	Pool     *sql.DB
}

// redis type
type RedisConfig struct {
	host     string
	password string
	prefix   string
}
