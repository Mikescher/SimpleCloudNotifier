package models

import (
	"context"
	"github.com/jmoiron/sqlx"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"time"
)

type ClientType string //@enum:type

const (
	ClientTypeAndroid ClientType = "ANDROID"
	ClientTypeIOS     ClientType = "IOS"
	ClientTypeLinux   ClientType = "LINUX"
	ClientTypeMacOS   ClientType = "MACOS"
	ClientTypeWindows ClientType = "WINDOWS"
)

type Client struct {
	ClientID         ClientID
	UserID           UserID
	Type             ClientType
	FCMToken         string
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
	ClientID         ClientID   `json:"client_id"`
	UserID           UserID     `json:"user_id"`
	Type             ClientType `json:"type"`
	FCMToken         string     `json:"fcm_token"`
	TimestampCreated string     `json:"timestamp_created"`
	AgentModel       string     `json:"agent_model"`
	AgentVersion     string     `json:"agent_version"`
}

type ClientDB struct {
	ClientID         ClientID   `db:"client_id"`
	UserID           UserID     `db:"user_id"`
	Type             ClientType `db:"type"`
	FCMToken         string     `db:"fcm_token"`
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
		TimestampCreated: timeFromMilli(c.TimestampCreated),
		AgentModel:       c.AgentModel,
		AgentVersion:     c.AgentVersion,
	}
}

func DecodeClient(ctx context.Context, q sq.Queryable, r *sqlx.Rows) (Client, error) {
	data, err := sq.ScanSingle[ClientDB](ctx, q, r, sq.SModeFast, sq.Safe, true)
	if err != nil {
		return Client{}, err
	}
	return data.Model(), nil
}

func DecodeClients(ctx context.Context, q sq.Queryable, r *sqlx.Rows) ([]Client, error) {
	data, err := sq.ScanAll[ClientDB](ctx, q, r, sq.SModeFast, sq.Safe, true)
	if err != nil {
		return nil, err
	}
	return langext.ArrMap(data, func(v ClientDB) Client { return v.Model() }), nil
}
