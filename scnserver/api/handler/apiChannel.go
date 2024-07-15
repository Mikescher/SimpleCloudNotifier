package handler

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/api/ginresp"
	ct "blackforestbytes.com/simplecloudnotifier/db/cursortoken"
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"errors"
	"fmt"
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/mathext"
	"net/http"
	"strings"
)

// ListChannels swaggerdoc
//
//	@Summary		List channels of a user (subscribed/owned/all)
//	@Description	The possible values for 'selector' are:
//	@Description	- "owned"          Return all channels of the user
//	@Description	- "subscribed"     Return all channels that the user is subscribing to
//	@Description	- "all"            Return channels that the user owns or is subscribing
//	@Description	- "subscribed_any" Return all channels that the user is subscribing to (even unconfirmed)
//	@Description	- "all_any"        Return channels that the user owns or is subscribing (even unconfirmed)
//
//	@ID				api-channels-list
//	@Tags			API-v2
//
//	@Param			uid			path		string	true	"UserID"
//	@Param			selector	query		string	false	"Filter channels (default: owned)"	Enums(owned, subscribed, all, subscribed_any, all_any)
//
//	@Success		200			{object}	handler.ListChannels.response
//	@Failure		400			{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure		401			{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure		500			{object}	ginresp.apiError	"internal server error"
//
//	@Router			/api/v2/users/{uid}/channels [GET]
func (h APIHandler) ListChannels(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
		UserID models.UserID `uri:"uid" binding:"entityid"`
	}
	type query struct {
		Selector *string `json:"selector" form:"selector"  enums:"owned,subscribed_any,all_any,subscribed,all"`
	}
	type response struct {
		Channels []models.ChannelWithSubscriptionJSON `json:"channels"`
	}

	var u uri
	var q query
	ctx, g, errResp := pctx.URI(&u).Query(&q).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	ctx, errResp := h.app.StartRequest(g, &u, &q, nil, nil)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	if permResp := ctx.CheckPermissionUserRead(u.UserID); permResp != nil {
		return *permResp
	}

	sel := strings.ToLower(langext.Coalesce(q.Selector, "owned"))

	var res []models.ChannelWithSubscriptionJSON

	if sel == "owned" {

		channels, err := h.database.ListChannelsByOwner(ctx, u.UserID, u.UserID)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query channels", err)
		}
		res = langext.ArrMap(channels, func(v models.ChannelWithSubscription) models.ChannelWithSubscriptionJSON { return v.JSON(true) })

	} else if sel == "subscribed_any" {

		channels, err := h.database.ListChannelsBySubscriber(ctx, u.UserID, nil)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query channels", err)
		}
		res = langext.ArrMap(channels, func(v models.ChannelWithSubscription) models.ChannelWithSubscriptionJSON { return v.JSON(false) })

	} else if sel == "all_any" {

		channels, err := h.database.ListChannelsByAccess(ctx, u.UserID, nil)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query channels", err)
		}
		res = langext.ArrMap(channels, func(v models.ChannelWithSubscription) models.ChannelWithSubscriptionJSON { return v.JSON(false) })

	} else if sel == "subscribed" {

		channels, err := h.database.ListChannelsBySubscriber(ctx, u.UserID, langext.Ptr(true))
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query channels", err)
		}
		res = langext.ArrMap(channels, func(v models.ChannelWithSubscription) models.ChannelWithSubscriptionJSON { return v.JSON(false) })

	} else if sel == "all" {

		channels, err := h.database.ListChannelsByAccess(ctx, u.UserID, langext.Ptr(true))
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query channels", err)
		}
		res = langext.ArrMap(channels, func(v models.ChannelWithSubscription) models.ChannelWithSubscriptionJSON { return v.JSON(false) })

	} else {

		return ginresp.APIError(g, 400, apierr.INVALID_ENUM_VALUE, "Invalid value for the [selector] parameter", nil)

	}

	return ctx.FinishSuccess(ginext.JSON(http.StatusOK, response{Channels: res}))
}

// GetChannel swaggerdoc
//
//	@Summary	Get a single channel
//	@ID			api-channels-get
//	@Tags		API-v2
//
//	@Param		uid	path		string	true	"UserID"
//	@Param		cid	path		string	true	"ChannelID"
//
//	@Success	200	{object}	models.ChannelWithSubscriptionJSON
//	@Failure	400	{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure	401	{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure	404	{object}	ginresp.apiError	"channel not found"
//	@Failure	500	{object}	ginresp.apiError	"internal server error"
//
//	@Router		/api/v2/users/{uid}/channels/{cid} [GET]
func (h APIHandler) GetChannel(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
		UserID    models.UserID    `uri:"uid" binding:"entityid"`
		ChannelID models.ChannelID `uri:"cid" binding:"entityid"`
	}

	var u uri
	ctx, g, errResp := h.app.StartRequest(pctx.URI(&u).Start())
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	if permResp := ctx.CheckPermissionUserRead(u.UserID); permResp != nil {
		return *permResp
	}

	channel, err := h.database.GetChannel(ctx, u.UserID, u.ChannelID, true)
	if errors.Is(err, sql.ErrNoRows) {
		return ginresp.APIError(g, 404, apierr.CHANNEL_NOT_FOUND, "Channel not found", err)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query channel", err)
	}

	return ctx.FinishSuccess(ginext.JSON(http.StatusOK, channel.JSON(true)))
}

// CreateChannel swaggerdoc
//
//	@Summary	Create a new (empty) channel
//	@ID			api-channels-create
//	@Tags		API-v2
//
//	@Param		uid			path		string						true	"UserID"
//	@Param		post_body	body		handler.CreateChannel.body	false	" "
//
//	@Success	200			{object}	models.ChannelWithSubscriptionJSON
//	@Failure	400			{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure	401			{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure	409			{object}	ginresp.apiError	"channel already exists"
//	@Failure	500			{object}	ginresp.apiError	"internal server error"
//
//	@Router		/api/v2/users/{uid}/channels [POST]
func (h APIHandler) CreateChannel(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
		UserID models.UserID `uri:"uid" binding:"entityid"`
	}
	type body struct {
		Name      string `json:"name"`
		Subscribe *bool  `json:"subscribe"`

		Description *string `json:"description"`
	}

	var u uri
	var b body
	ctx, g, errResp := h.app.StartRequest(pctx.URI(&u).Body(&b).Start())
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	if permResp := ctx.CheckPermissionUserAdmin(u.UserID); permResp != nil {
		return *permResp
	}

	if b.Name == "" {
		return ginresp.APIError(g, 400, apierr.INVALID_BODY_PARAM, "Missing parameter: name", nil)
	}

	channelDisplayName := h.app.NormalizeChannelDisplayName(b.Name)
	channelInternalName := h.app.NormalizeChannelInternalName(b.Name)

	channelExisting, err := h.database.GetChannelByName(ctx, u.UserID, channelInternalName)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query channel", err)
	}

	user, err := h.database.GetUser(ctx, u.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		return ginresp.APIError(g, 400, apierr.USER_NOT_FOUND, "User not found", nil)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query user", err)
	}

	if len(channelDisplayName) > user.MaxChannelNameLength() {
		return ginresp.APIError(g, 400, apierr.CHANNEL_TOO_LONG, fmt.Sprintf("Channel too long (max %d characters)", user.MaxChannelNameLength()), nil)
	}
	if len(strings.TrimSpace(channelDisplayName)) == 0 {
		return ginresp.APIError(g, 400, apierr.CHANNEL_NAME_EMPTY, fmt.Sprintf("Channel displayname cannot be empty"), nil)
	}
	if len(channelInternalName) > user.MaxChannelNameLength() {
		return ginresp.APIError(g, 400, apierr.CHANNEL_TOO_LONG, fmt.Sprintf("Channel too long (max %d characters)", user.MaxChannelNameLength()), nil)
	}
	if len(strings.TrimSpace(channelInternalName)) == 0 {
		return ginresp.APIError(g, 400, apierr.CHANNEL_NAME_EMPTY, fmt.Sprintf("Channel internalname cannot be empty"), nil)
	}

	if channelExisting != nil {
		return ginresp.APIError(g, 409, apierr.CHANNEL_ALREADY_EXISTS, "Channel with this name already exists", nil)
	}

	subscribeKey := h.app.GenerateRandomAuthKey()

	channel, err := h.database.CreateChannel(ctx, u.UserID, channelDisplayName, channelInternalName, subscribeKey, b.Description)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to create channel", err)
	}

	if langext.Coalesce(b.Subscribe, true) {

		sub, err := h.database.CreateSubscription(ctx, u.UserID, channel, true)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to create subscription", err)
		}

		return ctx.FinishSuccess(ginext.JSON(http.StatusOK, channel.WithSubscription(langext.Ptr(sub)).JSON(true)))

	} else {

		return ctx.FinishSuccess(ginext.JSON(http.StatusOK, channel.WithSubscription(nil).JSON(true)))

	}

}

// UpdateChannel swaggerdoc
//
//	@Summary	(Partially) update a channel
//	@ID			api-channels-update
//	@Tags		API-v2
//
//	@Param		uid				path		string	true	"UserID"
//	@Param		cid				path		string	true	"ChannelID"
//
//	@Param		subscribe_key	body		string	false	"Send `true` to create a new subscribe_key"
//	@Param		send_key		body		string	false	"Send `true` to create a new send_key"
//	@Param		display_name	body		string	false	"Change the cahnnel display-name (only chnages to lowercase/uppercase are allowed - internal_name must stay the same)"
//
//	@Success	200				{object}	models.ChannelWithSubscriptionJSON
//	@Failure	400				{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure	401				{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure	404				{object}	ginresp.apiError	"channel not found"
//	@Failure	500				{object}	ginresp.apiError	"internal server error"
//
//	@Router		/api/v2/users/{uid}/channels/{cid} [PATCH]
func (h APIHandler) UpdateChannel(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
		UserID    models.UserID    `uri:"uid" binding:"entityid"`
		ChannelID models.ChannelID `uri:"cid" binding:"entityid"`
	}
	type body struct {
		RefreshSubscribeKey *bool   `json:"subscribe_key"`
		DisplayName         *string `json:"display_name"`
		DescriptionName     *string `json:"description_name"`
	}

	var u uri
	var b body
	ctx, g, errResp := h.app.StartRequest(pctx.URI(&u).Body(&b).Start())
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	if permResp := ctx.CheckPermissionUserAdmin(u.UserID); permResp != nil {
		return *permResp
	}

	_, err := h.database.GetChannel(ctx, u.UserID, u.ChannelID, true)
	if errors.Is(err, sql.ErrNoRows) {
		return ginresp.APIError(g, 404, apierr.CHANNEL_NOT_FOUND, "Channel not found", err)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query channel", err)
	}

	user, err := h.database.GetUser(ctx, u.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		return ginresp.APIError(g, 400, apierr.USER_NOT_FOUND, "User not found", nil)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query user", err)
	}

	if langext.Coalesce(b.RefreshSubscribeKey, false) {
		newkey := h.app.GenerateRandomAuthKey()

		err := h.database.UpdateChannelSubscribeKey(ctx, u.ChannelID, newkey)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to update channel", err)
		}
	}

	if b.DisplayName != nil {

		newDisplayName := h.app.NormalizeChannelDisplayName(*b.DisplayName)

		if len(newDisplayName) > user.MaxChannelNameLength() {
			return ginresp.APIError(g, 400, apierr.CHANNEL_TOO_LONG, fmt.Sprintf("Channel too long (max %d characters)", user.MaxChannelNameLength()), nil)
		}

		if len(strings.TrimSpace(newDisplayName)) == 0 {
			return ginresp.APIError(g, 400, apierr.CHANNEL_NAME_EMPTY, fmt.Sprintf("Channel displayname cannot be empty"), nil)
		}

		err := h.database.UpdateChannelDisplayName(ctx, u.ChannelID, newDisplayName)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to update channel", err)
		}

	}

	if b.DescriptionName != nil {

		var descName *string = nil
		if strings.TrimSpace(*b.DescriptionName) != "" {
			descName = langext.Ptr(strings.TrimSpace(*b.DescriptionName))
		}

		if descName != nil && len(*descName) > user.MaxChannelDescriptionLength() {
			return ginresp.APIError(g, 400, apierr.CHANNEL_DESCRIPTION_TOO_LONG, fmt.Sprintf("Channel-Description too long (max %d characters)", user.MaxChannelDescriptionLength()), nil)
		}

		err := h.database.UpdateChannelDescriptionName(ctx, u.ChannelID, descName)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to update channel", err)
		}

	}

	channel, err := h.database.GetChannel(ctx, u.UserID, u.ChannelID, true)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query (updated) channel", err)
	}

	return ctx.FinishSuccess(ginext.JSON(http.StatusOK, channel.JSON(true)))
}

// ListChannelMessages swaggerdoc
//
//	@Summary		List messages of a channel
//	@Description	The next_page_token is an opaque token, the special value "@start" (or empty-string) is the beginning and "@end" is the end
//	@Description	Simply start the pagination without a next_page_token and get the next page by calling this endpoint with the returned next_page_token of the last query
//	@Description	If there are no more entries the token "@end" will be returned
//	@Description	By default we return long messages with a trimmed body, if trimmed=false is supplied we return full messages (this reduces the max page_size)
//	@ID				api-channel-messages
//	@Tags			API-v2
//
//	@Param			query_data	query		handler.ListChannelMessages.query	false	" "
//	@Param			uid			path		string								true	"UserID"
//	@Param			cid			path		string								true	"ChannelID"
//
//	@Success		200			{object}	handler.ListChannelMessages.response
//	@Failure		400			{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure		401			{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure		404			{object}	ginresp.apiError	"channel not found"
//	@Failure		500			{object}	ginresp.apiError	"internal server error"
//
//	@Router			/api/v2/users/{uid}/channels/{cid}/messages [GET]
func (h APIHandler) ListChannelMessages(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
		ChannelUserID models.UserID    `uri:"uid" binding:"entityid"`
		ChannelID     models.ChannelID `uri:"cid" binding:"entityid"`
	}
	type query struct {
		PageSize      *int    `json:"page_size"       form:"page_size"`
		NextPageToken *string `json:"next_page_token" form:"next_page_token"`
		Filter        *string `json:"filter"          form:"filter"`
		Trimmed       *bool   `json:"trimmed"         form:"trimmed"`
	}
	type response struct {
		Messages      []models.MessageJSON `json:"messages"`
		NextPageToken string               `json:"next_page_token"`
		PageSize      int                  `json:"page_size"`
	}

	var u uri
	var q query
	ctx, g, errResp := h.app.StartRequest(pctx.URI(&u).Query(&q).Start())
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	trimmed := langext.Coalesce(q.Trimmed, true)

	maxPageSize := langext.Conditional(trimmed, 16, 256)

	pageSize := mathext.Clamp(langext.Coalesce(q.PageSize, 64), 1, maxPageSize)

	channel, err := h.database.GetChannel(ctx, u.ChannelUserID, u.ChannelID, false)
	if errors.Is(err, sql.ErrNoRows) {
		return ginresp.APIError(g, 404, apierr.CHANNEL_NOT_FOUND, "Channel not found", err)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query channel", err)
	}

	if permResp := ctx.CheckPermissionChanMessagesRead(channel.Channel); permResp != nil {
		return *permResp
	}

	tok, err := ct.Decode(langext.Coalesce(q.NextPageToken, ""))
	if err != nil {
		return ginresp.APIError(g, 400, apierr.PAGETOKEN_ERROR, "Failed to decode next_page_token", err)
	}

	filter := models.MessageFilter{
		ChannelID: langext.Ptr([]models.ChannelID{channel.ChannelID}),
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

	return ctx.FinishSuccess(ginext.JSON(http.StatusOK, response{Messages: res, NextPageToken: npt.Token(), PageSize: pageSize}))
}
