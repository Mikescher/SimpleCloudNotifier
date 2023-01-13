package primary

import (
	"database/sql"
	"errors"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
)

func (db *Database) CreateCompatID(ctx TxContext, idtype string, newid string) (int64, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return 0, err
	}

	rows, err := tx.Query(ctx, "SELECT COALESCE(MAX(old), 0) FROM compat_ids", sq.PP{})
	if err != nil {
		return 0, err
	}

	if !rows.Next() {
		return 0, errors.New("failed to query MAX(old)")
	}

	var oldid int64
	err = rows.Scan(&oldid)
	if err != nil {
		return 0, err
	}

	oldid++

	_, err = tx.Exec(ctx, "INSERT INTO compat_ids (old, new, type) VALUES (:old, :new, :typ)", sq.PP{
		"old": oldid,
		"new": newid,
		"typ": idtype,
	})
	if err != nil {
		return 0, err
	}

	return oldid, nil
}

func (db *Database) ConvertCompatID(ctx TxContext, oldid int64, idtype string) (*string, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, "SELECT new FROM compat_ids WHERE old = :old AND type = :typ", sq.PP{
		"old": oldid,
		"typ": idtype,
	})
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, nil
	}

	var newid string
	err = rows.Scan(&newid)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &newid, nil
}

func (db *Database) ConvertToCompatID(ctx TxContext, newid string) (*int64, *string, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, nil, err
	}

	rows, err := tx.Query(ctx, "SELECT old, type FROM compat_ids WHERE new = :new", sq.PP{"new": newid})
	if err != nil {
		return nil, nil, err
	}

	if !rows.Next() {
		return nil, nil, nil
	}

	var oldid int64
	var idtype string
	err = rows.Scan(&oldid, &idtype)
	if err == sql.ErrNoRows {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}

	return &oldid, &idtype, nil
}
