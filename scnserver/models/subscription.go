package models

import (
	"context"
	"github.com/jmoiron/sqlx"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"time"
)

// [!] subscriptions are read-access to channels,
//
// The set of subscriptions specifies which messages the ListMessages() API call returns
// also single messages/channels that are subscribed can be queries
//
// (use keytokens for write-access)

type Subscription struct {
	SubscriptionID      SubscriptionID
	SubscriberUserID    UserID
	ChannelOwnerUserID  UserID
	ChannelID           ChannelID
	ChannelInternalName string
	TimestampCreated    time.Time
	Confirmed           bool
}

func (s Subscription) JSON() SubscriptionJSON {
	return SubscriptionJSON{
		SubscriptionID:      s.SubscriptionID,
		SubscriberUserID:    s.SubscriberUserID,
		ChannelOwnerUserID:  s.ChannelOwnerUserID,
		ChannelID:           s.ChannelID,
		ChannelInternalName: s.ChannelInternalName,
		TimestampCreated:    s.TimestampCreated.Format(time.RFC3339Nano),
		Confirmed:           s.Confirmed,
	}
}

type SubscriptionJSON struct {
	SubscriptionID      SubscriptionID `json:"subscription_id"`
	SubscriberUserID    UserID         `json:"subscriber_user_id"`
	ChannelOwnerUserID  UserID         `json:"channel_owner_user_id"`
	ChannelID           ChannelID      `json:"channel_id"`
	ChannelInternalName string         `json:"channel_internal_name"`
	TimestampCreated    string         `json:"timestamp_created"`
	Confirmed           bool           `json:"confirmed"`
}

type SubscriptionDB struct {
	SubscriptionID      SubscriptionID `db:"subscription_id"`
	SubscriberUserID    UserID         `db:"subscriber_user_id"`
	ChannelOwnerUserID  UserID         `db:"channel_owner_user_id"`
	ChannelID           ChannelID      `db:"channel_id"`
	ChannelInternalName string         `db:"channel_internal_name"`
	TimestampCreated    int64          `db:"timestamp_created"`
	Confirmed           int            `db:"confirmed"`
}

func (s SubscriptionDB) Model() Subscription {
	return Subscription{
		SubscriptionID:      s.SubscriptionID,
		SubscriberUserID:    s.SubscriberUserID,
		ChannelOwnerUserID:  s.ChannelOwnerUserID,
		ChannelID:           s.ChannelID,
		ChannelInternalName: s.ChannelInternalName,
		TimestampCreated:    timeFromMilli(s.TimestampCreated),
		Confirmed:           s.Confirmed != 0,
	}
}

func DecodeSubscription(ctx context.Context, q sq.Queryable, r *sqlx.Rows) (Subscription, error) {
	data, err := sq.ScanSingle[SubscriptionDB](ctx, q, r, sq.SModeFast, sq.Safe, true)
	if err != nil {
		return Subscription{}, err
	}
	return data.Model(), nil
}

func DecodeSubscriptions(ctx context.Context, q sq.Queryable, r *sqlx.Rows) ([]Subscription, error) {
	data, err := sq.ScanAll[SubscriptionDB](ctx, q, r, sq.SModeFast, sq.Safe, true)
	if err != nil {
		return nil, err
	}
	return langext.ArrMap(data, func(v SubscriptionDB) Subscription { return v.Model() }), nil
}
