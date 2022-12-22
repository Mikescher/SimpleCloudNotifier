package db

import (
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"time"
)

func (db *Database) GetChannelByName(ctx TxContext, userid models.UserID, chanName string) (*models.Channel, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, "SELECT * FROM channels WHERE owner_user_id = :uid AND internal_name = :nam LIMIT 1", sq.PP{
		"uid": userid,
		"nam": chanName,
	})
	if err != nil {
		return nil, err
	}

	channel, err := models.DecodeChannel(rows)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &channel, nil
}

func (db *Database) GetChannelByNameAndSendKey(ctx TxContext, chanName string, sendKey string) (*models.Channel, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, "SELECT * FROM channels WHERE internal_name = :chan_name OR send_key = :send_key LIMIT 1", sq.PP{
		"chan_name": chanName,
		"send_key":  sendKey,
	})
	if err != nil {
		return nil, err
	}

	channel, err := models.DecodeChannel(rows)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &channel, nil
}

func (db *Database) CreateChannel(ctx TxContext, userid models.UserID, dispName string, intName string, subscribeKey string, sendKey string) (models.Channel, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Channel{}, err
	}

	now := time.Now().UTC()

	res, err := tx.Exec(ctx, "INSERT INTO channels (owner_user_id, display_name, internal_name, subscribe_key, send_key, timestamp_created) VALUES (:ouid, :dnam, :inam, :subkey, :sendkey, :ts)", sq.PP{
		"ouid":    userid,
		"dnam":    dispName,
		"inam":    intName,
		"subkey":  subscribeKey,
		"sendkey": sendKey,
		"ts":      time2DB(now),
	})
	if err != nil {
		return models.Channel{}, err
	}

	liid, err := res.LastInsertId()
	if err != nil {
		return models.Channel{}, err
	}

	return models.Channel{
		ChannelID:         models.ChannelID(liid),
		OwnerUserID:       userid,
		DisplayName:       dispName,
		InternalName:      intName,
		SubscribeKey:      subscribeKey,
		SendKey:           sendKey,
		TimestampCreated:  now,
		TimestampLastSent: nil,
		MessagesSent:      0,
	}, nil
}

func (db *Database) ListChannelsByOwner(ctx TxContext, userid models.UserID, subUserID models.UserID) ([]models.ChannelWithSubscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, "SELECT channels.*, sub.* FROM channels LEFT JOIN subscriptions AS sub ON channels.channel_id = sub.channel_id AND sub.subscriber_user_id = :subuid WHERE owner_user_id = :ouid", sq.PP{
		"ouid":   userid,
		"subuid": subUserID,
	})
	if err != nil {
		return nil, err
	}

	data, err := models.DecodeChannelsWithSubscription(rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (db *Database) ListChannelsBySubscriber(ctx TxContext, userid models.UserID, confirmed *bool) ([]models.ChannelWithSubscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	confCond := ""
	if confirmed != nil && *confirmed {
		confCond = " AND sub.confirmed = 1"
	} else if confirmed != nil && !*confirmed {
		confCond = " AND sub.confirmed = 0"
	}

	rows, err := tx.Query(ctx, "SELECT channels.*, sub.* FROM channels LEFT JOIN subscriptions AS sub on channels.channel_id = sub.channel_id AND sub.subscriber_user_id = :subuid WHERE sub.subscription_id IS NOT NULL "+confCond, sq.PP{
		"subuid": userid,
	})
	if err != nil {
		return nil, err
	}

	data, err := models.DecodeChannelsWithSubscription(rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (db *Database) ListChannelsByAccess(ctx TxContext, userid models.UserID, confirmed *bool) ([]models.ChannelWithSubscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	confCond := ""
	if confirmed != nil && *confirmed {
		confCond = "OR sub.confirmed = 1"
	} else if confirmed != nil && !*confirmed {
		confCond = "OR sub.confirmed = 0"
	}

	rows, err := tx.Query(ctx, "SELECT channels.*, sub.* FROM channels LEFT JOIN subscriptions AS sub on channels.channel_id = sub.channel_id AND sub.subscriber_user_id = :subuid WHERE owner_user_id = :ouid "+confCond, sq.PP{
		"ouid":   userid,
		"subuid": userid,
	})
	if err != nil {
		return nil, err
	}

	data, err := models.DecodeChannelsWithSubscription(rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (db *Database) GetChannel(ctx TxContext, userid models.UserID, channelid models.ChannelID) (models.ChannelWithSubscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.ChannelWithSubscription{}, err
	}

	rows, err := tx.Query(ctx, "SELECT channels.*, sub.* FROM channels LEFT JOIN subscriptions AS sub on channels.channel_id = sub.channel_id AND sub.subscriber_user_id = :subuid WHERE owner_user_id = :ouid AND channels.channel_id = :cid LIMIT 1", sq.PP{
		"ouid":   userid,
		"cid":    channelid,
		"subuid": userid,
	})
	if err != nil {
		return models.ChannelWithSubscription{}, err
	}

	channel, err := models.DecodeChannelWithSubscription(rows)
	if err != nil {
		return models.ChannelWithSubscription{}, err
	}

	return channel, nil
}

func (db *Database) IncChannelMessageCounter(ctx TxContext, channel models.Channel) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE channels SET messages_sent = :ctr, timestamp_lastsent = :ts WHERE channel_id = :cid", sq.PP{
		"ctr": channel.MessagesSent + 1,
		"cid": time2DB(time.Now()),
		"ts":  channel.ChannelID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateChannelSendKey(ctx TxContext, channelid models.ChannelID, newkey string) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE channels SET send_key = :key WHERE channel_id = :cid", sq.PP{
		"key": newkey,
		"cid": channelid,
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateChannelSubscribeKey(ctx TxContext, channelid models.ChannelID, newkey string) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE channels SET subscribe_key = :key WHERE channel_id = :cid", sq.PP{
		"key": newkey,
		"cid": channelid,
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateChannelDisplayName(ctx TxContext, channelid models.ChannelID, dispname string) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE channels SET display_name = :nam WHERE channel_id = :cid", sq.PP{
		"nam": dispname,
		"cid": channelid,
	})
	if err != nil {
		return err
	}

	return nil
}
