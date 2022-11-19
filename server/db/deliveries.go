package db

import (
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
	next := now.Add(5 * time.Second)

	res, err := tx.ExecContext(ctx, "INSERT INTO deliveries (scn_message_id, receiver_user_id, receiver_client_id, timestamp_created, timestamp_finalized, status, fcm_message_id, next_delivery) VALUES (?, ?, ?, ?, ?, ?, ?)",
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
		DeliveryID:         liid,
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

	res, err := tx.ExecContext(ctx, "INSERT INTO deliveries (scn_message_id, receiver_user_id, receiver_client_id, timestamp_created, timestamp_finalized, status, fcm_message_id, next_delivery) VALUES (?, ?, ?, ?, ?, ?, ?)",
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
		DeliveryID:         liid,
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
