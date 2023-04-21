package primary

import (
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"strings"
	"time"
)

func (db *Database) CreateKeyToken(ctx TxContext, name string, owner models.UserID, allChannels bool, channels []models.ChannelID, permissions models.TokenPermissionList, token string) (models.KeyToken, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.KeyToken{}, err
	}

	now := time.Now().UTC()

	keyTokenid := models.NewKeyTokenID()

	_, err = tx.Exec(ctx, "INSERT INTO keytokens (keytoken_id, name, timestamp_created, owner_user_id, all_channels, channels, token, permissions) VALUES (:tid, :nam, :tsc, :owr, :all, :cha, :tok, :prm)", sq.PP{
		"tid": keyTokenid,
		"nam": name,
		"tsc": time2DB(now),
		"owr": owner.String(),
		"all": bool2DB(allChannels),
		"cha": strings.Join(langext.ArrMap(channels, func(v models.ChannelID) string { return v.String() }), ";"),
		"tok": token,
		"prm": permissions.String(),
	})
	if err != nil {
		return models.KeyToken{}, err
	}

	return models.KeyToken{
		KeyTokenID:        keyTokenid,
		Name:              name,
		TimestampCreated:  now,
		TimestampLastUsed: nil,
		OwnerUserID:       owner,
		AllChannels:       allChannels,
		Channels:          channels,
		Token:             token,
		Permissions:       permissions,
		MessagesSent:      0,
	}, nil
}

func (db *Database) ListKeyTokens(ctx TxContext, ownerID models.UserID) ([]models.KeyToken, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, "SELECT * FROM keytokens WHERE owner_user_id = :uid ORDER BY keytokens.timestamp_created DESC, keytokens.keytoken_id ASC", sq.PP{"uid": ownerID})
	if err != nil {
		return nil, err
	}

	data, err := models.DecodeKeyTokens(rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (db *Database) GetKeyToken(ctx TxContext, userid models.UserID, keyTokenid models.KeyTokenID) (models.KeyToken, error) {
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

	keyToken, err := models.DecodeKeyToken(rows)
	if err != nil {
		return models.KeyToken{}, err
	}

	return keyToken, nil
}

func (db *Database) GetKeyTokenByToken(ctx TxContext, key string) (*models.KeyToken, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, "SELECT * FROM keytokens WHERE token = :key LIMIT 1", sq.PP{"key": key})
	if err != nil {
		return nil, err
	}

	user, err := models.DecodeKeyToken(rows)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *Database) DeleteKeyToken(ctx TxContext, keyTokenid models.KeyTokenID) error {
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

func (db *Database) UpdateKeyTokenName(ctx TxContext, keyTokenid models.KeyTokenID, name string) error {
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

func (db *Database) UpdateKeyTokenPermissions(ctx TxContext, keyTokenid models.KeyTokenID, perm models.TokenPermissionList) error {
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

func (db *Database) UpdateKeyTokenAllChannels(ctx TxContext, keyTokenid models.KeyTokenID, allChannels bool) error {
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

func (db *Database) UpdateKeyTokenChannels(ctx TxContext, keyTokenid models.KeyTokenID, channels []models.ChannelID) error {
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

func (db *Database) IncKeyTokenMessageCounter(ctx TxContext, keyTokenid models.KeyTokenID) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE keytokens SET messages_sent = messages_sent + 1, timestamp_lastused = :ts WHERE keytoken_id = :tid", sq.PP{
		"ts":  time2DB(time.Now()),
		"tid": keyTokenid,
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateKeyTokenLastUsed(ctx TxContext, keyTokenid models.KeyTokenID) error {
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
