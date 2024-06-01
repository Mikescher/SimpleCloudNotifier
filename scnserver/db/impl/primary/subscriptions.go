package primary

import (
	"blackforestbytes.com/simplecloudnotifier/db"
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"errors"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"time"
)

func (db *Database) CreateSubscription(ctx db.TxContext, subscriberUID models.UserID, channel models.Channel, confirmed bool) (models.Subscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Subscription{}, err
	}

	entity := models.SubscriptionDB{
		SubscriptionID:      models.NewSubscriptionID(),
		SubscriberUserID:    subscriberUID,
		ChannelOwnerUserID:  channel.OwnerUserID,
		ChannelID:           channel.ChannelID,
		ChannelInternalName: channel.InternalName,
		TimestampCreated:    time2DB(time.Now()),
		Confirmed:           bool2DB(confirmed),
	}

	_, err = sq.InsertSingle(ctx, tx, "subscriptions", entity)
	if err != nil {
		return models.Subscription{}, err
	}

	return entity.Model(), nil
}

func (db *Database) ListSubscriptions(ctx db.TxContext, filter models.SubscriptionFilter) ([]models.Subscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	filterCond, filterJoin, prepParams, err := filter.SQL()

	orderClause := " ORDER BY subscriptions.timestamp_created DESC, subscriptions.subscription_id DESC "

	sqlQuery := "SELECT " + "subscriptions.*" + " FROM subscriptions " + filterJoin + " WHERE ( " + filterCond + " ) " + orderClause

	rows, err := tx.Query(ctx, sqlQuery, prepParams)
	if err != nil {
		return nil, err
	}

	data, err := models.DecodeSubscriptions(ctx, tx, rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (db *Database) GetSubscription(ctx db.TxContext, subid models.SubscriptionID) (models.Subscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Subscription{}, err
	}

	rows, err := tx.Query(ctx, "SELECT * FROM subscriptions WHERE subscription_id = :sid LIMIT 1", sq.PP{"sid": subid})
	if err != nil {
		return models.Subscription{}, err
	}

	sub, err := models.DecodeSubscription(ctx, tx, rows)
	if err != nil {
		return models.Subscription{}, err
	}

	return sub, nil
}

func (db *Database) GetSubscriptionBySubscriber(ctx db.TxContext, subscriberId models.UserID, channelId models.ChannelID) (*models.Subscription, error) {
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

	user, err := models.DecodeSubscription(ctx, tx, rows)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *Database) DeleteSubscription(ctx db.TxContext, subid models.SubscriptionID) error {
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

func (db *Database) UpdateSubscriptionConfirmed(ctx db.TxContext, subscriptionID models.SubscriptionID, confirmed bool) error {
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
