package db

import (
	"blackforestbytes.com/simplecloudnotifier/models"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"time"
)

func (db *Database) CreateClient(ctx TxContext, userid models.UserID, ctype models.ClientType, fcmToken string, agentModel string, agentVersion string) (models.Client, error) {
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
		ClientID:         models.ClientID(liid),
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

func (db *Database) ListClients(ctx TxContext, userid models.UserID) ([]models.Client, error) {
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

func (db *Database) GetClient(ctx TxContext, userid models.UserID, clientid models.ClientID) (models.Client, error) {
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

func (db *Database) DeleteClient(ctx TxContext, clientid models.ClientID) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM clients WHERE client_id = ?", clientid)
	if err != nil {
		return err
	}

	return nil
}
