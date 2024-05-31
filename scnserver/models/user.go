package models

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"context"
	"github.com/jmoiron/sqlx"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"time"
)

type User struct {
	UserID            UserID
	Username          *string
	TimestampCreated  time.Time
	TimestampLastRead *time.Time
	TimestampLastSent *time.Time
	MessagesSent      int
	QuotaUsed         int
	QuotaUsedDay      *string
	IsPro             bool
	ProToken          *string
}

func (u User) JSON() UserJSON {
	return UserJSON{
		UserID:                      u.UserID,
		Username:                    u.Username,
		TimestampCreated:            u.TimestampCreated.Format(time.RFC3339Nano),
		TimestampLastRead:           timeOptFmt(u.TimestampLastRead, time.RFC3339Nano),
		TimestampLastSent:           timeOptFmt(u.TimestampLastSent, time.RFC3339Nano),
		MessagesSent:                u.MessagesSent,
		QuotaUsed:                   u.QuotaUsedToday(),
		QuotaPerDay:                 u.QuotaPerDay(),
		QuotaRemaining:              u.QuotaRemainingToday(),
		IsPro:                       u.IsPro,
		DefaultChannel:              u.DefaultChannel(),
		MaxBodySize:                 u.MaxContentLength(),
		MaxTitleLength:              u.MaxTitleLength(),
		DefaultPriority:             u.DefaultPriority(),
		MaxChannelNameLength:        u.MaxChannelNameLength(),
		MaxChannelDescriptionLength: u.MaxChannelDescriptionLength(),
		MaxSenderNameLength:         u.MaxSenderNameLength(),
		MaxUserMessageIDLength:      u.MaxUserMessageIDLength(),
	}
}

func (u User) JSONWithClients(clients []Client, ak string, sk string, rk string) UserJSONWithClientsAndKeys {
	return UserJSONWithClientsAndKeys{
		UserJSON: u.JSON(),
		Clients:  langext.ArrMap(clients, func(v Client) ClientJSON { return v.JSON() }),
		SendKey:  sk,
		ReadKey:  rk,
		AdminKey: ak,
	}
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

type UserJSON struct {
	UserID                      UserID  `json:"user_id"`
	Username                    *string `json:"username"`
	TimestampCreated            string  `json:"timestamp_created"`
	TimestampLastRead           *string `json:"timestamp_lastread"`
	TimestampLastSent           *string `json:"timestamp_lastsent"`
	MessagesSent                int     `json:"messages_sent"`
	QuotaUsed                   int     `json:"quota_used"`
	QuotaRemaining              int     `json:"quota_remaining"`
	QuotaPerDay                 int     `json:"quota_max"`
	IsPro                       bool    `json:"is_pro"`
	DefaultChannel              string  `json:"default_channel"`
	MaxBodySize                 int     `json:"max_body_size"`
	MaxTitleLength              int     `json:"max_title_length"`
	DefaultPriority             int     `json:"default_priority"`
	MaxChannelNameLength        int     `json:"max_channel_name_length"`
	MaxChannelDescriptionLength int     `json:"max_channel_description_length"`
	MaxSenderNameLength         int     `json:"max_sender_name_length"`
	MaxUserMessageIDLength      int     `json:"max_user_message_id_length"`
}

type UserJSONWithClientsAndKeys struct {
	UserJSON
	Clients  []ClientJSON `json:"clients"`
	SendKey  string       `json:"send_key"`
	ReadKey  string       `json:"read_key"`
	AdminKey string       `json:"admin_key"`
}

type UserDB struct {
	UserID            UserID  `db:"user_id"`
	Username          *string `db:"username"`
	TimestampCreated  int64   `db:"timestamp_created"`
	TimestampLastRead *int64  `db:"timestamp_lastread"`
	TimestampLastSent *int64  `db:"timestamp_lastsent"`
	MessagesSent      int     `db:"messages_sent"`
	QuotaUsed         int     `db:"quota_used"`
	QuotaUsedDay      *string `db:"quota_used_day"`
	IsPro             bool    `db:"is_pro"`
	ProToken          *string `db:"pro_token"`
}

func (u UserDB) Model() User {
	return User{
		UserID:            u.UserID,
		Username:          u.Username,
		TimestampCreated:  timeFromMilli(u.TimestampCreated),
		TimestampLastRead: timeOptFromMilli(u.TimestampLastRead),
		TimestampLastSent: timeOptFromMilli(u.TimestampLastSent),
		MessagesSent:      u.MessagesSent,
		QuotaUsed:         u.QuotaUsed,
		QuotaUsedDay:      u.QuotaUsedDay,
		IsPro:             u.IsPro,
	}
}

func DecodeUser(ctx context.Context, q sq.Queryable, r *sqlx.Rows) (User, error) {
	data, err := sq.ScanSingle[UserDB](ctx, q, r, sq.SModeFast, sq.Safe, true)
	if err != nil {
		return User{}, err
	}
	return data.Model(), nil
}

func DecodeUsers(ctx context.Context, q sq.Queryable, r *sqlx.Rows) ([]User, error) {
	data, err := sq.ScanAll[UserDB](ctx, q, r, sq.SModeFast, sq.Safe, true)
	if err != nil {
		return nil, err
	}
	return langext.ArrMap(data, func(v UserDB) User { return v.Model() }), nil
}
