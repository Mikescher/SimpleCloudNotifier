package models

import (
	"context"
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

func (c Channel) JSONPreview() ChannelPreviewJSON {
	return ChannelPreviewJSON{
		ChannelID:       c.ChannelID,
		OwnerUserID:     c.OwnerUserID,
		InternalName:    c.InternalName,
		DisplayName:     c.DisplayName,
		DescriptionName: c.DescriptionName,
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
	TimestampCreated  string    `json:"timestamp_created"`
	TimestampLastSent *string   `json:"timestamp_lastsent"`
	MessagesSent      int       `json:"messages_sent"`
}

type ChannelWithSubscriptionJSON struct {
	ChannelJSON
	Subscription *SubscriptionJSON `json:"subscription"`
}

type ChannelPreviewJSON struct {
	ChannelID       ChannelID `json:"channel_id"`
	OwnerUserID     UserID    `json:"owner_user_id"`
	InternalName    string    `json:"internal_name"`
	DisplayName     string    `json:"display_name"`
	DescriptionName *string   `json:"description_name"`
}

type ChannelDB struct {
	ChannelID         ChannelID `db:"channel_id"`
	OwnerUserID       UserID    `db:"owner_user_id"`
	InternalName      string    `db:"internal_name"`
	DisplayName       string    `db:"display_name"`
	DescriptionName   *string   `db:"description_name"`
	SubscribeKey      string    `db:"subscribe_key"`
	TimestampCreated  int64     `db:"timestamp_created"`
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
		TimestampCreated:  timeFromMilli(c.TimestampCreated),
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

func DecodeChannel(ctx context.Context, q sq.Queryable, r *sqlx.Rows) (Channel, error) {
	data, err := sq.ScanSingle[ChannelDB](ctx, q, r, sq.SModeFast, sq.Safe, true)
	if err != nil {
		return Channel{}, err
	}
	return data.Model(), nil
}

func DecodeChannels(ctx context.Context, q sq.Queryable, r *sqlx.Rows) ([]Channel, error) {
	data, err := sq.ScanAll[ChannelDB](ctx, q, r, sq.SModeFast, sq.Safe, true)
	if err != nil {
		return nil, err
	}
	return langext.ArrMap(data, func(v ChannelDB) Channel { return v.Model() }), nil
}

func DecodeChannelWithSubscription(ctx context.Context, q sq.Queryable, r *sqlx.Rows) (ChannelWithSubscription, error) {
	data, err := sq.ScanSingle[ChannelWithSubscriptionDB](ctx, q, r, sq.SModeExtended, sq.Safe, true)
	if err != nil {
		return ChannelWithSubscription{}, err
	}
	return data.Model(), nil
}

func DecodeChannelsWithSubscription(ctx context.Context, q sq.Queryable, r *sqlx.Rows) ([]ChannelWithSubscription, error) {
	data, err := sq.ScanAll[ChannelWithSubscriptionDB](ctx, q, r, sq.SModeExtended, sq.Safe, true)
	if err != nil {
		return nil, err
	}
	return langext.ArrMap(data, func(v ChannelWithSubscriptionDB) ChannelWithSubscription { return v.Model() }), nil
}
