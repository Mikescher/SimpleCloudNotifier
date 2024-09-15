package handler

import (
	"blackforestbytes.com/simplecloudnotifier/logic"
	"database/sql"
	"errors"
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
	"net/http"
	"strings"
	"time"

	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/api/ginresp"
	ct "blackforestbytes.com/simplecloudnotifier/db/cursortoken"
	"blackforestbytes.com/simplecloudnotifier/models"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/mathext"
)

// ListMessages swaggerdoc
//
//	@Summary		List all (subscribed) messages
//	@Description	The next_page_token is an opaque token, the special value "@start" (or empty-string) is the beginning and "@end" is the end
//	@Description	Simply start the pagination without a next_page_token and get the next page by calling this endpoint with the returned next_page_token of the last query
//	@Description	If there are no more entries the token "@end" will be returned
//	@Description	By default we return long messages with a trimmed body, if trimmed=false is supplied we return full messages (this reduces the max page_size)
//	@ID				api-messages-list
//	@Tags			API-v2
//
//	@Param			query_data	query		handler.ListMessages.query	false	" "
//
//	@Success		200			{object}	handler.ListMessages.response
//	@Failure		400			{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure		401			{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure		500			{object}	ginresp.apiError	"internal server error"
//
//	@Router			/api/v2/messages [GET]
func (h APIHandler) ListMessages(pctx ginext.PreContext) ginext.HTTPResponse {
	type query struct {
		PageSize      *int     `json:"page_size"       form:"page_size"`
		NextPageToken *string  `json:"next_page_token" form:"next_page_token"`
		Filter        *string  `json:"filter"          form:"filter"`
		Trimmed       *bool    `json:"trimmed"         form:"trimmed"`
		Channels      []string `json:"channel"         form:"channel"`
		ChannelIDs    []string `json:"channel_id"      form:"channel_id"`
		Senders       []string `json:"sender"          form:"sender"`
		TimeBefore    *string  `json:"before"          form:"before"` // RFC3339
		TimeAfter     *string  `json:"after"           form:"after"`  // RFC3339
		Priority      []int    `json:"priority"        form:"priority"`
		KeyTokens     []string `json:"used_key"        form:"used_key"`
	}
	type response struct {
		Messages      []models.Message `json:"messages"`
		NextPageToken string           `json:"next_page_token"`
		PageSize      int              `json:"page_size"`
	}

	var q query
	ctx, g, errResp := pctx.Query(&q).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	return h.app.DoRequest(ctx, g, func(ctx *logic.AppContext, finishSuccess func(r ginext.HTTPResponse) ginext.HTTPResponse) ginext.HTTPResponse {

		trimmed := langext.Coalesce(q.Trimmed, true)

		maxPageSize := langext.Conditional(trimmed, 16, 256)

		pageSize := mathext.Clamp(langext.Coalesce(q.PageSize, 64), 1, maxPageSize)

		if permResp := ctx.CheckPermissionSelfAllMessagesRead(); permResp != nil {
			return *permResp
		}

		userid := *ctx.GetPermissionUserID()

		tok, err := ct.Decode(langext.Coalesce(q.NextPageToken, ""))
		if err != nil {
			return ginresp.APIError(g, 400, apierr.PAGETOKEN_ERROR, "Failed to decode next_page_token", err)
		}

		err = h.database.UpdateUserLastRead(ctx, userid)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to update last-read", err)
		}

		filter := models.MessageFilter{
			ConfirmedSubscriptionBy: langext.Ptr(userid),
		}

		if q.Filter != nil && strings.TrimSpace(*q.Filter) != "" {
			filter.SearchString = langext.Ptr([]string{strings.TrimSpace(*q.Filter)})
		}

		if len(q.Channels) != 0 {
			filter.ChannelNameCS = langext.Ptr(q.Channels)
		}

		if len(q.ChannelIDs) != 0 {
			cids := make([]models.ChannelID, 0, len(q.ChannelIDs))
			for _, v := range q.ChannelIDs {
				cid := models.ChannelID(v)
				if err = cid.Valid(); err != nil {
					return ginresp.APIError(g, 400, apierr.BINDFAIL_QUERY_PARAM, "Invalid channel-id", err)
				}
				cids = append(cids, cid)
			}
			filter.ChannelID = &cids
		}

		if len(q.Senders) != 0 {
			filter.SenderNameCS = langext.Ptr(q.Senders)
		}

		if q.TimeBefore != nil {
			t0, err := time.Parse(time.RFC3339, *q.TimeBefore)
			if err != nil {
				return ginresp.APIError(g, 400, apierr.BINDFAIL_QUERY_PARAM, "Invalid before-time", err)
			}
			filter.TimestampCoalesceBefore = &t0
		}

		if q.TimeAfter != nil {
			t0, err := time.Parse(time.RFC3339, *q.TimeAfter)
			if err != nil {
				return ginresp.APIError(g, 400, apierr.BINDFAIL_QUERY_PARAM, "Invalid after-time", err)
			}
			filter.TimestampCoalesceAfter = &t0
		}

		if len(q.Priority) != 0 {
			filter.Priority = langext.Ptr(q.Priority)
		}

		if len(q.KeyTokens) != 0 {
			tids := make([]models.KeyTokenID, 0, len(q.KeyTokens))
			for _, v := range q.KeyTokens {
				tid := models.KeyTokenID(v)
				if err = tid.Valid(); err != nil {
					return ginresp.APIError(g, 400, apierr.BINDFAIL_QUERY_PARAM, "Invalid keytoken-id", err)
				}
				tids = append(tids, tid)
			}
			filter.UsedKeyID = &tids
		}

		messages, npt, err := h.database.ListMessages(ctx, filter, &pageSize, tok)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query messages", err)
		}

		if trimmed {
			res := langext.ArrMap(messages, func(v models.Message) models.Message { return v.PreMarshal().Trim() })
			return finishSuccess(ginext.JSON(http.StatusOK, response{Messages: res, NextPageToken: npt.Token(), PageSize: pageSize}))
		} else {
			res := langext.ArrMap(messages, func(v models.Message) models.Message { return v.PreMarshal() })
			return finishSuccess(ginext.JSON(http.StatusOK, response{Messages: res, NextPageToken: npt.Token(), PageSize: pageSize}))
		}
	})
}

// GetMessage swaggerdoc
//
//	@Summary		Get a single message (untrimmed)
//	@Description	The user must either own the message and request the resource with the READ or ADMIN Key
//	@Description	Or the user must subscribe to the corresponding channel (and be confirmed) and request the resource with the READ or ADMIN Key
//	@Description	The returned message is never trimmed
//	@ID				api-messages-get
//	@Tags			API-v2
//
//	@Param			mid	path		string	true	"MessageID"
//
//	@Success		200	{object}	models.Message
//	@Failure		400	{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure		401	{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure		404	{object}	ginresp.apiError	"message not found"
//	@Failure		500	{object}	ginresp.apiError	"internal server error"
//
//	@Router			/api/v2/messages/{mid} [GET]
func (h APIHandler) GetMessage(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
		MessageID models.MessageID `uri:"mid" binding:"entityid"`
	}

	var u uri
	ctx, g, errResp := pctx.URI(&u).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	return h.app.DoRequest(ctx, g, func(ctx *logic.AppContext, finishSuccess func(r ginext.HTTPResponse) ginext.HTTPResponse) ginext.HTTPResponse {

		if permResp := ctx.CheckPermissionAny(); permResp != nil {
			return *permResp
		}

		msg, err := h.database.GetMessage(ctx, u.MessageID, false)
		if errors.Is(err, sql.ErrNoRows) {
			return ginresp.APIError(g, 404, apierr.MESSAGE_NOT_FOUND, "message not found", err)
		}
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query message", err)
		}

		// either we have direct read permissions (it is our message + read/admin key)
		// or we subscribe (+confirmed) to the channel and have read/admin key

		if ctx.CheckPermissionMessageRead(msg) {
			return finishSuccess(ginext.JSON(http.StatusOK, msg.PreMarshal()))
		}

		if uid := ctx.GetPermissionUserID(); uid != nil && ctx.CheckPermissionUserRead(*uid) == nil {
			sub, err := h.database.GetSubscriptionBySubscriber(ctx, *uid, msg.ChannelID)
			if err != nil {
				return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query subscription", err)
			}
			if sub == nil {
				// not subbed
				return ginresp.APIError(g, 401, apierr.USER_AUTH_FAILED, "You are not authorized for this action", nil)
			}
			if !sub.Confirmed {
				// sub not confirmed
				return ginresp.APIError(g, 401, apierr.USER_AUTH_FAILED, "You are not authorized for this action", nil)
			}

			// => perm okay
			return finishSuccess(ginext.JSON(http.StatusOK, msg.PreMarshal()))
		}

		return ginresp.APIError(g, 401, apierr.USER_AUTH_FAILED, "You are not authorized for this action", nil)

	})
}

// DeleteMessage swaggerdoc
//
//	@Summary		Delete a single message
//	@Description	The user must own the message and request the resource with the ADMIN Key
//	@ID				api-messages-delete
//	@Tags			API-v2
//
//	@Param			mid	path		string	true	"MessageID"
//
//	@Success		200	{object}	models.Message
//	@Failure		400	{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure		401	{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure		404	{object}	ginresp.apiError	"message not found"
//	@Failure		500	{object}	ginresp.apiError	"internal server error"
//
//	@Router			/api/v2/messages/{mid} [DELETE]
func (h APIHandler) DeleteMessage(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
		MessageID models.MessageID `uri:"mid" binding:"entityid"`
	}

	var u uri
	ctx, g, errResp := pctx.URI(&u).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	return h.app.DoRequest(ctx, g, func(ctx *logic.AppContext, finishSuccess func(r ginext.HTTPResponse) ginext.HTTPResponse) ginext.HTTPResponse {

		if permResp := ctx.CheckPermissionAny(); permResp != nil {
			return *permResp
		}

		msg, err := h.database.GetMessage(ctx, u.MessageID, false)
		if errors.Is(err, sql.ErrNoRows) {
			return ginresp.APIError(g, 404, apierr.MESSAGE_NOT_FOUND, "message not found", err)
		}
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query message", err)
		}

		if !ctx.CheckPermissionMessageDelete(msg) {
			return ginresp.APIError(g, 401, apierr.USER_AUTH_FAILED, "You are not authorized for this action", nil)
		}

		err = h.database.DeleteMessage(ctx, msg.MessageID)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to delete message", err)
		}

		err = h.database.CancelPendingDeliveries(ctx, msg.MessageID)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to cancel deliveries", err)
		}

		return finishSuccess(ginext.JSON(http.StatusOK, msg.PreMarshal()))

	})
}
