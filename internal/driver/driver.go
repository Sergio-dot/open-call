package driver

import (
	"database/sql"
	"time"
)

// DB holds the database connection pool
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbLifetime = 5 * time.Minute

func ConnectSQL(dsn string) (*DB, error) {
	// TODO create connection to database | 13.2 4:00
}

func NewDatabase(dsn string) (*sql.DB, error) {

}
