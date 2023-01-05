package primary

import (
	"blackforestbytes.com/simplecloudnotifier/db/cursortoken"
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

func (db *Database) GetMessage(ctx TxContext, scnMessageID models.SCNMessageID, allowDeleted bool) (models.Message, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Message{}, err
	}

	var sqlcmd string
	if allowDeleted {
		sqlcmd = "SELECT * FROM messages WHERE scn_message_id = :mid LIMIT 1"
	} else {
		sqlcmd = "SELECT * FROM messages WHERE scn_message_id = :mid AND deleted=0 LIMIT 1"
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

func (db *Database) CreateMessage(ctx TxContext, senderUserID models.UserID, channel models.Channel, timestampSend *time.Time, title string, content *string, priority int, userMsgId *string, senderIP string, senderName *string) (models.Message, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Message{}, err
	}

	now := time.Now().UTC()

	res, err := tx.Exec(ctx, "INSERT INTO messages (sender_user_id, owner_user_id, channel_internal_name, channel_id, timestamp_real, timestamp_client, title, content, priority, usr_message_id, sender_ip, sender_name) VALUES (:suid, :ouid, :cnam, :cid, :tsr, :tsc, :tit, :cnt, :prio, :umid, :ip, :snam)", sq.PP{
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
	})
	if err != nil {
		return models.Message{}, err
	}

	liid, err := res.LastInsertId()
	if err != nil {
		return models.Message{}, err
	}

	return models.Message{
		SCNMessageID:        models.SCNMessageID(liid),
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
	}, nil
}

func (db *Database) DeleteMessage(ctx TxContext, scnMessageID models.SCNMessageID) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE messages SET deleted=1 WHERE scn_message_id = :mid AND deleted=0", sq.PP{"mid": scnMessageID})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) ListMessages(ctx TxContext, filter models.MessageFilter, pageSize int, inTok cursortoken.CursorToken) ([]models.Message, cursortoken.CursorToken, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, cursortoken.CursorToken{}, err
	}

	if inTok.Mode == cursortoken.CTMEnd {
		return make([]models.Message, 0), cursortoken.End(), nil
	}

	pageCond := "1=1"
	if inTok.Mode == cursortoken.CTMNormal {
		pageCond = "timestamp_real < :tokts OR (timestamp_real = :tokts AND scn_message_id < :tokid )"
	}

	filterCond, filterJoin, prepParams, err := filter.SQL()

	orderClause := "ORDER BY COALESCE(timestamp_client, timestamp_real) DESC LIMIT :lim"

	sqlQuery := "SELECT " + "messages.*" + " FROM messages " + filterJoin + " WHERE ( " + pageCond + " ) AND ( " + filterCond + " ) " + orderClause

	prepParams["lim"] = pageSize + 1
	prepParams["tokts"] = inTok.Timestamp
	prepParams["tokid"] = inTok.Id

	rows, err := tx.Query(ctx, sqlQuery, prepParams)
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
		outToken := cursortoken.Normal(data[pageSize-1].Timestamp(), data[pageSize-1].SCNMessageID.IntID(), "DESC", filter.Hash())
		return data[0:pageSize], outToken, nil
	}
}
