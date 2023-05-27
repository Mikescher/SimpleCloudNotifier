package handler

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/api/ginresp"
	ct "blackforestbytes.com/simplecloudnotifier/db/cursortoken"
	primarydb "blackforestbytes.com/simplecloudnotifier/db/impl/primary"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/mathext"
	"net/http"
	"strings"
)

type APIHandler struct {
	app      *logic.Application
	database *primarydb.Database
}

func NewAPIHandler(app *logic.Application) APIHandler {
	return APIHandler{
		app:      app,
		database: app.Database.Primary,
	}
}

// CreateUser swaggerdoc
//
//	@Summary	Create a new user
//	@ID			api-user-create
//	@Tags		API-v2
//
//	@Param		post_body	body		handler.CreateUser.body	false	" "
//
//	@Success	200			{object}	models.UserJSONWithClientsAndKeys
//	@Failure	400			{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure	500			{object}	ginresp.apiError	"internal server error"
//
//	@Router		/api/v2/users [POST]
func (h APIHandler) CreateUser(g *gin.Context) ginresp.HTTPResponse {
	type body struct {
		FCMToken     string  `json:"fcm_token"`
		ProToken     *string `json:"pro_token"`
		Username     *string `json:"username"`
		AgentModel   string  `json:"agent_model"`
		AgentVersion string  `json:"agent_version"`
		ClientType   string  `json:"client_type"`
		NoClient     bool    `json:"no_client"`
	}

	var b body
	ctx, errResp := h.app.StartRequest(g, nil, nil, &b, nil)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	var clientType models.ClientType
	if !b.NoClient {
		if b.FCMToken == "" {
			return ginresp.APIError(g, 400, apierr.INVALID_CLIENTTYPE, "Missing FCMToken", nil)
		}
		if b.AgentVersion == "" {
			return ginresp.APIError(g, 400, apierr.INVALID_CLIENTTYPE, "Missing AgentVersion", nil)
		}
		if b.ClientType == "" {
			return ginresp.APIError(g, 400, apierr.INVALID_CLIENTTYPE, "Missing ClientType", nil)
		}
		if b.ClientType == string(models.ClientTypeAndroid) {
			clientType = models.ClientTypeAndroid
		} else if b.ClientType == string(models.ClientTypeIOS) {
			clientType = models.ClientTypeIOS
		} else {
			return ginresp.APIError(g, 400, apierr.BINDFAIL_BODY_PARAM, "Invalid ClientType", nil)
		}
	}

	if b.ProToken != nil {
		ptok, err := h.app.VerifyProToken(ctx, *b.ProToken)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.FAILED_VERIFY_PRO_TOKEN, "Failed to query purchase status", err)
		}

		if !ptok {
			return ginresp.APIError(g, 400, apierr.INVALID_PRO_TOKEN, "Purchase token could not be verified", nil)
		}
	}

	readKey := h.app.GenerateRandomAuthKey()
	sendKey := h.app.GenerateRandomAuthKey()
	adminKey := h.app.GenerateRandomAuthKey()

	err := h.database.ClearFCMTokens(ctx, b.FCMToken)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to clear existing fcm tokens", err)
	}

	if b.ProToken != nil {
		err := h.database.ClearProTokens(ctx, *b.ProToken)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to clear existing pro tokens", err)
		}
	}

	username := b.Username
	if username != nil {
		username = langext.Ptr(h.app.NormalizeUsername(*username))
	}

	userobj, err := h.database.CreateUser(ctx, b.ProToken, username)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to create user in db", err)
	}

	_, err = h.database.CreateKeyToken(ctx, "AdminKey (default)", userobj.UserID, true, make([]models.ChannelID, 0), models.TokenPermissionList{models.PermAdmin}, adminKey)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to create admin-key in db", err)
	}

	_, err = h.database.CreateKeyToken(ctx, "SendKey (default)", userobj.UserID, true, make([]models.ChannelID, 0), models.TokenPermissionList{models.PermChannelSend}, sendKey)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to create send-key in db", err)
	}

	_, err = h.database.CreateKeyToken(ctx, "ReadKey (default)", userobj.UserID, true, make([]models.ChannelID, 0), models.TokenPermissionList{models.PermUserRead, models.PermChannelRead}, readKey)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to create read-key in db", err)
	}

	if b.NoClient {
		return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, userobj.JSONWithClients(make([]models.Client, 0), adminKey, sendKey, readKey)))
	} else {
		err := h.database.DeleteClientsByFCM(ctx, b.FCMToken)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to delete existing clients in db", err)
		}

		client, err := h.database.CreateClient(ctx, userobj.UserID, clientType, b.FCMToken, b.AgentModel, b.AgentVersion)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to create client in db", err)
		}

		return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, userobj.JSONWithClients([]models.Client{client}, adminKey, sendKey, readKey)))
	}

}

// GetUser swaggerdoc
//
//	@Summary	Get a user
//	@ID			api-user-get
//	@Tags		API-v2
//
//	@Param		uid	path		int	true	"UserID"
//
//	@Success	200	{object}	models.UserJSON
//	@Failure	400	{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure	401	{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure	404	{object}	ginresp.apiError	"user not found"
//	@Failure	500	{object}	ginresp.apiError	"internal server error"
//
//	@Router		/api/v2/users/{uid} [GET]
func (h APIHandler) GetUser(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID models.UserID `uri:"uid" binding:"entityid"`
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

	user, err := h.database.GetUser(ctx, u.UserID)
	if err == sql.ErrNoRows {
		return ginresp.APIError(g, 404, apierr.USER_NOT_FOUND, "User not found", err)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query user", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, user.JSON()))
}

// UpdateUser swaggerdoc
//
//	@Summary		(Partially) update a user
//	@Description	The body-values are optional, only send the ones you want to update
//	@ID				api-user-update
//	@Tags			API-v2
//
//	@Param			uid			path		int		true	"UserID"
//
//	@Param			username	body		string	false	"Change the username (send an empty string to clear it)"
//	@Param			pro_token	body		string	false	"Send a verification of premium purchase"
//
//	@Success		200			{object}	models.UserJSON
//	@Failure		400			{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure		401			{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure		404			{object}	ginresp.apiError	"user not found"
//	@Failure		500			{object}	ginresp.apiError	"internal server error"
//
//	@Router			/api/v2/users/{uid} [PATCH]
func (h APIHandler) UpdateUser(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID models.UserID `uri:"uid" binding:"entityid"`
	}
	type body struct {
		Username *string `json:"username"`
		ProToken *string `json:"pro_token"`
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

	if b.Username != nil {
		username := langext.Ptr(h.app.NormalizeUsername(*b.Username))
		if *username == "" {
			username = nil
		}

		err := h.database.UpdateUserUsername(ctx, u.UserID, username)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to update user", err)
		}
	}

	if b.ProToken != nil {
		if *b.ProToken == "" {
			err := h.database.UpdateUserProToken(ctx, u.UserID, nil)
			if err != nil {
				return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to update user", err)
			}
		} else {
			ptok, err := h.app.VerifyProToken(ctx, *b.ProToken)
			if err != nil {
				return ginresp.APIError(g, 500, apierr.FAILED_VERIFY_PRO_TOKEN, "Failed to query purchase status", err)
			}

			if !ptok {
				return ginresp.APIError(g, 400, apierr.INVALID_PRO_TOKEN, "Purchase token could not be verified", nil)
			}

			err = h.database.ClearProTokens(ctx, *b.ProToken)
			if err != nil {
				return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to clear existing fcm tokens", err)
			}

			err = h.database.UpdateUserProToken(ctx, u.UserID, b.ProToken)
			if err != nil {
				return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to update user", err)
			}
		}
	}

	user, err := h.database.GetUser(ctx, u.UserID)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query (updated) user", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, user.JSON()))
}

// ListClients swaggerdoc
//
//	@Summary	List all clients
//	@ID			api-clients-list
//	@Tags		API-v2
//
//	@Param		uid	path		int	true	"UserID"
//
//	@Success	200	{object}	handler.ListClients.response
//	@Failure	400	{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure	401	{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure	500	{object}	ginresp.apiError	"internal server error"
//
//	@Router		/api/v2/users/{uid}/clients [GET]
func (h APIHandler) ListClients(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID models.UserID `uri:"uid" binding:"entityid"`
	}
	type response struct {
		Clients []models.ClientJSON `json:"clients"`
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

	clients, err := h.database.ListClients(ctx, u.UserID)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query clients", err)
	}

	res := langext.ArrMap(clients, func(v models.Client) models.ClientJSON { return v.JSON() })

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, response{Clients: res}))
}

// GetClient swaggerdoc
//
//	@Summary	Get a single client
//	@ID			api-clients-get
//	@Tags		API-v2
//
//	@Param		uid	path		int	true	"UserID"
//	@Param		cid	path		int	true	"ClientID"
//
//	@Success	200	{object}	models.ClientJSON
//	@Failure	400	{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure	401	{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure	404	{object}	ginresp.apiError	"client not found"
//	@Failure	500	{object}	ginresp.apiError	"internal server error"
//
//	@Router		/api/v2/users/{uid}/clients/{cid} [GET]
func (h APIHandler) GetClient(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID   models.UserID   `uri:"uid" binding:"entityid"`
		ClientID models.ClientID `uri:"cid" binding:"entityid"`
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

	client, err := h.database.GetClient(ctx, u.UserID, u.ClientID)
	if err == sql.ErrNoRows {
		return ginresp.APIError(g, 404, apierr.CLIENT_NOT_FOUND, "Client not found", err)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query client", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, client.JSON()))
}

// AddClient swaggerdoc
//
//	@Summary	Add a new clients
//	@ID			api-clients-create
//	@Tags		API-v2
//
//	@Param		uid			path		int						true	"UserID"
//
//	@Param		post_body	body		handler.AddClient.body	false	" "
//
//	@Success	200			{object}	models.ClientJSON
//	@Failure	400			{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure	401			{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure	500			{object}	ginresp.apiError	"internal server error"
//
//	@Router		/api/v2/users/{uid}/clients [POST]
func (h APIHandler) AddClient(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID models.UserID `uri:"uid" binding:"entityid"`
	}
	type body struct {
		FCMToken     string `json:"fcm_token" binding:"required"`
		AgentModel   string `json:"agent_model" binding:"required"`
		AgentVersion string `json:"agent_version" binding:"required"`
		ClientType   string `json:"client_type" binding:"required"`
	}

	var u uri
	var b body
	ctx, errResp := h.app.StartRequest(g, &u, nil, &b, nil)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	var clientType models.ClientType
	if b.ClientType == string(models.ClientTypeAndroid) {
		clientType = models.ClientTypeAndroid
	} else if b.ClientType == string(models.ClientTypeIOS) {
		clientType = models.ClientTypeIOS
	} else {
		return ginresp.APIError(g, 400, apierr.INVALID_CLIENTTYPE, "Invalid ClientType", nil)
	}

	if permResp := ctx.CheckPermissionUserAdmin(u.UserID); permResp != nil {
		return *permResp
	}

	err := h.database.DeleteClientsByFCM(ctx, b.FCMToken)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to delete existing clients in db", err)
	}

	client, err := h.database.CreateClient(ctx, u.UserID, clientType, b.FCMToken, b.AgentModel, b.AgentVersion)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to create client in db", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, client.JSON()))
}

// DeleteClient swaggerdoc
//
//	@Summary	Delete a client
//	@ID			api-clients-delete
//	@Tags		API-v2
//
//	@Param		uid	path		int	true	"UserID"
//	@Param		cid	path		int	true	"ClientID"
//
//	@Success	200	{object}	models.ClientJSON
//	@Failure	400	{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure	401	{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure	404	{object}	ginresp.apiError	"client not found"
//	@Failure	500	{object}	ginresp.apiError	"internal server error"
//
//	@Router		/api/v2/users/{uid}/clients/{cid} [DELETE]
func (h APIHandler) DeleteClient(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID   models.UserID   `uri:"uid" binding:"entityid"`
		ClientID models.ClientID `uri:"cid" binding:"entityid"`
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

	client, err := h.database.GetClient(ctx, u.UserID, u.ClientID)
	if err == sql.ErrNoRows {
		return ginresp.APIError(g, 404, apierr.CLIENT_NOT_FOUND, "Client not found", err)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query client", err)
	}

	err = h.database.DeleteClient(ctx, u.ClientID)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to delete client", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, client.JSON()))
}

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
//	@Param			uid			path		int		true	"UserID"
//	@Param			selector	query		string	false	"Filter channels (default: owned)"	Enums(owned, subscribed, all, subscribed_any, all_any)
//
//	@Success		200			{object}	handler.ListChannels.response
//	@Failure		400			{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure		401			{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure		500			{object}	ginresp.apiError	"internal server error"
//
//	@Router			/api/v2/users/{uid}/channels [GET]
func (h APIHandler) ListChannels(g *gin.Context) ginresp.HTTPResponse {
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

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, response{Channels: res}))
}

// GetChannel swaggerdoc
//
//	@Summary	Get a single channel
//	@ID			api-channels-get
//	@Tags		API-v2
//
//	@Param		uid	path		int	true	"UserID"
//	@Param		cid	path		int	true	"ChannelID"
//
//	@Success	200	{object}	models.ChannelWithSubscriptionJSON
//	@Failure	400	{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure	401	{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure	404	{object}	ginresp.apiError	"channel not found"
//	@Failure	500	{object}	ginresp.apiError	"internal server error"
//
//	@Router		/api/v2/users/{uid}/channels/{cid} [GET]
func (h APIHandler) GetChannel(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID    models.UserID    `uri:"uid" binding:"entityid"`
		ChannelID models.ChannelID `uri:"cid" binding:"entityid"`
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

	channel, err := h.database.GetChannel(ctx, u.UserID, u.ChannelID)
	if err == sql.ErrNoRows {
		return ginresp.APIError(g, 404, apierr.CHANNEL_NOT_FOUND, "Channel not found", err)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query channel", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, channel.JSON(true)))
}

// CreateChannel swaggerdoc
//
//	@Summary	Create a new (empty) channel
//	@ID			api-channels-create
//	@Tags		API-v2
//
//	@Param		uid			path		int							true	"UserID"
//	@Param		post_body	body		handler.CreateChannel.body	false	" "
//
//	@Success	200			{object}	models.ChannelWithSubscriptionJSON
//	@Failure	400			{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure	401			{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure	409			{object}	ginresp.apiError	"channel already exists"
//	@Failure	500			{object}	ginresp.apiError	"internal server error"
//
//	@Router		/api/v2/users/{uid}/channels [POST]
func (h APIHandler) CreateChannel(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID models.UserID `uri:"uid" binding:"entityid"`
	}
	type body struct {
		Name      string `json:"name"`
		Subscribe *bool  `json:"subscribe"`
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
	if err == sql.ErrNoRows {
		return ginresp.APIError(g, 400, apierr.USER_NOT_FOUND, "User not found", nil)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query user", err)
	}

	if len(channelDisplayName) > user.MaxChannelNameLength() {
		return ginresp.APIError(g, 400, apierr.CHANNEL_TOO_LONG, fmt.Sprintf("Channel too long (max %d characters)", user.MaxChannelNameLength()), nil)
	}
	if len(channelInternalName) > user.MaxChannelNameLength() {
		return ginresp.APIError(g, 400, apierr.CHANNEL_TOO_LONG, fmt.Sprintf("Channel too long (max %d characters)", user.MaxChannelNameLength()), nil)
	}

	if channelExisting != nil {
		return ginresp.APIError(g, 409, apierr.CHANNEL_ALREADY_EXISTS, "Channel with this name already exists", nil)
	}

	subscribeKey := h.app.GenerateRandomAuthKey()

	channel, err := h.database.CreateChannel(ctx, u.UserID, channelDisplayName, channelInternalName, subscribeKey)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to create channel", err)
	}

	if langext.Coalesce(b.Subscribe, true) {

		sub, err := h.database.CreateSubscription(ctx, u.UserID, channel, true)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to create subscription", err)
		}

		return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, channel.WithSubscription(langext.Ptr(sub)).JSON(true)))

	} else {

		return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, channel.WithSubscription(nil).JSON(true)))

	}

}

// UpdateChannel swaggerdoc
//
//	@Summary	(Partially) update a channel
//	@ID			api-channels-update
//	@Tags		API-v2
//
//	@Param		uid				path		int		true	"UserID"
//	@Param		cid				path		int		true	"ChannelID"
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
func (h APIHandler) UpdateChannel(g *gin.Context) ginresp.HTTPResponse {
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
	ctx, errResp := h.app.StartRequest(g, &u, nil, &b, nil)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	if permResp := ctx.CheckPermissionUserAdmin(u.UserID); permResp != nil {
		return *permResp
	}

	oldChannel, err := h.database.GetChannel(ctx, u.UserID, u.ChannelID)
	if err == sql.ErrNoRows {
		return ginresp.APIError(g, 404, apierr.CHANNEL_NOT_FOUND, "Channel not found", err)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query channel", err)
	}

	user, err := h.database.GetUser(ctx, u.UserID)
	if err == sql.ErrNoRows {
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
		newInternalName := h.app.NormalizeChannelInternalName(*b.DisplayName)

		if newInternalName != oldChannel.InternalName {
			return ginresp.APIError(g, 400, apierr.CHANNEL_NAME_WOULD_CHANGE, "Cannot substantially change the channel name", err)
		}

		if len(newDisplayName) > user.MaxChannelNameLength() {
			return ginresp.APIError(g, 400, apierr.CHANNEL_TOO_LONG, fmt.Sprintf("Channel too long (max %d characters)", user.MaxChannelNameLength()), nil)
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

		if descName != nil && len(*descName) > user.MaxChannelDescriptionNameLength() {
			return ginresp.APIError(g, 400, apierr.CHANNEL_DESCRIPTION_TOO_LONG, fmt.Sprintf("Channel-Description too long (max %d characters)", user.MaxChannelNameLength()), nil)
		}

		err := h.database.UpdateChannelDescriptionName(ctx, u.ChannelID, descName)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to update channel", err)
		}

	}

	channel, err := h.database.GetChannel(ctx, u.UserID, u.ChannelID)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query (updated) channel", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, channel.JSON(true)))
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
//	@Param			uid			path		int									true	"UserID"
//	@Param			cid			path		int									true	"ChannelID"
//
//	@Success		200			{object}	handler.ListChannelMessages.response
//	@Failure		400			{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure		401			{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure		404			{object}	ginresp.apiError	"channel not found"
//	@Failure		500			{object}	ginresp.apiError	"internal server error"
//
//	@Router			/api/v2/users/{uid}/channels/{cid}/messages [GET]
func (h APIHandler) ListChannelMessages(g *gin.Context) ginresp.HTTPResponse {
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
	ctx, errResp := h.app.StartRequest(g, &u, &q, nil, nil)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	trimmed := langext.Coalesce(q.Trimmed, true)

	maxPageSize := langext.Conditional(trimmed, 16, 256)

	pageSize := mathext.Clamp(langext.Coalesce(q.PageSize, 64), 1, maxPageSize)

	channel, err := h.database.GetChannel(ctx, u.ChannelUserID, u.ChannelID)
	if err == sql.ErrNoRows {
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
		return ginresp.APIError(g, 500, apierr.PAGETOKEN_ERROR, "Failed to decode next_page_token", err)
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

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, response{Messages: res, NextPageToken: npt.Token(), PageSize: pageSize}))
}

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
//	@Param			selector	query		string	true	"Filter subscriptions (default: owner_all)"	Enums(outgoing_all, outgoing_confirmed, outgoing_unconfirmed, incoming_all, incoming_confirmed, incoming_unconfirmed)
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
		Selector *string `json:"selector" form:"selector"  enums:"owner_all,owner_confirmed,owner_unconfirmed,incoming_all,incoming_confirmed,incoming_unconfirmed"`
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

	sel := strings.ToLower(langext.Coalesce(q.Selector, "owner_all"))

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
//	@Param		uid	path		int	true	"UserID"
//	@Param		cid	path		int	true	"ChannelID"
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

	_, err := h.database.GetChannel(ctx, u.UserID, u.ChannelID)
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
//	@Param		uid	path		int	true	"UserID"
//	@Param		sid	path		int	true	"SubscriptionID"
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
//	@Param		uid	path		int	true	"UserID"
//	@Param		sid	path		int	true	"SubscriptionID"
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
		return ginresp.APIError(g, 500, apierr.PAGETOKEN_ERROR, "Failed to decode next_page_token", err)
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

// ListUserKeys swaggerdoc
//
//	@Summary		List keys of the user
//	@Description	The request must be done with an ADMIN key, the returned keys are without their token.
//	@ID				api-tokenkeys-list
//	@Tags			API-v2
//
//	@Param			uid	path		int	true	"UserID"
//
//	@Success		200	{object}	handler.ListUserKeys.response
//	@Failure		400	{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure		401	{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure		404	{object}	ginresp.apiError	"message not found"
//	@Failure		500	{object}	ginresp.apiError	"internal server error"
//
//	@Router			/api/v2/users/:uid/keys [GET]
func (h APIHandler) ListUserKeys(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID models.UserID `uri:"uid" binding:"entityid"`
	}
	type response struct {
		Keys []models.KeyTokenJSON `json:"keys"`
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

	clients, err := h.database.ListKeyTokens(ctx, u.UserID)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query keys", err)
	}

	res := langext.ArrMap(clients, func(v models.KeyToken) models.KeyTokenJSON { return v.JSON() })

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, response{Keys: res}))
}

// GetUserKey swaggerdoc
//
//	@Summary		Get a single key
//	@Description	The request must be done with an ADMIN key, the returned key does not include its token.
//	@ID				api-tokenkeys-get
//	@Tags			API-v2
//
//	@Param			uid	path		int	true	"UserID"
//	@Param			kid	path		int	true	"TokenKeyID"
//
//	@Success		200	{object}	models.KeyTokenJSON
//	@Failure		400	{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure		401	{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure		404	{object}	ginresp.apiError	"message not found"
//	@Failure		500	{object}	ginresp.apiError	"internal server error"
//
//	@Router			/api/v2/users/:uid/keys/:kid [GET]
func (h APIHandler) GetUserKey(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID models.UserID     `uri:"uid" binding:"entityid"`
		KeyID  models.KeyTokenID `uri:"kid" binding:"entityid"`
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

	keytoken, err := h.database.GetKeyToken(ctx, u.UserID, u.KeyID)
	if err == sql.ErrNoRows {
		return ginresp.APIError(g, 404, apierr.KEY_NOT_FOUND, "Key not found", err)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query client", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, keytoken.JSON()))
}

// UpdateUserKey swaggerdoc
//
//	@Summary	Update a key
//	@ID			api-tokenkeys-update
//	@Tags		API-v2
//
//	@Param		uid	path		int	true	"UserID"
//	@Param		kid	path		int	true	"TokenKeyID"
//
//	@Param			post_body	body		handler.UpdateUserKey.body	false	" "
//
//	@Success	200	{object}	models.KeyTokenJSON
//	@Failure	400	{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure	401	{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure	404	{object}	ginresp.apiError	"message not found"
//	@Failure	500	{object}	ginresp.apiError	"internal server error"
//
//	@Router		/api/v2/users/:uid/keys/:kid [PATCH]
func (h APIHandler) UpdateUserKey(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID models.UserID     `uri:"uid" binding:"entityid"`
		KeyID  models.KeyTokenID `uri:"kid" binding:"entityid"`
	}
	type body struct {
		Name        *string             `json:"name"`
		AllChannels *bool               `json:"all_channels"`
		Channels    *[]models.ChannelID `json:"channels"`
		Permissions *string             `json:"permissions"`
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

	keytoken, err := h.database.GetKeyToken(ctx, u.UserID, u.KeyID)
	if err == sql.ErrNoRows {
		return ginresp.APIError(g, 404, apierr.KEY_NOT_FOUND, "Key not found", err)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query client", err)
	}

	if b.Name != nil {
		err := h.database.UpdateKeyTokenName(ctx, u.KeyID, *b.Name)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to update name", err)
		}
		keytoken.Name = *b.Name
	}

	if b.Permissions != nil {
		if keytoken.KeyTokenID == *ctx.GetPermissionKeyTokenID() {
			return ginresp.APIError(g, 400, apierr.CANNOT_SELFUPDATE_KEY, "Cannot update the currently used key", err)
		}

		permlist := models.ParseTokenPermissionList(*b.Permissions)
		err := h.database.UpdateKeyTokenPermissions(ctx, u.KeyID, permlist)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to update permissions", err)
		}
		keytoken.Permissions = permlist
	}

	if b.AllChannels != nil {
		if keytoken.KeyTokenID == *ctx.GetPermissionKeyTokenID() {
			return ginresp.APIError(g, 400, apierr.CANNOT_SELFUPDATE_KEY, "Cannot update the currently used key", err)
		}

		err := h.database.UpdateKeyTokenAllChannels(ctx, u.KeyID, *b.AllChannels)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to update all_channels", err)
		}
		keytoken.AllChannels = *b.AllChannels
	}

	if b.Channels != nil {
		if keytoken.KeyTokenID == *ctx.GetPermissionKeyTokenID() {
			return ginresp.APIError(g, 400, apierr.CANNOT_SELFUPDATE_KEY, "Cannot update the currently used key", err)
		}

		err := h.database.UpdateKeyTokenChannels(ctx, u.KeyID, *b.Channels)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to update channels", err)
		}
		keytoken.Channels = *b.Channels
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, keytoken.JSON()))
}

// CreateUserKey swaggerdoc
//
//	@Summary	Create a new key
//	@ID			api-tokenkeys-create
//	@Tags		API-v2
//
//	@Param		uid			path		int							true	"UserID"
//
//	@Param		post_body	body		handler.CreateUserKey.body	false	" "
//
//	@Success	200			{object}	models.KeyTokenJSON
//	@Failure	400			{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure	401			{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure	404			{object}	ginresp.apiError	"message not found"
//	@Failure	500			{object}	ginresp.apiError	"internal server error"
//
//	@Router		/api/v2/users/:uid/keys [POST]
func (h APIHandler) CreateUserKey(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID models.UserID `uri:"uid" binding:"entityid"`
	}
	type body struct {
		Name        string              `json:"name"         binding:"required"`
		AllChannels *bool               `json:"all_channels" binding:"required"`
		Channels    *[]models.ChannelID `json:"channels"     binding:"required"`
		Permissions *string             `json:"permissions"  binding:"required"`
	}

	var u uri
	var b body
	ctx, errResp := h.app.StartRequest(g, &u, nil, &b, nil)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	for _, c := range *b.Channels {
		if err := c.Valid(); err != nil {
			return ginresp.APIError(g, 400, apierr.INVALID_BODY_PARAM, "Invalid ChannelID", err)
		}
	}

	if permResp := ctx.CheckPermissionUserAdmin(u.UserID); permResp != nil {
		return *permResp
	}

	token := h.app.GenerateRandomAuthKey()

	perms := models.ParseTokenPermissionList(*b.Permissions)

	keytok, err := h.database.CreateKeyToken(ctx, b.Name, *ctx.GetPermissionUserID(), *b.AllChannels, *b.Channels, perms, token)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to create keytoken in db", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, keytok.JSON().WithToken(token)))
}

// DeleteUserKey swaggerdoc
//
//	@Summary		Delete a key
//	@Description	Cannot be used to delete the key used in the request itself
//	@ID				api-tokenkeys-delete
//	@Tags			API-v2
//
//	@Param			uid	path		int	true	"UserID"
//	@Param			kid	path		int	true	"TokenKeyID"
//
//	@Success		200	{object}	models.KeyTokenJSON
//	@Failure		400	{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure		401	{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure		404	{object}	ginresp.apiError	"message not found"
//	@Failure		500	{object}	ginresp.apiError	"internal server error"
//
//	@Router			/api/v2/users/:uid/keys/:kid [DELETE]
func (h APIHandler) DeleteUserKey(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID models.UserID     `uri:"uid" binding:"entityid"`
		KeyID  models.KeyTokenID `uri:"kid" binding:"entityid"`
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

	client, err := h.database.GetKeyToken(ctx, u.UserID, u.KeyID)
	if err == sql.ErrNoRows {
		return ginresp.APIError(g, 404, apierr.KEY_NOT_FOUND, "Key not found", err)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query client", err)
	}

	if u.KeyID == *ctx.GetPermissionKeyTokenID() {
		return ginresp.APIError(g, 400, apierr.CANNOT_SELFDELETE_KEY, "Cannot delete the currently used key", err)
	}

	err = h.database.DeleteKeyToken(ctx, u.KeyID)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to delete client", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, client.JSON()))
}
