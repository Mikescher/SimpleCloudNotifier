package models

type RequestLog struct {
	RequestID           RequestID   `db:"request_id"             json:"requestLog_id"`
	Method              string      `db:"method"                 json:"method"`
	URI                 string      `db:"uri"                    json:"uri"`
	UserAgent           *string     `db:"user_agent"             json:"user_agent"`
	Authentication      *string     `db:"authentication"         json:"authentication"`
	RequestBody         *string     `db:"request_body"           json:"request_body"`
	RequestBodySize     int64       `db:"request_body_size"      json:"request_body_size"`
	RequestContentType  string      `db:"request_content_type"   json:"request_content_type"`
	RemoteIP            string      `db:"remote_ip"              json:"remote_ip"`
	KeyID               *KeyTokenID `db:"key_id"                 json:"key_id"`
	UserID              *UserID     `db:"userid"                 json:"userid"`
	Permissions         *string     `db:"permissions"            json:"permissions"`
	ResponseStatuscode  *int64      `db:"response_statuscode"    json:"response_statuscode"`
	ResponseBodySize    *int64      `db:"response_body_size"     json:"response_body_size"`
	ResponseBody        *string     `db:"response_body"          json:"response_body"`
	ResponseContentType string      `db:"response_content_type"  json:"response_content_type"`
	RetryCount          int64       `db:"retry_count"            json:"retry_count"`
	Panicked            bool        `db:"panicked"               json:"panicked"`
	PanicStr            *string     `db:"panic_str"              json:"panic_str"`
	ProcessingTime      SCNDuration `db:"processing_time"        json:"processing_time"`
	TimestampCreated    SCNTime     `db:"timestamp_created"      json:"timestamp_created"`
	TimestampStart      SCNTime     `db:"timestamp_start"        json:"timestamp_start"`
	TimestampFinish     SCNTime     `db:"timestamp_finish"       json:"timestamp_finish"`
}
