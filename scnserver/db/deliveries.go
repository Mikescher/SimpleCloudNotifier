package db

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/models"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"time"
)

func (db *Database) CreateRetryDelivery(ctx TxContext, client models.Client, msg models.Message) (models.Delivery, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Delivery{}, err
	}

	now := time.Now().UTC()
	next := scn.NextDeliveryTimestamp(now)

	res, err := tx.Exec(ctx, "INSERT INTO deliveries (scn_message_id, receiver_user_id, receiver_client_id, timestamp_created, timestamp_finalized, status, fcm_message_id, next_delivery) VALUES (:mid, :ruid, :rcid, :tsc, :tsf, :stat, :fcm, :next)", sq.PP{
		"mid":  msg.SCNMessageID,
		"ruid": client.UserID,
		"rcid": client.ClientID,
		"tsc":  time2DB(now),
		"tsf":  nil,
		"stat": models.DeliveryStatusRetry,
		"fcm":  nil,
		"next": time2DB(next),
	})
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

	res, err := tx.Exec(ctx, "INSERT INTO deliveries (scn_message_id, receiver_user_id, receiver_client_id, timestamp_created, timestamp_finalized, status, fcm_message_id, next_delivery) VALUES (:mid, :ruid, :rcid, :tsc, :tsf, :stat, :fcm, :next)", sq.PP{
		"mid":  msg.SCNMessageID,
		"ruid": client.UserID,
		"rcid": client.ClientID,
		"tsc":  time2DB(now),
		"tsf":  time2DB(now),
		"stat": models.DeliveryStatusSuccess,
		"fcm":  fcmDelivID,
		"next": nil,
	})
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

	rows, err := tx.Query(ctx, "SELECT * FROM deliveries WHERE status = 'RETRY' AND next_delivery < :next LIMIT :lim", sq.PP{
		"next": time2DB(time.Now()),
		"lim":  pageSize,
	})
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

	_, err = tx.Exec(ctx, "UPDATE deliveries SET status = 'SUCCESS', next_delivery = NULL, retry_count = :rc, timestamp_finalized = :ts, fcm_message_id = :fcm WHERE delivery_id = :did", sq.PP{
		"rc":  delivery.RetryCount + 1,
		"ts":  time2DB(time.Now()),
		"fcm": fcmDelivID,
		"did": delivery.DeliveryID,
	})
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

	_, err = tx.Exec(ctx, "UPDATE deliveries SET status = 'FAILED', next_delivery = NULL, retry_count = :rc, timestamp_finalized = :ts WHERE delivery_id = :did",
		sq.PP{
			"rc":  delivery.RetryCount + 1,
			"ts":  time2DB(time.Now()),
			"did": delivery.DeliveryID,
		})
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

	_, err = tx.Exec(ctx, "UPDATE deliveries SET status = 'RETRY', next_delivery = :next, retry_count = :rc WHERE delivery_id = :did", sq.PP{
		"next": scn.NextDeliveryTimestamp(time.Now()),
		"rc":   delivery.RetryCount + 1,
		"did":  delivery.DeliveryID,
	})
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

	_, err = tx.Exec(ctx, "UPDATE deliveries SET status = 'FAILED', next_delivery = NULL, timestamp_finalized = :ts WHERE scn_message_id = :mid AND status = 'RETRY'", sq.PP{
		"ts":  time.Now(),
		"mid": scnMessageID,
	})
	if err != nil {
		return err
	}

	return nil
}
