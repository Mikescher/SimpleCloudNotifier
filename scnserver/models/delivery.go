package models

import (
	"github.com/jmoiron/sqlx"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"time"
)

type DeliveryStatus string

const (
	DeliveryStatusRetry   DeliveryStatus = "RETRY"
	DeliveryStatusSuccess DeliveryStatus = "SUCCESS"
	DeliveryStatusFailed  DeliveryStatus = "FAILED"
)

type Delivery struct {
	DeliveryID         DeliveryID
	MessageID          MessageID
	ReceiverUserID     UserID
	ReceiverClientID   ClientID
	TimestampCreated   time.Time
	TimestampFinalized *time.Time
	Status             DeliveryStatus
	RetryCount         int
	NextDelivery       *time.Time
	FCMMessageID       *string
}

func (d Delivery) JSON() DeliveryJSON {
	return DeliveryJSON{
		DeliveryID:         d.DeliveryID,
		MessageID:          d.MessageID,
		ReceiverUserID:     d.ReceiverUserID,
		ReceiverClientID:   d.ReceiverClientID,
		TimestampCreated:   d.TimestampCreated.Format(time.RFC3339Nano),
		TimestampFinalized: timeOptFmt(d.TimestampFinalized, time.RFC3339Nano),
		Status:             d.Status,
		RetryCount:         d.RetryCount,
		NextDelivery:       timeOptFmt(d.NextDelivery, time.RFC3339Nano),
		FCMMessageID:       d.FCMMessageID,
	}
}

func (d Delivery) MaxRetryCount() int {
	return 5
}

type DeliveryJSON struct {
	DeliveryID         DeliveryID     `json:"delivery_id"`
	MessageID          MessageID      `json:"message_id"`
	ReceiverUserID     UserID         `json:"receiver_user_id"`
	ReceiverClientID   ClientID       `json:"receiver_client_id"`
	TimestampCreated   string         `json:"timestamp_created"`
	TimestampFinalized *string        `json:"timestamp_finalized"`
	Status             DeliveryStatus `json:"status"`
	RetryCount         int            `json:"retry_count"`
	NextDelivery       *string        `json:"next_delivery"`
	FCMMessageID       *string        `json:"fcm_message_id"`
}

type DeliveryDB struct {
	DeliveryID         DeliveryID     `db:"delivery_id"`
	MessageID          MessageID      `db:"message_id"`
	ReceiverUserID     UserID         `db:"receiver_user_id"`
	ReceiverClientID   ClientID       `db:"receiver_client_id"`
	TimestampCreated   int64          `db:"timestamp_created"`
	TimestampFinalized *int64         `db:"timestamp_finalized"`
	Status             DeliveryStatus `db:"status"`
	RetryCount         int            `db:"retry_count"`
	NextDelivery       *int64         `db:"next_delivery"`
	FCMMessageID       *string        `db:"fcm_message_id"`
}

func (d DeliveryDB) Model() Delivery {
	return Delivery{
		DeliveryID:         d.DeliveryID,
		MessageID:          d.MessageID,
		ReceiverUserID:     d.ReceiverUserID,
		ReceiverClientID:   d.ReceiverClientID,
		TimestampCreated:   time.UnixMilli(d.TimestampCreated),
		TimestampFinalized: timeOptFromMilli(d.TimestampFinalized),
		Status:             d.Status,
		RetryCount:         d.RetryCount,
		NextDelivery:       timeOptFromMilli(d.NextDelivery),
		FCMMessageID:       d.FCMMessageID,
	}
}

func DecodeDelivery(r *sqlx.Rows) (Delivery, error) {
	data, err := sq.ScanSingle[DeliveryDB](r, sq.SModeFast, sq.Safe, true)
	if err != nil {
		return Delivery{}, err
	}
	return data.Model(), nil
}

func DecodeDeliveries(r *sqlx.Rows) ([]Delivery, error) {
	data, err := sq.ScanAll[DeliveryDB](r, sq.SModeFast, sq.Safe, true)
	if err != nil {
		return nil, err
	}
	return langext.ArrMap(data, func(v DeliveryDB) Delivery { return v.Model() }), nil
}
