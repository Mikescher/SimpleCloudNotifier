package models

import (
	"database/sql"
	"github.com/blockloop/scan"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"time"
)

type Subscription struct {
	SubscriptionID     int64
	SubscriberUserID   int64
	ChannelOwnerUserID int64
	ChannelID          int64
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
	SubscriptionID     int64  `json:"subscription_id"`
	SubscriberUserID   int64  `json:"subscriber_user_id"`
	ChannelOwnerUserID int64  `json:"channel_owner_user_id"`
	ChannelID          int64  `json:"channel_id"`
	ChannelName        string `json:"channel_name"`
	TimestampCreated   string `json:"timestamp_created"`
	Confirmed          bool   `json:"confirmed"`
}

type SubscriptionDB struct {
	SubscriptionID     int64  `db:"subscription_id"`
	SubscriberUserID   int64  `db:"subscriber_user_id"`
	ChannelOwnerUserID int64  `db:"channel_owner_user_id"`
	ChannelID          int64  `db:"channel_id"`
	ChannelName        string `db:"channel_name"`
	TimestampCreated   int64  `db:"timestamp_created"`
	Confirmed          int    `db:"confirmed"`
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

func DecodeSubscription(r *sql.Rows) (Subscription, error) {
	var data SubscriptionDB
	err := scan.RowStrict(&data, r)
	if err != nil {
		return Subscription{}, err
	}
	return data.Model(), nil
}

func DecodeSubscriptions(r *sql.Rows) ([]Subscription, error) {
	var data []SubscriptionDB
	err := scan.RowsStrict(&data, r)
	if err != nil {
		return nil, err
	}
	return langext.ArrMap(data, func(v SubscriptionDB) Subscription { return v.Model() }), nil
}
