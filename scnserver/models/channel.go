package models

import (
	"github.com/jmoiron/sqlx"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"time"
)

type Channel struct {
	ChannelID         ChannelID
	OwnerUserID       UserID
	InternalName      string
	DisplayName       string
	DescriptionName   *string
	SubscribeKey      string
	SendKey           string
	TimestampCreated  time.Time
	TimestampLastSent *time.Time
	MessagesSent      int
}

func (c Channel) JSON(includeKey bool) ChannelJSON {
	return ChannelJSON{
		ChannelID:         c.ChannelID,
		OwnerUserID:       c.OwnerUserID,
		InternalName:      c.InternalName,
		DisplayName:       c.DisplayName,
		DescriptionName:   c.DescriptionName,
		SubscribeKey:      langext.Conditional(includeKey, langext.Ptr(c.SubscribeKey), nil),
		SendKey:           langext.Conditional(includeKey, langext.Ptr(c.SendKey), nil),
		TimestampCreated:  c.TimestampCreated.Format(time.RFC3339Nano),
		TimestampLastSent: timeOptFmt(c.TimestampLastSent, time.RFC3339Nano),
		MessagesSent:      c.MessagesSent,
	}
}

func (c Channel) WithSubscription(sub *Subscription) ChannelWithSubscription {
	return ChannelWithSubscription{
		Channel:      c,
		Subscription: sub,
	}
}

type ChannelWithSubscription struct {
	Channel
	Subscription *Subscription
}

func (c ChannelWithSubscription) JSON(includeChannelKey bool) ChannelWithSubscriptionJSON {
	var sub *SubscriptionJSON = nil
	if c.Subscription != nil {
		sub = langext.Ptr(c.Subscription.JSON())
	}
	return ChannelWithSubscriptionJSON{
		ChannelJSON:  c.Channel.JSON(includeChannelKey),
		Subscription: sub,
	}
}

type ChannelJSON struct {
	ChannelID         ChannelID `json:"channel_id"`
	OwnerUserID       UserID    `json:"owner_user_id"`
	InternalName      string    `json:"internal_name"`
	DisplayName       string    `json:"display_name"`
	DescriptionName   *string   `json:"description_name"`
	SubscribeKey      *string   `json:"subscribe_key"` // can be nil, depending on endpoint
	SendKey           *string   `json:"send_key"`      // can be nil, depending on endpoint
	TimestampCreated  string    `json:"timestamp_created"`
	TimestampLastSent *string   `json:"timestamp_lastsent"`
	MessagesSent      int       `json:"messages_sent"`
}

type ChannelWithSubscriptionJSON struct {
	ChannelJSON
	Subscription *SubscriptionJSON `json:"subscription"`
}

type ChannelDB struct {
	ChannelID         ChannelID `db:"channel_id"`
	OwnerUserID       UserID    `db:"owner_user_id"`
	InternalName      string    `db:"internal_name"`
	DisplayName       string    `db:"display_name"`
	DescriptionName   *string   `db:"description_name"`
	SubscribeKey      string    `db:"subscribe_key"`
	SendKey           string    `db:"send_key"`
	TimestampCreated  int64     `db:"timestamp_created"`
	TimestampLastRead *int64    `db:"timestamp_lastread"`
	TimestampLastSent *int64    `db:"timestamp_lastsent"`
	MessagesSent      int       `db:"messages_sent"`
}

func (c ChannelDB) Model() Channel {
	return Channel{
		ChannelID:         c.ChannelID,
		OwnerUserID:       c.OwnerUserID,
		InternalName:      c.InternalName,
		DisplayName:       c.DisplayName,
		DescriptionName:   c.DescriptionName,
		SubscribeKey:      c.SubscribeKey,
		SendKey:           c.SendKey,
		TimestampCreated:  time.UnixMilli(c.TimestampCreated),
		TimestampLastSent: timeOptFromMilli(c.TimestampLastSent),
		MessagesSent:      c.MessagesSent,
	}
}

type ChannelWithSubscriptionDB struct {
	ChannelDB
	Subscription *SubscriptionDB `db:"sub"`
}

func (c ChannelWithSubscriptionDB) Model() ChannelWithSubscription {
	var sub *Subscription = nil
	if c.Subscription != nil {
		sub = langext.Ptr(c.Subscription.Model())
	}
	return ChannelWithSubscription{
		Channel:      c.ChannelDB.Model(),
		Subscription: sub,
	}
}

func DecodeChannel(r *sqlx.Rows) (Channel, error) {
	data, err := sq.ScanSingle[ChannelDB](r, sq.SModeFast, sq.Safe, true)
	if err != nil {
		return Channel{}, err
	}
	return data.Model(), nil
}

func DecodeChannels(r *sqlx.Rows) ([]Channel, error) {
	data, err := sq.ScanAll[ChannelDB](r, sq.SModeFast, sq.Safe, true)
	if err != nil {
		return nil, err
	}
	return langext.ArrMap(data, func(v ChannelDB) Channel { return v.Model() }), nil
}

func DecodeChannelWithSubscription(r *sqlx.Rows) (ChannelWithSubscription, error) {
	data, err := sq.ScanSingle[ChannelWithSubscriptionDB](r, sq.SModeExtended, sq.Safe, true)
	if err != nil {
		return ChannelWithSubscription{}, err
	}
	return data.Model(), nil
}

func DecodeChannelsWithSubscription(r *sqlx.Rows) ([]ChannelWithSubscription, error) {
	data, err := sq.ScanAll[ChannelWithSubscriptionDB](r, sq.SModeExtended, sq.Safe, true)
	if err != nil {
		return nil, err
	}
	return langext.ArrMap(data, func(v ChannelWithSubscriptionDB) ChannelWithSubscription { return v.Model() }), nil
}
