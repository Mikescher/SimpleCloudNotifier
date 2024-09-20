package primary

import (
	"blackforestbytes.com/simplecloudnotifier/db"
	ct "blackforestbytes.com/simplecloudnotifier/db/cursortoken"
	"blackforestbytes.com/simplecloudnotifier/models"
	"errors"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"time"
)

func (db *Database) GetMessageByUserMessageID(ctx db.TxContext, usrMsgId string) (*models.Message, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	return sq.QuerySingleOpt[models.Message](ctx, tx, "SELECT * FROM messages WHERE usr_message_id = :umid LIMIT 1", sq.PP{"umid": usrMsgId}, sq.SModeExtended, sq.Safe)
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

	return sq.QuerySingle[models.Message](ctx, tx, sqlcmd, sq.PP{"mid": scnMessageID}, sq.SModeExtended, sq.Safe)
}

func (db *Database) CreateMessage(ctx db.TxContext, senderUserID models.UserID, channel models.Channel, timestampSend *time.Time, title string, content *string, priority int, userMsgId *string, senderIP string, senderName *string, usedKeyID models.KeyTokenID) (models.Message, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Message{}, err
	}

	entity := models.Message{
		MessageID:           models.NewMessageID(),
		SenderUserID:        senderUserID,
		ChannelInternalName: channel.InternalName,
		ChannelID:           channel.ChannelID,
		SenderIP:            senderIP,
		SenderName:          senderName,
		TimestampReal:       models.NowSCNTime(),
		TimestampClient:     models.NewSCNTimePtr(timestampSend),
		Title:               title,
		Content:             content,
		Priority:            priority,
		UserMessageID:       userMsgId,
		UsedKeyID:           usedKeyID,
		Deleted:             false,
		MessageExtra:        models.MessageExtra{},
	}

	_, err = sq.InsertSingle(ctx, tx, "messages", entity)
	if err != nil {
		return models.Message{}, err
	}

	return entity, nil
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

func (db *Database) ListMessages(ctx db.TxContext, filter models.MessageFilter, pageSize *int, inTok ct.CursorToken) ([]models.Message, ct.CursorToken, int64, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, ct.CursorToken{}, 0, err
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

	sqlQueryList := "SELECT " + "messages.*" + " FROM messages " + filterJoin + " WHERE ( " + pageCond + " ) AND ( " + filterCond + " ) " + orderClause
	sqlQueryCount := "SELECT " + " COUNT(*) AS count FROM messages " + filterJoin + " WHERE  ( " + filterCond + " ) "

	prepParams["tokts"] = inTok.Timestamp
	prepParams["tokid"] = inTok.Id

	if inTok.Mode == ct.CTMEnd {

		dataCount, err := sq.QuerySingle[CountResponse](ctx, tx, sqlQueryCount, prepParams, sq.SModeFast, sq.Safe)
		if err != nil {
			return nil, ct.CursorToken{}, 0, err
		}

		return make([]models.Message, 0), ct.End(), dataCount.Count, nil
	}

	dataList, err := sq.QueryAll[models.Message](ctx, tx, sqlQueryList, prepParams, sq.SModeExtended, sq.Safe)
	if err != nil {
		return nil, ct.CursorToken{}, 0, err
	}

	if pageSize == nil || len(dataList) <= *pageSize {
		return dataList, ct.End(), int64(len(dataList)), nil
	} else {

		dataCount, err := sq.QuerySingle[CountResponse](ctx, tx, sqlQueryCount, prepParams, sq.SModeFast, sq.Safe)
		if err != nil {
			return nil, ct.CursorToken{}, 0, err
		}

		outToken := ct.Normal(dataList[*pageSize-1].Timestamp(), dataList[*pageSize-1].MessageID.String(), "DESC", filter.Hash())

		return dataList[0:*pageSize], outToken, dataCount.Count, nil
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
