package primary

import (
	server "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/db/dbtools"
	"blackforestbytes.com/simplecloudnotifier/db/schema"
	"blackforestbytes.com/simplecloudnotifier/db/simplectx"
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
		schemastr := schema.PrimarySchema[schema.PrimarySchemaVersion].SQL
		schemahash := schema.PrimarySchema[schema.PrimarySchemaVersion].Hash

		_, err = tx.Exec(tctx, schemastr, sq.PP{})
		if err != nil {
			return err
		}

		err = db.WriteMetaInt(tctx, "schema", int64(schema.PrimarySchemaVersion))
		if err != nil {
			return err
		}

		err = db.WriteMetaString(tctx, "schema_hash", schemahash)
		if err != nil {
			return err
		}

		ppReInit = true

		currschema = schema.PrimarySchemaVersion
	}

	if currschema == 1 {
		return errors.New("cannot autom. upgrade schema 1")
	}

	if currschema == 2 {
		return errors.New("cannot autom. upgrade schema 2")
	}

	if currschema == 3 {

		schemaHashMeta, err := db.ReadMetaString(tctx, "schema_hash")
		if err != nil {
			return err
		}

		schemHashDB, err := sq.HashSqliteDatabase(tctx, tx)
		if err != nil {
			return err
		}

		if schemHashDB != langext.Coalesce(schemaHashMeta, "") || langext.Coalesce(schemaHashMeta, "") != schema.PrimarySchema[currschema].Hash {
			log.Debug().Str("schemHashDB", schemHashDB).Msg("Schema (primary db)")
			log.Debug().Str("schemaHashMeta", langext.Coalesce(schemaHashMeta, "")).Msg("Schema (primary db)")
			log.Debug().Str("schemaHashAsset", schema.PrimarySchema[currschema].Hash).Msg("Schema (primary db)")
			return errors.New("database schema does not match (primary db)")
		} else {
			log.Debug().Str("schemHash", schemHashDB).Msg("Verified Schema consistency (primary db)")
		}

		log.Info().Int("currschema", currschema).Msg("Upgrade schema from 3 -> 4")

		_, err = tx.Exec(tctx, schema.PrimaryMigration_3_4, sq.PP{})
		if err != nil {
			return err
		}

		currschema = 4

		err = db.WriteMetaInt(tctx, "schema", int64(currschema))
		if err != nil {
			return err
		}

		err = db.WriteMetaString(tctx, "schema_hash", schema.PrimarySchema[currschema].Hash)
		if err != nil {
			return err
		}

		log.Info().Int("currschema", currschema).Msg("Upgrade schema from 3 -> 4 succesfuly")

		ppReInit = true
	}

	if currschema == 4 {

		schemaHashMeta, err := db.ReadMetaString(tctx, "schema_hash")
		if err != nil {
			return err
		}

		schemHashDB, err := sq.HashSqliteDatabase(tctx, tx)
		if err != nil {
			return err
		}

		if schemHashDB != langext.Coalesce(schemaHashMeta, "") || langext.Coalesce(schemaHashMeta, "") != schema.PrimarySchema[currschema].Hash {
			log.Debug().Str("schemHashDB", schemHashDB).Msg("Schema (primary db)")
			log.Debug().Str("schemaHashMeta", langext.Coalesce(schemaHashMeta, "")).Msg("Schema (primary db)")
			log.Debug().Str("schemaHashAsset", schema.PrimarySchema[currschema].Hash).Msg("Schema (primary db)")
			return errors.New("database schema does not match (primary db)")
		} else {
			log.Debug().Str("schemHash", schemHashDB).Msg("Verified Schema consistency (primary db)")
		}
	}

	if currschema != schema.PrimarySchemaVersion {
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
