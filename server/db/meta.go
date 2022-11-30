package db

import (
	"context"
	"errors"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
)

func (db *Database) ReadSchema(ctx context.Context) (retval int, reterr error) {

	r1, err := db.db.QueryContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name='meta'")
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

	r2, err := db.db.QueryContext(ctx, "SELECT value_int FROM meta WHERE meta_key='schema'")
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

	return dbschema, nil
}

func (db *Database) WriteMetaString(ctx context.Context, key string, value string) error {
	_, err := db.db.ExecContext(ctx, "INSERT INTO meta (meta_key, value_txt) VALUES (?, ?) ON CONFLICT(meta_key) DO UPDATE SET value_txt = ?",
		key,
		value,
		value)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) WriteMetaInt(ctx context.Context, key string, value int64) error {
	_, err := db.db.ExecContext(ctx, "INSERT INTO meta (meta_key, value_int) VALUES (?, ?) ON CONFLICT(meta_key) DO UPDATE SET value_int = ?",
		key,
		value,
		value)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) WriteMetaReal(ctx context.Context, key string, value float64) error {
	_, err := db.db.ExecContext(ctx, "INSERT INTO meta (meta_key, value_real) VALUES (?, ?) ON CONFLICT(meta_key) DO UPDATE SET value_real = ?",
		key,
		value,
		value)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) WriteMetaBlob(ctx context.Context, key string, value []byte) error {
	_, err := db.db.ExecContext(ctx, "INSERT INTO meta (meta_key, value_blob) VALUES (?, ?) ON CONFLICT(meta_key) DO UPDATE SET value_blob = ?",
		key,
		value,
		value)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) ReadMetaString(ctx context.Context, key string) (retval *string, reterr error) {
	r2, err := db.db.QueryContext(ctx, "SELECT value_txt FROM meta WHERE meta_key=?", key)
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

	return langext.Ptr(value), nil
}

func (db *Database) ReadMetaInt(ctx context.Context, key string) (retval *int64, reterr error) {
	r2, err := db.db.QueryContext(ctx, "SELECT value_int FROM meta WHERE meta_key=?", key)
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

	return langext.Ptr(value), nil
}

func (db *Database) ReadMetaReal(ctx context.Context, key string) (retval *float64, reterr error) {
	r2, err := db.db.QueryContext(ctx, "SELECT value_real FROM meta WHERE meta_key=?", key)
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

	return langext.Ptr(value), nil
}

func (db *Database) ReadMetaBlob(ctx context.Context, key string) (retval *[]byte, reterr error) {
	r2, err := db.db.QueryContext(ctx, "SELECT value_blob FROM meta WHERE meta_key=?", key)
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

	return langext.Ptr(value), nil
}

func (db *Database) DeleteMeta(ctx context.Context, key string) error {
	_, err := db.db.ExecContext(ctx, "DELETE FROM meta WHERE meta_key = ?", key)
	if err != nil {
		return err
	}
	return nil
}
