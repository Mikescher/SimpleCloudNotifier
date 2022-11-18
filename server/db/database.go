package db

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

//go:embed schema_1.0.ddl
var schema_1_0 string

//go:embed schema_2.0.ddl
var schema_2_0 string

//go:embed schema_3.0.ddl
var schema_3_0 string

type Database struct {
	db *sql.DB
}

func NewDatabase(conf scn.Config) (*Database, error) {
	db, err := sql.Open("sqlite3", conf.DBFile)
	if err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

func (db *Database) Migrate(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 24*time.Second)
	defer cancel()

	schema, err := db.ReadSchema(ctx)
	if schema == 0 {

		_, err = db.db.ExecContext(ctx, schema_3_0)
		if err != nil {
			return err
		}

		return nil

	} else if schema == 1 {
		return errors.New("cannot autom. upgrade schema 1")
	} else if schema == 2 {
		return errors.New("cannot autom. upgrade schema 2") //TODO
	} else if schema == 3 {
		return nil // current
	} else {
		return errors.New(fmt.Sprintf("Unknown DB schema: %d", schema))
	}

}

func (db *Database) ReadSchema(ctx context.Context) (int, error) {

	r1, err := db.db.QueryContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name='meta'")
	if err != nil {
		return 0, err
	}

	if !r1.Next() {
		return 0, nil
	}

	r2, err := db.db.QueryContext(ctx, "SELECT value_int FROM meta WHERE meta_key='schema'")
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

func (db *Database) Ping() error {
	return db.db.Ping()
}

func (db *Database) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return db.db.BeginTx(ctx, nil)
}
