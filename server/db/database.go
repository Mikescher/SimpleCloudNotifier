package db

import (
	server "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/db/schema"
	"blackforestbytes.com/simplecloudnotifier/sq"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"time"
)

type Database struct {
	db sq.DB
}

func NewDatabase(conf server.Config) (*Database, error) {
	url := fmt.Sprintf("file:%s?_journal=%s&_timeout=%d&_fk=%s", conf.DBFile, conf.DBJournal, conf.DBTimeout.Milliseconds(), langext.FormatBool(conf.DBCheckForeignKeys, "true", "false"))

	xdb, err := sqlx.Open("sqlite3", url)
	if err != nil {
		return nil, err
	}

	xdb.SetMaxOpenConns(5)
	xdb.SetMaxIdleConns(5)
	xdb.SetConnMaxLifetime(60 * time.Minute)
	xdb.SetConnMaxIdleTime(60 * time.Minute)

	qqdb := sq.NewDB(xdb)

	return &Database{qqdb}, nil
}

func (db *Database) Migrate(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 24*time.Second)
	defer cancel()

	currschema, err := db.ReadSchema(ctx)
	if currschema == 0 {

		_, err = db.db.Exec(ctx, schema.Schema3, sq.PP{})
		if err != nil {
			return err
		}

		err = db.WriteMetaInt(ctx, "schema", 3)
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

func (db *Database) Ping(ctx context.Context) error {
	return db.db.Ping(ctx)
}

func (db *Database) BeginTx(ctx context.Context) (sq.Tx, error) {
	return db.db.BeginTransaction(ctx, sql.LevelDefault)
}
