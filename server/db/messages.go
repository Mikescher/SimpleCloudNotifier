package db

import (
	"blackforestbytes.com/simplecloudnotifier/db/cursortoken"
	"blackforestbytes.com/simplecloudnotifier/models"
	"blackforestbytes.com/simplecloudnotifier/sq"
	"database/sql"
	"fmt"
	"time"
)

func (db *Database) GetMessageByUserMessageID(ctx TxContext, usrMsgId string) (*models.Message, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, "SELECT * FROM messages WHERE usr_message_id = :umid LIMIT 1", sq.PP{"umid": usrMsgId})
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

func (db *Database) GetMessage(ctx TxContext, scnMessageID models.SCNMessageID) (models.Message, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Message{}, err
	}

	rows, err := tx.Query(ctx, "SELECT * FROM messages WHERE scn_message_id = :mid LIMIT 1", sq.PP{"mid": scnMessageID})
	if err != nil {
		return models.Message{}, err
	}

	msg, err := models.DecodeMessage(rows)
	if err != nil {
		return models.Message{}, err
	}

	return msg, nil
}

func (db *Database) CreateMessage(ctx TxContext, senderUserID models.UserID, channel models.Channel, timestampSend *time.Time, title string, content *string, priority int, userMsgId *string, senderIP string, senderName *string) (models.Message, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Message{}, err
	}

	now := time.Now().UTC()

	res, err := tx.Exec(ctx, "INSERT INTO messages (sender_user_id, owner_user_id, channel_name, channel_id, timestamp_real, timestamp_client, title, content, priority, usr_message_id, sender_ip, sender_name) VALUES (:suid, :ouid, :cnam, :cid, :tsr, :tsc, :tit, :cnt, :prio, :umid, :ip, :snam)", sq.PP{
		"suid": senderUserID,
		"ouid": channel.OwnerUserID,
		"cnam": channel.Name,
		"cid":  channel.ChannelID,
		"tsr":  time2DB(now),
		"tsc":  time2DBOpt(timestampSend),
		"tit":  title,
		"cnt":  content,
		"prio": priority,
		"umid": userMsgId,
		"ip":   senderIP,
		"snam": senderName,
	})
	if err != nil {
		return models.Message{}, err
	}

	liid, err := res.LastInsertId()
	if err != nil {
		return models.Message{}, err
	}

	return models.Message{
		SCNMessageID:    models.SCNMessageID(liid),
		SenderUserID:    senderUserID,
		OwnerUserID:     channel.OwnerUserID,
		ChannelName:     channel.Name,
		ChannelID:       channel.ChannelID,
		SenderIP:        senderIP,
		SenderName:      senderName,
		TimestampReal:   now,
		TimestampClient: timestampSend,
		Title:           title,
		Content:         content,
		Priority:        priority,
		UserMessageID:   userMsgId,
	}, nil
}

func (db *Database) DeleteMessage(ctx TxContext, scnMessageID models.SCNMessageID) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "DELETE FROM messages WHERE scn_message_id = :mid", sq.PP{"mid": scnMessageID})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) ListMessages(ctx TxContext, userid models.UserID, pageSize int, inTok cursortoken.CursorToken) ([]models.Message, cursortoken.CursorToken, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, cursortoken.CursorToken{}, err
	}

	if inTok.Mode == cursortoken.CTMEnd {
		return make([]models.Message, 0), cursortoken.End(), nil
	}

	pageCond := ""
	if inTok.Mode == cursortoken.CTMNormal {
		pageCond = fmt.Sprintf("AND ( timestamp_real < %d OR (timestamp_real = %d AND scn_message_id < %d ) )", inTok.Timestamp, inTok.Timestamp, inTok.Id)
	}

	rows, err := tx.Query(ctx, "SELECT messages.* FROM messages LEFT JOIN subscriptions subs on messages.channel_id = subs.channel_id WHERE subs.subscriber_user_id = :uid AND subs.confirmed = 1 "+pageCond+" ORDER BY timestamp_real DESC LIMIT :lim", sq.PP{
		"uid": userid,
		"lim": pageSize + 1,
	})
	if err != nil {
		return nil, cursortoken.CursorToken{}, err
	}

	data, err := models.DecodeMessages(rows)
	if err != nil {
		return nil, cursortoken.CursorToken{}, err
	}

	if len(data) <= pageSize {
		return data, cursortoken.End(), nil
	} else {
		outToken := cursortoken.Normal(data[pageSize-1].TimestampReal, data[pageSize-1].SCNMessageID.IntID(), "DESC")
		return data[0:pageSize], outToken, nil
	}
}

func (db *Database) ListChannelMessages(ctx TxContext, channelid models.ChannelID, pageSize int, inTok cursortoken.CursorToken) ([]models.Message, cursortoken.CursorToken, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, cursortoken.CursorToken{}, err
	}

	if inTok.Mode == cursortoken.CTMEnd {
		return make([]models.Message, 0), cursortoken.End(), nil
	}

	pageCond := ""
	if inTok.Mode == cursortoken.CTMNormal {
		pageCond = "AND ( timestamp_real < :tokts OR (timestamp_real = :tokts AND scn_message_id < :tokid ) )"
	}

	rows, err := tx.Query(ctx, "SELECT * FROM messages WHERE channel_id = :cid "+pageCond+" ORDER BY timestamp_real DESC LIMIT :lim", sq.PP{
		"cid":   channelid,
		"lim":   pageSize + 1,
		"tokts": inTok.Timestamp,
		"tokid": inTok.Timestamp,
	})
	if err != nil {
		return nil, cursortoken.CursorToken{}, err
	}

	data, err := models.DecodeMessages(rows)
	if err != nil {
		return nil, cursortoken.CursorToken{}, err
	}

	if len(data) <= pageSize {
		return data, cursortoken.End(), nil
	} else {
		outToken := cursortoken.Normal(data[pageSize-1].TimestampReal, data[pageSize-1].SCNMessageID.IntID(), "DESC")
		return data[0:pageSize], outToken, nil
	}
}
