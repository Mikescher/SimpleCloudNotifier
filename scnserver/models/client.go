package models

type ClientType string //@enum:type

const (
	ClientTypeAndroid ClientType = "ANDROID"
	ClientTypeIOS     ClientType = "IOS"
	ClientTypeLinux   ClientType = "LINUX"
	ClientTypeMacOS   ClientType = "MACOS"
	ClientTypeWindows ClientType = "WINDOWS"
)

type Client struct {
	ClientID         ClientID   `db:"client_id"         json:"client_id"`
	UserID           UserID     `db:"user_id"           json:"user_id"`
	Type             ClientType `db:"type"              json:"type"`
	FCMToken         string     `db:"fcm_token"         json:"fcm_token"`
	TimestampCreated SCNTime    `db:"timestamp_created" json:"timestamp_created"`
	AgentModel       string     `db:"agent_model"       json:"agent_model"`
	AgentVersion     string     `db:"agent_version"     json:"agent_version"`
	Name             *string    `db:"name"              json:"name"`
	Deleted          bool       `db:"deleted"           json:"-"`
}
