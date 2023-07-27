package models

import (
	"github.com/jmoiron/sqlx"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"time"
)

const (
	ContentLengthTrim  = 1900
	ContentLengthShort = 200
)

type Message struct {
	MessageID           MessageID
	SenderUserID        UserID // user that sent the message (this is also the owner of the channel that contains it)
	ChannelInternalName string
	ChannelID           ChannelID
	SenderName          *string
	SenderIP            string
	TimestampReal       time.Time
	TimestampClient     *time.Time
	Title               string
	Content             *string
	Priority            int
	UserMessageID       *string
	UsedKeyID           KeyTokenID
	Deleted             bool
}

func (m Message) FullJSON() MessageJSON {
	return MessageJSON{
		MessageID:           m.MessageID,
		SenderUserID:        m.SenderUserID,
		ChannelInternalName: m.ChannelInternalName,
		ChannelID:           m.ChannelID,
		SenderName:          m.SenderName,
		SenderIP:            m.SenderIP,
		Timestamp:           m.Timestamp().Format(time.RFC3339Nano),
		Title:               m.Title,
		Content:             m.Content,
		Priority:            m.Priority,
		UserMessageID:       m.UserMessageID,
		UsedKeyID:           m.UsedKeyID,
		Trimmed:             false,
	}
}

func (m Message) TrimmedJSON() MessageJSON {
	return MessageJSON{
		MessageID:           m.MessageID,
		SenderUserID:        m.SenderUserID,
		ChannelInternalName: m.ChannelInternalName,
		ChannelID:           m.ChannelID,
		SenderName:          m.SenderName,
		SenderIP:            m.SenderIP,
		Timestamp:           m.Timestamp().Format(time.RFC3339Nano),
		Title:               m.Title,
		Content:             m.TrimmedContent(),
		Priority:            m.Priority,
		UserMessageID:       m.UserMessageID,
		UsedKeyID:           m.UsedKeyID,
		Trimmed:             m.NeedsTrim(),
	}
}

func (m Message) Timestamp() time.Time {
	return langext.Coalesce(m.TimestampClient, m.TimestampReal)
}

func (m Message) NeedsTrim() bool {
	return m.Content != nil && len(*m.Content) > ContentLengthTrim
}

func (m Message) TrimmedContent() *string {
	if m.Content == nil {
		return nil
	}
	if !m.NeedsTrim() {
		return m.Content
	}
	return langext.Ptr(langext.Coalesce(m.Content, "")[0:ContentLengthTrim-3] + "...")
}

func (m Message) ShortContent() string {
	if m.Content == nil {
		return ""
	}
	if len(*m.Content) < ContentLengthShort {
		return *m.Content
	}
	return (*m.Content)[0:ContentLengthShort-3] + "..."
}

type MessageJSON struct {
	MessageID           MessageID  `json:"message_id"`
	SenderUserID        UserID     `json:"sender_user_id"`
	ChannelInternalName string     `json:"channel_internal_name"`
	ChannelID           ChannelID  `json:"channel_id"`
	SenderName          *string    `json:"sender_name"`
	SenderIP            string     `json:"sender_ip"`
	Timestamp           string     `json:"timestamp"`
	Title               string     `json:"title"`
	Content             *string    `json:"content"`
	Priority            int        `json:"priority"`
	UserMessageID       *string    `json:"usr_message_id"`
	UsedKeyID           KeyTokenID `json:"used_key_id"`
	Trimmed             bool       `json:"trimmed"`
}

type MessageDB struct {
	MessageID           MessageID  `db:"message_id"`
	SenderUserID        UserID     `db:"sender_user_id"`
	ChannelInternalName string     `db:"channel_internal_name"`
	ChannelID           ChannelID  `db:"channel_id"`
	SenderName          *string    `db:"sender_name"`
	SenderIP            string     `db:"sender_ip"`
	TimestampReal       int64      `db:"timestamp_real"`
	TimestampClient     *int64     `db:"timestamp_client"`
	Title               string     `db:"title"`
	Content             *string    `db:"content"`
	Priority            int        `db:"priority"`
	UserMessageID       *string    `db:"usr_message_id"`
	UsedKeyID           KeyTokenID `db:"used_key_id"`
	Deleted             int        `db:"deleted"`
}

func (m MessageDB) Model() Message {
	return Message{
		MessageID:           m.MessageID,
		SenderUserID:        m.SenderUserID,
		ChannelInternalName: m.ChannelInternalName,
		ChannelID:           m.ChannelID,
		SenderName:          m.SenderName,
		SenderIP:            m.SenderIP,
		TimestampReal:       timeFromMilli(m.TimestampReal),
		TimestampClient:     timeOptFromMilli(m.TimestampClient),
		Title:               m.Title,
		Content:             m.Content,
		Priority:            m.Priority,
		UserMessageID:       m.UserMessageID,
		UsedKeyID:           m.UsedKeyID,
		Deleted:             m.Deleted != 0,
	}
}

func DecodeMessage(r *sqlx.Rows) (Message, error) {
	data, err := sq.ScanSingle[MessageDB](r, sq.SModeFast, sq.Safe, true)
	if err != nil {
		return Message{}, err
	}
	return data.Model(), nil
}

func DecodeMessages(r *sqlx.Rows) ([]Message, error) {
	data, err := sq.ScanAll[MessageDB](r, sq.SModeFast, sq.Safe, true)
	if err != nil {
		return nil, err
	}
	return langext.ArrMap(data, func(v MessageDB) Message { return v.Model() }), nil
}
