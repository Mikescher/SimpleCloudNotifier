package db

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"time"
)

func (db *Database) CreateUser(ctx TxContext, readKey string, sendKey string, adminKey string, protoken *string, username *string) (models.User, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.User{}, err
	}

	now := time.Now().UTC()

	res, err := tx.ExecContext(ctx, "INSERT INTO users (username, read_key, send_key, admin_key, is_pro, pro_token, timestamp_created) VALUES (?, ?, ?, ?, ?, ?, ?)",
		username,
		readKey,
		sendKey,
		adminKey,
		bool2DB(protoken != nil),
		protoken,
		time2DB(now))
	if err != nil {
		return models.User{}, err
	}

	liid, err := res.LastInsertId()
	if err != nil {
		return models.User{}, err
	}

	return models.User{
		UserID:            liid,
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

	_, err = tx.ExecContext(ctx, "UPDATE users SET is_pro=0, pro_token=NULL WHERE pro_token = ?", protoken)
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

	rows, err := tx.QueryContext(ctx, "SELECT * FROM users WHERE admin_key = ? OR send_key = ? OR read_key = ? LIMIT 1", key, key, key)
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

func (db *Database) GetUser(ctx TxContext, userid int64) (models.User, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.User{}, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT * FROM users WHERE user_id = ? LIMIT 1", userid)
	if err != nil {
		return models.User{}, err
	}

	user, err := models.DecodeUser(rows)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (db *Database) UpdateUserUsername(ctx TxContext, userid int64, username *string) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE users SET username = ? WHERE user_id = ?",
		username,
		userid)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateUserProToken(ctx TxContext, userid int64, protoken *string) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE users SET pro_token = ?, is_pro = ? WHERE user_id = ?",
		protoken,
		bool2DB(protoken != nil),
		userid)
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

	_, err = tx.ExecContext(ctx, "UPDATE users SET timestamp_lastsent = ?, messages_sent = ?, quota_used = ?, quota_used_day = ? WHERE user_id = ?",
		time2DB(time.Now()),
		user.MessagesSent+1,
		quota,
		scn.QuotaDayString(),
		user.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateUserLastRead(ctx TxContext, userid int64) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE users SET timestamp_lastread = ? WHERE user_id = ?",
		time2DB(time.Now()),
		userid)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateUserKeys(ctx TxContext, userid int64, sendKey string, readKey string, adminKey string) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE users SET send_key = ?, read_key = ?, admin_key = ? WHERE user_id = ?",
		sendKey,
		readKey,
		adminKey,
		userid)
	if err != nil {
		return err
	}

	return nil
}
