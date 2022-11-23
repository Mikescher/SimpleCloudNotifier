package push

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/models"
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"io"
	"net/http"
	"strconv"
	"time"
)

// https://firebase.google.com/docs/cloud-messaging/send-message#rest
// https://firebase.google.com/docs/cloud-messaging/auth-server

type FirebaseConnector struct {
	fbProject string
	client    http.Client
	auth      *FirebaseOAuth2
}

func NewFirebaseConn(conf scn.Config) (NotificationClient, error) {

	fbauth, err := NewAuth(conf.FirebaseTokenURI, conf.FirebaseProjectID, conf.FirebaseClientMail, conf.FirebasePrivateKey)
	if err != nil {
		return nil, err
	}

	return &FirebaseConnector{
		fbProject: conf.FirebaseProjectID,
		client:    http.Client{Timeout: 5 * time.Second},
		auth:      fbauth,
	}, nil
}

type Notification struct {
	Id       string
	Token    string
	Platform string
	Title    string
	Body     string
	Priority int
}

func (fb FirebaseConnector) SendNotification(ctx context.Context, client models.Client, msg models.Message) (string, error) {

	uri := "https://fcm.googleapis.com/v1/projects/" + fb.fbProject + "/messages:send"

	jsonBody := gin.H{
		"data": gin.H{
			"scn_msg_id": msg.SCNMessageID.String(),
			"usr_msg_id": langext.Coalesce(msg.UserMessageID, ""),
			"client_id":  client.ClientID.String(),
			"timestamp":  strconv.FormatInt(msg.Timestamp().Unix(), 10),
			"priority":   strconv.Itoa(msg.Priority),
			"trimmed":    langext.Conditional(msg.NeedsTrim(), "true", "false"),
			"title":      msg.Title,
			"body":       langext.Coalesce(msg.TrimmedContent(), ""),
		},
		"token": *client.FCMToken,
		"android": gin.H{
			"priority": "high",
		},
		"apns": gin.H{},
	}
	if client.Type == models.ClientTypeIOS {
		jsonBody["notification"] = gin.H{
			"title": msg.Title,
			"body":  msg.ShortContent(),
		}
	}

	bytesBody, err := json.Marshal(gin.H{"message": jsonBody})
	if err != nil {
		return "", err
	}

	request, err := http.NewRequestWithContext(ctx, "POST", uri, bytes.NewBuffer(bytesBody))
	if err != nil {
		return "", err
	}

	tok, err := fb.auth.Token(ctx)
	if err != nil {
		log.Err(err).Msg("Refreshing FB token failed")
		return "", err
	}

	request.Header.Set("Authorization", "Bearer "+tok)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := fb.client.Do(request)
	if err != nil {
		return "", err
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		if bstr, err := io.ReadAll(response.Body); err == nil {
			return "", errors.New(fmt.Sprintf("FCM-Request returned %d: %s", response.StatusCode, string(bstr)))
		} else {
			return "", errors.New(fmt.Sprintf("FCM-Request returned %d", response.StatusCode))
		}
	}

	respBodyBin, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var respBody struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(respBodyBin, &respBody); err != nil {
		return "", err
	}

	return respBody.Name, nil
}
