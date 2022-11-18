package models

import "time"

type ClientType string

const (
	ClientTypeAndroid ClientType = "ANDROID"
	ClientTypeIOS     ClientType = "IOS"
)

type Client struct {
	ClientID         int64
	UserID           int64
	Type             ClientType
	FCMToken         *string
	TimestampCreated time.Time
	AgentModel       string
	AgentVersion     string
}

type ClientJSON struct {
}
