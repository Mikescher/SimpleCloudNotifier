package handler

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/api/ginresp"
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"github.com/gin-gonic/gin"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"net/http"
	"strings"
)

// ListUserSubscriptions swaggerdoc
//
//	@Summary		List all subscriptions of a user (incoming/owned)
//	@Description	The possible values for 'selector' are:
//	@Description	- "outgoing_all"         All subscriptions (confirmed/unconfirmed) with the user as subscriber (= subscriptions he can use to read channels)
//	@Description	- "outgoing_confirmed"   Confirmed subscriptions with the user as subscriber
//	@Description	- "outgoing_unconfirmed" Unconfirmed (Pending) subscriptions with the user as subscriber
//	@Description	- "incoming_all"         All subscriptions (confirmed/unconfirmed) from other users to channels of this user (= incoming subscriptions and subscription requests)
//	@Description	- "incoming_confirmed"   Confirmed subscriptions from other users to channels of this user
//	@Description	- "incoming_unconfirmed" Unconfirmed subscriptions from other users to channels of this user (= requests)
//
//	@ID				api-user-subscriptions-list
//	@Tags			API-v2
//
//	@Param			uid			path		int		true	"UserID"
//	@Param			selector	query		string	true	"Filter subscriptions (default: outgoing_all)"	Enums(outgoing_all, outgoing_confirmed, outgoing_unconfirmed, incoming_all, incoming_confirmed, incoming_unconfirmed)
//
//	@Success		200			{object}	handler.ListUserSubscriptions.response
//	@Failure		400			{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure		401			{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure		500			{object}	ginresp.apiError	"internal server error"
//
//	@Router			/api/v2/users/{uid}/subscriptions [GET]
func (h APIHandler) ListUserSubscriptions(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID models.UserID `uri:"uid" binding:"entityid"`
	}
	type query struct {
		Selector *string `json:"selector" form:"selector"  enums:"outgoing_all,outgoing_confirmed,outgoing_unconfirmed,incoming_all,incoming_confirmed,incoming_unconfirmed"`
	}
	type response struct {
		Subscriptions []models.SubscriptionJSON `json:"subscriptions"`
	}

	var u uri
	var q query
	ctx, errResp := h.app.StartRequest(g, &u, &q, nil, nil)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	if permResp := ctx.CheckPermissionUserRead(u.UserID); permResp != nil {
		return *permResp
	}

	sel := strings.ToLower(langext.Coalesce(q.Selector, "outgoing_all"))

	var res []models.Subscription
	var err error

	if sel == "outgoing_all" {

		res, err = h.database.ListSubscriptionsBySubscriber(ctx, u.UserID, nil)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query subscriptions", err)
		}

	} else if sel == "outgoing_confirmed" {

		res, err = h.database.ListSubscriptionsBySubscriber(ctx, u.UserID, langext.Ptr(true))
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query subscriptions", err)
		}

	} else if sel == "outgoing_unconfirmed" {

		res, err = h.database.ListSubscriptionsBySubscriber(ctx, u.UserID, langext.Ptr(false))
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query subscriptions", err)
		}

	} else if sel == "incoming_all" {

		res, err = h.database.ListSubscriptionsByChannelOwner(ctx, u.UserID, nil)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query subscriptions", err)
		}

	} else if sel == "incoming_confirmed" {

		res, err = h.database.ListSubscriptionsByChannelOwner(ctx, u.UserID, langext.Ptr(true))
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query subscriptions", err)
		}

	} else if sel == "incoming_unconfirmed" {

		res, err = h.database.ListSubscriptionsByChannelOwner(ctx, u.UserID, langext.Ptr(false))
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query subscriptions", err)
		}

	} else {

		return ginresp.APIError(g, 400, apierr.INVALID_ENUM_VALUE, "Invalid value for the [selector] parameter", nil)

	}

	jsonres := langext.ArrMap(res, func(v models.Subscription) models.SubscriptionJSON { return v.JSON() })

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, response{Subscriptions: jsonres}))
}

// ListChannelSubscriptions swaggerdoc
//
//	@Summary	List all subscriptions of a channel
//	@ID			api-chan-subscriptions-list
//	@Tags		API-v2
//
//	@Param		uid	path		string	true	"UserID"
//	@Param		cid	path		string	true	"ChannelID"
//
//	@Success	200	{object}	handler.ListChannelSubscriptions.response
//	@Failure	400	{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure	401	{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure	404	{object}	ginresp.apiError	"channel not found"
//	@Failure	500	{object}	ginresp.apiError	"internal server error"
//
//	@Router		/api/v2/users/{uid}/channels/{cid}/subscriptions [GET]
func (h APIHandler) ListChannelSubscriptions(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID    models.UserID    `uri:"uid" binding:"entityid"`
		ChannelID models.ChannelID `uri:"cid" binding:"entityid"`
	}
	type response struct {
		Subscriptions []models.SubscriptionJSON `json:"subscriptions"`
	}

	var u uri
	ctx, errResp := h.app.StartRequest(g, &u, nil, nil, nil)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	if permResp := ctx.CheckPermissionUserRead(u.UserID); permResp != nil {
		return *permResp
	}

	_, err := h.database.GetChannel(ctx, u.UserID, u.ChannelID, true)
	if err == sql.ErrNoRows {
		return ginresp.APIError(g, 404, apierr.CHANNEL_NOT_FOUND, "Channel not found", err)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query channel", err)
	}

	clients, err := h.database.ListSubscriptionsByChannel(ctx, u.ChannelID)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query subscriptions", err)
	}

	res := langext.ArrMap(clients, func(v models.Subscription) models.SubscriptionJSON { return v.JSON() })

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, response{Subscriptions: res}))
}

// GetSubscription swaggerdoc
//
//	@Summary	Get a single subscription
//	@ID			api-subscriptions-get
//	@Tags		API-v2
//
//	@Param		uid	path		string	true	"UserID"
//	@Param		sid	path		string	true	"SubscriptionID"
//
//	@Success	200	{object}	models.SubscriptionJSON
//	@Failure	400	{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure	401	{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure	404	{object}	ginresp.apiError	"subscription not found"
//	@Failure	500	{object}	ginresp.apiError	"internal server error"
//
//	@Router		/api/v2/users/{uid}/subscriptions/{sid} [GET]
func (h APIHandler) GetSubscription(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID         models.UserID         `uri:"uid" binding:"entityid"`
		SubscriptionID models.SubscriptionID `uri:"sid" binding:"entityid"`
	}

	var u uri
	ctx, errResp := h.app.StartRequest(g, &u, nil, nil, nil)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	if permResp := ctx.CheckPermissionUserRead(u.UserID); permResp != nil {
		return *permResp
	}

	subscription, err := h.database.GetSubscription(ctx, u.SubscriptionID)
	if err == sql.ErrNoRows {
		return ginresp.APIError(g, 404, apierr.SUBSCRIPTION_NOT_FOUND, "Subscription not found", err)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query subscription", err)
	}
	if subscription.SubscriberUserID != u.UserID && subscription.ChannelOwnerUserID != u.UserID {
		return ginresp.APIError(g, 404, apierr.SUBSCRIPTION_USER_MISMATCH, "Subscription not found", nil)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, subscription.JSON()))
}

// CancelSubscription swaggerdoc
//
//	@Summary	Cancel (delete) subscription
//	@ID			api-subscriptions-delete
//	@Tags		API-v2
//
//	@Param		uid	path		string	true	"UserID"
//	@Param		sid	path		string	true	"SubscriptionID"
//
//	@Success	200	{object}	models.SubscriptionJSON
//	@Failure	400	{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure	401	{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure	404	{object}	ginresp.apiError	"subscription not found"
//	@Failure	500	{object}	ginresp.apiError	"internal server error"
//
//	@Router		/api/v2/users/{uid}/subscriptions/{sid} [DELETE]
func (h APIHandler) CancelSubscription(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID         models.UserID         `uri:"uid" binding:"entityid"`
		SubscriptionID models.SubscriptionID `uri:"sid" binding:"entityid"`
	}

	var u uri
	ctx, errResp := h.app.StartRequest(g, &u, nil, nil, nil)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	if permResp := ctx.CheckPermissionUserAdmin(u.UserID); permResp != nil {
		return *permResp
	}

	subscription, err := h.database.GetSubscription(ctx, u.SubscriptionID)
	if err == sql.ErrNoRows {
		return ginresp.APIError(g, 404, apierr.SUBSCRIPTION_NOT_FOUND, "Subscription not found", err)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query subscription", err)
	}
	if subscription.SubscriberUserID != u.UserID && subscription.ChannelOwnerUserID != u.UserID {
		return ginresp.APIError(g, 404, apierr.SUBSCRIPTION_USER_MISMATCH, "Subscription not found", nil)
	}

	err = h.database.DeleteSubscription(ctx, u.SubscriptionID)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to delete subscription", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, subscription.JSON()))
}

// CreateSubscription swaggerdoc
//
//	@Summary		Create/Request a subscription
//	@Description	Either [channel_owner_user_id, channel_internal_name] or [channel_id] must be supplied in the request body
//	@ID				api-subscriptions-create
//	@Tags			API-v2
//
//	@Param			uid			path		int									true	"UserID"
//	@Param			query_data	query		handler.CreateSubscription.query	false	" "
//	@Param			post_data	body		handler.CreateSubscription.body		false	" "
//
//	@Success		200			{object}	models.SubscriptionJSON
//	@Failure		400			{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure		401			{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure		500			{object}	ginresp.apiError	"internal server error"
//
//	@Router			/api/v2/users/{uid}/subscriptions [POST]
func (h APIHandler) CreateSubscription(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID models.UserID `uri:"uid" binding:"entityid"`
	}
	type body struct {
		ChannelOwnerUserID  *models.UserID    `json:"channel_owner_user_id" binding:"entityid"`
		ChannelInternalName *string           `json:"channel_internal_name"`
		ChannelID           *models.ChannelID `json:"channel_id" binding:"entityid"`
	}
	type query struct {
		ChanSubscribeKey *string `json:"chan_subscribe_key" form:"chan_subscribe_key"`
	}

	var u uri
	var q query
	var b body
	ctx, errResp := h.app.StartRequest(g, &u, &q, &b, nil)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	if permResp := ctx.CheckPermissionUserAdmin(u.UserID); permResp != nil {
		return *permResp
	}

	var channel models.Channel

	if b.ChannelOwnerUserID != nil && b.ChannelInternalName != nil && b.ChannelID == nil {

		channelInternalName := h.app.NormalizeChannelInternalName(*b.ChannelInternalName)

		outchannel, err := h.database.GetChannelByName(ctx, *b.ChannelOwnerUserID, channelInternalName)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query channel", err)
		}
		if outchannel == nil {
			return ginresp.APIError(g, 400, apierr.CHANNEL_NOT_FOUND, "Channel not found", err)
		}

		channel = *outchannel

	} else if b.ChannelOwnerUserID == nil && b.ChannelInternalName == nil && b.ChannelID != nil {

		outchannel, err := h.database.GetChannelByID(ctx, *b.ChannelID)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query channel", err)
		}
		if outchannel == nil {
			return ginresp.APIError(g, 400, apierr.CHANNEL_NOT_FOUND, "Channel not found", err)
		}

		channel = *outchannel

	} else {

		return ginresp.APIError(g, 400, apierr.INVALID_BODY_PARAM, "Must either supply [channel_owner_user_id, channel_internal_name] or [channel_id]", nil)

	}

	if channel.OwnerUserID != u.UserID && (q.ChanSubscribeKey == nil || *q.ChanSubscribeKey != channel.SubscribeKey) {
		return ginresp.APIError(g, 401, apierr.USER_AUTH_FAILED, "You are not authorized for this action", nil)
	}

	sub, err := h.database.CreateSubscription(ctx, u.UserID, channel, channel.OwnerUserID == u.UserID)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to create subscription", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, sub.JSON()))
}

// UpdateSubscription swaggerdoc
//
//	@Summary	Update a subscription (e.g. confirm)
//	@ID			api-subscriptions-update
//	@Tags		API-v2
//
//	@Param		uid			path		int								true	"UserID"
//	@Param		sid			path		int								true	"SubscriptionID"
//	@Param		post_data	body		handler.UpdateSubscription.body	false	" "
//
//	@Success	200			{object}	models.SubscriptionJSON
//	@Failure	400			{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure	401			{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure	404			{object}	ginresp.apiError	"subscription not found"
//	@Failure	500			{object}	ginresp.apiError	"internal server error"
//
//	@Router		/api/v2/users/{uid}/subscriptions/{sid} [PATCH]
func (h APIHandler) UpdateSubscription(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID         models.UserID         `uri:"uid" binding:"entityid"`
		SubscriptionID models.SubscriptionID `uri:"sid" binding:"entityid"`
	}
	type body struct {
		Confirmed *bool `form:"confirmed"`
	}

	var u uri
	var b body
	ctx, errResp := h.app.StartRequest(g, &u, nil, &b, nil)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	if permResp := ctx.CheckPermissionUserAdmin(u.UserID); permResp != nil {
		return *permResp
	}

	userid := *ctx.GetPermissionUserID()

	subscription, err := h.database.GetSubscription(ctx, u.SubscriptionID)
	if err == sql.ErrNoRows {
		return ginresp.APIError(g, 404, apierr.SUBSCRIPTION_NOT_FOUND, "Subscription not found", err)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query subscription", err)
	}
	if subscription.SubscriberUserID != u.UserID && subscription.ChannelOwnerUserID != u.UserID {
		return ginresp.APIError(g, 404, apierr.SUBSCRIPTION_USER_MISMATCH, "Subscription not found", nil)
	}

	if b.Confirmed != nil {
		if subscription.ChannelOwnerUserID != userid {
			return ginresp.APIError(g, 401, apierr.USER_AUTH_FAILED, "You are not authorized for this action", nil)
		}
		err = h.database.UpdateSubscriptionConfirmed(ctx, u.SubscriptionID, *b.Confirmed)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to update subscription", err)
		}
	}

	subscription, err = h.database.GetSubscription(ctx, u.SubscriptionID)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query subscription", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, subscription.JSON()))
}
