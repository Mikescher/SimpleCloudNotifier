package logic

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	hl "blackforestbytes.com/simplecloudnotifier/api/apihighlight"
	"blackforestbytes.com/simplecloudnotifier/api/ginresp"
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/mathext"
	"gogs.mikescher.com/BlackForestBytes/goext/timeext"
	"strings"
	"time"
)

type SendMessageResponse struct {
	User            models.User
	Message         models.Message
	MessageIsOld    bool
	CompatMessageID int64
}

func (app *Application) SendMessage(g *gin.Context, ctx *AppContext, UserID *models.UserID, Key *string, Channel *string, Title *string, Content *string, Priority *int, UserMessageID *string, SendTimestamp *float64, SenderName *string) (*SendMessageResponse, *ginext.HTTPResponse) {
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

	user, err := app.Database.Primary.GetUser(ctx, *UserID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 400, apierr.USER_NOT_FOUND, hl.USER_ID, "User not found", err))
	}
	if err != nil {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query user", err))
	}

	channelDisplayName := user.DefaultChannel()
	channelInternalName := user.DefaultChannel()
	if Channel != nil {
		channelDisplayName = app.NormalizeChannelDisplayName(*Channel)
		channelInternalName = app.NormalizeChannelInternalName(*Channel)
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
		msg, err := app.Database.Primary.GetMessageByUserMessageID(ctx, *UserMessageID)
		if err != nil {
			return nil, langext.Ptr(ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query existing message", err))
		}
		if msg != nil {

			existingCompID, _, err := app.Database.Primary.ConvertToCompatID(ctx, msg.MessageID.String())
			if err != nil {
				return nil, langext.Ptr(ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query compat-id", err))
			}

			if existingCompID == nil {
				v, err := app.Database.Primary.CreateCompatID(ctx, "messageid", msg.MessageID.String())
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

	channel, err := app.GetOrCreateChannel(ctx, *UserID, channelDisplayName, channelInternalName)
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

	msg, err := app.Database.Primary.CreateMessage(ctx, *UserID, channel, sendTimestamp, *Title, Content, priority, UserMessageID, clientIP, SenderName, keytok.KeyTokenID)
	if err != nil {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to create message in db", err))
	}

	compatMsgID, err := app.Database.Primary.CreateCompatID(ctx, "messageid", msg.MessageID.String())
	if err != nil {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to create compat-id", err))
	}

	subFilter := models.SubscriptionFilter{ChannelID: langext.Ptr([]models.ChannelID{channel.ChannelID}), Confirmed: langext.PTrue}
	activeSubscriptions, err := app.Database.Primary.ListSubscriptions(ctx, subFilter)
	if err != nil {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query subscriptions", err))
	}

	err = app.Database.Primary.IncUserMessageCounter(ctx, &user)
	if err != nil {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to inc user msg-counter", err))
	}

	err = app.Database.Primary.IncChannelMessageCounter(ctx, &channel)
	if err != nil {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to inc channel msg-counter", err))
	}

	err = app.Database.Primary.IncKeyTokenMessageCounter(ctx, keytok)
	if err != nil {
		return nil, langext.Ptr(ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to inc token msg-counter", err))
	}

	log.Info().Msg(fmt.Sprintf("Sending new notification %s for user %s", msg.MessageID, UserID))

	for _, sub := range activeSubscriptions {
		clients, err := app.Database.Primary.ListClients(ctx, sub.SubscriberUserID)
		if err != nil {
			return nil, langext.Ptr(ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query clients", err))
		}

		for _, client := range clients {

			fcmDelivID, err := app.DeliverMessage(ctx, user, client, channel, msg)
			if err != nil {
				_, err = app.Database.Primary.CreateRetryDelivery(ctx, client, msg)
				if err != nil {
					return nil, langext.Ptr(ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to create delivery", err))
				}
			} else {
				_, err = app.Database.Primary.CreateSuccessDelivery(ctx, client, msg, fcmDelivID)
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
