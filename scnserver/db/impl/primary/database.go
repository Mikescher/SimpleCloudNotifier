package primary

import (
	server "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/db/dbtools"
	"blackforestbytes.com/simplecloudnotifier/db/schema"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"time"
)

type Database struct {
	db  sq.DB
	pp  *dbtools.DBPreprocessor
	wal bool
}

func NewPrimaryDatabase(cfg server.Config) (*Database, error) {
	conf := cfg.DBMain

	url := fmt.Sprintf("file:%s?_journal=%s&_timeout=%d&_fk=%s&_busy_timeout=%d", conf.File, conf.Journal, conf.Timeout.Milliseconds(), langext.FormatBool(conf.CheckForeignKeys, "true", "false"), conf.BusyTimeout.Milliseconds())

	xdb, err := sqlx.Open("sqlite3", url)
	if err != nil {
		return nil, err
	}

	if conf.SingleConn {
		xdb.SetMaxOpenConns(1)
	} else {
		xdb.SetMaxOpenConns(5)
		xdb.SetMaxIdleConns(5)
		xdb.SetConnMaxLifetime(60 * time.Minute)
		xdb.SetConnMaxIdleTime(60 * time.Minute)
	}

	qqdb := sq.NewDB(xdb)

	if conf.EnableLogger {
		qqdb.AddListener(dbtools.DBLogger{})
	}

	pp, err := dbtools.NewDBPreprocessor(qqdb)
	if err != nil {
		return nil, err
	}

	qqdb.AddListener(pp)

	scndb := &Database{db: qqdb, pp: pp, wal: conf.Journal == "WAL"}

	return scndb, nil
}

func (db *Database) DB() sq.DB {
	return db.db
}

func (db *Database) Migrate(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 24*time.Second)
	defer cancel()

	currschema, err := db.ReadSchema(ctx)
	if err != nil {
		return err
	}

	if currschema == 0 {

		schemastr := schema.PrimarySchema3

		schemahash, err := sq.HashSqliteSchema(ctx, schemastr)
		if err != nil {
			return err
		}

		_, err = db.db.Exec(ctx, schemastr, sq.PP{})
		if err != nil {
			return err
		}

		err = db.WriteMetaInt(ctx, "schema", 3)
		if err != nil {
			return err
		}

		err = db.WriteMetaString(ctx, "schema_hash", schemahash)
		if err != nil {
			return err
		}

		err = db.pp.Init(ctx) // Re-Init
		if err != nil {
			return err
		}

		return nil

	} else if currschema == 1 {
		return errors.New("cannot autom. upgrade schema 1")
	} else if currschema == 2 {
		return errors.New("cannot autom. upgrade schema 2")
	} else if currschema == 3 {

		schemHashDB, err := sq.HashSqliteDatabase(ctx, db.db)
		if err != nil {
			return err
		}

		schemaHashMeta, err := db.ReadMetaString(ctx, "schema_hash")
		if err != nil {
			return err
		}

		schemHashAsset := schema.PrimaryHash3
		if err != nil {
			return err
		}

		if schemHashDB != langext.Coalesce(schemaHashMeta, "") || langext.Coalesce(schemaHashMeta, "") != schemHashAsset {
			log.Debug().Str("schemHashDB", schemHashDB).Msg("Schema (primary db)")
			log.Debug().Str("schemaHashMeta", langext.Coalesce(schemaHashMeta, "")).Msg("Schema (primary db)")
			log.Debug().Str("schemHashAsset", schemHashAsset).Msg("Schema (primary db)")
			return errors.New("database schema does not match (primary db)")
		} else {
			log.Debug().Str("schemHash", schemHashDB).Msg("Verified Schema consistency (primary db)")
		}

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

func (db *Database) Stop(ctx context.Context) error {
	if db.wal {
		_, err := db.db.Exec(ctx, "PRAGMA wal_checkpoint;", sq.PP{})
		if err != nil {
			return err
		}
	}
	err := db.db.Exit()
	if err != nil {
		return err
	}
	return nil
}
