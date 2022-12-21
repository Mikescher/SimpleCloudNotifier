package push

import (
	"blackforestbytes.com/simplecloudnotifier/models"
	"context"
)

type NotificationClient interface {
	SendNotification(ctx context.Context, client models.Client, msg models.Message) (string, error)
}
