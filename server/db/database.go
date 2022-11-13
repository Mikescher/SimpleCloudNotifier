package db

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema_1.0.sql
var schema_1_0 string

//go:embed schema_2.0.sql
var schema_2_0 string

func NewDatabase(ctx context.Context, conf scn.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", conf.DBFile)
	if err != nil {
		return nil, err
	}

	schema, err := getSchemaFromDB(ctx, db)
	if schema == 0 {

		_, err = db.ExecContext(ctx, schema_1_0)
		if err != nil {
			return nil, err
		}

		return db, nil

	} else if schema == 1 {
		return nil, errors.New("cannot autom. upgrade schema 1")
	} else if schema == 2 {
		return db, nil
	} else {
		return nil, errors.New(fmt.Sprintf("Unknown DB schema: %d", schema))
	}

}

func getSchemaFromDB(ctx context.Context, db *sql.DB) (int, error) {

	r1, err := db.QueryContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name='meta'")
	if err != nil {
		return 0, err
	}

	if !r1.Next() {
		return 0, nil
	}

	r2, err := db.QueryContext(ctx, "SELECT value_int FROM meta WHERE key='schema'")
	if err != nil {
		return 0, err
	}
	if !r2.Next() {
		return 0, errors.New("no schema entry in meta table")
	}

	var schema int
	err = r2.Scan(&schema)
	if err != nil {
		return 0, err
	}

	return schema, nil
}
