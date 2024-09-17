package test

import (
	"blackforestbytes.com/simplecloudnotifier/db/impl/logs"
	"blackforestbytes.com/simplecloudnotifier/db/impl/primary"
	"blackforestbytes.com/simplecloudnotifier/db/impl/requests"
	"blackforestbytes.com/simplecloudnotifier/db/schema"
	"blackforestbytes.com/simplecloudnotifier/db/simplectx"
	tt "blackforestbytes.com/simplecloudnotifier/test/util"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"testing"
)

func TestPrimaryDB_Current(t *testing.T) {
	dbf1, dbf2, dbf3, conf, stop := tt.StartSimpleTestspace(t)
	defer stop()

	ctx := context.Background()

	tt.AssertAny(dbf1)
	tt.AssertAny(dbf2)
	tt.AssertAny(dbf3)
	tt.AssertAny(conf)

	{
		db1, err := primary.NewPrimaryDatabase(conf)
		tt.TestFailIfErr(t, err)

		{
			tctx := simplectx.CreateSimpleContext(ctx, nil)

			schema1, err := db1.ReadSchema(tctx)
			tt.TestFailIfErr(t, err)
			tt.AssertEqual(t, "schema1", 0, schema1)

			err = tctx.CommitTransaction()
			tt.TestFailIfErr(t, err)
		}

		{
			err = db1.Migrate(ctx)
			tt.TestFailIfErr(t, err)
		}

		{
			tctx := simplectx.CreateSimpleContext(ctx, nil)

			schema2, err := db1.ReadSchema(tctx)
			tt.TestFailIfErr(t, err)
			tt.AssertEqual(t, "schema2", schema.PrimarySchemaVersion, schema2)

			err = tctx.CommitTransaction()
			tt.TestFailIfErr(t, err)
		}

		{
			err = db1.Migrate(ctx)
			tt.TestFailIfErr(t, err)
		}

		{
			tctx := simplectx.CreateSimpleContext(ctx, nil)

			schema2, err := db1.ReadSchema(tctx)
			tt.TestFailIfErr(t, err)
			tt.AssertEqual(t, "schema2", schema.PrimarySchemaVersion, schema2)

			err = tctx.CommitTransaction()
			tt.TestFailIfErr(t, err)
		}

		err = db1.Stop(ctx)
		tt.TestFailIfErr(t, err)
	}

	{
		db1New, err := primary.NewPrimaryDatabase(conf)
		tt.TestFailIfErr(t, err)

		{
			tctx := simplectx.CreateSimpleContext(ctx, nil)

			schema3, err := db1New.ReadSchema(tctx)
			tt.TestFailIfErr(t, err)
			tt.AssertEqual(t, "schema3", schema.PrimarySchemaVersion, schema3)

			err = tctx.CommitTransaction()
			tt.TestFailIfErr(t, err)
		}

		err = db1New.Migrate(ctx)
		tt.TestFailIfErr(t, err)

		{
			tctx := simplectx.CreateSimpleContext(ctx, nil)

			schema4, err := db1New.ReadSchema(tctx)
			tt.TestFailIfErr(t, err)
			tt.AssertEqual(t, "schema4", schema.PrimarySchemaVersion, schema4)

			err = tctx.CommitTransaction()
			tt.TestFailIfErr(t, err)
		}
	}
}

func TestLogsDB_Current(t *testing.T) {
	dbf1, dbf2, dbf3, conf, stop := tt.StartSimpleTestspace(t)
	defer stop()

	ctx := context.Background()

	tt.AssertAny(dbf1)
	tt.AssertAny(dbf2)
	tt.AssertAny(dbf3)
	tt.AssertAny(conf)

	{
		db1, err := logs.NewLogsDatabase(conf)
		tt.TestFailIfErr(t, err)

		{
			tctx := simplectx.CreateSimpleContext(ctx, nil)

			schema1, err := db1.ReadSchema(tctx)
			tt.TestFailIfErr(t, err)
			tt.AssertEqual(t, "schema1", 0, schema1)

			err = tctx.CommitTransaction()
			tt.TestFailIfErr(t, err)
		}

		{
			err = db1.Migrate(ctx)
			tt.TestFailIfErr(t, err)
		}

		{
			tctx := simplectx.CreateSimpleContext(ctx, nil)

			schema2, err := db1.ReadSchema(tctx)
			tt.TestFailIfErr(t, err)
			tt.AssertEqual(t, "schema2", schema.LogsSchemaVersion, schema2)

			err = tctx.CommitTransaction()
			tt.TestFailIfErr(t, err)
		}

		{
			err = db1.Migrate(ctx)
			tt.TestFailIfErr(t, err)
		}

		{
			tctx := simplectx.CreateSimpleContext(ctx, nil)

			schema2, err := db1.ReadSchema(tctx)
			tt.TestFailIfErr(t, err)
			tt.AssertEqual(t, "schema2", schema.LogsSchemaVersion, schema2)

			err = tctx.CommitTransaction()
			tt.TestFailIfErr(t, err)
		}

		err = db1.Stop(ctx)
		tt.TestFailIfErr(t, err)
	}

	{
		db1New, err := logs.NewLogsDatabase(conf)
		tt.TestFailIfErr(t, err)

		{
			tctx := simplectx.CreateSimpleContext(ctx, nil)

			schema3, err := db1New.ReadSchema(tctx)
			tt.TestFailIfErr(t, err)
			tt.AssertEqual(t, "schema3", schema.LogsSchemaVersion, schema3)

			err = tctx.CommitTransaction()
			tt.TestFailIfErr(t, err)
		}

		err = db1New.Migrate(ctx)
		tt.TestFailIfErr(t, err)

		{
			tctx := simplectx.CreateSimpleContext(ctx, nil)

			schema4, err := db1New.ReadSchema(tctx)
			tt.TestFailIfErr(t, err)
			tt.AssertEqual(t, "schema4", schema.LogsSchemaVersion, schema4)

			err = tctx.CommitTransaction()
			tt.TestFailIfErr(t, err)
		}
	}
}

func TestRequestsDB_Current(t *testing.T) {
	dbf1, dbf2, dbf3, conf, stop := tt.StartSimpleTestspace(t)
	defer stop()

	ctx := context.Background()

	tt.AssertAny(dbf1)
	tt.AssertAny(dbf2)
	tt.AssertAny(dbf3)
	tt.AssertAny(conf)

	{
		db1, err := requests.NewRequestsDatabase(conf)
		tt.TestFailIfErr(t, err)

		{
			tctx := simplectx.CreateSimpleContext(ctx, nil)

			schema1, err := db1.ReadSchema(tctx)
			tt.TestFailIfErr(t, err)
			tt.AssertEqual(t, "schema1", 0, schema1)

			err = tctx.CommitTransaction()
			tt.TestFailIfErr(t, err)
		}

		{
			err = db1.Migrate(ctx)
			tt.TestFailIfErr(t, err)
		}

		{
			tctx := simplectx.CreateSimpleContext(ctx, nil)

			schema2, err := db1.ReadSchema(tctx)
			tt.TestFailIfErr(t, err)
			tt.AssertEqual(t, "schema2", schema.RequestsSchemaVersion, schema2)

			err = tctx.CommitTransaction()
			tt.TestFailIfErr(t, err)
		}

		{
			err = db1.Migrate(ctx)
			tt.TestFailIfErr(t, err)
		}

		{
			tctx := simplectx.CreateSimpleContext(ctx, nil)

			schema2, err := db1.ReadSchema(tctx)
			tt.TestFailIfErr(t, err)
			tt.AssertEqual(t, "schema2", schema.RequestsSchemaVersion, schema2)

			err = tctx.CommitTransaction()
			tt.TestFailIfErr(t, err)
		}

		err = db1.Stop(ctx)
		tt.TestFailIfErr(t, err)
	}

	{
		db1New, err := requests.NewRequestsDatabase(conf)
		tt.TestFailIfErr(t, err)

		{
			tctx := simplectx.CreateSimpleContext(ctx, nil)

			schema3, err := db1New.ReadSchema(tctx)
			tt.TestFailIfErr(t, err)
			tt.AssertEqual(t, "schema3", schema.RequestsSchemaVersion, schema3)

			err = tctx.CommitTransaction()
			tt.TestFailIfErr(t, err)
		}

		err = db1New.Migrate(ctx)
		tt.TestFailIfErr(t, err)

		{
			tctx := simplectx.CreateSimpleContext(ctx, nil)

			schema4, err := db1New.ReadSchema(tctx)
			tt.TestFailIfErr(t, err)
			tt.AssertEqual(t, "schema4", schema.RequestsSchemaVersion, schema4)

			err = tctx.CommitTransaction()
			tt.TestFailIfErr(t, err)
		}
	}
}

func TestPrimaryDB_Migrate_from_3_to_latest(t *testing.T) {
	dbf1, dbf2, dbf3, conf, stop := tt.StartSimpleTestspace(t)
	defer stop()

	ctx := context.Background()

	tt.AssertAny(dbf1)
	tt.AssertAny(dbf2)
	tt.AssertAny(dbf3)
	tt.AssertAny(conf)

	{
		url := fmt.Sprintf("file:%s", dbf1)

		xdb, err := sqlx.Open("sqlite3", url)
		tt.TestFailIfErr(t, err)

		qqdb := sq.NewDB(xdb, sq.DBOptions{})

		schemavers := 3

		dbschema := schema.PrimarySchema[schemavers]

		_, err = qqdb.Exec(ctx, dbschema.SQL, sq.PP{})
		tt.TestFailIfErr(t, err)

		_, err = qqdb.Exec(ctx, "INSERT INTO meta (meta_key, value_int) VALUES (:key, :val) ON CONFLICT(meta_key) DO UPDATE SET value_int = :val", sq.PP{
			"key": "schema",
			"val": schemavers,
		})

		_, err = qqdb.Exec(ctx, "INSERT INTO meta (meta_key, value_txt) VALUES (:key, :val) ON CONFLICT(meta_key) DO UPDATE SET value_txt = :val", sq.PP{
			"key": "schema_hash",
			"val": dbschema.Hash,
		})

		{
			tctx := simplectx.CreateSimpleContext(ctx, nil)
			schemHashDB, err := sq.HashSqliteDatabase(tctx, qqdb)
			tt.TestFailIfErr(t, err)
			tt.AssertEqual(t, "schemHashDB", dbschema.Hash, schemHashDB)
			err = tctx.CommitTransaction()
			tt.TestFailIfErr(t, err)
		}

		err = qqdb.Exit()
		tt.TestFailIfErr(t, err)
	}

	{
		db1, err := primary.NewPrimaryDatabase(conf)
		tt.TestFailIfErr(t, err)

		{
			tctx := simplectx.CreateSimpleContext(ctx, nil)

			schema1, err := db1.ReadSchema(tctx)
			tt.TestFailIfErr(t, err)
			tt.AssertEqual(t, "schema1", 3, schema1)

			err = tctx.CommitTransaction()
			tt.TestFailIfErr(t, err)
		}

		//================================================
		{
			err = db1.Migrate(ctx)
			tt.TestFailIfErr(t, err)
		}
		//================================================

		{
			tctx := simplectx.CreateSimpleContext(ctx, nil)

			schema2, err := db1.ReadSchema(tctx)
			tt.TestFailIfErr(t, err)
			tt.AssertEqual(t, "schema2", schema.PrimarySchemaVersion, schema2)

			err = tctx.CommitTransaction()
			tt.TestFailIfErr(t, err)
		}

		{
			tctx := simplectx.CreateSimpleContext(ctx, nil)
			schemHashDB, err := sq.HashSqliteDatabase(tctx, db1.DB())
			tt.TestFailIfErr(t, err)
			tt.AssertEqual(t, "schemHashDB", schema.PrimarySchema[schema.PrimarySchemaVersion].Hash, schemHashDB)
			err = tctx.CommitTransaction()
			tt.TestFailIfErr(t, err)
		}

		err = db1.Stop(ctx)
		tt.TestFailIfErr(t, err)
	}
}

func TestPrimaryDB_Migrate_from_4_to_latest(t *testing.T) {
	dbf1, dbf2, dbf3, conf, stop := tt.StartSimpleTestspace(t)
	defer stop()

	ctx := context.Background()

	tt.AssertAny(dbf1)
	tt.AssertAny(dbf2)
	tt.AssertAny(dbf3)
	tt.AssertAny(conf)

	{
		url := fmt.Sprintf("file:%s", dbf1)

		xdb, err := sqlx.Open("sqlite3", url)
		tt.TestFailIfErr(t, err)

		qqdb := sq.NewDB(xdb, sq.DBOptions{})

		schemavers := 4

		dbschema := schema.PrimarySchema[schemavers]

		_, err = qqdb.Exec(ctx, dbschema.SQL, sq.PP{})
		tt.TestFailIfErr(t, err)

		_, err = qqdb.Exec(ctx, "INSERT INTO meta (meta_key, value_int) VALUES (:key, :val) ON CONFLICT(meta_key) DO UPDATE SET value_int = :val", sq.PP{
			"key": "schema",
			"val": schemavers,
		})

		_, err = qqdb.Exec(ctx, "INSERT INTO meta (meta_key, value_txt) VALUES (:key, :val) ON CONFLICT(meta_key) DO UPDATE SET value_txt = :val", sq.PP{
			"key": "schema_hash",
			"val": dbschema.Hash,
		})

		{
			tctx := simplectx.CreateSimpleContext(ctx, nil)
			schemHashDB, err := sq.HashSqliteDatabase(tctx, qqdb)
			tt.TestFailIfErr(t, err)
			tt.AssertEqual(t, "schemHashDB", dbschema.Hash, schemHashDB)
			err = tctx.CommitTransaction()
			tt.TestFailIfErr(t, err)
		}

		err = qqdb.Exit()
		tt.TestFailIfErr(t, err)
	}

	{
		db1, err := primary.NewPrimaryDatabase(conf)
		tt.TestFailIfErr(t, err)

		{
			tctx := simplectx.CreateSimpleContext(ctx, nil)

			schema1, err := db1.ReadSchema(tctx)
			tt.TestFailIfErr(t, err)
			tt.AssertEqual(t, "schema1", 4, schema1)

			err = tctx.CommitTransaction()
			tt.TestFailIfErr(t, err)
		}

		//================================================
		{
			err = db1.Migrate(ctx)
			tt.TestFailIfErr(t, err)
		}
		//================================================

		{
			tctx := simplectx.CreateSimpleContext(ctx, nil)

			schema2, err := db1.ReadSchema(tctx)
			tt.TestFailIfErr(t, err)
			tt.AssertEqual(t, "schema2", schema.PrimarySchemaVersion, schema2)

			err = tctx.CommitTransaction()
			tt.TestFailIfErr(t, err)
		}

		{
			tctx := simplectx.CreateSimpleContext(ctx, nil)
			schemHashDB, err := sq.HashSqliteDatabase(tctx, db1.DB())
			tt.TestFailIfErr(t, err)
			tt.AssertEqual(t, "schemHashDB", schema.PrimarySchema[schema.PrimarySchemaVersion].Hash, schemHashDB)
			err = tctx.CommitTransaction()
			tt.TestFailIfErr(t, err)
		}

		err = db1.Stop(ctx)
		tt.TestFailIfErr(t, err)
	}
}
