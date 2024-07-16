package handler

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/api/ginresp"
	primarydb "blackforestbytes.com/simplecloudnotifier/db/impl/primary"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"blackforestbytes.com/simplecloudnotifier/models"
	"fmt"
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"net/http"
	"time"
)

type ExternalHandler struct {
	app      *logic.Application
	database *primarydb.Database
}

func NewExternalHandler(app *logic.Application) ExternalHandler {
	return ExternalHandler{
		app:      app,
		database: app.Database.Primary,
	}
}

// UptimeKuma swaggerdoc
//
//	@Summary		Send a new message
//	@Description	All parameter can be set via query-parameter or the json body. Only UserID, UserKey and Title are required
//	@Tags			External
//
//	@Param			query_data	query		handler.UptimeKuma.query	false	" "
//	@Param			post_body	body		handler.UptimeKuma.body		false	" "
//
//	@Success		200			{object}	handler.UptimeKuma.response
//	@Failure		400			{object}	ginresp.apiError
//	@Failure		401			{object}	ginresp.apiError	"The user_id was not found or the user_key is wrong"
//	@Failure		403			{object}	ginresp.apiError	"The user has exceeded its daily quota - wait 24 hours or upgrade your account"
//	@Failure		500			{object}	ginresp.apiError	"An internal server error occurred - try again later"
//
//	@Router			/external/v1/uptime-kuma [POST]
func (h ExternalHandler) UptimeKuma(pctx ginext.PreContext) ginext.HTTPResponse {
	type query struct {
		UserID       *models.UserID `form:"user_id"     example:"7725"`
		KeyToken     *string        `form:"key"         example:"P3TNH8mvv14fm"`
		Channel      *string        `form:"channel"`
		ChannelUp    *string        `form:"channel_up"`
		ChannelDown  *string        `form:"channel_down"`
		Priority     *int           `form:"priority"`
		PriorityUp   *int           `form:"priority_up"`
		PriorityDown *int           `form:"priority_down"`
		SenderName   *string        `form:"senderName"`
	}
	type body struct {
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
		Msg *string `json:"msg"`
	}
	type response struct {
		MessageID models.MessageID `json:"message_id"`
	}

	var b body
	var q query
	ctx, g, errResp := pctx.Query(&q).Body(&b).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	return h.app.DoRequest(ctx, g, func(ctx *logic.AppContext, finishSuccess func(r ginext.HTTPResponse) ginext.HTTPResponse) ginext.HTTPResponse {

		if b.Heartbeat == nil {
			return ginresp.APIError(g, 400, apierr.BINDFAIL_BODY_PARAM, "missing field 'heartbeat' in request body", nil)
		}
		if b.Monitor == nil {
			return ginresp.APIError(g, 400, apierr.BINDFAIL_BODY_PARAM, "missing field 'monitor' in request body", nil)
		}
		if b.Msg == nil {
			return ginresp.APIError(g, 400, apierr.BINDFAIL_BODY_PARAM, "missing field 'msg' in request body", nil)
		}

		title := langext.Conditional(b.Heartbeat.Status == 1, fmt.Sprintf("Monitor %v is back online", b.Monitor.Name), fmt.Sprintf("Monitor %v went down!", b.Monitor.Name))

		content := b.Heartbeat.Msg

		var timestamp *float64 = nil
		if tz, err := time.LoadLocation(b.Heartbeat.Timezone); err == nil {
			if ts, err := time.ParseInLocation("2006-01-02 15:04:05", b.Heartbeat.LocalDateTime, tz); err == nil {
				timestamp = langext.Ptr(float64(ts.Unix()))
			}
		}

		var channel *string = nil
		if q.Channel != nil {
			channel = q.Channel
		}
		if q.ChannelUp != nil && b.Heartbeat.Status == 1 {
			channel = q.ChannelUp
		}
		if q.ChannelDown != nil && b.Heartbeat.Status != 1 {
			channel = q.ChannelDown
		}

		var priority *int = nil
		if q.Priority != nil {
			priority = q.Priority
		}
		if q.PriorityUp != nil && b.Heartbeat.Status == 1 {
			priority = q.PriorityUp
		}
		if q.PriorityDown != nil && b.Heartbeat.Status != 1 {
			priority = q.PriorityDown
		}

		okResp, errResp := h.app.SendMessage(g, ctx, q.UserID, q.KeyToken, channel, &title, &content, priority, nil, timestamp, q.SenderName)
		if errResp != nil {
			return *errResp
		}

		return finishSuccess(ginext.JSON(http.StatusOK, response{
			MessageID: okResp.Message.MessageID,
		}))

	})
}
