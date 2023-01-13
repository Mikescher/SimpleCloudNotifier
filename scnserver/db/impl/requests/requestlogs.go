package requests

import (
	"blackforestbytes.com/simplecloudnotifier/models"
	"context"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"time"
)

func (db *Database) InsertRequestLog(ctx context.Context, data models.RequestLogDB) (models.RequestLogDB, error) {

	now := time.Now()

	res, err := db.db.Exec(ctx, "INSERT INTO requests (method, uri, user_agent, authentication, request_body, request_body_size, request_content_type, remote_ip, userid, permissions, response_statuscode, response_body_size, response_body, response_content_type, retry_count, panicked, panic_str, processing_time, timestamp_created, timestamp_start, timestamp_finish) VALUES (:method, :uri, :user_agent, :authentication, :request_body, :request_body_size, :request_content_type, :remote_ip, :userid, :permissions, :response_statuscode, :response_body_size, :response_body, :response_content_type, :retry_count, :panicked, :panic_str, :processing_time, :timestamp_created, :timestamp_start, :timestamp_finish)", sq.PP{
		"method":                data.Method,
		"uri":                   data.URI,
		"user_agent":            data.UserAgent,
		"authentication":        data.Authentication,
		"request_body":          data.RequestBody,
		"request_body_size":     data.RequestBodySize,
		"request_content_type":  data.RequestContentType,
		"remote_ip":             data.RemoteIP,
		"userid":                data.UserID,
		"permissions":           data.Permissions,
		"response_statuscode":   data.ResponseStatuscode,
		"response_body_size":    data.ResponseBodySize,
		"response_body":         data.ResponseBody,
		"response_content_type": data.ResponseContentType,
		"retry_count":           data.RetryCount,
		"panicked":              data.Panicked,
		"panic_str":             data.PanicStr,
		"processing_time":       data.ProcessingTime,
		"timestamp_created":     now.UnixMilli(),
		"timestamp_start":       data.TimestampStart,
		"timestamp_finish":      data.TimestampFinish,
	})
	if err != nil {
		return models.RequestLogDB{}, err
	}

	liid, err := res.LastInsertId()
	if err != nil {
		return models.RequestLogDB{}, err
	}

	return models.RequestLogDB{
		RequestID:           models.RequestID(liid),
		Method:              data.Method,
		URI:                 data.URI,
		UserAgent:           data.UserAgent,
		Authentication:      data.Authentication,
		RequestBody:         data.RequestBody,
		RequestBodySize:     data.RequestBodySize,
		RequestContentType:  data.RequestContentType,
		RemoteIP:            data.RemoteIP,
		UserID:              data.UserID,
		Permissions:         data.Permissions,
		ResponseStatuscode:  data.ResponseStatuscode,
		ResponseBodySize:    data.ResponseBodySize,
		ResponseBody:        data.ResponseBody,
		ResponseContentType: data.ResponseContentType,
		RetryCount:          data.RetryCount,
		Panicked:            data.Panicked,
		PanicStr:            data.PanicStr,
		ProcessingTime:      data.ProcessingTime,
		TimestampCreated:    now.UnixMilli(),
		TimestampStart:      data.TimestampStart,
		TimestampFinish:     data.TimestampFinish,
	}, nil
}

func (db *Database) Cleanup(ctx context.Context, count int, duration time.Duration) (int64, error) {
	res1, err := db.db.Exec(ctx, "DELETE FROM requests WHERE request_id NOT IN ( SELECT request_id FROM requests ORDER BY timestamp_created DESC LIMIT :lim ) ", sq.PP{
		"lim": count,
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
