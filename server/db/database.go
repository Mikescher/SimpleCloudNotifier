package db

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func NewDatabase(conf scn.Config) (*sql.DB, error) {
	return sql.Open("sqlite3", conf.DBFile)
}
