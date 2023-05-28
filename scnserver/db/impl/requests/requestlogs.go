package requests

import (
	ct "blackforestbytes.com/simplecloudnotifier/db/cursortoken"
	"blackforestbytes.com/simplecloudnotifier/models"
	"context"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"time"
)

func (db *Database) InsertRequestLog(ctx context.Context, requestid models.RequestID, data models.RequestLog) (models.RequestLog, error) {

	entity := data.DB()
	entity.RequestID = requestid
	entity.TimestampCreated = time2DB(time.Now())

	_, err := sq.InsertSingle(ctx, db.db, "requests", entity)
	if err != nil {
		return models.RequestLog{}, err
	}

	return entity.Model(), nil
}

func (db *Database) Cleanup(ctx context.Context, count int, duration time.Duration) (int64, error) {
	res1, err := db.db.Exec(ctx, "DELETE FROM requests WHERE request_id NOT IN ( SELECT request_id FROM requests ORDER BY timestamp_created DESC LIMIT :keep ) ", sq.PP{
		"keep": count,
	})
	if err != nil {
		return 0, err
	}
	affected1, err := res1.RowsAffected()
	if err != nil {
		return 0, err
	}

	res2, err := db.db.Exec(ctx, "DELETE FROM requests WHERE timestamp_created < :tslim", sq.PP{
		"tslim": time.Now().Add(-duration).UnixMilli(),
	})
	if err != nil {
		return 0, err
	}
	affected2, err := res2.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected1 + affected2, nil
}

func (db *Database) ListRequestLogs(ctx context.Context, filter models.RequestLogFilter, pageSize *int, inTok ct.CursorToken) ([]models.RequestLog, ct.CursorToken, error) {
	if inTok.Mode == ct.CTMEnd {
		return make([]models.RequestLog, 0), ct.End(), nil
	}

	pageCond := "1=1"
	if inTok.Mode == ct.CTMNormal {
		pageCond = "timestamp_created < :tokts OR (timestamp_created = :tokts AND request_id < :tokid )"
	}

	filterCond, filterJoin, prepParams, err := filter.SQL()

	orderClause := ""
	if pageSize != nil {
		orderClause = "ORDER BY timestamp_created DESC, request_id DESC LIMIT :lim"
		prepParams["lim"] = *pageSize + 1
	} else {
		orderClause = "ORDER BY timestamp_created DESC, request_id DESC"
	}

	sqlQuery := "SELECT " + "requests.*" + " FROM requests " + filterJoin + " WHERE ( " + pageCond + " ) AND ( " + filterCond + " ) " + orderClause

	prepParams["tokts"] = inTok.Timestamp
	prepParams["tokid"] = inTok.Id

	rows, err := db.db.Query(ctx, sqlQuery, prepParams)
	if err != nil {
		return nil, ct.CursorToken{}, err
	}

	data, err := models.DecodeRequestLogs(rows)
	if err != nil {
		return nil, ct.CursorToken{}, err
	}

	if pageSize == nil || len(data) <= *pageSize {
		return data, ct.End(), nil
	} else {
		outToken := ct.Normal(data[*pageSize-1].TimestampCreated, data[*pageSize-1].RequestID.String(), "DESC", filter.Hash())
		return data[0:*pageSize], outToken, nil
	}
}
