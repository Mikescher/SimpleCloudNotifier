package requests

import (
	server "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/db/dbtools"
	"blackforestbytes.com/simplecloudnotifier/db/schema"
	"blackforestbytes.com/simplecloudnotifier/db/simplectx"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
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

func NewRequestsDatabase(cfg server.Config) (*Database, error) {
	conf := cfg.DBRequests

	url := fmt.Sprintf("file:%s?_pragma=journal_mode(%s)&_pragma=timeout(%d)&_pragma=foreign_keys(%s)&_pragma=busy_timeout(%d)",
		conf.File,
		conf.Journal,
		conf.Timeout.Milliseconds(),
		langext.FormatBool(conf.CheckForeignKeys, "true", "false"),
		conf.BusyTimeout.Milliseconds())

	if !langext.InArray("sqlite3", sql.Drivers()) {
		sqlite.RegisterAsSQLITE3()
	}

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

	qqdb := sq.NewDB(xdb, sq.DBOptions{})

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

func (db *Database) Migrate(outerctx context.Context) error {
	innerctx, cancel := context.WithTimeout(outerctx, 24*time.Second)
	tctx := simplectx.CreateSimpleContext(innerctx, cancel)

	tx, err := tctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}
	defer func() {
		if tx.Status() == sq.TxStatusInitial || tx.Status() == sq.TxStatusActive {
			err = tx.Rollback()
			if err != nil {
				log.Err(err).Msg("failed to rollback transaction")
			}
		}
	}()

	ppReInit := false

	currschema, err := db.ReadSchema(tctx)
	if err != nil {
		return err
	}

	if currschema == 0 {
		schemastr := schema.RequestsSchema[schema.RequestsSchemaVersion].SQL
		schemahash := schema.RequestsSchema[schema.RequestsSchemaVersion].Hash

		schemahash, err := sq.HashGoSqliteSchema(tctx, schemastr)
		if err != nil {
			return err
		}

		_, err = tx.Exec(tctx, schemastr, sq.PP{})
		if err != nil {
			return err
		}

		err = db.WriteMetaInt(tctx, "schema", int64(schema.RequestsSchemaVersion))
		if err != nil {
			return err
		}

		err = db.WriteMetaString(tctx, "schema_hash", schemahash)
		if err != nil {
			return err
		}

		ppReInit = true

		currschema = schema.LogsSchemaVersion
	}

	if currschema == 1 {
		schemHashDB, err := sq.HashSqliteDatabase(tctx, tx)
		if err != nil {
			return err
		}

		schemaHashMeta, err := db.ReadMetaString(tctx, "schema_hash")
		if err != nil {
			return err
		}

		if schemHashDB != langext.Coalesce(schemaHashMeta, "") || langext.Coalesce(schemaHashMeta, "") != schema.RequestsSchema[currschema].Hash {
			log.Debug().Str("schemHashDB", schemHashDB).Msg("Schema (requests db)")
			log.Debug().Str("schemaHashMeta", langext.Coalesce(schemaHashMeta, "")).Msg("Schema (requests db)")
			log.Debug().Str("schemaHashAsset", schema.RequestsSchema[currschema].Hash).Msg("Schema (requests db)")
			return errors.New("database schema does not match (requests db)")
		} else {
			log.Debug().Str("schemHash", schemHashDB).Msg("Verified Schema consistency (requests db)")
		}
	}

	if currschema != schema.RequestsSchemaVersion {
		return errors.New(fmt.Sprintf("Unknown DB schema: %d", currschema))
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	if ppReInit {
		log.Debug().Msg("Re-Init preprocessor")
		err = db.pp.Init(outerctx) // Re-Init
		if err != nil {
			return err
		}
	}

	return nil
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
