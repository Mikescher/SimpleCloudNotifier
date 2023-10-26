package models

import (
	"gogs.mikescher.com/BlackForestBytes/goext/rext"
)

//go:generate go run ../_gen/id-generate.go -- ids_gen.go

type EntityID interface {
	String() string
	Valid() error
	Prefix() string
	Raw() string
	CheckString() string
	Regex() rext.Regex
}

type UserID string         //@csid:type [USR]
type ChannelID string      //@csid:type [CHA]
type DeliveryID string     //@csid:type [DEL]
type MessageID string      //@csid:type [MSG]
type SubscriptionID string //@csid:type [SUB]
type ClientID string       //@csid:type [CLN]
type RequestID string      //@csid:type [REQ]
type KeyTokenID string     //@csid:type [TOK]
