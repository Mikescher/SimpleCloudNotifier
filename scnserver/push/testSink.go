package push

import (
	"blackforestbytes.com/simplecloudnotifier/models"
	"context"
	_ "embed"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
)

type SinkData struct {
	Message models.Message
	Client  models.Client
}

type TestSink struct {
	Data []SinkData
}

func NewTestSink() NotificationClient {
	return &TestSink{}
}

func (d *TestSink) Last() SinkData {
	return d.Data[len(d.Data)-1]
}

func (d *TestSink) SendNotification(ctx context.Context, client models.Client, msg models.Message) (string, error) {
	id, err := langext.NewHexUUID()
	if err != nil {
		return "", err
	}

	key := "TestSink[" + id + "]"

	d.Data = append(d.Data, SinkData{
		Message: msg,
		Client:  client,
	})

	return key, nil
}
