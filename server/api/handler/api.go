package handler

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/common/ginresp"
	"blackforestbytes.com/simplecloudnotifier/db"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"github.com/gin-gonic/gin"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"net/http"
	"regexp"
)

type APIHandler struct {
	app      *logic.Application
	database *db.Database
}

func NewAPIHandler(app *logic.Application) APIHandler {
	return APIHandler{
		app:      app,
		database: app.Database,
	}
}

// CreateUser swaggerdoc
//
// @Summary Create a new user
// @ID      api-user-create
//
// @Param   post_body body     handler.CreateUser.body false " "
//
// @Success 200       {object} models.UserJSON
// @Failure 400       {object} ginresp.apiError
// @Failure 500       {object} ginresp.apiError
//
// @Router  /api-v2/users/ [POST]
func (h APIHandler) CreateUser(g *gin.Context) ginresp.HTTPResponse {
	type body struct {
		FCMToken     string  `json:"fcm_token"`
		ProToken     *string `json:"pro_token"`
		Username     *string `json:"username"`
		AgentModel   string  `json:"agent_model"`
		AgentVersion string  `json:"agent_version"`
		ClientType   string  `json:"client_type"`
	}

	var b body
	ctx, errResp := h.app.StartRequest(g, nil, nil, &b)
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
		return ginresp.InternAPIError(400, apierr.INVALID_CLIENTTYPE, "Invalid ClientType", nil)
	}

	if b.ProToken != nil {
		ptok, err := h.app.VerifyProToken(*b.ProToken)
		if err != nil {
			return ginresp.InternAPIError(500, apierr.FAILED_VERIFY_PRO_TOKEN, "Failed to query purchase status", err)
		}

		if !ptok {
			return ginresp.InternAPIError(400, apierr.INVALID_PRO_TOKEN, "Purchase token could not be verified", nil)
		}
	}

	readKey := h.app.GenerateRandomAuthKey()
	sendKey := h.app.GenerateRandomAuthKey()
	adminKey := h.app.GenerateRandomAuthKey()

	err := h.database.ClearFCMTokens(ctx, b.FCMToken)
	if err != nil {
		return ginresp.InternAPIError(500, apierr.DATABASE_ERROR, "Failed to clear existing fcm tokens", err)
	}

	if b.ProToken != nil {
		err := h.database.ClearProTokens(ctx, *b.ProToken)
		if err != nil {
			return ginresp.InternAPIError(500, apierr.DATABASE_ERROR, "Failed to clear existing fcm tokens", err)
		}
	}

	userobj, err := h.database.CreateUser(ctx, readKey, sendKey, adminKey, b.ProToken, b.Username)
	if err != nil {
		return ginresp.InternAPIError(500, apierr.DATABASE_ERROR, "Failed to create user in db", err)
	}

	_, err = h.database.CreateClient(ctx, userobj.UserID, clientType, b.FCMToken, b.AgentModel, b.AgentVersion)
	if err != nil {
		return ginresp.InternAPIError(500, apierr.DATABASE_ERROR, "Failed to create user in db", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, userobj.JSON()))
}

// GetUser swaggerdoc
//
// @Summary Get a user
// @ID      api-user-get
//
// @Param   uid path     int true "UserID"
//
// @Success 200 {object} models.UserJSON
// @Failure 400 {object} ginresp.apiError
// @Failure 401 {object} ginresp.apiError
// @Failure 404 {object} ginresp.apiError
// @Failure 500 {object} ginresp.apiError
//
// @Router  /api-v2/users/{uid} [GET]
func (h APIHandler) GetUser(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID int64 `uri:"uid"`
	}

	var u uri
	ctx, errResp := h.app.StartRequest(g, &u, nil, nil)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	if permResp := ctx.CheckPermissionUserRead(u.UserID); permResp != nil {
		return *permResp
	}

	user, err := h.database.GetUser(ctx, u.UserID)
	if err == sql.ErrNoRows {
		return ginresp.InternAPIError(404, apierr.USER_NOT_FOUND, "User not found", err)
	}
	if err != nil {
		return ginresp.InternAPIError(500, apierr.DATABASE_ERROR, "Failed to query user", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, user.JSON()))
}

// UpdateUser swaggerdoc
//
// @Summary     (Partially) update a user
// @Description The body-values are optional, only send the ones you want to update
// @ID          api-user-update
//
// @Param       post_body body     handler.UpdateUser.body false " "
//
// @Success     200       {object} models.UserJSON
// @Failure     400       {object} ginresp.apiError
// @Failure     401       {object} ginresp.apiError
// @Failure     404       {object} ginresp.apiError
// @Failure     500       {object} ginresp.apiError
//
// @Router      /api-v2/users/{uid} [PATCH]
func (h APIHandler) UpdateUser(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID int64 `uri:"uid"`
	}
	type body struct {
		Username *string `json:"username"`
		ProToken *string `json:"pro_token"`
	}

	var u uri
	var b body
	ctx, errResp := h.app.StartRequest(g, &u, nil, &b)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	if permResp := ctx.CheckPermissionUserAdmin(u.UserID); permResp != nil {
		return *permResp
	}

	if b.Username != nil {
		username := langext.Ptr(regexp.MustCompile(`[[:alnum:]\-_]`).ReplaceAllString(*b.Username, ""))
		if *username == "" {
			username = nil
		}

		err := h.database.UpdateUserUsername(ctx, u.UserID, b.Username)
		if err != nil {
			return ginresp.InternAPIError(500, apierr.DATABASE_ERROR, "Failed to update user", err)
		}
	}

	if b.ProToken != nil {
		ptok, err := h.app.VerifyProToken(*b.ProToken)
		if err != nil {
			return ginresp.InternAPIError(500, apierr.FAILED_VERIFY_PRO_TOKEN, "Failed to query purchase status", err)
		}

		if !ptok {
			return ginresp.InternAPIError(400, apierr.INVALID_PRO_TOKEN, "Purchase token could not be verified", nil)
		}

		err = h.database.ClearProTokens(ctx, *b.ProToken)
		if err != nil {
			return ginresp.InternAPIError(500, apierr.DATABASE_ERROR, "Failed to clear existing fcm tokens", err)
		}

		err = h.database.UpdateUserProToken(ctx, u.UserID, b.ProToken)
		if err != nil {
			return ginresp.InternAPIError(500, apierr.DATABASE_ERROR, "Failed to update user", err)
		}
	}

	user, err := h.database.GetUser(ctx, u.UserID)
	if err != nil {
		return ginresp.InternAPIError(500, apierr.DATABASE_ERROR, "Failed to query (updated) user", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, user.JSON()))
}

// ListClients swaggerdoc
//
// @Summary List all clients
// @ID      api-clients-list
//
// @Param   uid path     int true "UserID"
//
// @Success 200 {object} handler.ListClients.response
// @Failure 400 {object} ginresp.apiError
// @Failure 401 {object} ginresp.apiError
// @Failure 404 {object} ginresp.apiError
// @Failure 500 {object} ginresp.apiError
//
// @Router  /api-v2/users/{uid}/clients [GET]
func (h APIHandler) ListClients(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID int64 `uri:"uid"`
	}
	type response struct {
		Clients []models.ClientJSON `json:"clients"`
	}

	var u uri
	ctx, errResp := h.app.StartRequest(g, &u, nil, nil)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	if permResp := ctx.CheckPermissionUserRead(u.UserID); permResp != nil {
		return *permResp
	}

	clients, err := h.database.ListClients(ctx, u.UserID)
	if err != nil {
		return ginresp.InternAPIError(500, apierr.DATABASE_ERROR, "Failed to query clients", err)
	}

	res := langext.ArrMap(clients, func(v models.Client) models.ClientJSON { return v.JSON() })

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, response{Clients: res}))
}

// GetClient swaggerdoc
//
// @Summary Get a single clients
// @ID      api-clients-get
//
// @Param   uid path     int true "UserID"
// @Param   cid path     int true "ClientID"
//
// @Success 200 {object} models.ClientJSON
// @Failure 400 {object} ginresp.apiError
// @Failure 401 {object} ginresp.apiError
// @Failure 404 {object} ginresp.apiError
// @Failure 500 {object} ginresp.apiError
//
// @Router  /api-v2/users/{uid}/clients/{cid} [GET]
func (h APIHandler) GetClient(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID   int64 `uri:"uid"`
		ClientID int64 `uri:"cid"`
	}

	var u uri
	ctx, errResp := h.app.StartRequest(g, &u, nil, nil)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	if permResp := ctx.CheckPermissionUserRead(u.UserID); permResp != nil {
		return *permResp
	}

	client, err := h.database.GetClient(ctx, u.UserID, u.ClientID)
	if err == sql.ErrNoRows {
		return ginresp.InternAPIError(404, apierr.CLIENT_NOT_FOUND, "Client not found", err)
	}
	if err != nil {
		return ginresp.InternAPIError(500, apierr.DATABASE_ERROR, "Failed to query client", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, client.JSON()))
}

// AddClient swaggerdoc
//
// @Summary Add a new clients
// @ID      api-clients-create
//
// @Param   uid       path     int                    true  "UserID"
//
// @Param   post_body body     handler.AddClient.body false " "
//
// @Success 200       {object} models.ClientJSON
// @Failure 400       {object} ginresp.apiError
// @Failure 401       {object} ginresp.apiError
// @Failure 404       {object} ginresp.apiError
// @Failure 500       {object} ginresp.apiError
//
// @Router  /api-v2/users/{uid}/clients [POST]
func (h APIHandler) AddClient(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID int64 `uri:"uid"`
	}
	type body struct {
		FCMToken     string `json:"fcm_token"`
		AgentModel   string `json:"agent_model"`
		AgentVersion string `json:"agent_version"`
		ClientType   string `json:"client_type"`
	}

	var u uri
	var b body
	ctx, errResp := h.app.StartRequest(g, &u, nil, &b)
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
		return ginresp.InternAPIError(400, apierr.INVALID_CLIENTTYPE, "Invalid ClientType", nil)
	}

	if permResp := ctx.CheckPermissionUserAdmin(u.UserID); permResp != nil {
		return *permResp
	}

	client, err := h.database.CreateClient(ctx, u.UserID, clientType, b.FCMToken, b.AgentModel, b.AgentVersion)
	if err != nil {
		return ginresp.InternAPIError(500, apierr.DATABASE_ERROR, "Failed to create user in db", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, client.JSON()))
}

// DeleteClient swaggerdoc
//
// @Summary Delete a client
// @ID      api-clients-delete
//
// @Param   uid path     int true "UserID"
// @Param   cid path     int true "ClientID"
//
// @Success 200 {object} models.ClientJSON
// @Failure 400 {object} ginresp.apiError
// @Failure 401 {object} ginresp.apiError
// @Failure 404 {object} ginresp.apiError
// @Failure 500 {object} ginresp.apiError
//
// @Router  /api-v2/users/{uid}/clients [POST]
func (h APIHandler) DeleteClient(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID   int64 `uri:"uid"`
		ClientID int64 `uri:"cid"`
	}

	var u uri
	ctx, errResp := h.app.StartRequest(g, &u, nil, nil)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	if permResp := ctx.CheckPermissionUserAdmin(u.UserID); permResp != nil {
		return *permResp
	}

	client, err := h.database.GetClient(ctx, u.UserID, u.ClientID)
	if err == sql.ErrNoRows {
		return ginresp.InternAPIError(404, apierr.CLIENT_NOT_FOUND, "Client not found", err)
	}
	if err != nil {
		return ginresp.InternAPIError(500, apierr.DATABASE_ERROR, "Failed to query client", err)
	}

	err = h.database.DeleteClient(ctx, u.ClientID)
	if err != nil {
		return ginresp.InternAPIError(500, apierr.DATABASE_ERROR, "Failed to delete client", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, client.JSON()))
}

// ListChannels swaggerdoc
//
// @Summary List all channels of a user
// @ID      api-channels-list
//
// @Param   uid path     int true "UserID"
//
// @Success 200 {object} handler.ListChannels.response
// @Failure 400 {object} ginresp.apiError
// @Failure 401 {object} ginresp.apiError
// @Failure 404 {object} ginresp.apiError
// @Failure 500 {object} ginresp.apiError
//
// @Router  /api-v2/users/{uid}/channels [GET]
func (h APIHandler) ListChannels(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID int64 `uri:"uid"`
	}
	type response struct {
		Channels []models.ChannelJSON `json:"channels"`
	}

	var u uri
	ctx, errResp := h.app.StartRequest(g, &u, nil, nil)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	if permResp := ctx.CheckPermissionUserRead(u.UserID); permResp != nil {
		return *permResp
	}

	clients, err := h.database.ListChannels(ctx, u.UserID)
	if err != nil {
		return ginresp.InternAPIError(500, apierr.DATABASE_ERROR, "Failed to query channels", err)
	}

	res := langext.ArrMap(clients, func(v models.Channel) models.ChannelJSON { return v.JSON() })

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, response{Channels: res}))
}

// GetChannel swaggerdoc
//
// @Summary List all channels of a user
// @ID      api-channels-get
//
// @Param   uid path     int true "UserID"
// @Param   cid path     int true "ChannelID"
//
// @Success 200 {object} models.ChannelJSON
// @Failure 400 {object} ginresp.apiError
// @Failure 401 {object} ginresp.apiError
// @Failure 404 {object} ginresp.apiError
// @Failure 500 {object} ginresp.apiError
//
// @Router  /api-v2/users/{uid}/channels/{cid} [GET]
func (h APIHandler) GetChannel(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID    int64 `uri:"uid"`
		ChannelID int64 `uri:"cid"`
	}

	var u uri
	ctx, errResp := h.app.StartRequest(g, &u, nil, nil)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	if permResp := ctx.CheckPermissionUserRead(u.UserID); permResp != nil {
		return *permResp
	}

	channel, err := h.database.GetChannel(ctx, u.UserID, u.ChannelID)
	if err == sql.ErrNoRows {
		return ginresp.InternAPIError(404, apierr.CLIENT_NOT_FOUND, "Channel not found", err)
	}
	if err != nil {
		return ginresp.InternAPIError(500, apierr.DATABASE_ERROR, "Failed to query channel", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, channel.JSON()))
}

func (h APIHandler) GetChannelMessages(g *gin.Context) ginresp.HTTPResponse {
	return ginresp.NotImplemented() //TODO
}

// ListUserSubscriptions swaggerdoc
//
// @Summary List all channels of a user
// @ID      api-user-subscriptions-list
//
// @Param   uid path     int true "UserID"
//
// @Success 200 {object} handler.ListUserSubscriptions.response
// @Failure 400 {object} ginresp.apiError
// @Failure 401 {object} ginresp.apiError
// @Failure 404 {object} ginresp.apiError
// @Failure 500 {object} ginresp.apiError
//
// @Router  /api-v2/users/{uid}/subscriptions [GET]
func (h APIHandler) ListUserSubscriptions(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID int64 `uri:"uid"`
	}
	type response struct {
		Subscriptions []models.SubscriptionJSON `json:"subscriptions"`
	}

	var u uri
	ctx, errResp := h.app.StartRequest(g, &u, nil, nil)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	if permResp := ctx.CheckPermissionUserRead(u.UserID); permResp != nil {
		return *permResp
	}

	clients, err := h.database.ListSubscriptionsByOwner(ctx, u.UserID)
	if err != nil {
		return ginresp.InternAPIError(500, apierr.DATABASE_ERROR, "Failed to query channels", err)
	}

	res := langext.ArrMap(clients, func(v models.Subscription) models.SubscriptionJSON { return v.JSON() })

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, response{Subscriptions: res}))
}

// ListChannelSubscriptions swaggerdoc
//
// @Summary List all subscriptions of a channel
// @ID      api-chan-subscriptions-list
//
// @Param   uid path     int true "UserID"
// @Param   cid path     int true "ChannelID"
//
// @Success 200 {object} handler.ListChannelSubscriptions.response
// @Failure 400 {object} ginresp.apiError
// @Failure 401 {object} ginresp.apiError
// @Failure 404 {object} ginresp.apiError
// @Failure 500 {object} ginresp.apiError
//
// @Router  /api-v2/users/{uid}/channels/{cid}/subscriptions [GET]
func (h APIHandler) ListChannelSubscriptions(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID    int64 `uri:"uid"`
		ChannelID int64 `uri:"cid"`
	}
	type response struct {
		Subscriptions []models.SubscriptionJSON `json:"subscriptions"`
	}

	var u uri
	ctx, errResp := h.app.StartRequest(g, &u, nil, nil)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	if permResp := ctx.CheckPermissionUserRead(u.UserID); permResp != nil {
		return *permResp
	}

	_, err := h.database.GetChannel(ctx, u.UserID, u.ChannelID)
	if err == sql.ErrNoRows {
		return ginresp.InternAPIError(404, apierr.CHANNEL_NOT_FOUND, "Channel not found", err)
	}
	if err != nil {
		return ginresp.InternAPIError(500, apierr.DATABASE_ERROR, "Failed to query channels", err)
	}

	clients, err := h.database.ListSubscriptionsByChannel(ctx, u.ChannelID)
	if err != nil {
		return ginresp.InternAPIError(500, apierr.DATABASE_ERROR, "Failed to query channels", err)
	}

	res := langext.ArrMap(clients, func(v models.Subscription) models.SubscriptionJSON { return v.JSON() })

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, response{Subscriptions: res}))
}

// GetSubscription swaggerdoc
//
// @Summary Get a single subscription
// @ID      api-subscriptions-get
//
// @Param   uid path     int true "UserID"
// @Param   sid path     int true "SubscriptionID"
//
// @Success 200 {object} models.SubscriptionJSON
// @Failure 400 {object} ginresp.apiError
// @Failure 401 {object} ginresp.apiError
// @Failure 404 {object} ginresp.apiError
// @Failure 500 {object} ginresp.apiError
//
// @Router  /api-v2/users/{uid}/subscriptions/{sid} [GET]
func (h APIHandler) GetSubscription(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID         int64 `uri:"uid"`
		SubscriptionID int64 `uri:"sid"`
	}

	var u uri
	ctx, errResp := h.app.StartRequest(g, &u, nil, nil)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	if permResp := ctx.CheckPermissionUserRead(u.UserID); permResp != nil {
		return *permResp
	}

	subscription, err := h.database.GetSubscription(ctx, u.SubscriptionID)
	if err == sql.ErrNoRows {
		return ginresp.InternAPIError(404, apierr.SUBSCRIPTION_NOT_FOUND, "Subscription not found", err)
	}
	if err != nil {
		return ginresp.InternAPIError(500, apierr.DATABASE_ERROR, "Failed to query channel", err)
	}

	if subscription.SubscriberUserID != u.UserID {
		return ginresp.InternAPIError(401, apierr.USER_AUTH_FAILED, "You are not authorized for this action", nil)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, subscription.JSON()))
}

// CancelSubscription swaggerdoc
//
// @Summary Cancel (delete) subscription
// @ID      api-subscriptions-delete
//
// @Param   uid path     int true "UserID"
// @Param   sid path     int true "SubscriptionID"
//
// @Success 200 {object} models.SubscriptionJSON
// @Failure 400 {object} ginresp.apiError
// @Failure 401 {object} ginresp.apiError
// @Failure 404 {object} ginresp.apiError
// @Failure 500 {object} ginresp.apiError
//
// @Router  /api-v2/users/{uid}/subscriptions/{sid} [DELETE]
func (h APIHandler) CancelSubscription(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID         int64 `uri:"uid"`
		SubscriptionID int64 `uri:"sid"`
	}

	var u uri
	ctx, errResp := h.app.StartRequest(g, &u, nil, nil)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	if permResp := ctx.CheckPermissionUserRead(u.UserID); permResp != nil {
		return *permResp
	}

	subscription, err := h.database.GetSubscription(ctx, u.SubscriptionID)
	if err == sql.ErrNoRows {
		return ginresp.InternAPIError(404, apierr.SUBSCRIPTION_NOT_FOUND, "Subscription not found", err)
	}
	if err != nil {
		return ginresp.InternAPIError(500, apierr.DATABASE_ERROR, "Failed to query channel", err)
	}

	if subscription.SubscriberUserID != u.UserID {
		return ginresp.InternAPIError(401, apierr.USER_AUTH_FAILED, "You are not authorized for this action", nil)
	}

	err = h.database.DeleteSubscription(ctx, u.SubscriptionID)
	if err != nil {
		return ginresp.InternAPIError(500, apierr.DATABASE_ERROR, "Failed to delete subscription", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, subscription.JSON()))
}

func (h APIHandler) CreateSubscription(g *gin.Context) ginresp.HTTPResponse {
	return ginresp.NotImplemented() //TODO
}

func (h APIHandler) UpdateSubscription(g *gin.Context) ginresp.HTTPResponse {
	return ginresp.NotImplemented() //TODO
}

func (h APIHandler) ListMessages(g *gin.Context) ginresp.HTTPResponse {
	return ginresp.NotImplemented() //TODO
}

func (h APIHandler) GetMessage(g *gin.Context) ginresp.HTTPResponse {
	return ginresp.NotImplemented() //TODO
}

func (h APIHandler) DeleteMessage(g *gin.Context) ginresp.HTTPResponse {
	return ginresp.NotImplemented() //TODO
}

func (h APIHandler) SendMessage(g *gin.Context) ginresp.HTTPResponse {
	return ginresp.NotImplemented() //TODO
}
