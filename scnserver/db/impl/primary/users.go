package primary

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/models"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"time"
)

func (db *Database) CreateUser(ctx TxContext, protoken *string, username *string) (models.User, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.User{}, err
	}

	entity := models.UserDB{
		UserID:            models.NewUserID(),
		Username:          username,
		TimestampCreated:  time2DB(time.Now()),
		TimestampLastRead: nil,
		TimestampLastSent: nil,
		MessagesSent:      0,
		QuotaUsed:         0,
		QuotaUsedDay:      nil,
		IsPro:             protoken != nil,
		ProToken:          protoken,
	}

	_, err = sq.InsertSingle(ctx, tx, "users", entity)
	if err != nil {
		return models.User{}, err
	}

	return entity.Model(), nil
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

func (db *Database) IncUserMessageCounter(ctx TxContext, user *models.User) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	now := time.Now()

	quota := user.QuotaUsedToday() + 1

	user.QuotaUsed = quota
	user.QuotaUsedDay = langext.Ptr(scn.QuotaDayString())

	_, err = tx.Exec(ctx, "UPDATE users SET timestamp_lastsent = :ts, messages_sent = messages_sent+1, quota_used = :qu, quota_used_day = :qd WHERE user_id = :uid", sq.PP{
		"ts":  time2DB(now),
		"qu":  user.QuotaUsed,
		"qd":  user.QuotaUsedDay,
		"uid": user.UserID,
	})
	if err != nil {
		return err
	}

	user.TimestampLastSent = &now
	user.MessagesSent = user.MessagesSent + 1

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
