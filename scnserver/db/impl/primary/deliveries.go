package primary

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/db"
	"blackforestbytes.com/simplecloudnotifier/models"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"time"
)

func (db *Database) CreateRetryDelivery(ctx db.TxContext, client models.Client, msg models.Message) (models.Delivery, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Delivery{}, err
	}

	now := time.Now()
	next := scn.NextDeliveryTimestamp(now)

	entity := models.Delivery{
		DeliveryID:         models.NewDeliveryID(),
		MessageID:          msg.MessageID,
		ReceiverUserID:     client.UserID,
		ReceiverClientID:   client.ClientID,
		TimestampCreated:   models.NewSCNTime(now),
		TimestampFinalized: nil,
		Status:             models.DeliveryStatusRetry,
		RetryCount:         0,
		NextDelivery:       models.NewSCNTimePtr(&next),
		FCMMessageID:       nil,
	}

	_, err = sq.InsertSingle(ctx, tx, "deliveries", entity)
	if err != nil {
		return models.Delivery{}, err
	}

	return entity, nil
}

func (db *Database) CreateSuccessDelivery(ctx db.TxContext, client models.Client, msg models.Message, fcmDelivID string) (models.Delivery, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Delivery{}, err
	}

	now := time.Now()

	entity := models.Delivery{
		DeliveryID:         models.NewDeliveryID(),
		MessageID:          msg.MessageID,
		ReceiverUserID:     client.UserID,
		ReceiverClientID:   client.ClientID,
		TimestampCreated:   models.NewSCNTime(now),
		TimestampFinalized: models.NewSCNTimePtr(&now),
		Status:             models.DeliveryStatusSuccess,
		RetryCount:         0,
		NextDelivery:       nil,
		FCMMessageID:       langext.Ptr(fcmDelivID),
	}

	_, err = sq.InsertSingle(ctx, tx, "deliveries", entity)
	if err != nil {
		return models.Delivery{}, err
	}

	return entity, nil
}

func (db *Database) ListRetrieableDeliveries(ctx db.TxContext, pageSize int) ([]models.Delivery, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	return sq.QueryAll[models.Delivery](ctx, tx, "SELECT * FROM deliveries WHERE status = 'RETRY' AND next_delivery < :next ORDER BY next_delivery ASC LIMIT :lim", sq.PP{
		"next": time2DB(time.Now()),
		"lim":  pageSize,
	}, sq.SModeExtended, sq.Safe)
}

func (db *Database) SetDeliverySuccess(ctx db.TxContext, delivery models.Delivery, fcmDelivID string) error {
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

func (db *Database) SetDeliveryFailed(ctx db.TxContext, delivery models.Delivery) error {
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

func (db *Database) SetDeliveryRetry(ctx db.TxContext, delivery models.Delivery) error {
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

func (db *Database) CancelPendingDeliveries(ctx db.TxContext, messageID models.MessageID) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE deliveries SET status = 'FAILED', next_delivery = NULL, timestamp_finalized = :ts WHERE message_id = :mid AND status = 'RETRY'", sq.PP{
		"ts":  time.Now(),
		"mid": messageID,
	})
	if err != nil {
		return err
	}

	return nil
}
