package models

import "strconv"

type EntityID interface {
	IntID() int64
	String() string
}

type UserID int64

func (id UserID) IntID() int64 {
	return int64(id)
}

func (id UserID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

type ChannelID int64

func (id ChannelID) IntID() int64 {
	return int64(id)
}

func (id ChannelID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

type DeliveryID int64

func (id DeliveryID) IntID() int64 {
	return int64(id)
}

func (id DeliveryID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

type SCNMessageID int64

func (id SCNMessageID) IntID() int64 {
	return int64(id)
}

func (id SCNMessageID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

type SubscriptionID int64

func (id SubscriptionID) IntID() int64 {
	return int64(id)
}

func (id SubscriptionID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

type ClientID int64

func (id ClientID) IntID() int64 {
	return int64(id)
}

func (id ClientID) String() string {
	return strconv.FormatInt(int64(id), 10)
}
