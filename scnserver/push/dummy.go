package push

import (
	"blackforestbytes.com/simplecloudnotifier/models"
	"context"
	_ "embed"
)

type DummyConnector struct{}

func NewDummy() NotificationClient {
	return &DummyConnector{}
}

func (d DummyConnector) SendNotification(ctx context.Context, client models.Client, msg models.Message, compatTitleOverride *string, compatMsgIDOverride *string) (string, error) {
	return "%DUMMY%", nil
}
