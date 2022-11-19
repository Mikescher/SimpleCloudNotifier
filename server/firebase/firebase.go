package firebase

import (
	"context"
	_ "embed"
	fb "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
)

//go:embed scnserviceaccountkey.json
var scnserviceaccountkey []byte

type App struct {
	app       *fb.App
	messaging *messaging.Client
}

func NewFirebaseApp() App {
	opt := option.WithCredentialsJSON(scnserviceaccountkey)
	app, err := fb.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Error().Err(err).Msg("failed to init firebase app")
	}
	msg, err := app.Messaging(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("failed to init messaging client")
	}
	log.Info().Msg("Initialized Firebase")
	return App{
		app:       app,
		messaging: msg,
	}
}

type Notification struct {
	Id       string
	Token    string
	Platform string
	Title    string
	Body     string
	Priority int
}

func (fb App) SendNotification(ctx context.Context, notification Notification) (string, error) {
	n := messaging.Message{
		Data: map[string]string{"scn_msg_id": notification.Id},
		Notification: &messaging.Notification{
			Title: notification.Title,
			Body:  notification.Body,
		},
		Android:    nil,
		APNS:       nil,
		Webpush:    nil,
		FCMOptions: nil,
		Token:      notification.Token,
		Topic:      "",
		Condition:  "",
	}
	if notification.Platform == "ios" {
		n.APNS = nil
	}

	if notification.Platform == "android" {
		n.Android = nil
	}

	res, err := fb.messaging.Send(ctx, &n)
	if err != nil {
		log.Error().Err(err).Msg("failed to send push")
	}
	return res, err
}
