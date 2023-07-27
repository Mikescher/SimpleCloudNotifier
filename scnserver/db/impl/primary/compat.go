package primary

import (
	"blackforestbytes.com/simplecloudnotifier/db"
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"errors"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
)

func (db *Database) CreateCompatID(ctx db.TxContext, idtype string, newid string) (int64, error) {
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

func (db *Database) ConvertCompatID(ctx db.TxContext, oldid int64, idtype string) (*string, error) {
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

func (db *Database) ConvertToCompatID(ctx db.TxContext, newid string) (*int64, *string, error) {
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

func (db *Database) ConvertToCompatIDOrCreate(ctx db.TxContext, idtype string, newid string) (int64, error) {
	id1, _, err := db.ConvertToCompatID(ctx, newid)
	if err != nil {
		return 0, err
	}
	if id1 != nil {
		return *id1, nil
	}

	id2, err := db.CreateCompatID(ctx, idtype, newid)
	if err != nil {
		return 0, err
	}
	return id2, nil
}

func (db *Database) GetAck(ctx db.TxContext, msgid models.MessageID) (bool, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return false, err
	}

	rows, err := tx.Query(ctx, "SELECT * FROM compat_acks WHERE message_id = :msgid LIMIT 1", sq.PP{
		"msgid": msgid,
	})
	if err != nil {
		return false, err
	}

	res := rows.Next()

	err = rows.Close()
	if err != nil {
		return false, err
	}

	return res, nil
}

func (db *Database) SetAck(ctx db.TxContext, userid models.UserID, msgid models.MessageID) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "INSERT INTO compat_acks (user_id, message_id) VALUES (:uid, :mid)", sq.PP{
		"uid": userid,
		"mid": msgid,
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) IsCompatClient(ctx db.TxContext, clientid models.ClientID) (bool, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return false, err
	}

	rows, err := tx.Query(ctx, "SELECT * FROM compat_clients WHERE client_id = :id LIMIT 1", sq.PP{
		"id": clientid,
	})
	if err != nil {
		return false, err
	}

	res := rows.Next()

	err = rows.Close()
	if err != nil {
		return false, err
	}

	return res, nil
}
