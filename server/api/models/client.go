package models

import (
	"database/sql"
	"github.com/blockloop/scan"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"time"
)

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

func (c Client) JSON() ClientJSON {
	return ClientJSON{
		ClientID:         c.ClientID,
		UserID:           c.UserID,
		Type:             c.Type,
		FCMToken:         c.FCMToken,
		TimestampCreated: c.TimestampCreated.Format(time.RFC3339Nano),
		AgentModel:       c.AgentModel,
		AgentVersion:     c.AgentVersion,
	}
}

type ClientJSON struct {
	ClientID         int64      `json:"client_id"`
	UserID           int64      `json:"user_id"`
	Type             ClientType `json:"type"`
	FCMToken         *string    `json:"fcm_token"`
	TimestampCreated string     `json:"timestamp_created"`
	AgentModel       string     `json:"agent_model"`
	AgentVersion     string     `json:"agent_version"`
}

type ClientDB struct {
	ClientID         int64      `db:"client_id"`
	UserID           int64      `db:"user_id"`
	Type             ClientType `db:"type"`
	FCMToken         *string    `db:"fcm_token"`
	TimestampCreated int64      `db:"timestamp_created"`
	AgentModel       string     `db:"agent_model"`
	AgentVersion     string     `db:"agent_version"`
}

func (c ClientDB) Model() Client {
	return Client{
		ClientID:         c.ClientID,
		UserID:           c.UserID,
		Type:             c.Type,
		FCMToken:         c.FCMToken,
		TimestampCreated: time.UnixMilli(c.TimestampCreated),
		AgentModel:       c.AgentModel,
		AgentVersion:     c.AgentVersion,
	}
}

func DecodeClient(r *sql.Rows) (Client, error) {
	var data ClientDB
	err := scan.RowStrict(&data, r)
	if err != nil {
		return Client{}, err
	}
	return data.Model(), nil
}

func DecodeClients(r *sql.Rows) ([]Client, error) {
	var data []ClientDB
	err := scan.RowsStrict(&data, r)
	if err != nil {
		return nil, err
	}
	return langext.ArrMap(data, func(v ClientDB) Client { return v.Model() }), nil
}
