package primary

import (
	"blackforestbytes.com/simplecloudnotifier/db"
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"errors"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"strings"
	"time"
)

func (db *Database) CreateKeyToken(ctx db.TxContext, name string, owner models.UserID, allChannels bool, channels []models.ChannelID, permissions models.TokenPermissionList, token string) (models.KeyToken, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.KeyToken{}, err
	}

	entity := models.KeyTokenDB{
		KeyTokenID:        models.NewKeyTokenID(),
		Name:              name,
		TimestampCreated:  time2DB(time.Now()),
		TimestampLastUsed: nil,
		OwnerUserID:       owner,
		AllChannels:       allChannels,
		Channels:          strings.Join(langext.ArrMap(channels, func(v models.ChannelID) string { return v.String() }), ";"),
		Token:             token,
		Permissions:       permissions.String(),
		MessagesSent:      0,
	}

	_, err = sq.InsertSingle(ctx, tx, "keytokens", entity)
	if err != nil {
		return models.KeyToken{}, err
	}

	return entity.Model(), nil
}

func (db *Database) ListKeyTokens(ctx db.TxContext, ownerID models.UserID) ([]models.KeyToken, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, "SELECT * FROM keytokens WHERE owner_user_id = :uid ORDER BY keytokens.timestamp_created DESC, keytokens.keytoken_id ASC", sq.PP{"uid": ownerID})
	if err != nil {
		return nil, err
	}

	data, err := models.DecodeKeyTokens(ctx, tx, rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (db *Database) GetKeyToken(ctx db.TxContext, userid models.UserID, keyTokenid models.KeyTokenID) (models.KeyToken, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.KeyToken{}, err
	}

	rows, err := tx.Query(ctx, "SELECT * FROM keytokens WHERE owner_user_id = :uid AND keytoken_id = :cid LIMIT 1", sq.PP{
		"uid": userid,
		"cid": keyTokenid,
	})
	if err != nil {
		return models.KeyToken{}, err
	}

	keyToken, err := models.DecodeKeyToken(ctx, tx, rows)
	if err != nil {
		return models.KeyToken{}, err
	}

	return keyToken, nil
}

func (db *Database) GetKeyTokenByID(ctx db.TxContext, keyTokenid models.KeyTokenID) (models.KeyToken, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.KeyToken{}, err
	}

	rows, err := tx.Query(ctx, "SELECT * FROM keytokens WHERE keytoken_id = :cid LIMIT 1", sq.PP{
		"cid": keyTokenid,
	})
	if err != nil {
		return models.KeyToken{}, err
	}

	keyToken, err := models.DecodeKeyToken(ctx, tx, rows)
	if err != nil {
		return models.KeyToken{}, err
	}

	return keyToken, nil
}

func (db *Database) GetKeyTokenByToken(ctx db.TxContext, key string) (*models.KeyToken, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, "SELECT * FROM keytokens WHERE token = :key LIMIT 1", sq.PP{"key": key})
	if err != nil {
		return nil, err
	}

	user, err := models.DecodeKeyToken(ctx, tx, rows)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *Database) DeleteKeyToken(ctx db.TxContext, keyTokenid models.KeyTokenID) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "DELETE FROM keytokens WHERE keytoken_id = :tid", sq.PP{"tid": keyTokenid})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateKeyTokenName(ctx db.TxContext, keyTokenid models.KeyTokenID, name string) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE keytokens SET name = :nam WHERE keytoken_id = :tid", sq.PP{
		"nam": name,
		"tid": keyTokenid,
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateKeyTokenPermissions(ctx db.TxContext, keyTokenid models.KeyTokenID, perm models.TokenPermissionList) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE keytokens SET permissions = :prm WHERE keytoken_id = :tid", sq.PP{
		"tid": keyTokenid,
		"prm": perm.String(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateKeyTokenAllChannels(ctx db.TxContext, keyTokenid models.KeyTokenID, allChannels bool) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE keytokens SET all_channels = :all WHERE keytoken_id = :tid", sq.PP{
		"tid": keyTokenid,
		"all": bool2DB(allChannels),
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateKeyTokenChannels(ctx db.TxContext, keyTokenid models.KeyTokenID, channels []models.ChannelID) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE keytokens SET channels = :cha WHERE keytoken_id = :tid", sq.PP{
		"tid": keyTokenid,
		"cha": strings.Join(langext.ArrMap(channels, func(v models.ChannelID) string { return v.String() }), ";"),
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) IncKeyTokenMessageCounter(ctx db.TxContext, keyToken *models.KeyToken) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	now := time.Now()

	_, err = tx.Exec(ctx, "UPDATE keytokens SET messages_sent = messages_sent+1, timestamp_lastused = :ts WHERE keytoken_id = :tid", sq.PP{
		"ts":  time2DB(now),
		"tid": keyToken.KeyTokenID,
	})
	if err != nil {
		return err
	}

	keyToken.TimestampLastUsed = &now
	keyToken.MessagesSent += 1

	return nil
}

func (db *Database) UpdateKeyTokenLastUsed(ctx db.TxContext, keyTokenid models.KeyTokenID) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE keytokens SET timestamp_lastused = :ts WHERE keytoken_id = :tid", sq.PP{
		"ts":  time2DB(time.Now()),
		"tid": keyTokenid,
	})
	if err != nil {
		return err
	}

	return nil
}
