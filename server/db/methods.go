package db

import (
	"blackforestbytes.com/simplecloudnotifier/api/models"
	"database/sql"
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

func (db *Database) GetChannelByKey(ctx TxContext, key string) (*models.Channel, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT * FROM channels WHERE subscribe_key = ? OR send_key = ? LIMIT 1", key, key)
	if err != nil {
		return nil, err
	}

	channel, err := models.DecodeChannel(rows)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &channel, nil
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

	_, err = tx.ExecContext(ctx, "UPDATE users SET username = ? WHERE user_id = ?", username, userid)
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

	_, err = tx.ExecContext(ctx, "UPDATE users SET pro_token = ? AND is_pro = ? WHERE user_id = ?", protoken, bool2DB(protoken != nil), userid)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) ListClients(ctx TxContext, userid int64) ([]models.Client, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT * FROM clients WHERE user_id = ?", userid)
	if err != nil {
		return nil, err
	}

	data, err := models.DecodeClients(rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (db *Database) GetClient(ctx TxContext, userid int64, clientid int64) (models.Client, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Client{}, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT * FROM clients WHERE user_id = ? AND client_id = ? LIMIT 1", userid, clientid)
	if err != nil {
		return models.Client{}, err
	}

	client, err := models.DecodeClient(rows)
	if err != nil {
		return models.Client{}, err
	}

	return client, nil
}
