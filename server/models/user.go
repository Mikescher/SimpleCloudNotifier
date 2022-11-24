package models

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"database/sql"
	"github.com/blockloop/scan"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"time"
)

type User struct {
	UserID            UserID
	Username          *string
	SendKey           string
	ReadKey           string
	AdminKey          string
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
		UserID:            u.UserID,
		Username:          u.Username,
		ReadKey:           u.ReadKey,
		SendKey:           u.SendKey,
		AdminKey:          u.AdminKey,
		TimestampCreated:  u.TimestampCreated.Format(time.RFC3339Nano),
		TimestampLastRead: timeOptFmt(u.TimestampLastRead, time.RFC3339Nano),
		TimestampLastSent: timeOptFmt(u.TimestampLastSent, time.RFC3339Nano),
		MessagesSent:      u.MessagesSent,
		QuotaUsed:         u.QuotaUsed,
		QuotaUsedDay:      u.QuotaUsedDay,
		IsPro:             u.IsPro,
		DefaultChannel:    u.DefaultChannel(),
	}
}

func (u User) JSONWithClients(clients []Client) UserJSONWithClients {
	return UserJSONWithClients{
		UserJSON: u.JSON(),
		Clients:  langext.ArrMap(clients, func(v Client) ClientJSON { return v.JSON() }),
	}
}

func (u User) MaxContentLength() int {
	if u.IsPro {
		return 16384
	} else {
		return 2048
	}
}

func (u User) MaxTitleLength() int {
	return 120
}

func (u User) QuotaPerDay() int {
	if u.IsPro {
		return 1000
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

func (u User) MaxChannelNameLength() int {
	return 120
}

type UserJSON struct {
	UserID            UserID  `json:"user_id"`
	Username          *string `json:"username"`
	ReadKey           string  `json:"read_key"`
	SendKey           string  `json:"send_key"`
	AdminKey          string  `json:"admin_key"`
	TimestampCreated  string  `json:"timestamp_created"`
	TimestampLastRead *string `json:"timestamp_last_read"`
	TimestampLastSent *string `json:"timestamp_last_sent"`
	MessagesSent      int     `json:"messages_sent"`
	QuotaUsed         int     `json:"quota_used"`
	QuotaUsedDay      *string `json:"quota_used_day"`
	IsPro             bool    `json:"is_pro"`
	DefaultChannel    string  `json:"default_channel"`
}

type UserJSONWithClients struct {
	UserJSON
	Clients []ClientJSON `json:"clients"`
}

type UserDB struct {
	UserID            UserID  `db:"user_id"`
	Username          *string `db:"username"`
	SendKey           string  `db:"send_key"`
	ReadKey           string  `db:"read_key"`
	AdminKey          string  `db:"admin_key"`
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
		SendKey:           u.SendKey,
		ReadKey:           u.ReadKey,
		AdminKey:          u.AdminKey,
		TimestampCreated:  time.UnixMilli(u.TimestampCreated),
		TimestampLastRead: timeOptFromMilli(u.TimestampLastRead),
		TimestampLastSent: timeOptFromMilli(u.TimestampLastSent),
		MessagesSent:      u.MessagesSent,
		QuotaUsed:         u.QuotaUsed,
		QuotaUsedDay:      u.QuotaUsedDay,
		IsPro:             u.IsPro,
	}
}

func DecodeUser(r *sql.Rows) (User, error) {
	var data UserDB
	err := scan.RowStrict(&data, r)
	if err != nil {
		return User{}, err
	}
	return data.Model(), nil
}

func DecodeUsers(r *sql.Rows) ([]User, error) {
	var data []UserDB
	err := scan.RowsStrict(&data, r)
	if err != nil {
		return nil, err
	}
	return langext.ArrMap(data, func(v UserDB) User { return v.Model() }), nil
}
