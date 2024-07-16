package handler

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/api/ginresp"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"errors"
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"net/http"
	"strings"
)

// ListUserSubscriptions swaggerdoc
//
//	@Summary		List all subscriptions of a user (incoming/owned)
//
//	@Description	The possible values for 'direction' are:
//	@Description	- "outgoing"       Subscriptions with the user as subscriber (= subscriptions he can use to read channels)
//	@Description	- "incoming"       Subscriptions to channels of this user (= incoming subscriptions and subscription requests)
//	@Description	- "both"           Combines "outgoing" and "incoming" (default)
//	@Description
//	@Description	The possible values for 'confirmation' are:
//	@Description	- "confirmed"      Confirmed (active) subscriptions
//	@Description	- "unconfirmed"    Unconfirmed (pending) subscriptions
//	@Description	- "all"            Combines "confirmed" and "unconfirmed" (default)
//	@Description
//	@Description	The possible values for 'external' are:
//	@Description	- "true"           Subscriptions with subscriber_user_id != channel_owner_user_id  (subscriptions from other users)
//	@Description	- "false"          Subscriptions with subscriber_user_id == channel_owner_user_id  (subscriptions from this user to his own channels)
//	@Description	- "all"            Combines "external" and "internal" (default)
//	@Description
//	@Description	The `subscriber_user_id` parameter can be used to additionally filter the subscriber_user_id (return subscribtions from a specific user)
//	@Description
//	@Description	The `channel_owner_user_id` parameter can be used to additionally filter the channel_owner_user_id (return subscribtions to a specific user)
//
//	@ID				api-user-subscriptions-list
//	@Tags			API-v2
//
//	@Param			uid			path		string	true	"UserID"
//	@Param			selector	query		string	true	"Filter subscriptions (default: outgoing_all)"	Enums(outgoing_all, outgoing_confirmed, outgoing_unconfirmed, incoming_all, incoming_confirmed, incoming_unconfirmed)
//
//	@Success		200			{object}	handler.ListUserSubscriptions.response
//	@Failure		400			{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure		401			{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure		500			{object}	ginresp.apiError	"internal server error"
//
//	@Router			/api/v2/users/{uid}/subscriptions [GET]
func (h APIHandler) ListUserSubscriptions(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
		UserID models.UserID `uri:"uid" binding:"entityid"`
	}
	type query struct {
		Direction          *string        `json:"direction"             form:"direction"                enums:"incoming,outgoing,both"`
		Confirmation       *string        `json:"confirmation"          form:"confirmation"             enums:"confirmed,unconfirmed,all"`
		External           *string        `json:"external"              form:"external"                 enums:"true,false,all"`
		SubscriberUserID   *models.UserID `json:"subscriber_user_id"    form:"subscriber_user_id"`
		ChannelOwnerUserID *models.UserID `json:"channel_owner_user_id" form:"channel_owner_user_id"`
	}
	type response struct {
		Subscriptions []models.SubscriptionJSON `json:"subscriptions"`
	}

	var u uri
	var q query
	ctx, g, errResp := pctx.URI(&u).Query(&q).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	return h.app.DoRequest(ctx, g, func(ctx *logic.AppContext, finishSuccess func(r ginext.HTTPResponse) ginext.HTTPResponse) ginext.HTTPResponse {

		if permResp := ctx.CheckPermissionUserRead(u.UserID); permResp != nil {
			return *permResp
		}

		filter := models.SubscriptionFilter{}
		filter.AnyUserID = langext.Ptr(u.UserID)

		if q.Direction != nil {
			if strings.EqualFold(*q.Direction, "incoming") {
				filter.ChannelOwnerUserID = langext.Ptr([]models.UserID{u.UserID})
			} else if strings.EqualFold(*q.Direction, "outgoing") {
				filter.SubscriberUserID = langext.Ptr([]models.UserID{u.UserID})
			} else if strings.EqualFold(*q.Direction, "both") {
				// both
			} else {
				return ginresp.APIError(g, 400, apierr.BINDFAIL_QUERY_PARAM, "Invalid value for param 'direction'", nil)
			}
		}

		if q.Confirmation != nil {
			if strings.EqualFold(*q.Confirmation, "confirmed") {
				filter.Confirmed = langext.PTrue
			} else if strings.EqualFold(*q.Confirmation, "unconfirmed") {
				filter.Confirmed = langext.PFalse
			} else if strings.EqualFold(*q.Confirmation, "all") {
				// both
			} else {
				return ginresp.APIError(g, 400, apierr.BINDFAIL_QUERY_PARAM, "Invalid value for param 'confirmation'", nil)
			}
		}

		if q.External != nil {
			if strings.EqualFold(*q.External, "true") {
				filter.SubscriberIsChannelOwner = langext.PFalse
			} else if strings.EqualFold(*q.External, "false") {
				filter.SubscriberIsChannelOwner = langext.PTrue
			} else if strings.EqualFold(*q.External, "all") {
				// both
			} else {
				return ginresp.APIError(g, 400, apierr.BINDFAIL_QUERY_PARAM, "Invalid value for param 'external'", nil)
			}
		}

		if q.SubscriberUserID != nil {
			filter.SubscriberUserID2 = langext.Ptr([]models.UserID{*q.SubscriberUserID})
		}

		if q.ChannelOwnerUserID != nil {
			filter.ChannelOwnerUserID2 = langext.Ptr([]models.UserID{*q.ChannelOwnerUserID})
		}

		res, err := h.database.ListSubscriptions(ctx, filter)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query subscriptions", err)
		}

		jsonres := langext.ArrMap(res, func(v models.Subscription) models.SubscriptionJSON { return v.JSON() })

		return finishSuccess(ginext.JSON(http.StatusOK, response{Subscriptions: jsonres}))

	})
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
func (h APIHandler) ListChannelSubscriptions(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
		UserID    models.UserID    `uri:"uid" binding:"entityid"`
		ChannelID models.ChannelID `uri:"cid" binding:"entityid"`
	}
	type response struct {
		Subscriptions []models.SubscriptionJSON `json:"subscriptions"`
	}

	var u uri
	ctx, g, errResp := pctx.URI(&u).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	return h.app.DoRequest(ctx, g, func(ctx *logic.AppContext, finishSuccess func(r ginext.HTTPResponse) ginext.HTTPResponse) ginext.HTTPResponse {

		if permResp := ctx.CheckPermissionUserRead(u.UserID); permResp != nil {
			return *permResp
		}

		_, err := h.database.GetChannel(ctx, u.UserID, u.ChannelID, true)
		if errors.Is(err, sql.ErrNoRows) {
			return ginresp.APIError(g, 404, apierr.CHANNEL_NOT_FOUND, "Channel not found", err)
		}
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query channel", err)
		}

		clients, err := h.database.ListSubscriptions(ctx, models.SubscriptionFilter{AnyUserID: langext.Ptr(u.UserID), ChannelID: langext.Ptr([]models.ChannelID{u.ChannelID})})
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query subscriptions", err)
		}

		res := langext.ArrMap(clients, func(v models.Subscription) models.SubscriptionJSON { return v.JSON() })

		return finishSuccess(ginext.JSON(http.StatusOK, response{Subscriptions: res}))

	})
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
func (h APIHandler) GetSubscription(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
		UserID         models.UserID         `uri:"uid" binding:"entityid"`
		SubscriptionID models.SubscriptionID `uri:"sid" binding:"entityid"`
	}

	var u uri
	ctx, g, errResp := pctx.URI(&u).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	return h.app.DoRequest(ctx, g, func(ctx *logic.AppContext, finishSuccess func(r ginext.HTTPResponse) ginext.HTTPResponse) ginext.HTTPResponse {

		if permResp := ctx.CheckPermissionUserRead(u.UserID); permResp != nil {
			return *permResp
		}

		subscription, err := h.database.GetSubscription(ctx, u.SubscriptionID)
		if errors.Is(err, sql.ErrNoRows) {
			return ginresp.APIError(g, 404, apierr.SUBSCRIPTION_NOT_FOUND, "Subscription not found", err)
		}
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query subscription", err)
		}
		if subscription.SubscriberUserID != u.UserID && subscription.ChannelOwnerUserID != u.UserID {
			return ginresp.APIError(g, 404, apierr.SUBSCRIPTION_USER_MISMATCH, "Subscription not found", nil)
		}

		return finishSuccess(ginext.JSON(http.StatusOK, subscription.JSON()))

	})
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
func (h APIHandler) CancelSubscription(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
		UserID         models.UserID         `uri:"uid" binding:"entityid"`
		SubscriptionID models.SubscriptionID `uri:"sid" binding:"entityid"`
	}

	var u uri
	ctx, g, errResp := pctx.URI(&u).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	return h.app.DoRequest(ctx, g, func(ctx *logic.AppContext, finishSuccess func(r ginext.HTTPResponse) ginext.HTTPResponse) ginext.HTTPResponse {

		if permResp := ctx.CheckPermissionUserAdmin(u.UserID); permResp != nil {
			return *permResp
		}

		subscription, err := h.database.GetSubscription(ctx, u.SubscriptionID)
		if errors.Is(err, sql.ErrNoRows) {
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

		return finishSuccess(ginext.JSON(http.StatusOK, subscription.JSON()))

	})
}

// CreateSubscription swaggerdoc
//
//	@Summary		Create/Request a subscription
//	@Description	Either [channel_owner_user_id, channel_internal_name] or [channel_id] must be supplied in the request body
//	@ID				api-subscriptions-create
//	@Tags			API-v2
//
//	@Param			uid			path		string								true	"UserID"
//	@Param			query_data	query		handler.CreateSubscription.query	false	" "
//	@Param			post_data	body		handler.CreateSubscription.body		false	" "
//
//	@Success		200			{object}	models.SubscriptionJSON
//	@Failure		400			{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure		401			{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure		500			{object}	ginresp.apiError	"internal server error"
//
//	@Router			/api/v2/users/{uid}/subscriptions [POST]
func (h APIHandler) CreateSubscription(pctx ginext.PreContext) ginext.HTTPResponse {
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
	ctx, g, errResp := pctx.URI(&u).Query(&q).Body(&b).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	return h.app.DoRequest(ctx, g, func(ctx *logic.AppContext, finishSuccess func(r ginext.HTTPResponse) ginext.HTTPResponse) ginext.HTTPResponse {

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

		existingSub, err := h.database.GetSubscriptionBySubscriber(ctx, u.UserID, channel.ChannelID)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query existing subscription", err)
		}
		if existingSub != nil {
			if !existingSub.Confirmed && channel.OwnerUserID == u.UserID {
				err = h.database.UpdateSubscriptionConfirmed(ctx, existingSub.SubscriptionID, true)
				if err != nil {
					return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to update subscription", err)
				}
				existingSub.Confirmed = true
			}

			return finishSuccess(ginext.JSON(http.StatusOK, existingSub.JSON()))
		}

		sub, err := h.database.CreateSubscription(ctx, u.UserID, channel, channel.OwnerUserID == u.UserID)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to create subscription", err)
		}

		return finishSuccess(ginext.JSON(http.StatusOK, sub.JSON()))

	})
}

// UpdateSubscription swaggerdoc
//
//	@Summary	Update a subscription (e.g. confirm)
//	@ID			api-subscriptions-update
//	@Tags		API-v2
//
//	@Param		uid			path		string							true	"UserID"
//	@Param		sid			path		string							true	"SubscriptionID"
//	@Param		post_data	body		handler.UpdateSubscription.body	false	" "
//
//	@Success	200			{object}	models.SubscriptionJSON
//	@Failure	400			{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure	401			{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure	404			{object}	ginresp.apiError	"subscription not found"
//	@Failure	500			{object}	ginresp.apiError	"internal server error"
//
//	@Router		/api/v2/users/{uid}/subscriptions/{sid} [PATCH]
func (h APIHandler) UpdateSubscription(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
		UserID         models.UserID         `uri:"uid" binding:"entityid"`
		SubscriptionID models.SubscriptionID `uri:"sid" binding:"entityid"`
	}
	type body struct {
		Confirmed *bool `form:"confirmed"`
	}

	var u uri
	var b body
	ctx, g, errResp := pctx.URI(&u).Body(&b).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	return h.app.DoRequest(ctx, g, func(ctx *logic.AppContext, finishSuccess func(r ginext.HTTPResponse) ginext.HTTPResponse) ginext.HTTPResponse {

		if permResp := ctx.CheckPermissionUserAdmin(u.UserID); permResp != nil {
			return *permResp
		}

		userid := *ctx.GetPermissionUserID()

		subscription, err := h.database.GetSubscription(ctx, u.SubscriptionID)
		if errors.Is(err, sql.ErrNoRows) {
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

		return finishSuccess(ginext.JSON(http.StatusOK, subscription.JSON()))

	})
}
