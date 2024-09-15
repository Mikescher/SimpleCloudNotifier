package primary

import (
	"blackforestbytes.com/simplecloudnotifier/db"
	"blackforestbytes.com/simplecloudnotifier/models"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
)

func (db *Database) CreateSubscription(ctx db.TxContext, subscriberUID models.UserID, channel models.Channel, confirmed bool) (models.Subscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Subscription{}, err
	}

	entity := models.Subscription{
		SubscriptionID:      models.NewSubscriptionID(),
		SubscriberUserID:    subscriberUID,
		ChannelOwnerUserID:  channel.OwnerUserID,
		ChannelID:           channel.ChannelID,
		ChannelInternalName: channel.InternalName,
		TimestampCreated:    models.NowSCNTime(),
		Confirmed:           confirmed,
	}

	_, err = sq.InsertSingle(ctx, tx, "subscriptions", entity)
	if err != nil {
		return models.Subscription{}, err
	}

	return entity, nil
}

func (db *Database) ListSubscriptions(ctx db.TxContext, filter models.SubscriptionFilter) ([]models.Subscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	filterCond, filterJoin, prepParams, err := filter.SQL()

	orderClause := " ORDER BY subscriptions.timestamp_created DESC, subscriptions.subscription_id DESC "

	sqlQuery := "SELECT " + "subscriptions.*" + " FROM subscriptions " + filterJoin + " WHERE ( " + filterCond + " ) " + orderClause

	return sq.QueryAll[models.Subscription](ctx, tx, sqlQuery, prepParams, sq.SModeExtended, sq.Safe)
}

func (db *Database) GetSubscription(ctx db.TxContext, subid models.SubscriptionID) (models.Subscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Subscription{}, err
	}

	return sq.QuerySingle[models.Subscription](ctx, tx, "SELECT * FROM subscriptions WHERE subscription_id = :sid LIMIT 1", sq.PP{"sid": subid}, sq.SModeExtended, sq.Safe)
}

func (db *Database) GetSubscriptionBySubscriber(ctx db.TxContext, subscriberId models.UserID, channelId models.ChannelID) (*models.Subscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	return sq.QuerySingleOpt[models.Subscription](ctx, tx, "SELECT * FROM subscriptions WHERE subscriber_user_id = :suid AND channel_id = :cid LIMIT 1", sq.PP{
		"suid": subscriberId,
		"cid":  channelId,
	}, sq.SModeExtended, sq.Safe)
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
