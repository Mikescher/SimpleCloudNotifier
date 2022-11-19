package db

import (
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"time"
)

func (db *Database) CreateSubscription(ctx TxContext, subscriberUID int64, channel models.Channel, confirmed bool) (models.Subscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Subscription{}, err
	}

	now := time.Now().UTC()

	res, err := tx.ExecContext(ctx, "INSERT INTO subscriptions (subscriber_user_id, channel_owner_user_id, channel_name, channel_id, timestamp_created, confirmed) VALUES (?, ?, ?, ?, ?, ?)",
		subscriberUID,
		channel.OwnerUserID,
		channel.Name,
		channel.ChannelID,
		time2DB(now),
		confirmed)
	if err != nil {
		return models.Subscription{}, err
	}

	liid, err := res.LastInsertId()
	if err != nil {
		return models.Subscription{}, err
	}

	return models.Subscription{
		SubscriptionID:     liid,
		SubscriberUserID:   subscriberUID,
		ChannelOwnerUserID: channel.OwnerUserID,
		ChannelID:          channel.ChannelID,
		ChannelName:        channel.Name,
		TimestampCreated:   now,
		Confirmed:          confirmed,
	}, nil
}

func (db *Database) ListSubscriptionsByChannel(ctx TxContext, channelID int64) ([]models.Subscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT * FROM subscriptions WHERE channel_id = ?", channelID)
	if err != nil {
		return nil, err
	}

	data, err := models.DecodeSubscriptions(rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (db *Database) ListSubscriptionsByOwner(ctx TxContext, ownerUserID int64) ([]models.Subscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT * FROM subscriptions WHERE channel_owner_user_id = ?", ownerUserID)
	if err != nil {
		return nil, err
	}

	data, err := models.DecodeSubscriptions(rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (db *Database) GetSubscription(ctx TxContext, subid int64) (models.Subscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Subscription{}, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT * FROM subscriptions WHERE subscription_id = ? LIMIT 1", subid)
	if err != nil {
		return models.Subscription{}, err
	}

	sub, err := models.DecodeSubscription(rows)
	if err != nil {
		return models.Subscription{}, err
	}

	return sub, nil
}

func (db *Database) GetSubscriptionBySubscriber(ctx TxContext, subscriberId int64, channelId int64) (*models.Subscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT * FROM subscriptions WHERE subscriber_user_id = ? AND channel_id = ? LIMIT 1", subscriberId, channelId)
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

func (db *Database) DeleteSubscription(ctx TxContext, subid int64) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM subscriptions WHERE subscription_id = ?", subid)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateSubscriptionConfirmed(ctx TxContext, subscriptionID int64, confirmed bool) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE subscriptions SET confirmed = ? WHERE subscription_id = ?", confirmed, subscriptionID)
	if err != nil {
		return err
	}

	return nil
}
