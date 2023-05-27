package models

import (
	"github.com/jmoiron/sqlx"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"strings"
	"time"
)

type TokenPerm string //@enum:type

const (
	PermAdmin       TokenPerm = "A"  // Edit userdata (+ includes all other permissions)
	PermChannelRead TokenPerm = "CR" // Read messages
	PermChannelSend TokenPerm = "CS" // Send messages
	PermUserRead    TokenPerm = "UR" // Read userdata
)

type TokenPermissionList []TokenPerm

func (e TokenPermissionList) Any(p ...TokenPerm) bool {
	for _, v1 := range e {
		for _, v2 := range p {
			if v1 == v2 {
				return true
			}
		}
	}
	return false
}

func (e TokenPermissionList) String() string {
	return strings.Join(langext.ArrMap(e, func(v TokenPerm) string { return string(v) }), ";")
}

func ParseTokenPermissionList(input string) TokenPermissionList {
	r := make([]TokenPerm, 0, len(input))
	for _, v := range strings.Split(input, ";") {
		if vv, ok := ParseTokenPerm(v); ok {
			r = append(r, vv)
		}
	}
	return r
}

type KeyToken struct {
	KeyTokenID        KeyTokenID
	Name              string
	TimestampCreated  time.Time
	TimestampLastUsed *time.Time
	OwnerUserID       UserID
	AllChannels       bool
	Channels          []ChannelID // can also be owned by other user (needs active subscription)
	Token             string
	Permissions       TokenPermissionList
	MessagesSent      int
}

func (k KeyToken) IsUserRead(uid UserID) bool {
	return k.OwnerUserID == uid && k.Permissions.Any(PermAdmin, PermUserRead)
}

func (k KeyToken) IsAllMessagesRead(uid UserID) bool {
	return k.OwnerUserID == uid && k.AllChannels == true && k.Permissions.Any(PermAdmin, PermChannelRead)
}

func (k KeyToken) IsChannelMessagesRead(cid ChannelID) bool {
	return (k.AllChannels == true || langext.InArray(cid, k.Channels)) && k.Permissions.Any(PermAdmin, PermChannelRead)
}

func (k KeyToken) IsAdmin(uid UserID) bool {
	return k.OwnerUserID == uid && k.Permissions.Any(PermAdmin)
}

func (k KeyToken) IsChannelMessagesSend(c Channel) bool {
	return (k.AllChannels == true || langext.InArray(c.ChannelID, k.Channels)) && k.OwnerUserID == c.OwnerUserID && k.Permissions.Any(PermAdmin, PermChannelSend)
}

func (k KeyToken) JSON() KeyTokenJSON {
	return KeyTokenJSON{
		KeyTokenID:        k.KeyTokenID,
		Name:              k.Name,
		TimestampCreated:  k.TimestampCreated,
		TimestampLastUsed: k.TimestampLastUsed,
		OwnerUserID:       k.OwnerUserID,
		AllChannels:       k.AllChannels,
		Channels:          k.Channels,
		Permissions:       k.Permissions.String(),
		MessagesSent:      k.MessagesSent,
	}
}

type KeyTokenJSON struct {
	KeyTokenID        KeyTokenID  `json:"keytoken_id"`
	Name              string      `json:"name"`
	TimestampCreated  time.Time   `json:"timestamp_created"`
	TimestampLastUsed *time.Time  `json:"timestamp_lastused"`
	OwnerUserID       UserID      `json:"owner_user_id"`
	AllChannels       bool        `json:"all_channels"`
	Channels          []ChannelID `json:"channels"`
	Permissions       string      `json:"permissions"`
	MessagesSent      int         `json:"messages_sent"`
}

type KeyTokenWithTokenJSON struct {
	KeyTokenJSON
	Token string `json:"token"`
}

func (j KeyTokenJSON) WithToken(tok string) any {
	return KeyTokenWithTokenJSON{
		KeyTokenJSON: j,
		Token:        tok,
	}
}

type KeyTokenDB struct {
	KeyTokenID        KeyTokenID `db:"keytoken_id"`
	Name              string     `db:"name"`
	TimestampCreated  int64      `db:"timestamp_created"`
	TimestampLastUsed *int64     `db:"timestamp_lastused"`
	OwnerUserID       UserID     `db:"owner_user_id"`
	AllChannels       bool       `db:"all_channels"`
	Channels          string     `db:"channels"`
	Token             string     `db:"token"`
	Permissions       string     `db:"permissions"`
	MessagesSent      int        `db:"messages_sent"`
}

func (k KeyTokenDB) Model() KeyToken {
	return KeyToken{
		KeyTokenID:        k.KeyTokenID,
		Name:              k.Name,
		TimestampCreated:  timeFromMilli(k.TimestampCreated),
		TimestampLastUsed: timeOptFromMilli(k.TimestampLastUsed),
		OwnerUserID:       k.OwnerUserID,
		AllChannels:       k.AllChannels,
		Channels:          langext.ArrMap(strings.Split(k.Channels, ";"), func(v string) ChannelID { return ChannelID(v) }),
		Token:             k.Token,
		Permissions:       ParseTokenPermissionList(k.Permissions),
		MessagesSent:      k.MessagesSent,
	}
}

func DecodeKeyToken(r *sqlx.Rows) (KeyToken, error) {
	data, err := sq.ScanSingle[KeyTokenDB](r, sq.SModeFast, sq.Safe, true)
	if err != nil {
		return KeyToken{}, err
	}
	return data.Model(), nil
}

func DecodeKeyTokens(r *sqlx.Rows) ([]KeyToken, error) {
	data, err := sq.ScanAll[KeyTokenDB](r, sq.SModeFast, sq.Safe, true)
	if err != nil {
		return nil, err
	}
	return langext.ArrMap(data, func(v KeyTokenDB) KeyToken { return v.Model() }), nil
}
