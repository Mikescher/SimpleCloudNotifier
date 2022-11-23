package db

import (
	"blackforestbytes.com/simplecloudnotifier/db/schema"
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(filename string) (*Database, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

func (db *Database) Migrate(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 24*time.Second)
	defer cancel()

	currschema, err := db.ReadSchema(ctx)
	if currschema == 0 {

		_, err = db.db.ExecContext(ctx, schema.Schema3)
		if err != nil {
			return err
		}

		return nil

	} else if currschema == 1 {
		return errors.New("cannot autom. upgrade schema 1")
	} else if currschema == 2 {
		return errors.New("cannot autom. upgrade schema 2") //TODO
	} else if currschema == 3 {
		return nil // current
	} else {
		return errors.New(fmt.Sprintf("Unknown DB schema: %d", currschema))
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

	var dbschema int
	err = r2.Scan(&dbschema)
	if err != nil {
		return 0, err
	}

	return dbschema, nil
}

func (db *Database) Ping() error {
	return db.db.Ping()
}

func (db *Database) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return db.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
}
