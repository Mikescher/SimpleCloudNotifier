package primary

import (
	"blackforestbytes.com/simplecloudnotifier/models"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"time"
)

func (db *Database) CreateClient(ctx TxContext, userid models.UserID, ctype models.ClientType, fcmToken string, agentModel string, agentVersion string) (models.Client, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Client{}, err
	}

	now := time.Now().UTC()

	clientid := models.NewClientID()

	_, err = tx.Exec(ctx, "INSERT INTO clients (client_id, user_id, type, fcm_token, timestamp_created, agent_model, agent_version) VALUES (:cid, :uid, :typ, :fcm, :ts, :am, :av)", sq.PP{
		"cid": clientid,
		"uid": userid,
		"typ": string(ctype),
		"fcm": fcmToken,
		"ts":  time2DB(now),
		"am":  agentModel,
		"av":  agentVersion,
	})
	if err != nil {
		return models.Client{}, err
	}

	return models.Client{
		ClientID:         clientid,
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

	_, err = tx.Exec(ctx, "DELETE FROM clients WHERE fcm_token = :fcm", sq.PP{"fcm": fcmtoken})
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

	rows, err := tx.Query(ctx, "SELECT * FROM clients WHERE user_id = :uid ORDER BY clients.timestamp_created DESC, clients.client_id ASC", sq.PP{"uid": userid})
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

	rows, err := tx.Query(ctx, "SELECT * FROM clients WHERE user_id = :uid AND client_id = :cid LIMIT 1", sq.PP{
		"uid": userid,
		"cid": clientid,
	})
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

	_, err = tx.Exec(ctx, "DELETE FROM clients WHERE client_id = :cid", sq.PP{"cid": clientid})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteClientsByFCM(ctx TxContext, fcmtoken string) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "DELETE FROM clients WHERE fcm_token = :fcm", sq.PP{"fcm": fcmtoken})
	if err != nil {
		return err
	}

	return nil
}
