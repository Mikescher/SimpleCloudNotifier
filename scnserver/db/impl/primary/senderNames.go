package primary

import (
	"blackforestbytes.com/simplecloudnotifier/db"
	"blackforestbytes.com/simplecloudnotifier/models"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
)

func (db *Database) ListSenderNames(ctx db.TxContext, userid models.UserID, includeForeignSubscribed bool) ([]models.SenderNameStatistics, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	var sqlStr string

	prepParams := sq.PP{"uid": userid}

	if includeForeignSubscribed {
		sqlStr = "SELECT sender_name AS name, MAX(timestamp_real) AS ts_last, MIN(timestamp_real) AS ts_first, COUNT(*) AS count FROM messages LEFT JOIN subscriptions AS subs on messages.channel_id = subs.channel_id WHERE (subs.subscriber_user_id = :uid AND subs.confirmed = 1) AND sender_NAME NOT NULL GROUP BY sender_name ORDER BY ts_last DESC"
	} else {
		sqlStr = "SELECT sender_name AS name, MAX(timestamp_real) AS ts_last, MIN(timestamp_real) AS ts_first, COUNT(*) AS count FROM messages WHERE sender_user_id = :uid AND sender_NAME NOT NULL GROUP BY sender_name ORDER BY ts_last DESC"
	}

	return sq.QueryAll[models.SenderNameStatistics](ctx, tx, sqlStr, prepParams, sq.SModeExtended, sq.Safe)
}
