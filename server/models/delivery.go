package models

import (
	"database/sql"
	"github.com/blockloop/scan"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"time"
)

type DeliveryStatus string

const (
	DeliveryStatusRetry   DeliveryStatus = "RETRY"
	DeliveryStatusSuccess DeliveryStatus = "SUCCESS"
	DeliveryStatusFailed  DeliveryStatus = "FAILED"
)

type Delivery struct {
	DeliveryID         int64
	SCNMessageID       int64
	ReceiverUserID     int64
	ReceiverClientID   int64
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
		SCNMessageID:       d.SCNMessageID,
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
	DeliveryID         int64          `json:"delivery_id"`
	SCNMessageID       int64          `json:"scn_message_id"`
	ReceiverUserID     int64          `json:"receiver_user_id"`
	ReceiverClientID   int64          `json:"receiver_client_id"`
	TimestampCreated   string         `json:"timestamp_created"`
	TimestampFinalized *string        `json:"tiestamp_finalized"`
	Status             DeliveryStatus `json:"status"`
	RetryCount         int            `json:"retry_count"`
	NextDelivery       *string        `json:"next_delivery"`
	FCMMessageID       *string        `json:"fcm_message_id"`
}

type DeliveryDB struct {
	DeliveryID         int64          `db:"delivery_id"`
	SCNMessageID       int64          `db:"scn_message_id"`
	ReceiverUserID     int64          `db:"receiver_user_id"`
	ReceiverClientID   int64          `db:"receiver_client_id"`
	TimestampCreated   int64          `db:"timestamp_created"`
	TimestampFinalized *int64         `db:"tiestamp_finalized"`
	Status             DeliveryStatus `db:"status"`
	RetryCount         int            `db:"retry_count"`
	NextDelivery       *int64         `db:"next_delivery"`
	FCMMessageID       *string        `db:"fcm_message_id"`
}

func (d DeliveryDB) Model() Delivery {
	return Delivery{
		DeliveryID:         d.DeliveryID,
		SCNMessageID:       d.SCNMessageID,
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

func DecodeDelivery(r *sql.Rows) (Delivery, error) {
	var data DeliveryDB
	err := scan.RowStrict(&data, r)
	if err != nil {
		return Delivery{}, err
	}
	return data.Model(), nil
}

func DecodeDeliveries(r *sql.Rows) ([]Delivery, error) {
	var data []DeliveryDB
	err := scan.RowsStrict(&data, r)
	if err != nil {
		return nil, err
	}
	return langext.ArrMap(data, func(v DeliveryDB) Delivery { return v.Model() }), nil
}
