package primary

import (
	"blackforestbytes.com/simplecloudnotifier/db"
	"blackforestbytes.com/simplecloudnotifier/models"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"time"
)

func (db *Database) CreateClient(ctx db.TxContext, userid models.UserID, ctype models.ClientType, fcmToken string, agentModel string, agentVersion string) (models.Client, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Client{}, err
	}

	entity := models.ClientDB{
		ClientID:         models.NewClientID(),
		UserID:           userid,
		Type:             ctype,
		FCMToken:         fcmToken,
		TimestampCreated: time2DB(time.Now()),
		AgentModel:       agentModel,
		AgentVersion:     agentVersion,
	}

	_, err = sq.InsertSingle(ctx, tx, "clients", entity)
	if err != nil {
		return models.Client{}, err
	}

	return entity.Model(), nil
}

func (db *Database) ClearFCMTokens(ctx db.TxContext, fcmtoken string) error {
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

func (db *Database) ListClients(ctx db.TxContext, userid models.UserID) ([]models.Client, error) {
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

func (db *Database) GetClient(ctx db.TxContext, userid models.UserID, clientid models.ClientID) (models.Client, error) {
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

func (db *Database) DeleteClient(ctx db.TxContext, clientid models.ClientID) error {
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

func (db *Database) DeleteClientsByFCM(ctx db.TxContext, fcmtoken string) error {
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

func (db *Database) UpdateClientFCMToken(ctx db.TxContext, clientid models.ClientID, fcmtoken string) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE clients SET fcm_token = :vvv WHERE client_id = :cid", sq.PP{
		"vvv": fcmtoken,
		"cid": clientid,
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateClientAgentModel(ctx db.TxContext, clientid models.ClientID, agentModel string) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE clients SET agent_model = :vvv WHERE client_id = :cid", sq.PP{
		"vvv": agentModel,
		"cid": clientid,
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateClientAgentVersion(ctx db.TxContext, clientid models.ClientID, agentVersion string) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE clients SET agent_version = :vvv WHERE client_id = :cid", sq.PP{
		"vvv": agentVersion,
		"cid": clientid,
	})
	if err != nil {
		return err
	}

	return nil
}
