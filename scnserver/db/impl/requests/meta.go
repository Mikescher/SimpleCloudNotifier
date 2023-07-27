package requests

import (
	"blackforestbytes.com/simplecloudnotifier/db"
	"errors"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
)

func (db *Database) ReadSchema(ctx db.TxContext) (retval int, reterr error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return 0, err
	}

	r1, err := tx.Query(ctx, "SELECT name FROM sqlite_master WHERE type = :typ AND name = :name", sq.PP{"typ": "table", "name": "meta"})
	if err != nil {
		return 0, err
	}
	defer func() {
		err = r1.Close()
		if err != nil {
			// overwrite return values
			retval = 0
			reterr = err
		}
	}()

	if !r1.Next() {
		return 0, nil
	}

	err = r1.Close()
	if err != nil {
		return 0, err
	}

	r2, err := tx.Query(ctx, "SELECT value_int FROM meta WHERE meta_key = :key", sq.PP{"key": "schema"})
	if err != nil {
		return 0, err
	}
	defer func() {
		err = r2.Close()
		if err != nil {
			// overwrite return values
			retval = 0
			reterr = err
		}
	}()

	if !r2.Next() {
		return 0, errors.New("no schema entry in meta table")
	}

	var dbschema int
	err = r2.Scan(&dbschema)
	if err != nil {
		return 0, err
	}

	err = r2.Close()
	if err != nil {
		return 0, err
	}

	return dbschema, nil
}

func (db *Database) WriteMetaString(ctx db.TxContext, key string, value string) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "INSERT INTO meta (meta_key, value_txt) VALUES (:key, :val) ON CONFLICT(meta_key) DO UPDATE SET value_txt = :val", sq.PP{
		"key": key,
		"val": value,
	})
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) WriteMetaInt(ctx db.TxContext, key string, value int64) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "INSERT INTO meta (meta_key, value_int) VALUES (:key, :val) ON CONFLICT(meta_key) DO UPDATE SET value_int = :val", sq.PP{
		"key": key,
		"val": value,
	})
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) WriteMetaReal(ctx db.TxContext, key string, value float64) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "INSERT INTO meta (meta_key, value_real) VALUES (:key, :val) ON CONFLICT(meta_key) DO UPDATE SET value_real = :val", sq.PP{
		"key": key,
		"val": value,
	})
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) WriteMetaBlob(ctx db.TxContext, key string, value []byte) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "INSERT INTO meta (meta_key, value_blob) VALUES (:key, :val) ON CONFLICT(meta_key) DO UPDATE SET value_blob = :val", sq.PP{
		"key": key,
		"val": value,
	})
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) ReadMetaString(ctx db.TxContext, key string) (retval *string, reterr error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	r2, err := tx.Query(ctx, "SELECT value_txt FROM meta WHERE meta_key = :key", sq.PP{"key": key})
	if err != nil {
		return nil, err
	}
	defer func() {
		err = r2.Close()
		if err != nil {
			// overwrite return values
			retval = nil
			reterr = err
		}
	}()
	if !r2.Next() {
		return nil, errors.New("no matching entry in meta table")
	}

	var value string
	err = r2.Scan(&value)
	if err != nil {
		return nil, err
	}

	err = r2.Close()
	if err != nil {
		return nil, err
	}

	return langext.Ptr(value), nil
}

func (db *Database) ReadMetaInt(ctx db.TxContext, key string) (retval *int64, reterr error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	r2, err := tx.Query(ctx, "SELECT value_int FROM meta WHERE meta_key = :key", sq.PP{"key": key})
	if err != nil {
		return nil, err
	}
	defer func() {
		err = r2.Close()
		if err != nil {
			// overwrite return values
			retval = nil
			reterr = err
		}
	}()

	if !r2.Next() {
		return nil, errors.New("no matching entry in meta table")
	}

	var value int64
	err = r2.Scan(&value)
	if err != nil {
		return nil, err
	}

	err = r2.Close()
	if err != nil {
		return nil, err
	}

	return langext.Ptr(value), nil
}

func (db *Database) ReadMetaReal(ctx db.TxContext, key string) (retval *float64, reterr error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	r2, err := tx.Query(ctx, "SELECT value_real FROM meta WHERE meta_key = :key", sq.PP{"key": key})
	if err != nil {
		return nil, err
	}
	defer func() {
		err = r2.Close()
		if err != nil {
			// overwrite return values
			retval = nil
			reterr = err
		}
	}()

	if !r2.Next() {
		return nil, errors.New("no matching entry in meta table")
	}

	var value float64
	err = r2.Scan(&value)
	if err != nil {
		return nil, err
	}

	err = r2.Close()
	if err != nil {
		return nil, err
	}

	return langext.Ptr(value), nil
}

func (db *Database) ReadMetaBlob(ctx db.TxContext, key string) (retval *[]byte, reterr error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	r2, err := tx.Query(ctx, "SELECT value_blob FROM meta WHERE meta_key = :key", sq.PP{"key": key})
	if err != nil {
		return nil, err
	}
	defer func() {
		err = r2.Close()
		if err != nil {
			// overwrite return values
			retval = nil
			reterr = err
		}
	}()

	if !r2.Next() {
		return nil, errors.New("no matching entry in meta table")
	}

	var value []byte
	err = r2.Scan(&value)
	if err != nil {
		return nil, err
	}

	err = r2.Close()
	if err != nil {
		return nil, err
	}

	return langext.Ptr(value), nil
}

func (db *Database) DeleteMeta(ctx db.TxContext, key string) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "DELETE FROM meta WHERE meta_key = :key", sq.PP{"key": key})
	if err != nil {
		return err
	}
	return nil
}
