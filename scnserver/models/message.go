package models

import (
	"fmt"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"time"
)

const (
	ContentLengthTrim  = 1900
	ContentLengthShort = 200
)

type Message struct {
	MessageID           MessageID  `db:"message_id"              json:"message_id"`
	SenderUserID        UserID     `db:"sender_user_id"          json:"sender_user_id"` // user that sent the message (this is also the owner of the channel that contains it)
	ChannelInternalName string     `db:"channel_internal_name"   json:"channel_internal_name"`
	ChannelID           ChannelID  `db:"channel_id"              json:"channel_id"`
	SenderName          *string    `db:"sender_name"             json:"sender_name"`
	SenderIP            string     `db:"sender_ip"               json:"sender_ip"`
	TimestampReal       SCNTime    `db:"timestamp_real"          json:"-"`
	TimestampClient     *SCNTime   `db:"timestamp_client"        json:"-"`
	Title               string     `db:"title"                   json:"title"`
	Content             *string    `db:"content"                 json:"content"`
	Priority            int        `db:"priority"                json:"priority"`
	UserMessageID       *string    `db:"usr_message_id"          json:"usr_message_id"`
	UsedKeyID           KeyTokenID `db:"used_key_id"             json:"used_key_id"`
	Deleted             bool       `db:"deleted"                 json:"-"`

	MessageExtra `db:"-"` // fields that are not in DB and are set on PreMarshal
}

type MessageExtra struct {
	Timestamp SCNTime `db:"-" json:"timestamp"`
	Trimmed   bool    `db:"-" json:"trimmed"`
}

func (u *Message) PreMarshal() Message {
	u.MessageExtra.Timestamp = NewSCNTime(u.Timestamp())
	return *u
}

func (m Message) Trim() Message {
	r := m
	if !r.Trimmed && r.NeedsTrim() {
		r.Content = r.TrimmedContent()
		r.MessageExtra.Trimmed = true
	}
	return r.PreMarshal()
}

func (m Message) Timestamp() time.Time {
	return langext.Coalesce(m.TimestampClient, m.TimestampReal).Time()
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

func (m Message) FormatNotificationTitle(user User, channel Channel) string {
	if m.ChannelInternalName == user.DefaultChannel() {
		return m.Title
	}

	return fmt.Sprintf("[%s] %s", channel.DisplayName, m.Title)
}
