package models

import "gogs.mikescher.com/BlackForestBytes/goext/rfctime"

type DeliveryStatus string //@enum:type

const (
	DeliveryStatusRetry   DeliveryStatus = "RETRY"
	DeliveryStatusSuccess DeliveryStatus = "SUCCESS"
	DeliveryStatusFailed  DeliveryStatus = "FAILED"
)

type Delivery struct {
	DeliveryID         DeliveryID               `db:"delivery_id"         json:"delivery_id"`
	MessageID          MessageID                `db:"message_id"          json:"message_id"`
	ReceiverUserID     UserID                   `db:"receiver_user_id"    json:"receiver_user_id"`
	ReceiverClientID   ClientID                 `db:"receiver_client_id"  json:"receiver_client_id"`
	TimestampCreated   SCNTime                  `db:"timestamp_created"   json:"timestamp_created"`
	TimestampFinalized *SCNTime                 `db:"timestamp_finalized" json:"timestamp_finalized"`
	Status             DeliveryStatus           `db:"status"              json:"status"`
	RetryCount         int                      `db:"retry_count"         json:"retry_count"`
	NextDelivery       *rfctime.RFC3339NanoTime `db:"next_delivery"       json:"next_delivery"`
	FCMMessageID       *string                  `db:"fcm_message_id"      json:"fcm_message_id"`
}

func (d Delivery) MaxRetryCount() int {
	return 5
}
