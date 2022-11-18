package db

import (
	"blackforestbytes.com/simplecloudnotifier/api/models"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
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
		QuotaToday:        0,
		QuotaDay:          nil,
		IsPro:             protoken != nil,
		ProToken:          protoken,
	}, nil
}

func (db *Database) CreateClient(ctx TxContext, userid int64, ctype models.ClientType, fcmToken string, agentModel string, agentVersion string) (models.Client, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Client{}, err
	}

	now := time.Now().UTC()

	res, err := tx.ExecContext(ctx, "INSERT INTO clients (user_id, type, fcm_token, timestamp_created, agent_model, agent_version) VALUES (?, ?, ?, ?, ?, ?)",
		userid,
		string(ctype),
		fcmToken,
		time2DB(now),
		agentModel,
		agentVersion)
	if err != nil {
		return models.Client{}, err
	}

	liid, err := res.LastInsertId()
	if err != nil {
		return models.Client{}, err
	}

	return models.Client{
		ClientID:         liid,
		UserID:           userid,
		Type:             ctype,
		FCMToken:         langext.Ptr(fcmToken),
		TimestampCreated: now,
		AgentModel:       agentModel,
		AgentVersion:     agentVersion,
	}, nil
}

func (db *Database) ClearFCMTokens(ctx TxContext, fcmtoken string) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM clients WHERE fcm_token = ?", fcmtoken)
	if err != nil {
		return err
	}

	return nil
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
