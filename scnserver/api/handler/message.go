package handler

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	primarydb "blackforestbytes.com/simplecloudnotifier/db/impl/primary"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"blackforestbytes.com/simplecloudnotifier/models"
	"gogs.mikescher.com/BlackForestBytes/goext/dataext"
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"net/http"
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
func (h MessageHandler) SendMessage(pctx ginext.PreContext) ginext.HTTPResponse {
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
	ctx, g, errResp := pctx.Form(&f).Query(&q).Body(&b).IgnoreWrongContentType().Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	return h.app.DoRequest(ctx, g, models.TLockReadWrite, func(ctx *logic.AppContext, finishSuccess func(r ginext.HTTPResponse) ginext.HTTPResponse) ginext.HTTPResponse {

		// query has highest prio, then form, then json
		data := dataext.ObjectMerge(dataext.ObjectMerge(b, f), q)

		okResp, errResp := h.app.SendMessage(g, ctx, data.UserID, data.KeyToken, data.Channel, data.Title, data.Content, data.Priority, data.UserMessageID, data.SendTimestamp, data.SenderName)
		if errResp != nil {
			return *errResp
		} else {
			return finishSuccess(ginext.JSON(http.StatusOK, response{
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

	})
}
