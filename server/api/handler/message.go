package handler

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/common/ginresp"
	"blackforestbytes.com/simplecloudnotifier/db"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"gogs.mikescher.com/BlackForestBytes/goext/dataext"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/mathext"
	"gogs.mikescher.com/BlackForestBytes/goext/timeext"
	"net/http"
	"strings"
	"time"
)

type MessageHandler struct {
	app      *logic.Application
	database *db.Database
}

func NewMessageHandler(app *logic.Application) MessageHandler {
	return MessageHandler{
		app:      app,
		database: app.Database,
	}
}

// SendMessage swaggerdoc
//
// @Summary     Send a new message
// @Description All parameter can be set via query-parameter or the json body. Only UserID, UserKey and Title are required
//
// @Param       query_data query    handler.SendMessage.query false " "
// @Param       post_body  body     handler.SendMessage.body  false " "
//
// @Success     200        {object} handler.SendMessage.response
// @Failure     400        {object} ginresp.apiError
// @Failure     401        {object} ginresp.apiError
// @Failure     403        {object} ginresp.apiError
// @Failure     404        {object} ginresp.apiError
// @Failure     500        {object} ginresp.apiError
//
// @Router      /     [POST]
// @Router      /send [POST]
func (h MessageHandler) SendMessage(g *gin.Context) ginresp.HTTPResponse {
	type query struct {
		UserID        *int64   `form:"user_id"`
		UserKey       *string  `form:"user_key"`
		Channel       *string  `form:"channel"`
		ChanKey       *string  `form:"chan_key"`
		Title         *string  `form:"message_title"`
		Content       *string  `form:"message_content"`
		Priority      *int     `form:"priority"`
		UserMessageID *string  `form:"msg_id"`
		SendTimestamp *float64 `form:"timestamp"`
	}
	type body struct {
		UserID        *int64   `json:"user_id"`
		UserKey       *string  `json:"user_key"`
		Channel       *string  `json:"channel"`
		ChanKey       *string  `form:"chan_key"`
		Title         *string  `json:"message_title"`
		Content       *string  `json:"message_content"`
		Priority      *int     `json:"priority"`
		UserMessageID *string  `json:"msg_id"`
		SendTimestamp *float64 `json:"timestamp"`
	}
	type response struct {
		Success        bool            `json:"success"`
		ErrorID        apierr.APIError `json:"error"`
		ErrorHighlight int             `json:"errhighlight"`
		Message        string          `json:"message"`
		SuppressSend   bool            `json:"suppress_send"`
		MessageCount   int             `json:"messagecount"`
		Quota          int             `json:"quota"`
		IsPro          bool            `json:"is_pro"`
		QuotaMax       int             `json:"quota_max"`
		SCNMessageID   int64           `json:"scn_msg_id"`
	}

	var b body
	var q query
	ctx, errResp := h.app.StartRequest(g, nil, &q, &b)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	data := dataext.ObjectMerge(b, q)

	if data.UserID == nil {
		return ginresp.SendAPIError(400, apierr.MISSING_UID, 101, "Missing parameter [[user_id]]")
	}
	if data.UserKey == nil {
		return ginresp.SendAPIError(400, apierr.MISSING_UID, 102, "Missing parameter [[user_token]]")
	}
	if data.Title == nil {
		return ginresp.SendAPIError(400, apierr.MISSING_UID, 103, "Missing parameter [[title]]")
	}
	if data.SendTimestamp != nil && mathext.Abs(*data.SendTimestamp-float64(time.Now().Unix())) > (24*time.Hour).Seconds() {
		return ginresp.SendAPIError(400, apierr.TIMESTAMP_OUT_OF_RANGE, -1, "The timestamp mus be within 24 hours of now()")
	}
	if data.Priority != nil && (*data.Priority != 0 && *data.Priority != 1 && *data.Priority != 2) {
		return ginresp.SendAPIError(400, apierr.INVALID_PRIO, 105, "Invalid priority")
	}
	if len(strings.TrimSpace(*data.Title)) == 0 {
		return ginresp.SendAPIError(400, apierr.NO_TITLE, 103, "No title specified")
	}
	if data.UserMessageID != nil && len(strings.TrimSpace(*data.UserMessageID)) > 64 {
		return ginresp.SendAPIError(400, apierr.USR_MSG_ID_TOO_LONG, -1, "MessageID too long (64 characters)")
	}

	channelName := "main"
	if data.Channel != nil {
		channelName = strings.ToLower(strings.TrimSpace(*data.Channel))
	}

	user, err := h.database.GetUser(ctx, *data.UserID)
	if err == sql.ErrNoRows {
		return ginresp.SendAPIError(400, apierr.USER_NOT_FOUND, -1, "User not found")
	}
	if err != nil {
		return ginresp.SendAPIError(500, apierr.DATABASE_ERROR, -1, "Failed to query user")
	}

	if len(strings.TrimSpace(*data.Title)) > 120 {
		return ginresp.SendAPIError(400, apierr.TITLE_TOO_LONG, 103, "Title too long (120 characters)")
	}
	if data.Content != nil && len(strings.TrimSpace(*data.Content)) > user.MaxContentLength() {
		return ginresp.SendAPIError(400, apierr.CONTENT_TOO_LONG, 104, fmt.Sprintf("Content too long (%d characters; max := %d characters)", len(strings.TrimSpace(*data.Content)), user.MaxContentLength()))
	}

	if data.UserMessageID != nil {
		msg, err := h.database.GetMessageByUserMessageID(ctx, *data.UserMessageID)
		if err != nil {
			return ginresp.SendAPIError(500, apierr.DATABASE_ERROR, -1, "Failed to query existing message")
		}
		if msg != nil {
			return ginresp.JSON(http.StatusOK, response{
				Success:        true,
				ErrorID:        apierr.NO_ERROR,
				ErrorHighlight: -1,
				Message:        "Message already sent",
				SuppressSend:   true,
				MessageCount:   user.MessagesSent,
				Quota:          user.QuotaUsedToday(),
				IsPro:          user.IsPro,
				QuotaMax:       user.QuotaPerDay(),
				SCNMessageID:   msg.SCNMessageID,
			})
		}
	}

	if user.QuotaRemainingToday() <= 0 {
		return ginresp.SendAPIError(403, apierr.QUOTA_REACHED, -1, fmt.Sprintf("Daily quota reached (%d)", user.QuotaPerDay()))
	}

	channel, err := h.app.GetOrCreateChannel(ctx, *data.UserID, channelName)
	if err != nil {
		return ginresp.SendAPIError(500, apierr.DATABASE_ERROR, -1, "Failed to query/create channel")
	}

	selfChanAdmin := *data.UserID == channel.OwnerUserID && *data.UserKey == user.AdminKey
	selfChanSend := *data.UserID == channel.OwnerUserID && *data.UserKey == user.SendKey
	forgChanSend := *data.UserID != channel.OwnerUserID && data.ChanKey != nil && *data.ChanKey == channel.SendKey

	if !selfChanAdmin && !selfChanSend && !forgChanSend {
		return ginresp.SendAPIError(401, apierr.USER_AUTH_FAILED, 102, fmt.Sprintf("Daily quota reached (%d)", user.QuotaPerDay()))
	}

	var sendTimestamp *time.Time = nil
	if data.SendTimestamp != nil {
		sendTimestamp = langext.Ptr(timeext.UnixFloatSeconds(*data.SendTimestamp))
	}

	priority := langext.Coalesce(data.Priority, 1)

	msg, err := h.database.CreateMessage(ctx, *data.UserID, channel, sendTimestamp, *data.Title, data.Content, priority, data.UserMessageID)
	if err != nil {
		return ginresp.SendAPIError(500, apierr.DATABASE_ERROR, -1, "Failed to create message in db")
	}

	subscriptions, err := h.database.ListSubscriptionsByChannel(ctx, channel.ChannelID)
	if err != nil {
		return ginresp.SendAPIError(500, apierr.DATABASE_ERROR, -1, "Failed to query subscriptions")
	}

	err = h.database.IncUserMessageCounter(ctx, user)
	if err != nil {
		return ginresp.SendAPIError(500, apierr.DATABASE_ERROR, -1, "Failed to inc user msg-counter")
	}

	err = h.database.IncChannelMessageCounter(ctx, channel)
	if err != nil {
		return ginresp.SendAPIError(500, apierr.DATABASE_ERROR, -1, "Failed to channel msg-counter")
	}

	for _, sub := range subscriptions {
		clients, err := h.database.ListClients(ctx, sub.SubscriberUserID)
		if err != nil {
			return ginresp.SendAPIError(500, apierr.DATABASE_ERROR, -1, "Failed to query clients")
		}

		if !sub.Confirmed {
			continue
		}

		for _, client := range clients {

			fcmDelivID, err := h.deliverMessage(ctx, client, msg)
			if err != nil {
				_, err = h.database.CreateRetryDelivery(ctx, client, msg)
				if err != nil {
					return ginresp.SendAPIError(500, apierr.DATABASE_ERROR, -1, "Failed to create delivery")
				}
			} else {
				_, err = h.database.CreateSuccessDelivery(ctx, client, msg, *fcmDelivID)
				if err != nil {
					return ginresp.SendAPIError(500, apierr.DATABASE_ERROR, -1, "Failed to create delivery")
				}
			}

		}
	}

	return ginresp.JSON(http.StatusOK, response{
		Success:        true,
		ErrorID:        apierr.NO_ERROR,
		ErrorHighlight: -1,
		Message:        "Message sent",
		SuppressSend:   false,
		MessageCount:   user.MessagesSent + 1,
		Quota:          user.QuotaUsedToday() + 1,
		IsPro:          user.IsPro,
		QuotaMax:       user.QuotaPerDay(),
		SCNMessageID:   msg.SCNMessageID,
	})
}

func (h MessageHandler) deliverMessage(ctx *logic.AppContext, client models.Client, msg models.Message) (*string, error) {
	if client.FCMToken != nil {
		fcmDelivID, err := h.app.Firebase.SendNotification(ctx, client, msg)
		if err != nil {
			return nil, err
		}
		return langext.Ptr(fcmDelivID), nil
	} else {
		return langext.Ptr(""), nil
	}
}
