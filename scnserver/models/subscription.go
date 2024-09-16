package models

// [!] subscriptions are read-access to channels,
//
// The set of subscriptions specifies which messages the ListMessages() API call returns
// also single messages/channels that are subscribed can be queries
//
// (use keytokens for write-access)

type Subscription struct {
	SubscriptionID      SubscriptionID `db:"subscription_id"          json:"subscription_id"`
	SubscriberUserID    UserID         `db:"subscriber_user_id"       json:"subscriber_user_id"`
	ChannelOwnerUserID  UserID         `db:"channel_owner_user_id"    json:"channel_owner_user_id"`
	ChannelID           ChannelID      `db:"channel_id"               json:"channel_id"`
	ChannelInternalName string         `db:"channel_internal_name"    json:"channel_internal_name"`
	TimestampCreated    SCNTime        `db:"timestamp_created"        json:"timestamp_created"`
	Confirmed           bool           `db:"confirmed"                json:"confirmed"`
}
