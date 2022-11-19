package db

import (
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"time"
)

func (db *Database) GetMessageByUserMessageID(ctx TxContext, usrMsgId string) (*models.Message, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT * FROM messages WHERE usr_message_id = ? LIMIT 1", usrMsgId)
	if err != nil {
		return nil, err
	}

	msg, err := models.DecodeMessage(rows)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (db *Database) GetMessage(ctx TxContext, scnMessageID int64) (models.Message, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Message{}, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT * FROM messages WHERE scn_message_id = ? LIMIT 1", scnMessageID)
	if err != nil {
		return models.Message{}, err
	}

	msg, err := models.DecodeMessage(rows)
	if err != nil {
		return models.Message{}, err
	}

	return msg, nil
}

func (db *Database) CreateMessage(ctx TxContext, senderUserID int64, channel models.Channel, timestampSend *time.Time, title string, content *string, priority int, userMsgId *string) (models.Message, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Message{}, err
	}

	now := time.Now().UTC()

	res, err := tx.ExecContext(ctx, "INSERT INTO messages (sender_user_id, owner_user_id, channel_name, channel_id, timestamp_real, timestamp_client, title, content, priority, usr_message_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		senderUserID,
		channel.OwnerUserID,
		channel.Name,
		channel.ChannelID,
		time2DB(now),
		time2DBOpt(timestampSend),
		title,
		content,
		priority,
		userMsgId)
	if err != nil {
		return models.Message{}, err
	}

	liid, err := res.LastInsertId()
	if err != nil {
		return models.Message{}, err
	}

	return models.Message{
		SCNMessageID:    liid,
		SenderUserID:    senderUserID,
		OwnerUserID:     channel.OwnerUserID,
		ChannelName:     channel.Name,
		ChannelID:       channel.ChannelID,
		TimestampReal:   now,
		TimestampClient: timestampSend,
		Title:           title,
		Content:         content,
		Priority:        priority,
		UserMessageID:   userMsgId,
	}, nil
}

func (db *Database) DeleteMessage(ctx TxContext, scnMessageID int64) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM messages WHERE scn_message_id = ?", scnMessageID)
	if err != nil {
		return err
	}

	return nil
}
