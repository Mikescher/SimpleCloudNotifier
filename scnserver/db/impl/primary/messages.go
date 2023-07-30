package primary

import (
	"blackforestbytes.com/simplecloudnotifier/db"
	ct "blackforestbytes.com/simplecloudnotifier/db/cursortoken"
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"errors"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"time"
)

func (db *Database) GetMessageByUserMessageID(ctx db.TxContext, usrMsgId string) (*models.Message, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, "SELECT * FROM messages WHERE usr_message_id = :umid LIMIT 1", sq.PP{"umid": usrMsgId})
	if err != nil {
		return nil, err
	}

	msg, err := models.DecodeMessage(rows)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (db *Database) GetMessage(ctx db.TxContext, scnMessageID models.MessageID, allowDeleted bool) (models.Message, error) {
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

func (db *Database) CreateMessage(ctx db.TxContext, senderUserID models.UserID, channel models.Channel, timestampSend *time.Time, title string, content *string, priority int, userMsgId *string, senderIP string, senderName *string, usedKeyID models.KeyTokenID) (models.Message, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Message{}, err
	}

	entity := models.MessageDB{
		MessageID:           models.NewMessageID(),
		SenderUserID:        senderUserID,
		ChannelInternalName: channel.InternalName,
		ChannelID:           channel.ChannelID,
		SenderIP:            senderIP,
		SenderName:          senderName,
		TimestampReal:       time2DB(time.Now()),
		TimestampClient:     time2DBOpt(timestampSend),
		Title:               title,
		Content:             content,
		Priority:            priority,
		UserMessageID:       userMsgId,
		UsedKeyID:           usedKeyID,
		Deleted:             bool2DB(false),
	}

	_, err = sq.InsertSingle(ctx, tx, "messages", entity)
	if err != nil {
		return models.Message{}, err
	}

	return entity.Model(), nil
}

func (db *Database) DeleteMessage(ctx db.TxContext, messageID models.MessageID) error {
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

func (db *Database) ListMessages(ctx db.TxContext, filter models.MessageFilter, pageSize *int, inTok ct.CursorToken) ([]models.Message, ct.CursorToken, error) {
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

func (db *Database) CountMessages(ctx db.TxContext, filter models.MessageFilter) (int64, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return 0, err
	}

	filterCond, filterJoin, prepParams, err := filter.SQL()

	sqlQuery := "SELECT " + "COUNT(*)" + " FROM messages " + filterJoin + " WHERE  ( " + filterCond + " ) "

	rows, err := tx.Query(ctx, sqlQuery, prepParams)
	if err != nil {
		return 0, err
	}

	if !rows.Next() {
		return 0, errors.New("COUNT query returned no results")
	}

	var countRes int64
	err = rows.Scan(&countRes)
	if err != nil {
		return 0, err
	}

	return countRes, nil
}
