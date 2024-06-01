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

func (d DummyConnector) SendNotification(ctx context.Context, user models.User, client models.Client, channel models.Channel, msg models.Message) (string, error) {
	return "%DUMMY%", nil
}
