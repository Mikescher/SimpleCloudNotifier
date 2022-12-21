package db

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"time"
)

func (db *Database) CreateUser(ctx TxContext, readKey string, sendKey string, adminKey string, protoken *string, username *string) (models.User, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.User{}, err
	}

	now := time.Now().UTC()

	res, err := tx.Exec(ctx, "INSERT INTO users (username, read_key, send_key, admin_key, is_pro, pro_token, timestamp_created) VALUES (:un, :rk, :sk, :ak, :pro, :tok, :ts)", sq.PP{
		"un":  username,
		"rk":  readKey,
		"sk":  sendKey,
		"ak":  adminKey,
		"pro": bool2DB(protoken != nil),
		"tok": protoken,
		"ts":  time2DB(now),
	})
	if err != nil {
		return models.User{}, err
	}

	liid, err := res.LastInsertId()
	if err != nil {
		return models.User{}, err
	}

	return models.User{
		UserID:            models.UserID(liid),
		Username:          username,
		ReadKey:           readKey,
		SendKey:           sendKey,
		AdminKey:          adminKey,
		TimestampCreated:  now,
		TimestampLastRead: nil,
		TimestampLastSent: nil,
		MessagesSent:      0,
		QuotaUsed:         0,
		QuotaUsedDay:      nil,
		IsPro:             protoken != nil,
		ProToken:          protoken,
	}, nil
}

func (db *Database) ClearProTokens(ctx TxContext, protoken string) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE users SET is_pro=0, pro_token=NULL WHERE pro_token = :tok", sq.PP{"tok": protoken})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetUserByKey(ctx TxContext, key string) (*models.User, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, "SELECT * FROM users WHERE admin_key = :key OR send_key = :key OR read_key = :key LIMIT 1", sq.PP{"key": key})
	if err != nil {
		return nil, err
	}

	user, err := models.DecodeUser(rows)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *Database) GetUser(ctx TxContext, userid models.UserID) (models.User, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.User{}, err
	}

	rows, err := tx.Query(ctx, "SELECT * FROM users WHERE user_id = :uid LIMIT 1", sq.PP{"uid": userid})
	if err != nil {
		return models.User{}, err
	}

	user, err := models.DecodeUser(rows)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (db *Database) UpdateUserUsername(ctx TxContext, userid models.UserID, username *string) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE users SET username = :nam WHERE user_id = :uid", sq.PP{
		"nam": username,
		"uid": userid,
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateUserProToken(ctx TxContext, userid models.UserID, protoken *string) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE users SET pro_token = :tok, is_pro = :pro WHERE user_id = :uid", sq.PP{
		"tok": protoken,
		"pro": bool2DB(protoken != nil),
		"uid": userid,
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) IncUserMessageCounter(ctx TxContext, user models.User) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	quota := user.QuotaUsedToday() + 1

	_, err = tx.Exec(ctx, "UPDATE users SET timestamp_lastsent = :ts, messages_sent = :ctr, quota_used = :qu, quota_used_day = :qd WHERE user_id = :uid", sq.PP{
		"ts":  time2DB(time.Now()),
		"ctr": user.MessagesSent + 1,
		"qu":  quota,
		"qd":  scn.QuotaDayString(),
		"uid": user.UserID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateUserLastRead(ctx TxContext, userid models.UserID) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE users SET timestamp_lastread = :ts WHERE user_id = :uid", sq.PP{
		"ts":  time2DB(time.Now()),
		"uid": userid,
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateUserKeys(ctx TxContext, userid models.UserID, sendKey string, readKey string, adminKey string) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE users SET send_key = :sk, read_key = :rk, admin_key = :ak WHERE user_id = :uid", sq.PP{
		"sk":  sendKey,
		"rk":  readKey,
		"ak":  adminKey,
		"uid": userid,
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateUserSendKey(ctx TxContext, userid models.UserID, newkey string) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE users SET send_key = :sk WHERE user_id = :uid", sq.PP{
		"sk":  newkey,
		"uid": userid,
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateUserReadKey(ctx TxContext, userid models.UserID, newkey string) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE users SET read_key = :rk WHERE user_id = :uid", sq.PP{
		"rk":  newkey,
		"uid": userid,
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateUserAdminKey(ctx TxContext, userid models.UserID, newkey string) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE users SET admin_key = :ak WHERE user_id = :uid", sq.PP{
		"ak":  newkey,
		"uid": userid,
	})
	if err != nil {
		return err
	}

	return nil
}
