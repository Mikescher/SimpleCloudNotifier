package models

import (
	"database/sql"
	"github.com/blockloop/scan"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"time"
)

type Message struct {
	SCNMessageID    int64
	SenderUserID    int64
	OwnerUserID     int64
	ChannelName     string
	ChannelID       int64
	TimestampReal   time.Time
	TimestampClient *time.Time
	Title           string
	Content         *string
	Priority        int
	UserMessageID   *string
}

func (m Message) JSON() MessageJSON {
	return MessageJSON{
		SCNMessageID:  m.SCNMessageID,
		SenderUserID:  m.SenderUserID,
		OwnerUserID:   m.OwnerUserID,
		ChannelName:   m.ChannelName,
		ChannelID:     m.ChannelID,
		Timestamp:     langext.Coalesce(m.TimestampClient, m.TimestampReal).Format(time.RFC3339Nano),
		Title:         m.Title,
		Content:       m.Content,
		Priority:      m.Priority,
		UserMessageID: m.UserMessageID,
		Trimmed:       false,
	}
}

type MessageJSON struct {
	SCNMessageID  int64   `json:"scn_message_id"`
	SenderUserID  int64   `json:"sender_user_id"`
	OwnerUserID   int64   `json:"owner_user_id"`
	ChannelName   string  `json:"channel_name"`
	ChannelID     int64   `json:"channel_id"`
	Timestamp     string  `json:"timestamp"`
	Title         string  `json:"title"`
	Content       *string `json:"body"`
	Priority      int     `json:"priority"`
	UserMessageID *string `json:"usr_message_id"`
	Trimmed       bool    `json:"trimmed"`
}

type MessageDB struct {
	SCNMessageID    int64   `db:"scn_message_id"`
	SenderUserID    int64   `db:"sender_user_id"`
	OwnerUserID     int64   `db:"owner_user_id"`
	ChannelName     string  `db:"channel_name"`
	ChannelID       int64   `db:"channel_id"`
	TimestampReal   int64   `db:"timestamp_real"`
	TimestampClient *int64  `db:"timestamp_client"`
	Title           string  `db:"title"`
	Content         *string `db:"content"`
	Priority        int     `db:"priority"`
	UserMessageID   *string `db:"usr_message_id"`
}

func (m MessageDB) Model() Message {
	return Message{
		SCNMessageID:    m.SCNMessageID,
		SenderUserID:    m.SenderUserID,
		OwnerUserID:     m.OwnerUserID,
		ChannelName:     m.ChannelName,
		ChannelID:       m.ChannelID,
		TimestampReal:   time.UnixMilli(m.TimestampReal),
		TimestampClient: timeOptFromMilli(m.TimestampClient),
		Title:           m.Title,
		Content:         m.Content,
		Priority:        m.Priority,
		UserMessageID:   m.UserMessageID,
	}
}

func DecodeMessage(r *sql.Rows) (Message, error) {
	var data MessageDB
	err := scan.RowStrict(&data, r)
	if err != nil {
		return Message{}, err
	}
	return data.Model(), nil
}

func DecodeMessages(r *sql.Rows) ([]Message, error) {
	var data []MessageDB
	err := scan.RowsStrict(&data, r)
	if err != nil {
		return nil, err
	}
	return langext.ArrMap(data, func(v MessageDB) Message { return v.Model() }), nil
}
