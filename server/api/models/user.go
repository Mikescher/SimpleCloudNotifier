package models

import "time"

type User struct {
	UserID            int64
	Username          *string
	ReadKey           string
	SendKey           string
	AdminKey          string
	TimestampCreated  time.Time
	TimestampLastRead *time.Time
	TimestampLastSent *time.Time
	MessagesSent      int
	QuotaToday        int
	QuotaDay          *string
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
		QuotaToday:        u.QuotaToday,
		QuotaDay:          u.QuotaDay,
		IsPro:             u.IsPro,
	}
}

type UserJSON struct {
	UserID            int64   `json:"user_id"`
	Username          *string `json:"username"`
	ReadKey           string  `json:"read_key"`
	SendKey           string  `json:"send_key"`
	AdminKey          string  `json:"admin_key"`
	TimestampCreated  string  `json:"timestamp_created"`
	TimestampLastRead *string `json:"timestamp_last_read"`
	TimestampLastSent *string `json:"timestamp_last_sent"`
	MessagesSent      int     `json:"messages_sent"`
	QuotaToday        int     `json:"quota_today"`
	QuotaDay          *string `json:"quota_day"`
	IsPro             bool    `json:"is_pro"`
}
