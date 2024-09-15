package models

import (
	"encoding/json"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"strings"
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

func (e TokenPermissionList) MarshalToDB(v TokenPermissionList) (string, error) {
	return v.String(), nil
}

func (e TokenPermissionList) UnmarshalToModel(v string) (TokenPermissionList, error) {
	return ParseTokenPermissionList(v), nil
}

func (t TokenPermissionList) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

type ChannelIDArr []ChannelID

func (t ChannelIDArr) MarshalToDB(v ChannelIDArr) (string, error) {
	return strings.Join(langext.ArrMap(v, func(v ChannelID) string { return v.String() }), ";"), nil
}

func (t ChannelIDArr) UnmarshalToModel(v string) (ChannelIDArr, error) {
	channels := make([]ChannelID, 0)
	if strings.TrimSpace(v) != "" {
		channels = langext.ArrMap(strings.Split(v, ";"), func(v string) ChannelID { return ChannelID(v) })
	}

	return channels, nil
}

type KeyToken struct {
	KeyTokenID        KeyTokenID          `db:"keytoken_id"          json:"keytoken_id"`
	Name              string              `db:"name"                 json:"name"`
	TimestampCreated  SCNTime             `db:"timestamp_created"    json:"timestamp_created"`
	TimestampLastUsed *SCNTime            `db:"timestamp_lastused"   json:"timestamp_lastused"`
	OwnerUserID       UserID              `db:"owner_user_id"        json:"owner_user_id"`
	AllChannels       bool                `db:"all_channels"         json:"all_channels"`
	Channels          ChannelIDArr        `db:"channels"             json:"channels"`
	Token             string              `db:"token"                json:"token"               jsonfilter:"INCLUDE_TOKEN"`
	Permissions       TokenPermissionList `db:"permissions"          json:"permissions"`
	MessagesSent      int                 `db:"messages_sent"        json:"messages_sent"`
}

type KeyTokenPreview struct {
	KeyTokenID  KeyTokenID  `json:"keytoken_id"`
	Name        string      `json:"name"`
	OwnerUserID UserID      `json:"owner_user_id"`
	AllChannels bool        `json:"all_channels"`
	Channels    []ChannelID `json:"channels"`
	Permissions string      `json:"permissions"`
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

func (k KeyToken) Preview() KeyTokenPreview {
	return KeyTokenPreview{
		KeyTokenID:  k.KeyTokenID,
		Name:        k.Name,
		OwnerUserID: k.OwnerUserID,
		AllChannels: k.AllChannels,
		Channels:    k.Channels,
		Permissions: k.Permissions.String(),
	}
}
