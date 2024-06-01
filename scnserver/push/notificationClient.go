package push

import (
	"blackforestbytes.com/simplecloudnotifier/models"
	"context"
)

type NotificationClient interface {
	SendNotification(ctx context.Context, user models.User, client models.Client, channel models.Channel, msg models.Message) (string, error)
}
