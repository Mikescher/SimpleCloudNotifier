package primary

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/db"
	"blackforestbytes.com/simplecloudnotifier/models"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"time"
)

func (db *Database) CreateUser(ctx db.TxContext, protoken *string, username *string) (models.User, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.User{}, err
	}

	entity := models.User{
		UserID:            models.NewUserID(),
		Username:          username,
		TimestampCreated:  models.NowSCNTime(),
		TimestampLastRead: nil,
		TimestampLastSent: nil,
		MessagesSent:      0,
		QuotaUsed:         0,
		QuotaUsedDay:      nil,
		IsPro:             protoken != nil,
		ProToken:          protoken,
		UserExtra:         models.UserExtra{},
	}

	entity.PreMarshal()

	_, err = sq.InsertSingle(ctx, tx, "users", entity)
	if err != nil {
		return models.User{}, err
	}

	return entity, nil
}

func (db *Database) ClearProTokens(ctx db.TxContext, protoken string) error {
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

func (db *Database) GetUser(ctx db.TxContext, userid models.UserID) (models.User, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.User{}, err
	}

	return sq.QuerySingle[models.User](ctx, tx, "SELECT * FROM users WHERE user_id = :uid LIMIT 1", sq.PP{"uid": userid}, sq.SModeExtended, sq.Safe)
}

func (db *Database) GetUserOpt(ctx db.TxContext, userid models.UserID) (*models.User, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	return sq.QuerySingleOpt[models.User](ctx, tx, "SELECT * FROM users WHERE user_id = :uid LIMIT 1", sq.PP{"uid": userid}, sq.SModeExtended, sq.Safe)
}

func (db *Database) UpdateUserUsername(ctx db.TxContext, userid models.UserID, username *string) error {
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

func (db *Database) UpdateUserProToken(ctx db.TxContext, userid models.UserID, protoken *string) error {
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

func (db *Database) IncUserMessageCounter(ctx db.TxContext, user *models.User) error {
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

	user.TimestampLastSent = models.NewSCNTimePtr(&now)
	user.MessagesSent = user.MessagesSent + 1

	return nil
}

func (db *Database) UpdateUserLastRead(ctx db.TxContext, userid models.UserID) error {
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
