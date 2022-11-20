package db

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/models"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"time"
)

func (db *Database) CreateRetryDelivery(ctx TxContext, client models.Client, msg models.Message) (models.Delivery, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Delivery{}, err
	}

	now := time.Now().UTC()
	next := scn.NextDeliveryTimestamp(now)

	res, err := tx.ExecContext(ctx, "INSERT INTO deliveries (scn_message_id, receiver_user_id, receiver_client_id, timestamp_created, timestamp_finalized, status, fcm_message_id, next_delivery) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		msg.SCNMessageID,
		client.UserID,
		client.ClientID,
		time2DB(now),
		nil,
		models.DeliveryStatusRetry,
		nil,
		time2DB(next))
	if err != nil {
		return models.Delivery{}, err
	}

	liid, err := res.LastInsertId()
	if err != nil {
		return models.Delivery{}, err
	}

	return models.Delivery{
		DeliveryID:         models.DeliveryID(liid),
		SCNMessageID:       msg.SCNMessageID,
		ReceiverUserID:     client.UserID,
		ReceiverClientID:   client.ClientID,
		TimestampCreated:   now,
		TimestampFinalized: nil,
		Status:             models.DeliveryStatusRetry,
		RetryCount:         0,
		NextDelivery:       langext.Ptr(next),
		FCMMessageID:       nil,
	}, nil
}

func (db *Database) CreateSuccessDelivery(ctx TxContext, client models.Client, msg models.Message, fcmDelivID string) (models.Delivery, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Delivery{}, err
	}

	now := time.Now().UTC()

	res, err := tx.ExecContext(ctx, "INSERT INTO deliveries (scn_message_id, receiver_user_id, receiver_client_id, timestamp_created, timestamp_finalized, status, fcm_message_id, next_delivery) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		msg.SCNMessageID,
		client.UserID,
		client.ClientID,
		time2DB(now),
		time2DB(now),
		models.DeliveryStatusSuccess,
		fcmDelivID,
		nil)
	if err != nil {
		return models.Delivery{}, err
	}

	liid, err := res.LastInsertId()
	if err != nil {
		return models.Delivery{}, err
	}

	return models.Delivery{
		DeliveryID:         models.DeliveryID(liid),
		SCNMessageID:       msg.SCNMessageID,
		ReceiverUserID:     client.UserID,
		ReceiverClientID:   client.ClientID,
		TimestampCreated:   now,
		TimestampFinalized: langext.Ptr(now),
		Status:             models.DeliveryStatusSuccess,
		RetryCount:         0,
		NextDelivery:       nil,
		FCMMessageID:       langext.Ptr(fcmDelivID),
	}, nil
}

func (db *Database) ListRetrieableDeliveries(ctx TxContext, pageSize int) ([]models.Delivery, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT * FROM deliveries WHERE status = 'RETRY' AND next_delivery < ? LIMIT ?",
		time2DB(time.Now()),
		pageSize)
	if err != nil {
		return nil, err
	}

	data, err := models.DecodeDeliveries(rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (db *Database) SetDeliverySuccess(ctx TxContext, delivery models.Delivery, fcmDelivID string) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE deliveries SET status = 'SUCCESS', next_delivery = NULL, retry_count = ?, timestamp_finalized = ?, fcm_message_id = ? WHERE delivery_id = ?",
		delivery.RetryCount+1,
		time2DB(time.Now()),
		fcmDelivID,
		delivery.DeliveryID)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) SetDeliveryFailed(ctx TxContext, delivery models.Delivery) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE deliveries SET status = 'FAILED', next_delivery = NULL, retry_count = ?, timestamp_finalized = ? WHERE delivery_id = ?",
		delivery.RetryCount+1,
		time2DB(time.Now()),
		delivery.DeliveryID)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) SetDeliveryRetry(ctx TxContext, delivery models.Delivery) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE deliveries SET status = 'RETRY', next_delivery = ?, retry_count = ? WHERE delivery_id = ?",
		scn.NextDeliveryTimestamp(time.Now()),
		delivery.RetryCount+1,
		delivery.DeliveryID)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) CancelPendingDeliveries(ctx TxContext, scnMessageID models.SCNMessageID) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE deliveries SET status = 'FAILED', next_delivery = NULL, timestamp_finalized = ? WHERE scn_message_id = ? AND status = 'RETRY'",
		time.Now(),
		scnMessageID)
	if err != nil {
		return err
	}

	return nil
}
