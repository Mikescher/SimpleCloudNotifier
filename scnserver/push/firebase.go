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
	"strings"
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

	pkey := strings.ReplaceAll(conf.FirebasePrivateKey, "\\n", "\n")

	fbauth, err := NewAuth(conf.FirebaseTokenURI, conf.FirebaseProjectID, conf.FirebaseClientMail, pkey)
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

func (fb FirebaseConnector) SendNotification(ctx context.Context, user models.User, client models.Client, channel models.Channel, msg models.Message) (string, string, error) {

	uri := "https://fcm.googleapis.com/v1/projects/" + fb.fbProject + "/messages:send"

	jsonBody := gin.H{}

	if client.Type == models.ClientTypeIOS {
		jsonBody = gin.H{
			"token": client.FCMToken,
			"notification": gin.H{
				"title": msg.Title,
				"body":  msg.ShortContent(),
			},
			"apns": gin.H{},
		}
	} else if client.Type == models.ClientTypeAndroid {
		jsonBody = gin.H{
			"token": client.FCMToken,
			"android": gin.H{
				"priority":    "high",
				"fcm_options": gin.H{},
			},
			"data": gin.H{
				"scn_msg_id": msg.MessageID.String(),
				"usr_msg_id": langext.Coalesce(msg.UserMessageID, ""),
				"client_id":  client.ClientID.String(),
				"timestamp":  strconv.FormatInt(msg.Timestamp().Unix(), 10),
				"priority":   strconv.Itoa(msg.Priority),
				"trimmed":    langext.Conditional(msg.NeedsTrim(), "true", "false"),
				"title":      msg.Title,
				"channel":    channel.DisplayName,
				"channel_id": channel.ChannelID,
				"body":       langext.Coalesce(msg.TrimmedContent(), ""),
			},
		}
	} else {
		jsonBody = gin.H{
			"token": client.FCMToken,
			"notification": gin.H{
				"title": msg.FormatNotificationTitle(user, channel),
				"body":  msg.ShortContent(),
			},
		}
	}

	bytesBody, err := json.Marshal(gin.H{"message": jsonBody})
	if err != nil {
		return "", "", err
	}

	request, err := http.NewRequestWithContext(ctx, "POST", uri, bytes.NewBuffer(bytesBody))
	if err != nil {
		return "", "", err
	}

	tok, err := fb.auth.Token(ctx)
	if err != nil {
		log.Err(err).Msg("Refreshing FB token failed")
		return "", "", err
	}

	request.Header.Set("Authorization", "Bearer "+tok)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := fb.client.Do(request)
	if err != nil {
		return "", "", err
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		if bstr, err := io.ReadAll(response.Body); err == nil {

			var errRespBody struct {
				Error struct {
					Code    int    `json:"code"`
					Message string `json:"message"`
					Status  string `json:"status"`
					Details []struct {
						AtType string `json:"@type"`
						ECode  string `json:"errorCode"`
					} `json:"details"`
				} `json:"error"`
			}

			if err := json.Unmarshal(bstr, &errRespBody); err == nil {
				for _, v := range errRespBody.Error.Details {
					return "", v.ECode, errors.New(fmt.Sprintf("FCM-Request returned %d [UNREGISTERED]: %s", response.StatusCode, string(bstr)))
				}
			}

			return "", "", errors.New(fmt.Sprintf("FCM-Request returned %d: %s", response.StatusCode, string(bstr)))

		} else {

			return "", "", errors.New(fmt.Sprintf("FCM-Request returned %d", response.StatusCode))

		}
	}

	respBodyBin, err := io.ReadAll(response.Body)
	if err != nil {
		return "", "", err
	}

	var respBody struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(respBodyBin, &respBody); err != nil {
		return "", "", err
	}

	log.Info().Msg(fmt.Sprintf("Sucessfully pushed notification %s", msg.MessageID))

	return respBody.Name, "", nil
}
