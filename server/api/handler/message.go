package handler

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/common/ginresp"
	"blackforestbytes.com/simplecloudnotifier/db"
	"blackforestbytes.com/simplecloudnotifier/logic"
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

// SendMessageCompat swaggerdoc
//
// @Deprecated
//
// @Summary     Send a new message (compatibility)
// @Description All parameter can be set via query-parameter or form-data body. Only UserID, UserKey and Title are required
//
// @Param       query_data query    handler.SendMessageCompat.query false " "
// @Param       form_data  formData handler.SendMessageCompat.form  false " "
//
// @Success     200        {object} handler.sendMessageInternal.response
// @Failure     400        {object} ginresp.apiError
// @Failure     401        {object} ginresp.apiError
// @Failure     403        {object} ginresp.apiError
// @Failure     404        {object} ginresp.apiError
// @Failure     500        {object} ginresp.apiError
//
// @Router      /send.php [POST]
func (h MessageHandler) SendMessageCompat(g *gin.Context) ginresp.HTTPResponse {
	type query struct {
		UserID        *int64   `form:"user_id"`
		UserKey       *string  `form:"user_key"`
		Title         *string  `form:"title"`
		Content       *string  `form:"content"`
		Priority      *int     `form:"priority"`
		UserMessageID *string  `form:"msg_id"`
		SendTimestamp *float64 `form:"timestamp"`
	}
	type form struct {
		UserID        *int64   `form:"user_id"`
		UserKey       *string  `form:"user_key"`
		Title         *string  `form:"title"`
		Content       *string  `form:"content"`
		Priority      *int     `form:"priority"`
		UserMessageID *string  `form:"msg_id"`
		SendTimestamp *float64 `form:"timestamp"`
	}

	var f form
	var q query
	ctx, errResp := h.app.StartRequest(g, nil, &q, nil, &f)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	data := dataext.ObjectMerge(f, q)

	return h.sendMessageInternal(g, ctx, data.UserID, data.UserKey, langext.Ptr(h.app.DefaultChannel), nil, data.Title, data.Content, data.Priority, data.UserMessageID, data.SendTimestamp)

}

// SendMessage swaggerdoc
//
// @Summary     Send a new message
// @Description All parameter can be set via query-parameter or the json body. Only UserID, UserKey and Title are required
//
// @Param       query_data query    handler.SendMessage.query false " "
// @Param       post_body  body     handler.SendMessage.body  false " "
//
// @Success     200        {object} handler.sendMessageInternal.response
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
		Title         *string  `form:"title"`
		Content       *string  `form:"content"`
		Priority      *int     `form:"priority"`
		UserMessageID *string  `form:"msg_id"`
		SendTimestamp *float64 `form:"timestamp"`
	}
	type body struct {
		UserID        *int64   `json:"user_id"`
		UserKey       *string  `json:"user_key"`
		Channel       *string  `json:"channel"`
		ChanKey       *string  `form:"chan_key"`
		Title         *string  `json:"title"`
		Content       *string  `json:"content"`
		Priority      *int     `json:"priority"`
		UserMessageID *string  `json:"msg_id"`
		SendTimestamp *float64 `json:"timestamp"`
	}

	var b body
	var q query
	ctx, errResp := h.app.StartRequest(g, nil, &q, &b, nil)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	data := dataext.ObjectMerge(b, q)

	return h.sendMessageInternal(g, ctx, data.UserID, data.UserKey, data.Channel, data.ChanKey, data.Title, data.Content, data.Priority, data.UserMessageID, data.SendTimestamp)

}

func (h MessageHandler) sendMessageInternal(g *gin.Context, ctx *logic.AppContext, UserID *int64, UserKey *string, Channel *string, ChanKey *string, Title *string, Content *string, Priority *int, UserMessageID *string, SendTimestamp *float64) ginresp.HTTPResponse {
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

	if Title != nil {
		Title = langext.Ptr(strings.TrimSpace(*Title))
	}
	if UserMessageID != nil {
		UserMessageID = langext.Ptr(strings.TrimSpace(*UserMessageID))
	}

	if UserID == nil {
		return ginresp.SendAPIError(g, 400, apierr.MISSING_UID, 101, "Missing parameter [[user_id]]", nil)
	}
	if UserKey == nil {
		return ginresp.SendAPIError(g, 400, apierr.MISSING_UID, 102, "Missing parameter [[user_token]]", nil)
	}
	if Title == nil {
		return ginresp.SendAPIError(g, 400, apierr.MISSING_UID, 103, "Missing parameter [[title]]", nil)
	}
	if SendTimestamp != nil && mathext.Abs(*SendTimestamp-float64(time.Now().Unix())) > (24*time.Hour).Seconds() {
		return ginresp.SendAPIError(g, 400, apierr.TIMESTAMP_OUT_OF_RANGE, -1, "The timestamp mus be within 24 hours of now()", nil)
	}
	if Priority != nil && (*Priority != 0 && *Priority != 1 && *Priority != 2) {
		return ginresp.SendAPIError(g, 400, apierr.INVALID_PRIO, 105, "Invalid priority", nil)
	}
	if len(*Title) == 0 {
		return ginresp.SendAPIError(g, 400, apierr.NO_TITLE, 103, "No title specified", nil)
	}
	if UserMessageID != nil && len(*UserMessageID) > 64 {
		return ginresp.SendAPIError(g, 400, apierr.USR_MSG_ID_TOO_LONG, -1, "MessageID too long (64 characters)", nil)
	}

	channelName := h.app.DefaultChannel
	if Channel != nil {
		channelName = h.app.NormalizeChannelName(*Channel)
	}

	user, err := h.database.GetUser(ctx, *UserID)
	if err == sql.ErrNoRows {
		return ginresp.SendAPIError(g, 400, apierr.USER_NOT_FOUND, -1, "User not found", nil)
	}
	if err != nil {
		return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, -1, "Failed to query user", err)
	}

	if len(*Title) > 120 {
		return ginresp.SendAPIError(g, 400, apierr.TITLE_TOO_LONG, 103, "Title too long (120 characters)", nil)
	}
	if Content != nil && len(*Content) > user.MaxContentLength() {
		return ginresp.SendAPIError(g, 400, apierr.CONTENT_TOO_LONG, 104, fmt.Sprintf("Content too long (%d characters; max := %d characters)", len(*Content), user.MaxContentLength()), nil)
	}

	if UserMessageID != nil {
		msg, err := h.database.GetMessageByUserMessageID(ctx, *UserMessageID)
		if err != nil {
			return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, -1, "Failed to query existing message", err)
		}
		if msg != nil {
			return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, response{
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
			}))
		}
	}

	if user.QuotaRemainingToday() <= 0 {
		return ginresp.SendAPIError(g, 403, apierr.QUOTA_REACHED, -1, fmt.Sprintf("Daily quota reached (%d)", user.QuotaPerDay()), nil)
	}

	channel, err := h.app.GetOrCreateChannel(ctx, *UserID, channelName)
	if err != nil {
		return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, -1, "Failed to query/create channel", err)
	}

	selfChanAdmin := *UserID == channel.OwnerUserID && *UserKey == user.AdminKey
	selfChanSend := *UserID == channel.OwnerUserID && *UserKey == user.SendKey
	forgChanSend := *UserID != channel.OwnerUserID && ChanKey != nil && *ChanKey == channel.SendKey

	if !selfChanAdmin && !selfChanSend && !forgChanSend {
		return ginresp.SendAPIError(g, 401, apierr.USER_AUTH_FAILED, 102, fmt.Sprintf("Daily quota reached (%d)", user.QuotaPerDay()), nil)
	}

	var sendTimestamp *time.Time = nil
	if SendTimestamp != nil {
		sendTimestamp = langext.Ptr(timeext.UnixFloatSeconds(*SendTimestamp))
	}

	priority := langext.Coalesce(Priority, 1)

	msg, err := h.database.CreateMessage(ctx, *UserID, channel, sendTimestamp, *Title, Content, priority, UserMessageID)
	if err != nil {
		return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, -1, "Failed to create message in db", err)
	}

	subscriptions, err := h.database.ListSubscriptionsByChannel(ctx, channel.ChannelID)
	if err != nil {
		return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, -1, "Failed to query subscriptions", err)
	}

	err = h.database.IncUserMessageCounter(ctx, user)
	if err != nil {
		return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, -1, "Failed to inc user msg-counter", err)
	}

	err = h.database.IncChannelMessageCounter(ctx, channel)
	if err != nil {
		return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, -1, "Failed to inc channel msg-counter", err)
	}

	for _, sub := range subscriptions {
		clients, err := h.database.ListClients(ctx, sub.SubscriberUserID)
		if err != nil {
			return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, -1, "Failed to query clients", err)
		}

		if !sub.Confirmed {
			continue
		}

		for _, client := range clients {

			fcmDelivID, err := h.app.DeliverMessage(ctx, client, msg)
			if err != nil {
				_, err = h.database.CreateRetryDelivery(ctx, client, msg)
				if err != nil {
					return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, -1, "Failed to create delivery", err)
				}
			} else {
				_, err = h.database.CreateSuccessDelivery(ctx, client, msg, *fcmDelivID)
				if err != nil {
					return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, -1, "Failed to create delivery", err)
				}
			}

		}
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, response{
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
	}))
}
