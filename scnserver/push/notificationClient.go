package push

import (
	"blackforestbytes.com/simplecloudnotifier/models"
	"context"
)

type NotificationClient interface {
	SendNotification(ctx context.Context, client models.Client, msg models.Message, compatTitleOverride *string, compatMsgIDOverride *string) (string, error)
}
