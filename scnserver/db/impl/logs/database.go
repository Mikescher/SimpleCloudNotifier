package logs

import (
	server "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/db/dbtools"
	"blackforestbytes.com/simplecloudnotifier/db/schema"
	"blackforestbytes.com/simplecloudnotifier/db/simplectx"
	"blackforestbytes.com/simplecloudnotifier/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/exerr"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"os"
	"path/filepath"
	"time"
)

type Database struct {
	db            sq.DB
	pp            *dbtools.DBPreprocessor
	wal           bool
	name          string
	schemaVersion int
	schema        map[int]schema.Def
}

func NewLogsDatabase(cfg server.Config) (*Database, error) {
	conf := cfg.DBLogs

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

	qqdb := sq.NewDB(xdb, sq.DBOptions{RegisterDefaultConverter: langext.PTrue, RegisterCommentTrimmer: langext.PTrue})
	models.RegisterConverter(qqdb)

	if conf.EnableLogger {
		qqdb.AddListener(dbtools.DBLogger{})
	}

	pp, err := dbtools.NewDBPreprocessor(qqdb)
	if err != nil {
		return nil, err
	}

	qqdb.AddListener(pp)

	scndb := &Database{
		db:            qqdb,
		pp:            pp,
		wal:           conf.Journal == "WAL",
		schemaVersion: schema.LogsSchemaVersion,
		schema:        schema.LogsSchema,
		name:          "logs",
	}

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

	if currschema == db.schemaVersion {
		log.Info().Msgf("Database [%s] is up-to-date (%d == %d)", db.name, currschema, db.schemaVersion)
	}

	for currschema < db.schemaVersion {

		if currschema == 0 {
			log.Info().Msgf("Migrate database (initialize) [%s] %d -> %d", db.name, currschema, db.schemaVersion)

			schemastr := db.schema[db.schemaVersion].SQL
			schemahash := db.schema[db.schemaVersion].Hash

			_, err = tx.Exec(tctx, schemastr, sq.PP{})
			if err != nil {
				return err
			}

			err = db.WriteMetaInt(tctx, "schema", int64(db.schemaVersion))
			if err != nil {
				return err
			}

			err = db.WriteMetaString(tctx, "schema_hash", schemahash)
			if err != nil {
				return err
			}

			ppReInit = true

			currschema = db.schemaVersion
		} else {
			log.Info().Msgf("Migrate database [%s] %d -> %d", db.name, currschema, currschema+1)

			err = db.migrateSingle(tctx, tx, currschema, currschema+1)
			if err != nil {
				return err
			}

			currschema = currschema + 1

			ppReInit = true
		}
	}

	if currschema != db.schemaVersion {
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

//goland:noinspection SqlConstantCondition,SqlWithoutWhere
func (db *Database) migrateSingle(tctx *simplectx.SimpleContext, tx sq.Tx, schemaFrom int, schemaTo int) error {

	if schemaFrom == schemaTo-1 {

		migSQL := db.schema[schemaTo].MigScript
		if migSQL == "" {
			return exerr.New(exerr.TypeInternal, fmt.Sprintf("missing %s migration from %d to %d", db.name, schemaFrom, schemaTo)).Build()
		}

		return db.migrateBySQL(tctx, tx, migSQL, schemaFrom, schemaTo, db.schema[schemaTo].Hash, nil)
	}

	return exerr.New(exerr.TypeInternal, fmt.Sprintf("missing %s migration from %d to %d", db.name, schemaFrom, schemaTo)).Build()
}

func (db *Database) migrateBySQL(tctx *simplectx.SimpleContext, tx sq.Tx, stmts string, currSchemaVers int, resultSchemVers int, resultHash string, post func(tctx *simplectx.SimpleContext, tx sq.Tx) error) error {

	schemaHashMeta, err := db.ReadMetaString(tctx, "schema_hash")
	if err != nil {
		return err
	}

	schemHashDBBefore, err := sq.HashSqliteDatabase(tctx, tx)
	if err != nil {
		return err
	}

	if schemHashDBBefore != langext.Coalesce(schemaHashMeta, "") || langext.Coalesce(schemaHashMeta, "") != db.schema[currSchemaVers].Hash {
		log.Debug().Str("schemHashDB", schemHashDBBefore).Msg("Schema (primary db)")
		log.Debug().Str("schemaHashMeta", langext.Coalesce(schemaHashMeta, "")).Msg("Schema (primary db)")
		log.Debug().Str("schemaHashAsset", db.schema[currSchemaVers].Hash).Msg("Schema (primary db)")
		return errors.New("database schema does not match (primary db)")
	} else {
		log.Debug().Str("schemHash", schemHashDBBefore).Msg("Verified Schema consistency (primary db)")
	}

	log.Info().Msgf("Upgrade schema from %d -> %d", currSchemaVers, resultSchemVers)

	_, err = tx.Exec(tctx, stmts, sq.PP{})
	if err != nil {
		return err
	}

	schemHashDBAfter, err := sq.HashSqliteDatabase(tctx, tx)
	if err != nil {
		return err
	}

	if schemHashDBAfter != resultHash {

		schemaDBStr := langext.Must(createSqliteDatabaseSchemaStringFromSQL(tctx, db.schema[resultSchemVers].SQL))
		resultDBStr := langext.Must(sq.CreateSqliteDatabaseSchemaString(tctx, tx))

		fmt.Printf("========================================= SQL SCHEMA-DUMP STR (CORRECT | FROM COMPILED SCHEMA):%s\n=========================================\n\n", schemaDBStr)
		fmt.Printf("========================================= SQL SCHEMA-DUMP STR (CURRNET | AFTER MIGRATION):%s\n=========================================\n\n", resultDBStr)

		return fmt.Errorf("database [%s] schema does not match after [%d -> %d] migration (expected: %s | actual: %s)", db.name, currSchemaVers, resultSchemVers, resultHash, schemHashDBBefore)
	}

	err = db.WriteMetaInt(tctx, "schema", int64(resultSchemVers))
	if err != nil {
		return err
	}

	err = db.WriteMetaString(tctx, "schema_hash", resultHash)
	if err != nil {
		return err
	}

	log.Info().Msgf("Upgrade schema from %d -> %d succesfully", currSchemaVers, resultSchemVers)

	return nil
}

func createSqliteDatabaseSchemaStringFromSQL(ctx context.Context, schemaStr string) (string, error) {
	dbdir := os.TempDir()
	dbfile1 := filepath.Join(dbdir, langext.MustHexUUID()+".sqlite3")
	defer func() { _ = os.Remove(dbfile1) }()

	err := os.MkdirAll(dbdir, os.ModePerm)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("file:%s?_pragma=journal_mode(%s)&_pragma=timeout(%d)&_pragma=foreign_keys(%s)&_pragma=busy_timeout(%d)", dbfile1, "DELETE", 1000, "true", 1000)

	xdb, err := sqlx.Open("sqlite", url)
	if err != nil {
		return "", err
	}

	db := sq.NewDB(xdb, sq.DBOptions{})

	_, err = db.Exec(ctx, schemaStr, sq.PP{})
	if err != nil {
		return "", err
	}

	return sq.CreateSqliteDatabaseSchemaString(ctx, db)
}

func (db *Database) Ping(ctx context.Context) error {
	return db.db.Ping(ctx)
}

func (db *Database) Version(ctx context.Context) (string, string, error) {
	type rt struct {
		Version  string `db:"version"`
		SourceID string `db:"sourceID"`
	}

	resp, err := sq.QuerySingle[rt](ctx, db.db, "SELECT sqlite_version() AS version, sqlite_source_id() AS sourceID", sq.PP{}, sq.SModeFast, sq.Safe)
	if err != nil {
		return "", "", err
	}

	return resp.Version, resp.SourceID, nil
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
