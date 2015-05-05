package model

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB(driverName, dataSourceName string) (err error) {
	db, err = sql.Open(driverName, dataSourceName)
	return
}
