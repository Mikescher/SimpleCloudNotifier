package handler

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	hl "blackforestbytes.com/simplecloudnotifier/api/apihighlight"
	"blackforestbytes.com/simplecloudnotifier/api/ginresp"
	primarydb "blackforestbytes.com/simplecloudnotifier/db/impl/primary"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/dataext"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/mathext"
	"gogs.mikescher.com/BlackForestBytes/goext/timeext"
	"net/http"
	"strings"
	"time"
)

type SendMessageResponse struct {
	User            models.User
	Message         models.Message
	MessageIsOld    bool
	CompatMessageID int64
}

type MessageHandler struct {
	app      *logic.Application
	database *primarydb.Database
}

func NewMessageHandler(app *logic.Application) MessageHandler {
	return MessageHandler{
		app:      app,
		database: app.Database.Primary,
	}
}

// SendMessage swaggerdoc
//
//	@Summary		Send a new message
//	@Description	All parameter can be set via query-parameter or the json body. Only UserID, UserKey and Title are required
//	@Tags			External
//
//	@Param			query_data	query		handler.SendMessage.combined	false	" "
//	@Param			post_body	body		handler.SendMessage.combined	false	" "
//	@Param			form_body	formData	handler.SendMessage.combined	false	" "
//
//	@Success		200			{object}	handler.SendMessage.response
//	@Failure		400			{object}	ginresp.apiError
//	@Failure		401			{object}	ginresp.apiError	"The user_id was not found or the user_key is wrong"
//	@Failure		403			{object}	ginresp.apiError	"The user has exceeded its daily quota - wait 24 hours or upgrade your account"
//	@Failure		500			{object}	ginresp.apiError	"An internal server error occurred - try again later"
//
//	@Router			/     [POST]
//	@Router			/send [POST]
func (h MessageHandler) SendMessage(g *gin.Context) ginresp.HTTPResponse {
	type combined struct {
		UserID        *models.UserID `json:"user_id"     form:"user_id"     example:"7725"                               `
		KeyToken      *string        `json:"key"         form:"key"         example:"P3TNH8mvv14fm"                      `
		Channel       *string        `json:"channel"     form:"channel"     example:"test"                               `
		Title         *string        `json:"title"       form:"title"       example:"Hello World"                        `
		Content       *string        `json:"content"     form:"content"     example:"This is a message"                  `
		Priority      *int           `json:"priority"    form:"priority"    example:"1"                   enums:"0,1,2"  `
		UserMessageID *string        `json:"msg_id"      form:"msg_id"      example:"db8b0e6a-a08c-4646"                 `
		SendTimestamp *float64       `json:"timestamp"   form:"timestamp"   example:"1669824037"                         `
		SenderName    *string        `json:"sender_name" form:"sender_name" example:"example-server"                     `
	}

	type response struct {
		Success        bool             `json:"success"`
		ErrorID        apierr.APIError  `json:"error"`
		ErrorHighlight int              `json:"errhighlight"`
		Message        string           `json:"message"`
		SuppressSend   bool             `json:"suppress_send"`
		MessageCount   int              `json:"messagecount"`
		Quota          int              `json:"quota"`
		IsPro          bool             `json:"is_pro"`
		QuotaMax       int              `json:"quota_max"`
		SCNMessageID   models.MessageID `json:"scn_msg_id"`
	}

	var b combined
	var q combined
	var f combined
	ctx, errResp := h.app.StartRequest(g, nil, &q, &b, &f, logic.RequestOptions{IgnoreWrongContentType: true})
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	// query has highest prio, then form, then json
	data := dataext.ObjectMerge(dataext.ObjectMerge(b, f), q)

	okResp, errResp := h.sendMessageInternal(g, ctx, data.UserID, data.KeyToken, data.Channel, data.Title, data.Content, data.Priority, data.UserMessageID, data.SendTimestamp, data.SenderName)
	if errResp != nil {
		return *errResp
	} else {
		return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, response{
			Success:        true,
			ErrorID:        apierr.NO_ERROR,
			ErrorHighlight: -1,
			Message:        langext.Conditional(okResp.MessageIsOld, "Message already sent", "Message sent"),
			SuppressSend:   okResp.MessageIsOld,
			MessageCount:   okResp.User.MessagesSent,
			Quota:          okResp.User.QuotaUsedToday(),
			IsPro:          okResp.User.IsPro,
			QuotaMax:       okResp.User.QuotaPerDay(),
			SCNMessageID:   okResp.Message.MessageID,
		}))
	}
}

func (h MessageHandler) sendMessageInternal(g *gin.Context, ctx *logic.AppContext, UserID *models.UserID, Key *string, Channel *string, Title *string, Content *string, Priority *int, UserMessageID *string, SendTimestamp *float64, SenderName *string) (*SendMessageResponse, *ginresp.HTTPResponse) {
	if Title != nil {
		Title = langext.Ptr(strings.TrimSpace(*Title))
	}
	if UserMessageID != nil {
		UserMessageID = langext.Ptr(strings.TrimSpace(*UserMessageID))
	}

	if UserID == nil {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 400, apierr.MISSING_UID, hl.USER_ID, "Missing parameter [[user_id]]", nil))
	}
	if Key == nil {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 400, apierr.MISSING_TOK, hl.USER_KEY, "Missing parameter [[key]]", nil))
	}
	if Title == nil {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 400, apierr.MISSING_TITLE, hl.TITLE, "Missing parameter [[title]]", nil))
	}
	if Priority != nil && (*Priority != 0 && *Priority != 1 && *Priority != 2) {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 400, apierr.INVALID_PRIO, hl.PRIORITY, "Invalid priority", nil))
	}
	if len(*Title) == 0 {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 400, apierr.NO_TITLE, hl.TITLE, "No title specified", nil))
	}

	user, err := h.database.GetUser(ctx, *UserID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 400, apierr.USER_NOT_FOUND, hl.USER_ID, "User not found", err))
	}
	if err != nil {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query user", err))
	}

	channelDisplayName := user.DefaultChannel()
	channelInternalName := user.DefaultChannel()
	if Channel != nil {
		channelDisplayName = h.app.NormalizeChannelDisplayName(*Channel)
		channelInternalName = h.app.NormalizeChannelInternalName(*Channel)
	}

	if len(*Title) > user.MaxTitleLength() {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 400, apierr.TITLE_TOO_LONG, hl.TITLE, fmt.Sprintf("Title too long (max %d characters)", user.MaxTitleLength()), nil))
	}
	if Content != nil && len(*Content) > user.MaxContentLength() {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 400, apierr.CONTENT_TOO_LONG, hl.CONTENT, fmt.Sprintf("Content too long (%d characters; max := %d characters)", len(*Content), user.MaxContentLength()), nil))
	}
	if len(channelDisplayName) > user.MaxChannelNameLength() {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 400, apierr.CHANNEL_TOO_LONG, hl.CHANNEL, fmt.Sprintf("Channel too long (max %d characters)", user.MaxChannelNameLength()), nil))
	}
	if len(strings.TrimSpace(channelDisplayName)) == 0 {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 400, apierr.CHANNEL_NAME_EMPTY, hl.CHANNEL, fmt.Sprintf("Channel displayname cannot be empty"), nil))
	}
	if len(channelInternalName) > user.MaxChannelNameLength() {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 400, apierr.CHANNEL_TOO_LONG, hl.CHANNEL, fmt.Sprintf("Channel too long (max %d characters)", user.MaxChannelNameLength()), nil))
	}
	if len(strings.TrimSpace(channelInternalName)) == 0 {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 400, apierr.CHANNEL_NAME_EMPTY, hl.CHANNEL, fmt.Sprintf("Channel internalname cannot be empty"), nil))
	}
	if SenderName != nil && len(*SenderName) > user.MaxSenderNameLength() {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 400, apierr.SENDERNAME_TOO_LONG, hl.SENDER_NAME, fmt.Sprintf("SenderName too long (max %d characters)", user.MaxSenderNameLength()), nil))
	}
	if UserMessageID != nil && len(*UserMessageID) > user.MaxUserMessageIDLength() {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 400, apierr.USR_MSG_ID_TOO_LONG, hl.USER_MESSAGE_ID, fmt.Sprintf("MessageID too long (max %d characters)", user.MaxUserMessageIDLength()), nil))
	}
	if SendTimestamp != nil && mathext.Abs(*SendTimestamp-float64(time.Now().Unix())) > timeext.FromHours(user.MaxTimestampDiffHours()).Seconds() {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 400, apierr.TIMESTAMP_OUT_OF_RANGE, hl.NONE, fmt.Sprintf("The timestamp mus be within %d hours of now()", user.MaxTimestampDiffHours()), nil))
	}

	if UserMessageID != nil {
		msg, err := h.database.GetMessageByUserMessageID(ctx, *UserMessageID)
		if err != nil {
			return nil, langext.Ptr(ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query existing message", err))
		}
		if msg != nil {

			existingCompID, _, err := h.database.ConvertToCompatID(ctx, msg.MessageID.String())
			if err != nil {
				return nil, langext.Ptr(ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query compat-id", err))
			}

			if existingCompID == nil {
				v, err := h.database.CreateCompatID(ctx, "messageid", msg.MessageID.String())
				if err != nil {
					return nil, langext.Ptr(ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to create compat-id", err))
				}
				existingCompID = &v
			}

			//the found message can be deleted (!), but we still return NO_ERROR here...
			return &SendMessageResponse{
				User:            user,
				Message:         *msg,
				MessageIsOld:    true,
				CompatMessageID: *existingCompID,
			}, nil
		}
	}

	if user.QuotaRemainingToday() <= 0 {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 403, apierr.QUOTA_REACHED, hl.NONE, fmt.Sprintf("Daily quota reached (%d)", user.QuotaPerDay()), nil))
	}

	channel, err := h.app.GetOrCreateChannel(ctx, *UserID, channelDisplayName, channelInternalName)
	if err != nil {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query/create (owned) channel", err))
	}

	keytok, permResp := ctx.CheckPermissionSend(channel, *Key)
	if permResp != nil {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 401, apierr.USER_AUTH_FAILED, hl.USER_KEY, "You are not authorized for this action", nil))
	}

	var sendTimestamp *time.Time = nil
	if SendTimestamp != nil {
		sendTimestamp = langext.Ptr(timeext.UnixFloatSeconds(*SendTimestamp))
	}

	priority := langext.Coalesce(Priority, user.DefaultPriority())

	clientIP := g.ClientIP()

	msg, err := h.database.CreateMessage(ctx, *UserID, channel, sendTimestamp, *Title, Content, priority, UserMessageID, clientIP, SenderName, keytok.KeyTokenID)
	if err != nil {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to create message in db", err))
	}

	compatMsgID, err := h.database.CreateCompatID(ctx, "messageid", msg.MessageID.String())
	if err != nil {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to create compat-id", err))
	}

	subFilter := models.SubscriptionFilter{ChannelID: langext.Ptr([]models.ChannelID{channel.ChannelID}), Confirmed: langext.PTrue}
	activeSubscriptions, err := h.database.ListSubscriptions(ctx, subFilter)
	if err != nil {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query subscriptions", err))
	}

	err = h.database.IncUserMessageCounter(ctx, &user)
	if err != nil {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to inc user msg-counter", err))
	}

	err = h.database.IncChannelMessageCounter(ctx, &channel)
	if err != nil {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to inc channel msg-counter", err))
	}

	err = h.database.IncKeyTokenMessageCounter(ctx, keytok)
	if err != nil {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to inc token msg-counter", err))
	}

	log.Info().Msg(fmt.Sprintf("Sending new notification %s for user %s", msg.MessageID, UserID))

	for _, sub := range activeSubscriptions {
		clients, err := h.database.ListClients(ctx, sub.SubscriberUserID)
		if err != nil {
			return nil, langext.Ptr(ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query clients", err))
		}

		for _, client := range clients {

			isCompatClient, err := h.database.IsCompatClient(ctx, client.ClientID)
			if err != nil {
				return nil, langext.Ptr(ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query compat_clients", err))
			}

			var titleOverride *string = nil
			var msgidOverride *string = nil
			if isCompatClient {
				titleOverride = langext.Ptr(h.app.CompatizeMessageTitle(ctx, msg))
				msgidOverride = langext.Ptr(fmt.Sprintf("%d", compatMsgID))
			}

			fcmDelivID, err := h.app.DeliverMessage(ctx, client, msg, titleOverride, msgidOverride)
			if err != nil {
				_, err = h.database.CreateRetryDelivery(ctx, client, msg)
				if err != nil {
					return nil, langext.Ptr(ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to create delivery", err))
				}
			} else {
				_, err = h.database.CreateSuccessDelivery(ctx, client, msg, fcmDelivID)
				if err != nil {
					return nil, langext.Ptr(ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to create delivery", err))
				}
			}

		}
	}

	return &SendMessageResponse{
		User:            user,
		Message:         msg,
		MessageIsOld:    false,
		CompatMessageID: compatMsgID,
	}, nil
}

// UptimeKumaWebHook swaggerdoc
//
//	@Summary		Send a new message
//	@Description	All parameter can be set via query-parameter or the json body. Only UserID, UserKey and Title are required
//	@Tags			External
//
//	@Param			query_data	query		handler.UptimeKumaWebHook.query					false	" "
//	@Param			post_body	body		handler.UptimeKumaWebHook.uptimeKumaWebhookBody	false	" "
//
//	@Success		200			{object}	any
//	@Failure		400			{object}	ginresp.apiError
//	@Failure		401			{object}	ginresp.apiError	"The user_id was not found or the user_key is wrong"
//	@Failure		403			{object}	ginresp.apiError	"The user has exceeded its daily quota - wait 24 hours or upgrade your account"
//	@Failure		500			{object}	ginresp.apiError	"An internal server error occurred - try again later"
//
//	@Router			/webhook/uptime-kuma [POST]
func (h MessageHandler) UptimeKumaWebHook(g *gin.Context) ginresp.HTTPResponse {
	type query struct {
		UserID   *models.UserID `form:"user_id"     example:"7725"`
		KeyToken *string        `form:"key"         example:"P3TNH8mvv14fm"`
	}

	type uptimeKumaWebhookBody struct {
		Heartbeat *struct {
			Time           string `json:"time"`
			Status         int    `json:"status"`
			Msg            string `json:"msg"`
			Timezone       string `json:"timezone"`
			TimezoneOffset string `json:"timezoneOffset"`
			LocalDateTime  string `json:"localDateTime"`
		} `json:"heartbeat"`
		Monitor *struct {
			Name string  `json:"name"`
			Url  *string `json:"url"`
		} `json:"monitor"`
		Msg string `json:"msg"`
	}

	var b uptimeKumaWebhookBody
	var q query

	ctx, httpErr := h.app.StartRequest(g, nil, &q, &b, nil)
	if httpErr != nil {
		return *httpErr
	}
	defer ctx.Cancel()

	var title = ""

	var content = ""
	content += fmt.Sprintf("%v\n", b.Msg)
	if b.Monitor != nil {
		content += fmt.Sprintf("%v\n", b.Monitor.Name)
		if b.Monitor.Url != nil {
			content += fmt.Sprintf("url: %v\n", *b.Monitor.Url)
		}

		if b.Heartbeat != nil {
			statusString := "down"

			if b.Heartbeat.Status == 1 {
				statusString = "up"
			}
			title = fmt.Sprintf("%v %v!", b.Monitor.Name, statusString)
		}

	}

	if b.Heartbeat != nil {
		content += "\n===== Heartbeat ======\n"
		content += fmt.Sprintf("msg: %v\n", b.Heartbeat.Msg)
		content += fmt.Sprintf("timestamp: %v\n", b.Heartbeat.Time)
		content += fmt.Sprintf("timezone: %v\n", b.Heartbeat.Timezone)
		content += fmt.Sprintf("timezone offset: %v\n", b.Heartbeat.TimezoneOffset)
		content += fmt.Sprintf("local date time: %v\n", b.Heartbeat.TimezoneOffset)
	}
	okResp, errResp := h.sendMessageInternal(g, ctx, q.UserID, q.KeyToken, nil, &title, &content, langext.Ptr(1), nil, nil, nil)

	if errResp != nil {
		return *errResp
	}
	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, okResp))
}
