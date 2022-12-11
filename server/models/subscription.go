package models

import (
	"github.com/jmoiron/sqlx"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"time"
)

type Subscription struct {
	SubscriptionID     SubscriptionID
	SubscriberUserID   UserID
	ChannelOwnerUserID UserID
	ChannelID          ChannelID
	ChannelName        string
	TimestampCreated   time.Time
	Confirmed          bool
}

func (s Subscription) JSON() SubscriptionJSON {
	return SubscriptionJSON{
		SubscriptionID:     s.SubscriptionID,
		SubscriberUserID:   s.SubscriberUserID,
		ChannelOwnerUserID: s.ChannelOwnerUserID,
		ChannelID:          s.ChannelID,
		ChannelName:        s.ChannelName,
		TimestampCreated:   s.TimestampCreated.Format(time.RFC3339Nano),
		Confirmed:          s.Confirmed,
	}
}

type SubscriptionJSON struct {
	SubscriptionID     SubscriptionID `json:"subscription_id"`
	SubscriberUserID   UserID         `json:"subscriber_user_id"`
	ChannelOwnerUserID UserID         `json:"channel_owner_user_id"`
	ChannelID          ChannelID      `json:"channel_id"`
	ChannelName        string         `json:"channel_name"`
	TimestampCreated   string         `json:"timestamp_created"`
	Confirmed          bool           `json:"confirmed"`
}

type SubscriptionDB struct {
	SubscriptionID     SubscriptionID `db:"subscription_id"`
	SubscriberUserID   UserID         `db:"subscriber_user_id"`
	ChannelOwnerUserID UserID         `db:"channel_owner_user_id"`
	ChannelID          ChannelID      `db:"channel_id"`
	ChannelName        string         `db:"channel_name"`
	TimestampCreated   int64          `db:"timestamp_created"`
	Confirmed          int            `db:"confirmed"`
}

func (s SubscriptionDB) Model() Subscription {
	return Subscription{
		SubscriptionID:     s.SubscriptionID,
		SubscriberUserID:   s.SubscriberUserID,
		ChannelOwnerUserID: s.ChannelOwnerUserID,
		ChannelID:          s.ChannelID,
		ChannelName:        s.ChannelName,
		TimestampCreated:   time.UnixMilli(s.TimestampCreated),
		Confirmed:          s.Confirmed != 0,
	}
}

func DecodeSubscription(r *sqlx.Rows) (Subscription, error) {
	data, err := sq.ScanSingle[SubscriptionDB](r, true)
	if err != nil {
		return Subscription{}, err
	}
	return data.Model(), nil
}

func DecodeSubscriptions(r *sqlx.Rows) ([]Subscription, error) {
	data, err := sq.ScanAll[SubscriptionDB](r, true)
	if err != nil {
		return nil, err
	}
	return langext.ArrMap(data, func(v SubscriptionDB) Subscription { return v.Model() }), nil
}
