package db

import (
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"time"
)

func (db *Database) GetChannelByKey(ctx TxContext, key string) (*models.Channel, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT * FROM channels WHERE subscribe_key = ? OR send_key = ? LIMIT 1", key, key)
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

func (db *Database) GetChannelByName(ctx TxContext, userid int64, chanName string) (*models.Channel, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT * FROM channels WHERE owner_user_id = ? OR name = ? LIMIT 1", userid, chanName)
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

func (db *Database) CreateChannel(ctx TxContext, userid int64, name string, subscribeKey string, sendKey string) (models.Channel, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Channel{}, err
	}

	now := time.Now().UTC()

	res, err := tx.ExecContext(ctx, "INSERT INTO channels (owner_user_id, name, subscribe_key, send_key, timestamp_created) VALUES (?, ?, ?, ?, ?)",
		userid,
		name,
		subscribeKey,
		sendKey,
		time2DB(now))
	if err != nil {
		return models.Channel{}, err
	}

	liid, err := res.LastInsertId()
	if err != nil {
		return models.Channel{}, err
	}

	return models.Channel{
		ChannelID:         liid,
		OwnerUserID:       userid,
		Name:              name,
		SubscribeKey:      subscribeKey,
		SendKey:           sendKey,
		TimestampCreated:  now,
		TimestampLastSent: nil,
		MessagesSent:      0,
	}, nil
}

func (db *Database) ListChannels(ctx TxContext, userid int64) ([]models.Channel, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT * FROM channels WHERE owner_user_id = ?", userid)
	if err != nil {
		return nil, err
	}

	data, err := models.DecodeChannels(rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (db *Database) GetChannel(ctx TxContext, userid int64, channelid int64) (models.Channel, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Channel{}, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT * FROM channels WHERE owner_user_id = ? AND channel_id = ? LIMIT 1", userid, channelid)
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

	_, err = tx.ExecContext(ctx, "UPDATE channels SET messages_sent = ? AND timestamp_lastsent = ? WHERE channel_id = ?",
		channel.MessagesSent+1,
		time2DB(time.Now()),
		channel.ChannelID)
	if err != nil {
		return err
	}

	return nil
}