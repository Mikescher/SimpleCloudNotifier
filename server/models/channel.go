package models

import (
	"database/sql"
	"github.com/blockloop/scan"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"time"
)

type Channel struct {
	ChannelID         int64
	OwnerUserID       int64
	Name              string
	SubscribeKey      string
	SendKey           string
	TimestampCreated  time.Time
	TimestampLastSent *time.Time
	MessagesSent      int
}

func (c Channel) JSON(includeKey bool) ChannelJSON {
	return ChannelJSON{
		ChannelID:         c.ChannelID,
		OwnerUserID:       c.OwnerUserID,
		Name:              c.Name,
		SubscribeKey:      langext.Conditional(includeKey, langext.Ptr(c.SubscribeKey), nil),
		SendKey:           langext.Conditional(includeKey, langext.Ptr(c.SendKey), nil),
		TimestampCreated:  c.TimestampCreated.Format(time.RFC3339Nano),
		TimestampLastSent: timeOptFmt(c.TimestampLastSent, time.RFC3339Nano),
		MessagesSent:      c.MessagesSent,
	}
}

type ChannelJSON struct {
	ChannelID         int64   `json:"channel_id"`
	OwnerUserID       int64   `json:"owner_user_id"`
	Name              string  `json:"name"`
	SubscribeKey      *string `json:"subscribe_key"` // can be nil, depending on endpoint
	SendKey           *string `json:"send_key"`      // can be nil, depending on endpoint
	TimestampCreated  string  `json:"timestamp_created"`
	TimestampLastSent *string `json:"timestamp_last_sent"`
	MessagesSent      int     `json:"messages_sent"`
}

type ChannelDB struct {
	ChannelID         int64  `db:"channel_id"`
	OwnerUserID       int64  `db:"owner_user_id"`
	Name              string `db:"name"`
	SubscribeKey      string `db:"subscribe_key"`
	SendKey           string `db:"send_key"`
	TimestampCreated  int64  `db:"timestamp_created"`
	TimestampLastRead *int64 `db:"timestamp_last_read"`
	TimestampLastSent *int64 `db:"timestamp_last_sent"`
	MessagesSent      int    `db:"messages_sent"`
}

func (c ChannelDB) Model() Channel {
	return Channel{
		ChannelID:         c.ChannelID,
		OwnerUserID:       c.OwnerUserID,
		Name:              c.Name,
		SubscribeKey:      c.SubscribeKey,
		SendKey:           c.SendKey,
		TimestampCreated:  time.UnixMilli(c.TimestampCreated),
		TimestampLastSent: timeOptFromMilli(c.TimestampLastSent),
		MessagesSent:      c.MessagesSent,
	}
}

func DecodeChannel(r *sql.Rows) (Channel, error) {
	var data ChannelDB
	err := scan.RowStrict(&data, r)
	if err != nil {
		return Channel{}, err
	}
	return data.Model(), nil
}

func DecodeChannels(r *sql.Rows) ([]Channel, error) {
	var data []ChannelDB
	err := scan.RowsStrict(&data, r)
	if err != nil {
		return nil, err
	}
	return langext.ArrMap(data, func(v ChannelDB) Channel { return v.Model() }), nil
}
