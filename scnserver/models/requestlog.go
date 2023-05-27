package models

import (
	"github.com/jmoiron/sqlx"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"gogs.mikescher.com/BlackForestBytes/goext/timeext"
	"time"
)

type RequestLog struct {
	RequestID           RequestID
	Method              string
	URI                 string
	UserAgent           *string
	Authentication      *string
	RequestBody         *string
	RequestBodySize     int64
	RequestContentType  string
	RemoteIP            string
	KeyID               *KeyTokenID
	UserID              *UserID
	Permissions         *string
	ResponseStatuscode  *int64
	ResponseBodySize    *int64
	ResponseBody        *string
	ResponseContentType string
	RetryCount          int64
	Panicked            bool
	PanicStr            *string
	ProcessingTime      time.Duration
	TimestampCreated    time.Time
	TimestampStart      time.Time
	TimestampFinish     time.Time
}

func (c RequestLog) JSON() RequestLogJSON {
	return RequestLogJSON{
		RequestID:           c.RequestID,
		Method:              c.Method,
		URI:                 c.URI,
		UserAgent:           c.UserAgent,
		Authentication:      c.Authentication,
		RequestBody:         c.RequestBody,
		RequestBodySize:     c.RequestBodySize,
		RequestContentType:  c.RequestContentType,
		RemoteIP:            c.RemoteIP,
		KeyID:               c.KeyID,
		UserID:              c.UserID,
		Permissions:         c.Permissions,
		ResponseStatuscode:  c.ResponseStatuscode,
		ResponseBodySize:    c.ResponseBodySize,
		ResponseBody:        c.ResponseBody,
		ResponseContentType: c.ResponseContentType,
		RetryCount:          c.RetryCount,
		Panicked:            c.Panicked,
		PanicStr:            c.PanicStr,
		ProcessingTime:      c.ProcessingTime.Seconds(),
		TimestampCreated:    c.TimestampCreated.Format(time.RFC3339Nano),
		TimestampStart:      c.TimestampStart.Format(time.RFC3339Nano),
		TimestampFinish:     c.TimestampFinish.Format(time.RFC3339Nano),
	}
}

func (c RequestLog) DB() RequestLogDB {
	return RequestLogDB{
		RequestID:           c.RequestID,
		Method:              c.Method,
		URI:                 c.URI,
		UserAgent:           c.UserAgent,
		Authentication:      c.Authentication,
		RequestBody:         c.RequestBody,
		RequestBodySize:     c.RequestBodySize,
		RequestContentType:  c.RequestContentType,
		RemoteIP:            c.RemoteIP,
		KeyID:               c.KeyID,
		UserID:              c.UserID,
		Permissions:         c.Permissions,
		ResponseStatuscode:  c.ResponseStatuscode,
		ResponseBodySize:    c.ResponseBodySize,
		ResponseBody:        c.ResponseBody,
		ResponseContentType: c.ResponseContentType,
		RetryCount:          c.RetryCount,
		Panicked:            langext.Conditional[int64](c.Panicked, 1, 0),
		PanicStr:            c.PanicStr,
		ProcessingTime:      c.ProcessingTime.Milliseconds(),
		TimestampCreated:    c.TimestampCreated.UnixMilli(),
		TimestampStart:      c.TimestampStart.UnixMilli(),
		TimestampFinish:     c.TimestampFinish.UnixMilli(),
	}
}

type RequestLogJSON struct {
	RequestID           RequestID   `json:"requestLog_id"`
	Method              string      `json:"method"`
	URI                 string      `json:"uri"`
	UserAgent           *string     `json:"user_agent"`
	Authentication      *string     `json:"authentication"`
	RequestBody         *string     `json:"request_body"`
	RequestBodySize     int64       `json:"request_body_size"`
	RequestContentType  string      `json:"request_content_type"`
	RemoteIP            string      `json:"remote_ip"`
	KeyID               *KeyTokenID `json:"key_id"`
	UserID              *UserID     `json:"userid"`
	Permissions         *string     `json:"permissions"`
	ResponseStatuscode  *int64      `json:"response_statuscode"`
	ResponseBodySize    *int64      `json:"response_body_size"`
	ResponseBody        *string     `json:"response_body"`
	ResponseContentType string      `json:"response_content_type"`
	RetryCount          int64       `json:"retry_count"`
	Panicked            bool        `json:"panicked"`
	PanicStr            *string     `json:"panic_str"`
	ProcessingTime      float64     `json:"processing_time"`
	TimestampCreated    string      `json:"timestamp_created"`
	TimestampStart      string      `json:"timestamp_start"`
	TimestampFinish     string      `json:"timestamp_finish"`
}

type RequestLogDB struct {
	RequestID           RequestID   `db:"request_id"`
	Method              string      `db:"method"`
	URI                 string      `db:"uri"`
	UserAgent           *string     `db:"user_agent"`
	Authentication      *string     `db:"authentication"`
	RequestBody         *string     `db:"request_body"`
	RequestBodySize     int64       `db:"request_body_size"`
	RequestContentType  string      `db:"request_content_type"`
	RemoteIP            string      `db:"remote_ip"`
	KeyID               *KeyTokenID `db:"key_id"`
	UserID              *UserID     `db:"userid"`
	Permissions         *string     `db:"permissions"`
	ResponseStatuscode  *int64      `db:"response_statuscode"`
	ResponseBodySize    *int64      `db:"response_body_size"`
	ResponseBody        *string     `db:"response_body"`
	ResponseContentType string      `db:"response_content_type"`
	RetryCount          int64       `db:"retry_count"`
	Panicked            int64       `db:"panicked"`
	PanicStr            *string     `db:"panic_str"`
	ProcessingTime      int64       `db:"processing_time"`
	TimestampCreated    int64       `db:"timestamp_created"`
	TimestampStart      int64       `db:"timestamp_start"`
	TimestampFinish     int64       `db:"timestamp_finish"`
}

func (c RequestLogDB) Model() RequestLog {
	return RequestLog{
		RequestID:           c.RequestID,
		Method:              c.Method,
		URI:                 c.URI,
		UserAgent:           c.UserAgent,
		Authentication:      c.Authentication,
		RequestBody:         c.RequestBody,
		RequestBodySize:     c.RequestBodySize,
		RequestContentType:  c.RequestContentType,
		RemoteIP:            c.RemoteIP,
		KeyID:               c.KeyID,
		UserID:              c.UserID,
		Permissions:         c.Permissions,
		ResponseStatuscode:  c.ResponseStatuscode,
		ResponseBodySize:    c.ResponseBodySize,
		ResponseBody:        c.ResponseBody,
		ResponseContentType: c.ResponseContentType,
		RetryCount:          c.RetryCount,
		Panicked:            c.Panicked != 0,
		PanicStr:            c.PanicStr,
		ProcessingTime:      timeext.FromMilliseconds(c.ProcessingTime),
		TimestampCreated:    timeFromMilli(c.TimestampCreated),
		TimestampStart:      timeFromMilli(c.TimestampStart),
		TimestampFinish:     timeFromMilli(c.TimestampFinish),
	}
}

func DecodeRequestLog(r *sqlx.Rows) (RequestLog, error) {
	data, err := sq.ScanSingle[RequestLogDB](r, sq.SModeFast, sq.Safe, true)
	if err != nil {
		return RequestLog{}, err
	}
	return data.Model(), nil
}

func DecodeRequestLogs(r *sqlx.Rows) ([]RequestLog, error) {
	data, err := sq.ScanAll[RequestLogDB](r, sq.SModeFast, sq.Safe, true)
	if err != nil {
		return nil, err
	}
	return langext.ArrMap(data, func(v RequestLogDB) RequestLog { return v.Model() }), nil
}
