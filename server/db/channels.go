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

	rows, err := tx.Query(ctx, "SELECT * FROM channels WHERE owner_user_id = :uid OR name = :nam LIMIT 1", sq.PP{
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

	rows, err := tx.Query(ctx, "SELECT * FROM channels WHERE name = :chan_name OR send_key = :send_key LIMIT 1", sq.PP{
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

func (db *Database) CreateChannel(ctx TxContext, userid models.UserID, name string, subscribeKey string, sendKey string) (models.Channel, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Channel{}, err
	}

	now := time.Now().UTC()

	res, err := tx.Exec(ctx, "INSERT INTO channels (owner_user_id, name, subscribe_key, send_key, timestamp_created) VALUES (:ouid, :nam, :subkey, :sendkey, :ts)", sq.PP{
		"ouid":    userid,
		"nam":     name,
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
		Name:              name,
		SubscribeKey:      subscribeKey,
		SendKey:           sendKey,
		TimestampCreated:  now,
		TimestampLastSent: nil,
		MessagesSent:      0,
	}, nil
}

func (db *Database) ListChannelsByOwner(ctx TxContext, userid models.UserID) ([]models.Channel, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, "SELECT * FROM channels WHERE owner_user_id = :ouid", sq.PP{"ouid": userid})
	if err != nil {
		return nil, err
	}

	data, err := models.DecodeChannels(rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (db *Database) ListChannelsBySubscriber(ctx TxContext, userid models.UserID, confirmed bool) ([]models.Channel, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	confCond := ""
	if confirmed {
		confCond = " AND sub.confirmed = 1"
	}

	rows, err := tx.Query(ctx, "SELECT * FROM channels LEFT JOIN subscriptions sub on channels.channel_id = sub.channel_id WHERE sub.subscriber_user_id = :suid "+confCond, sq.PP{
		"suid": userid,
	})
	if err != nil {
		return nil, err
	}

	data, err := models.DecodeChannels(rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (db *Database) ListChannelsByAccess(ctx TxContext, userid models.UserID, confirmed bool) ([]models.Channel, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	confCond := "OR sub.subscriber_user_id = ?"
	if confirmed {
		confCond = "OR (sub.subscriber_user_id = ? AND sub.confirmed = 1)"
	}

	rows, err := tx.Query(ctx, "SELECT * FROM channels LEFT JOIN subscriptions sub on channels.channel_id = sub.channel_id WHERE owner_user_id = :ouid "+confCond, sq.PP{
		"ouid": userid,
	})
	if err != nil {
		return nil, err
	}

	data, err := models.DecodeChannels(rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (db *Database) GetChannel(ctx TxContext, userid models.UserID, channelid models.ChannelID) (models.Channel, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Channel{}, err
	}

	rows, err := tx.Query(ctx, "SELECT * FROM channels WHERE owner_user_id = :ouid AND channel_id = :cid LIMIT 1", sq.PP{
		"ouid": userid,
		"cid":  channelid,
	})
	if err != nil {
		return models.Channel{}, err
	}

	client, err := models.DecodeChannel(rows)
	if err != nil {
		return models.Channel{}, err
	}

	return client, nil
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
