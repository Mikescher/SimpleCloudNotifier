package models

type Channel struct {
	ChannelID         ChannelID `db:"channel_id"         json:"channel_id"`
	OwnerUserID       UserID    `db:"owner_user_id"      json:"owner_user_id"`
	InternalName      string    `db:"internal_name"      json:"internal_name"`
	DisplayName       string    `db:"display_name"       json:"display_name"`
	DescriptionName   *string   `db:"description_name"   json:"description_name"`
	SubscribeKey      string    `db:"subscribe_key"      json:"subscribe_key"      jsonfilter:"INCLUDE_KEY"` // can be nil, depending on endpoint
	TimestampCreated  SCNTime   `db:"timestamp_created"  json:"timestamp_created"`
	TimestampLastSent *SCNTime  `db:"timestamp_lastsent" json:"timestamp_lastsent"`
	MessagesSent      int       `db:"messages_sent"      json:"messages_sent"`
}

type ChannelWithSubscription struct {
	Channel
	Subscription *Subscription `db:"sub" json:"subscription"`
}

type ChannelPreview struct {
	ChannelID       ChannelID `json:"channel_id"`
	OwnerUserID     UserID    `json:"owner_user_id"`
	InternalName    string    `json:"internal_name"`
	DisplayName     string    `json:"display_name"`
	DescriptionName *string   `json:"description_name"`
}

func (c Channel) WithSubscription(sub *Subscription) ChannelWithSubscription {
	return ChannelWithSubscription{
		Channel:      c,
		Subscription: sub,
	}
}

func (c Channel) Preview() ChannelPreview {
	return ChannelPreview{
		ChannelID:       c.ChannelID,
		OwnerUserID:     c.OwnerUserID,
		InternalName:    c.InternalName,
		DisplayName:     c.DisplayName,
		DescriptionName: c.DescriptionName,
	}
}
