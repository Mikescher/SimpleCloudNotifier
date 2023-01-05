package primary

import (
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"time"
)

func (db *Database) CreateSubscription(ctx TxContext, subscriberUID models.UserID, channel models.Channel, confirmed bool) (models.Subscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Subscription{}, err
	}

	now := time.Now().UTC()

	res, err := tx.Exec(ctx, "INSERT INTO subscriptions (subscriber_user_id, channel_owner_user_id, channel_internal_name, channel_id, timestamp_created, confirmed) VALUES (:suid, :ouid, :cnam, :cid, :ts, :conf)", sq.PP{
		"suid": subscriberUID,
		"ouid": channel.OwnerUserID,
		"cnam": channel.InternalName,
		"cid":  channel.ChannelID,
		"ts":   time2DB(now),
		"conf": confirmed,
	})
	if err != nil {
		return models.Subscription{}, err
	}

	liid, err := res.LastInsertId()
	if err != nil {
		return models.Subscription{}, err
	}

	return models.Subscription{
		SubscriptionID:      models.SubscriptionID(liid),
		SubscriberUserID:    subscriberUID,
		ChannelOwnerUserID:  channel.OwnerUserID,
		ChannelID:           channel.ChannelID,
		ChannelInternalName: channel.InternalName,
		TimestampCreated:    now,
		Confirmed:           confirmed,
	}, nil
}

func (db *Database) ListSubscriptionsByChannel(ctx TxContext, channelID models.ChannelID) ([]models.Subscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, "SELECT * FROM subscriptions WHERE channel_id = :cid", sq.PP{"cid": channelID})
	if err != nil {
		return nil, err
	}

	data, err := models.DecodeSubscriptions(rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (db *Database) ListSubscriptionsByChannelOwner(ctx TxContext, ownerUserID models.UserID, confirmed *bool) ([]models.Subscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	cond := ""
	if confirmed != nil && *confirmed {
		cond = " AND confirmed = 1"
	} else if confirmed != nil && !*confirmed {
		cond = " AND confirmed = 0"
	}

	rows, err := tx.Query(ctx, "SELECT * FROM subscriptions WHERE channel_owner_user_id = :ouid"+cond, sq.PP{"ouid": ownerUserID})
	if err != nil {
		return nil, err
	}

	data, err := models.DecodeSubscriptions(rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (db *Database) ListSubscriptionsBySubscriber(ctx TxContext, subscriberUserID models.UserID, confirmed *bool) ([]models.Subscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	cond := ""
	if confirmed != nil && *confirmed {
		cond = " AND confirmed = 1"
	} else if confirmed != nil && !*confirmed {
		cond = " AND confirmed = 0"
	}

	rows, err := tx.Query(ctx, "SELECT * FROM subscriptions WHERE subscriber_user_id = :suid"+cond, sq.PP{"suid": subscriberUserID})
	if err != nil {
		return nil, err
	}

	data, err := models.DecodeSubscriptions(rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (db *Database) GetSubscription(ctx TxContext, subid models.SubscriptionID) (models.Subscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Subscription{}, err
	}

	rows, err := tx.Query(ctx, "SELECT * FROM subscriptions WHERE subscription_id = :sid LIMIT 1", sq.PP{"sid": subid})
	if err != nil {
		return models.Subscription{}, err
	}

	sub, err := models.DecodeSubscription(rows)
	if err != nil {
		return models.Subscription{}, err
	}

	return sub, nil
}

func (db *Database) GetSubscriptionBySubscriber(ctx TxContext, subscriberId models.UserID, channelId models.ChannelID) (*models.Subscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, "SELECT * FROM subscriptions WHERE subscriber_user_id = :suid AND channel_id = :cid LIMIT 1", sq.PP{
		"suid": subscriberId,
		"cid":  channelId,
	})
	if err != nil {
		return nil, err
	}

	user, err := models.DecodeSubscription(rows)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *Database) DeleteSubscription(ctx TxContext, subid models.SubscriptionID) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "DELETE FROM subscriptions WHERE subscription_id = :sid", sq.PP{"sid": subid})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateSubscriptionConfirmed(ctx TxContext, subscriptionID models.SubscriptionID, confirmed bool) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE subscriptions SET confirmed = :conf WHERE subscription_id = :sid", sq.PP{
		"conf": confirmed,
		"sid":  subscriptionID,
	})
	if err != nil {
		return err
	}

	return nil
}
