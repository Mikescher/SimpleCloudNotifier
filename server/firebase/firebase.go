package firebase

import (
	"blackforestbytes.com/simplecloudnotifier/models"
	"context"
	_ "embed"
	fb "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"google.golang.org/api/option"
	"strconv"
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

func (fb App) SendNotification(ctx context.Context, client models.Client, msg models.Message) (string, error) {

	n := messaging.Message{
		Data: map[string]string{
			"scn_msg_id": strconv.FormatInt(msg.SCNMessageID, 10),
			"usr_msg_id": langext.Coalesce(msg.UserMessageID, ""),
			"client_id":  strconv.FormatInt(client.ClientID, 10),
			"timestamp":  strconv.FormatInt(msg.Timestamp().Unix(), 10),
			"priority":   strconv.Itoa(msg.Priority),
			"trimmed":    langext.Conditional(msg.NeedsTrim(), "true", "false"),
			"title":      msg.Title,
			"body":       langext.Coalesce(msg.TrimmedContent(), ""),
		},
		Notification: &messaging.Notification{
			Title: msg.Title,
			Body:  msg.ShortContent(),
		},
		Android:    nil,
		APNS:       nil,
		Webpush:    nil,
		FCMOptions: nil,
		Token:      *client.FCMToken,
		Topic:      "",
		Condition:  "",
	}

	if client.Type == models.ClientTypeIOS {
		n.APNS = nil
	} else if client.Type == models.ClientTypeAndroid {
		n.Notification = nil
		n.Android = &messaging.AndroidConfig{Priority: "high"}
	}

	res, err := fb.messaging.Send(ctx, &n)
	if err != nil {
		log.Error().Err(err).Msg("failed to send push")
	}
	return res, err
}
