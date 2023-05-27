package primary

import (
	ct "blackforestbytes.com/simplecloudnotifier/db/cursortoken"
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
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

func (db *Database) GetMessage(ctx TxContext, scnMessageID models.MessageID, allowDeleted bool) (models.Message, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Message{}, err
	}

	var sqlcmd string
	if allowDeleted {
		sqlcmd = "SELECT * FROM messages WHERE message_id = :mid LIMIT 1"
	} else {
		sqlcmd = "SELECT * FROM messages WHERE message_id = :mid AND deleted=0 LIMIT 1"
	}

	rows, err := tx.Query(ctx, sqlcmd, sq.PP{"mid": scnMessageID})
	if err != nil {
		return models.Message{}, err
	}

	msg, err := models.DecodeMessage(rows)
	if err != nil {
		return models.Message{}, err
	}

	return msg, nil
}

func (db *Database) CreateMessage(ctx TxContext, senderUserID models.UserID, channel models.Channel, timestampSend *time.Time, title string, content *string, priority int, userMsgId *string, senderIP string, senderName *string, usedKeyID models.KeyTokenID) (models.Message, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Message{}, err
	}

	now := time.Now().UTC()

	messageid := models.NewMessageID()

	_, err = tx.Exec(ctx, "INSERT INTO messages (message_id, sender_user_id, owner_user_id, channel_internal_name, channel_id, timestamp_real, timestamp_client, title, content, priority, usr_message_id, sender_ip, sender_name, used_key_id) VALUES (:mid, :suid, :ouid, :cnam, :cid, :tsr, :tsc, :tit, :cnt, :prio, :umid, :ip, :snam, :uk)", sq.PP{
		"mid":  messageid,
		"suid": senderUserID,
		"ouid": channel.OwnerUserID,
		"cnam": channel.InternalName,
		"cid":  channel.ChannelID,
		"tsr":  time2DB(now),
		"tsc":  time2DBOpt(timestampSend),
		"tit":  title,
		"cnt":  content,
		"prio": priority,
		"umid": userMsgId,
		"ip":   senderIP,
		"snam": senderName,
		"uk":   usedKeyID,
	})
	if err != nil {
		return models.Message{}, err
	}

	return models.Message{
		MessageID:           messageid,
		SenderUserID:        senderUserID,
		OwnerUserID:         channel.OwnerUserID,
		ChannelInternalName: channel.InternalName,
		ChannelID:           channel.ChannelID,
		SenderIP:            senderIP,
		SenderName:          senderName,
		TimestampReal:       now,
		TimestampClient:     timestampSend,
		Title:               title,
		Content:             content,
		Priority:            priority,
		UserMessageID:       userMsgId,
		UsedKeyID:           usedKeyID,
	}, nil
}

func (db *Database) DeleteMessage(ctx TxContext, messageID models.MessageID) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE messages SET deleted=1 WHERE message_id = :mid AND deleted=0", sq.PP{"mid": messageID})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) ListMessages(ctx TxContext, filter models.MessageFilter, pageSize *int, inTok ct.CursorToken) ([]models.Message, ct.CursorToken, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, ct.CursorToken{}, err
	}

	if inTok.Mode == ct.CTMEnd {
		return make([]models.Message, 0), ct.End(), nil
	}

	pageCond := "1=1"
	if inTok.Mode == ct.CTMNormal {
		pageCond = "timestamp_real < :tokts OR (timestamp_real = :tokts AND message_id < :tokid )"
	}

	filterCond, filterJoin, prepParams, err := filter.SQL()

	orderClause := ""
	if pageSize != nil {
		orderClause = "ORDER BY COALESCE(timestamp_client, timestamp_real) DESC, message_id DESC LIMIT :lim"
		prepParams["lim"] = *pageSize + 1
	} else {
		orderClause = "ORDER BY COALESCE(timestamp_client, timestamp_real) DESC, message_id DESC"
	}

	sqlQuery := "SELECT " + "messages.*" + " FROM messages " + filterJoin + " WHERE ( " + pageCond + " ) AND ( " + filterCond + " ) " + orderClause

	prepParams["tokts"] = inTok.Timestamp
	prepParams["tokid"] = inTok.Id

	rows, err := tx.Query(ctx, sqlQuery, prepParams)
	if err != nil {
		return nil, ct.CursorToken{}, err
	}

	data, err := models.DecodeMessages(rows)
	if err != nil {
		return nil, ct.CursorToken{}, err
	}

	if pageSize == nil || len(data) <= *pageSize {
		return data, ct.End(), nil
	} else {
		outToken := ct.Normal(data[*pageSize-1].Timestamp(), data[*pageSize-1].MessageID.String(), "DESC", filter.Hash())
		return data[0:*pageSize], outToken, nil
	}
}
