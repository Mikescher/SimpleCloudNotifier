package handler

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	hl "blackforestbytes.com/simplecloudnotifier/api/apihighlight"
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

// SendMessageCompat swaggerdoc
//
// @Deprecated
//
// @Summary     Send a new message (compatibility)
// @Description All parameter can be set via query-parameter or form-data body. Only UserID, UserKey and Title are required
// @Tags        External
//
// @Param       query_data query    handler.SendMessageCompat.combined false " "
// @Param       form_data  formData handler.SendMessageCompat.combined false " "
//
// @Success     200        {object} handler.sendMessageInternal.response
// @Failure     400        {object} ginresp.apiError
// @Failure     401        {object} ginresp.apiError
// @Failure     403        {object} ginresp.apiError
// @Failure     500        {object} ginresp.apiError
//
// @Router      /send.php [POST]
func (h MessageHandler) SendMessageCompat(g *gin.Context) ginresp.HTTPResponse {
	type combined struct {
		UserID        *models.UserID `json:"user_id"   form:"user_id"`
		UserKey       *string        `json:"user_key"  form:"user_key"`
		Title         *string        `json:"title"     form:"title"`
		Content       *string        `json:"content"   form:"content"`
		Priority      *int           `json:"priority"  form:"priority"`
		UserMessageID *string        `json:"msg_id"    form:"msg_id"`
		SendTimestamp *float64       `json:"timestamp" form:"timestamp"`
	}

	var f combined
	var q combined
	ctx, errResp := h.app.StartRequest(g, nil, &q, nil, &f)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	data := dataext.ObjectMerge(f, q)

	return h.sendMessageInternal(g, ctx, data.UserID, data.UserKey, nil, nil, data.Title, data.Content, data.Priority, data.UserMessageID, data.SendTimestamp, nil)
}

// SendMessage swaggerdoc
//
// @Summary     Send a new message
// @Description All parameter can be set via query-parameter or the json body. Only UserID, UserKey and Title are required
// @Tags        External
//
// @Param       query_data query    handler.SendMessage.combined false " "
// @Param       post_body  body     handler.SendMessage.combined false " "
// @Param       form_body  formData handler.SendMessage.combined false " "
//
// @Success     200        {object} handler.sendMessageInternal.response
// @Failure     400        {object} ginresp.apiError
// @Failure     401        {object} ginresp.apiError "The user_id was not found or the user_key is wrong"
// @Failure     403        {object} ginresp.apiError "The user has exceeded its daily quota - wait 24 hours or upgrade your account"
// @Failure     500        {object} ginresp.apiError "An internal server error occurred - try again later"
//
// @Router      /     [POST]
// @Router      /send [POST]
func (h MessageHandler) SendMessage(g *gin.Context) ginresp.HTTPResponse {
	type combined struct {
		UserID        *models.UserID `json:"user_id"     form:"user_id"     example:"7725"                               `
		UserKey       *string        `json:"user_key"    form:"user_key"    example:"P3TNH8mvv14fm"                      `
		Channel       *string        `json:"channel"     form:"channel"     example:"test"                               `
		ChanKey       *string        `json:"chan_key"    form:"chan_key"    example:"qhnUbKcLgp6tg"                      `
		Title         *string        `json:"title"       form:"title"       example:"Hello World"                        `
		Content       *string        `json:"content"     form:"content"     example:"This is a message"                  `
		Priority      *int           `json:"priority"    form:"priority"    example:"1"                   enums:"0,1,2"  `
		UserMessageID *string        `json:"msg_id"      form:"msg_id"      example:"db8b0e6a-a08c-4646"                 `
		SendTimestamp *float64       `json:"timestamp"   form:"timestamp"   example:"1669824037"                         `
		SenderName    *string        `json:"sender_name" form:"sender_name" example:"example-server"                     `
	}

	var b combined
	var q combined
	var f combined
	ctx, errResp := h.app.StartRequest(g, nil, &q, &b, &f)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	// query has highest prio, then form, then json
	data := dataext.ObjectMerge(dataext.ObjectMerge(b, f), q)

	return h.sendMessageInternal(g, ctx, data.UserID, data.UserKey, data.Channel, data.ChanKey, data.Title, data.Content, data.Priority, data.UserMessageID, data.SendTimestamp, data.SenderName)

}

func (h MessageHandler) sendMessageInternal(g *gin.Context, ctx *logic.AppContext, UserID *models.UserID, UserKey *string, Channel *string, ChanKey *string, Title *string, Content *string, Priority *int, UserMessageID *string, SendTimestamp *float64, SenderName *string) ginresp.HTTPResponse {
	type response struct {
		Success        bool                `json:"success"`
		ErrorID        apierr.APIError     `json:"error"`
		ErrorHighlight int                 `json:"errhighlight"`
		Message        string              `json:"message"`
		SuppressSend   bool                `json:"suppress_send"`
		MessageCount   int                 `json:"messagecount"`
		Quota          int                 `json:"quota"`
		IsPro          bool                `json:"is_pro"`
		QuotaMax       int                 `json:"quota_max"`
		SCNMessageID   models.SCNMessageID `json:"scn_msg_id"`
	}

	if Title != nil {
		Title = langext.Ptr(strings.TrimSpace(*Title))
	}
	if UserMessageID != nil {
		UserMessageID = langext.Ptr(strings.TrimSpace(*UserMessageID))
	}

	if UserID == nil {
		return ginresp.SendAPIError(g, 400, apierr.MISSING_UID, hl.USER_ID, "Missing parameter [[user_id]]", nil)
	}
	if UserKey == nil {
		return ginresp.SendAPIError(g, 400, apierr.MISSING_TOK, hl.USER_KEY, "Missing parameter [[user_token]]", nil)
	}
	if Title == nil {
		return ginresp.SendAPIError(g, 400, apierr.MISSING_TITLE, hl.TITLE, "Missing parameter [[title]]", nil)
	}
	if Priority != nil && (*Priority != 0 && *Priority != 1 && *Priority != 2) {
		return ginresp.SendAPIError(g, 400, apierr.INVALID_PRIO, hl.PRIORITY, "Invalid priority", nil)
	}
	if len(*Title) == 0 {
		return ginresp.SendAPIError(g, 400, apierr.NO_TITLE, hl.TITLE, "No title specified", nil)
	}

	user, err := h.database.GetUser(ctx, *UserID)
	if err == sql.ErrNoRows {
		return ginresp.SendAPIError(g, 400, apierr.USER_NOT_FOUND, hl.USER_ID, "User not found", nil)
	}
	if err != nil {
		return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query user", err)
	}

	channelName := user.DefaultChannel()
	if Channel != nil {
		channelName = h.app.NormalizeChannelName(*Channel)
	}

	if len(*Title) > user.MaxTitleLength() {
		return ginresp.SendAPIError(g, 400, apierr.TITLE_TOO_LONG, hl.TITLE, fmt.Sprintf("Title too long (max %d characters)", user.MaxTitleLength()), nil)
	}
	if Content != nil && len(*Content) > user.MaxContentLength() {
		return ginresp.SendAPIError(g, 400, apierr.CONTENT_TOO_LONG, hl.CONTENT, fmt.Sprintf("Content too long (%d characters; max := %d characters)", len(*Content), user.MaxContentLength()), nil)
	}
	if len(channelName) > user.MaxChannelNameLength() {
		return ginresp.SendAPIError(g, 400, apierr.CONTENT_TOO_LONG, hl.CHANNEL, fmt.Sprintf("Channel too long (max %d characters)", user.MaxChannelNameLength()), nil)
	}
	if SenderName != nil && len(*SenderName) > user.MaxSenderName() {
		return ginresp.SendAPIError(g, 400, apierr.SENDERNAME_TOO_LONG, hl.SENDER_NAME, fmt.Sprintf("SenderName too long (max %d characters)", user.MaxSenderName()), nil)
	}
	if UserMessageID != nil && len(*UserMessageID) > user.MaxUserMessageID() {
		return ginresp.SendAPIError(g, 400, apierr.USR_MSG_ID_TOO_LONG, hl.USER_MESSAGE_ID, fmt.Sprintf("MessageID too long (max %d characters)", user.MaxUserMessageID()), nil)
	}
	if SendTimestamp != nil && mathext.Abs(*SendTimestamp-float64(time.Now().Unix())) > timeext.FromHours(user.MaxTimestampDiffHours()).Seconds() {
		return ginresp.SendAPIError(g, 400, apierr.TIMESTAMP_OUT_OF_RANGE, hl.NONE, fmt.Sprintf("The timestamp mus be within %d hours of now()", user.MaxTimestampDiffHours()), nil)
	}

	if UserMessageID != nil {
		msg, err := h.database.GetMessageByUserMessageID(ctx, *UserMessageID)
		if err != nil {
			return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query existing message", err)
		}
		if msg != nil {
			//the found message can be deleted (!), but we still return NO_ERROR here...
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
		return ginresp.SendAPIError(g, 403, apierr.QUOTA_REACHED, hl.NONE, fmt.Sprintf("Daily quota reached (%d)", user.QuotaPerDay()), nil)
	}

	var channel models.Channel
	if ChanKey != nil {
		// foreign channel (+ channel send-key)

		foreignChan, err := h.database.GetChannelByNameAndSendKey(ctx, channelName, *ChanKey)
		if err != nil {
			return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query (foreign) channel", err)
		}
		if foreignChan == nil {
			return ginresp.SendAPIError(g, 400, apierr.CHANNEL_NOT_FOUND, hl.CHANNEL, "(Foreign) Channel not found", err)
		}
		channel = *foreignChan
	} else {
		// own channel

		channel, err = h.app.GetOrCreateChannel(ctx, *UserID, channelName)
		if err != nil {
			return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query/create (owned) channel", err)
		}
	}

	selfChanAdmin := *UserID == channel.OwnerUserID && *UserKey == user.AdminKey
	selfChanSend := *UserID == channel.OwnerUserID && *UserKey == user.SendKey
	forgChanSend := *UserID != channel.OwnerUserID && ChanKey != nil && *ChanKey == channel.SendKey

	if !selfChanAdmin && !selfChanSend && !forgChanSend {
		return ginresp.SendAPIError(g, 401, apierr.USER_AUTH_FAILED, hl.USER_KEY, "You are not authorized for this action", nil)
	}

	var sendTimestamp *time.Time = nil
	if SendTimestamp != nil {
		sendTimestamp = langext.Ptr(timeext.UnixFloatSeconds(*SendTimestamp))
	}

	priority := langext.Coalesce(Priority, user.DefaultPriority())

	clientIP := g.ClientIP()

	msg, err := h.database.CreateMessage(ctx, *UserID, channel, sendTimestamp, *Title, Content, priority, UserMessageID, clientIP, SenderName)
	if err != nil {
		return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to create message in db", err)
	}

	subscriptions, err := h.database.ListSubscriptionsByChannel(ctx, channel.ChannelID)
	if err != nil {
		return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query subscriptions", err)
	}

	err = h.database.IncUserMessageCounter(ctx, user)
	if err != nil {
		return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to inc user msg-counter", err)
	}

	err = h.database.IncChannelMessageCounter(ctx, channel)
	if err != nil {
		return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to inc channel msg-counter", err)
	}

	for _, sub := range subscriptions {
		clients, err := h.database.ListClients(ctx, sub.SubscriberUserID)
		if err != nil {
			return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query clients", err)
		}

		if !sub.Confirmed {
			continue
		}

		for _, client := range clients {

			fcmDelivID, err := h.app.DeliverMessage(ctx, client, msg)
			if err != nil {
				_, err = h.database.CreateRetryDelivery(ctx, client, msg)
				if err != nil {
					return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to create delivery", err)
				}
			} else {
				_, err = h.database.CreateSuccessDelivery(ctx, client, msg, *fcmDelivID)
				if err != nil {
					return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to create delivery", err)
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
