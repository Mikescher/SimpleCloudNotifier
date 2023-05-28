package handler

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/api/ginresp"
	ct "blackforestbytes.com/simplecloudnotifier/db/cursortoken"
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"github.com/gin-gonic/gin"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/mathext"
	"net/http"
	"strings"
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
func (h APIHandler) ListMessages(g *gin.Context) ginresp.HTTPResponse {
	type query struct {
		PageSize      *int    `json:"page_size"       form:"page_size"`
		NextPageToken *string `json:"next_page_token" form:"next_page_token"`
		Filter        *string `json:"filter"          form:"filter"`
		Trimmed       *bool   `json:"trimmed"         form:"trimmed"` //TODO more filter (sender-name, channel, timestamps, prio, )
	}
	type response struct {
		Messages      []models.MessageJSON `json:"messages"`
		NextPageToken string               `json:"next_page_token"`
		PageSize      int                  `json:"page_size"`
	}

	var q query
	ctx, errResp := h.app.StartRequest(g, nil, &q, nil, nil)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

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

	messages, npt, err := h.database.ListMessages(ctx, filter, &pageSize, tok)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query messages", err)
	}

	var res []models.MessageJSON
	if trimmed {
		res = langext.ArrMap(messages, func(v models.Message) models.MessageJSON { return v.TrimmedJSON() })
	} else {
		res = langext.ArrMap(messages, func(v models.Message) models.MessageJSON { return v.FullJSON() })
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, response{Messages: res, NextPageToken: npt.Token(), PageSize: pageSize}))
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
//	@Param			mid	path		int	true	"MessageID"
//
//	@Success		200	{object}	models.MessageJSON
//	@Failure		400	{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure		401	{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure		404	{object}	ginresp.apiError	"message not found"
//	@Failure		500	{object}	ginresp.apiError	"internal server error"
//
//	@Router			/api/v2/messages/{mid} [PATCH]
func (h APIHandler) GetMessage(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		MessageID models.MessageID `uri:"mid" binding:"entityid"`
	}

	var u uri
	ctx, errResp := h.app.StartRequest(g, &u, nil, nil, nil)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	if permResp := ctx.CheckPermissionAny(); permResp != nil {
		return *permResp
	}

	msg, err := h.database.GetMessage(ctx, u.MessageID, false)
	if err == sql.ErrNoRows {
		return ginresp.APIError(g, 404, apierr.MESSAGE_NOT_FOUND, "message not found", err)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query message", err)
	}

	// either we have direct read permissions (it is our message + read/admin key)
	// or we subscribe (+confirmed) to the channel and have read/admin key

	if ctx.CheckPermissionMessageRead(msg) {
		return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, msg.FullJSON()))
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
		return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, msg.FullJSON()))
	}

	return ginresp.APIError(g, 401, apierr.USER_AUTH_FAILED, "You are not authorized for this action", nil)
}

// DeleteMessage swaggerdoc
//
//	@Summary		Delete a single message
//	@Description	The user must own the message and request the resource with the ADMIN Key
//	@ID				api-messages-delete
//	@Tags			API-v2
//
//	@Param			mid	path		int	true	"MessageID"
//
//	@Success		200	{object}	models.MessageJSON
//	@Failure		400	{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure		401	{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure		404	{object}	ginresp.apiError	"message not found"
//	@Failure		500	{object}	ginresp.apiError	"internal server error"
//
//	@Router			/api/v2/messages/{mid} [DELETE]
func (h APIHandler) DeleteMessage(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		MessageID models.MessageID `uri:"mid" binding:"entityid"`
	}

	var u uri
	ctx, errResp := h.app.StartRequest(g, &u, nil, nil, nil)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	if permResp := ctx.CheckPermissionAny(); permResp != nil {
		return *permResp
	}

	msg, err := h.database.GetMessage(ctx, u.MessageID, false)
	if err == sql.ErrNoRows {
		return ginresp.APIError(g, 404, apierr.MESSAGE_NOT_FOUND, "message not found", err)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query message", err)
	}

	if !ctx.CheckPermissionMessageRead(msg) {
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

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, msg.FullJSON()))
}
