package models

import (
	scn "blackforestbytes.com/simplecloudnotifier"
)

type User struct {
	UserID            UserID   `db:"user_id"              json:"user_id"`
	Username          *string  `db:"username"             json:"username"`
	TimestampCreated  SCNTime  `db:"timestamp_created"    json:"timestamp_created"`
	TimestampLastRead *SCNTime `db:"timestamp_lastread"   json:"timestamp_lastread"`
	TimestampLastSent *SCNTime `db:"timestamp_lastsent"   json:"timestamp_lastsent"`
	MessagesSent      int      `db:"messages_sent"        json:"messages_sent"`
	QuotaUsed         int      `db:"quota_used"           json:"quota_used"`
	QuotaUsedDay      *string  `db:"quota_used_day"       json:"-"`
	IsPro             bool     `db:"is_pro"               json:"is_pro"`
	ProToken          *string  `db:"pro_token"            json:"-"`

	UserExtra `db:"-"` // fields that are not in DB and are set on PreMarshal
}

type UserExtra struct {
	QuotaRemaining              int    `json:"quota_remaining"`
	QuotaPerDay                 int    `json:"quota_max"`
	DefaultChannel              string `json:"default_channel"`
	MaxBodySize                 int    `json:"max_body_size"`
	MaxTitleLength              int    `json:"max_title_length"`
	DefaultPriority             int    `json:"default_priority"`
	MaxChannelNameLength        int    `json:"max_channel_name_length"`
	MaxChannelDescriptionLength int    `json:"max_channel_description_length"`
	MaxSenderNameLength         int    `json:"max_sender_name_length"`
	MaxUserMessageIDLength      int    `json:"max_user_message_id_length"`
}

type UserPreview struct {
	UserID   UserID  `json:"user_id"`
	Username *string `json:"username"`
}

type UserWithClientsAndKeys struct {
	User
	Clients  []Client `json:"clients"`
	SendKey  string   `json:"send_key"`
	ReadKey  string   `json:"read_key"`
	AdminKey string   `json:"admin_key"`
}

func (u User) WithClients(clients []Client, ak string, sk string, rk string) UserWithClientsAndKeys {
	return UserWithClientsAndKeys{
		User:     u.PreMarshal(),
		Clients:  clients,
		SendKey:  sk,
		ReadKey:  rk,
		AdminKey: ak,
	}
}

func (u *User) PreMarshal() User {
	u.UserExtra = UserExtra{
		QuotaPerDay:                 u.QuotaPerDay(),
		QuotaRemaining:              u.QuotaRemainingToday(),
		DefaultChannel:              u.DefaultChannel(),
		MaxBodySize:                 u.MaxContentLength(),
		MaxTitleLength:              u.MaxTitleLength(),
		DefaultPriority:             u.DefaultPriority(),
		MaxChannelNameLength:        u.MaxChannelNameLength(),
		MaxChannelDescriptionLength: u.MaxChannelDescriptionLength(),
		MaxSenderNameLength:         u.MaxSenderNameLength(),
		MaxUserMessageIDLength:      u.MaxUserMessageIDLength(),
	}
	return *u
}

func (u User) MaxContentLength() int {
	if u.IsPro {
		return 2 * 1024 * 1024 // 2 MB
	} else {
		return 2 * 1024 // 2 KB
	}
}

func (u User) MaxTitleLength() int {
	return 120
}

func (u User) QuotaPerDay() int {
	if u.IsPro {
		return 5000
	} else {
		return 50
	}
}

func (u User) QuotaUsedToday() int {
	now := scn.QuotaDayString()
	if u.QuotaUsedDay != nil && *u.QuotaUsedDay == now {
		return u.QuotaUsed
	} else {
		return 0
	}
}

func (u User) QuotaRemainingToday() int {
	return u.QuotaPerDay() - u.QuotaUsedToday()
}

func (u User) DefaultChannel() string {
	return "main"
}

func (u User) DefaultPriority() int {
	return 1
}

func (u User) MaxChannelNameLength() int {
	return 120
}

func (u User) MaxChannelDescriptionLength() int {
	return 300
}

func (u User) MaxSenderNameLength() int {
	return 120
}

func (u User) MaxUserMessageIDLength() int {
	return 64
}

func (u User) MaxTimestampDiffHours() int {
	return 24
}

func (u User) JSONPreview() UserPreview {
	return UserPreview{
		UserID:   u.UserID,
		Username: u.Username,
	}
}
