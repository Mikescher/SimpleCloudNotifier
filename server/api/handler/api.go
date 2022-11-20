package handler

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/common/ginresp"
	"blackforestbytes.com/simplecloudnotifier/db"
	"blackforestbytes.com/simplecloudnotifier/db/cursortoken"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/mathext"
	"gogs.mikescher.com/BlackForestBytes/goext/timeext"
	"net/http"
	"strings"
	"time"
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
		FCMToken     string  `json:"fcm_token" binding:"required"`
		ProToken     *string `json:"pro_token"`
		Username     *string `json:"username"`
		AgentModel   string  `json:"agent_model" binding:"required"`
		AgentVersion string  `json:"agent_version" binding:"required"`
		ClientType   string  `json:"client_type" binding:"required"`
	}

	var b body
	ctx, errResp := h.app.StartRequest(g, nil, nil, &b, nil)
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

	if b.ProToken != nil {
		ptok, err := h.app.VerifyProToken(*b.ProToken)
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
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to clear existing fcm tokens", err)
		}
	}

	username := b.Username
	if username != nil {
		username = langext.Ptr(h.app.NormalizeUsername(*username))
	}

	userobj, err := h.database.CreateUser(ctx, readKey, sendKey, adminKey, b.ProToken, username)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to create user in db", err)
	}

	_, err = h.database.CreateClient(ctx, userobj.UserID, clientType, b.FCMToken, b.AgentModel, b.AgentVersion)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to create user in db", err)
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
		UserID models.UserID `uri:"uid"`
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
// @Summary     (Partially) update a user
// @Description The body-values are optional, only send the ones you want to update
// @ID          api-user-update
//
// @Param       uid       path     int    true  "UserID"
//
// @Param       username  body     string false "Change the username (send an empty string to clear it)"
// @Param       pro_token body     string false "Send a verification of permium purchase"
// @Param       read_key  body     string false "Send `true` to create a new read_key"
// @Param       send_key  body     string false "Send `true` to create a new send_key"
// @Param       admin_key body     string false "Send `true` to create a new admin_key"
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
		UserID models.UserID `uri:"uid"`
	}
	type body struct {
		Username        *string `json:"username"`
		ProToken        *string `json:"pro_token"`
		RefreshReadKey  *bool   `json:"read_key"`
		RefreshSendKey  *bool   `json:"send_key"`
		RefreshAdminKey *bool   `json:"admin_key"`
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

		err := h.database.UpdateUserUsername(ctx, u.UserID, b.Username)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to update user", err)
		}
	}

	if b.ProToken != nil {
		ptok, err := h.app.VerifyProToken(*b.ProToken)
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

	if langext.Coalesce(b.RefreshSendKey, false) {
		newkey := h.app.GenerateRandomAuthKey()

		err := h.database.UpdateUserSendKey(ctx, u.UserID, newkey)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to update user", err)
		}
	}

	if langext.Coalesce(b.RefreshReadKey, false) {
		newkey := h.app.GenerateRandomAuthKey()

		err := h.database.UpdateUserReadKey(ctx, u.UserID, newkey)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to update user", err)
		}
	}

	if langext.Coalesce(b.RefreshAdminKey, false) {
		newkey := h.app.GenerateRandomAuthKey()

		err := h.database.UpdateUserAdminKey(ctx, u.UserID, newkey)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to update user", err)
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
		UserID models.UserID `uri:"uid"`
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
		UserID   models.UserID   `uri:"uid"`
		ClientID models.ClientID `uri:"cid"`
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
		UserID models.UserID `uri:"uid"`
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

	client, err := h.database.CreateClient(ctx, u.UserID, clientType, b.FCMToken, b.AgentModel, b.AgentVersion)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to create user in db", err)
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
		UserID   models.UserID   `uri:"uid"`
		ClientID models.ClientID `uri:"cid"`
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
// @Summary     List channels of a user (subscribed/owned)
// @Description The possible values for 'selector' are:
// @Description - "owned" Return all channels of the user
// @Description - "subscribed" Return all channels that the user is subscribing to
// @Description - "all" Return channels that the user owns or is subscribing
// @Description - "subscribed_any" Return all channels that the user is subscribing to (even unconfirmed)
// @Description - "all_any" Return channels that the user owns or is subscribing (even unconfirmed)
// @ID          api-channels-list
//
// @Param       uid      path     int    true "UserID"
// @Param       selector query    string true "Filter channels (default: owned)" Enums(owned, subscribed, all, subscribed_any, all_any)
//
// @Success     200      {object} handler.ListChannels.response
// @Failure     400      {object} ginresp.apiError
// @Failure     401      {object} ginresp.apiError
// @Failure     404      {object} ginresp.apiError
// @Failure     500      {object} ginresp.apiError
//
// @Router      /api-v2/users/{uid}/channels [GET]
func (h APIHandler) ListChannels(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID models.UserID `uri:"uid"`
	}
	type query struct {
		Selector *string `form:"selector"`
	}
	type response struct {
		Channels []models.ChannelJSON `json:"channels"`
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

	var res []models.ChannelJSON

	if sel == "owned" {
		channels, err := h.database.ListChannelsByOwner(ctx, u.UserID)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query channels", err)
		}
		res = langext.ArrMap(channels, func(v models.Channel) models.ChannelJSON { return v.JSON(true) })
	} else if sel == "subscribed_any" {
		channels, err := h.database.ListChannelsBySubscriber(ctx, u.UserID, false)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query channels", err)
		}
		res = langext.ArrMap(channels, func(v models.Channel) models.ChannelJSON { return v.JSON(false) })
	} else if sel == "all_any" {
		channels, err := h.database.ListChannelsByAccess(ctx, u.UserID, false)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query channels", err)
		}
		res = langext.ArrMap(channels, func(v models.Channel) models.ChannelJSON { return v.JSON(false) })
	} else if sel == "subscribed" {
		channels, err := h.database.ListChannelsBySubscriber(ctx, u.UserID, true)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query channels", err)
		}
		res = langext.ArrMap(channels, func(v models.Channel) models.ChannelJSON { return v.JSON(false) })
	} else if sel == "all" {
		channels, err := h.database.ListChannelsByAccess(ctx, u.UserID, true)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query channels", err)
		}
		res = langext.ArrMap(channels, func(v models.Channel) models.ChannelJSON { return v.JSON(false) })
	} else {
		return ginresp.APIError(g, 400, apierr.INVALID_ENUM_VALUE, "Invalid value for the [selector] parameter", nil)
	}

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
		UserID    models.UserID    `uri:"uid"`
		ChannelID models.ChannelID `uri:"cid"`
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

// UpdateChannel swaggerdoc
//
// @Summary (Partially) update a channel
// @ID      api-channels-update
//
// @Param   uid           path     int    true  "UserID"
// @Param   cid           path     int    true  "ChannelID"
//
// @Param   subscribe_key body     string false "Send `true` to create a new subscribe_key"
// @Param   send_key      body     string false "Send `true` to create a new send_key"
//
// @Success 200           {object} models.ChannelJSON
// @Failure 400           {object} ginresp.apiError
// @Failure 401           {object} ginresp.apiError
// @Failure 404           {object} ginresp.apiError
// @Failure 500           {object} ginresp.apiError
//
// @Router  /api-v2/users/{uid}/channels/{cid} [PATCH]
func (h APIHandler) UpdateChannel(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID    models.UserID    `uri:"uid"`
		ChannelID models.ChannelID `uri:"cid"`
	}
	type body struct {
		RefreshSubscribeKey *bool `json:"subscribe_key"`
		RefreshSendKey      *bool `json:"send_key"`
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

	_, err := h.database.GetChannel(ctx, u.UserID, u.ChannelID)
	if err == sql.ErrNoRows {
		return ginresp.APIError(g, 404, apierr.CHANNEL_NOT_FOUND, "Channel not found", err)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query channel", err)
	}

	if langext.Coalesce(b.RefreshSendKey, false) {
		newkey := h.app.GenerateRandomAuthKey()

		err := h.database.UpdateChannelSendKey(ctx, u.ChannelID, newkey)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to update user", err)
		}
	}

	if langext.Coalesce(b.RefreshSubscribeKey, false) {
		newkey := h.app.GenerateRandomAuthKey()

		err := h.database.UpdateChannelSubscribeKey(ctx, u.ChannelID, newkey)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to update user", err)
		}
	}

	user, err := h.database.GetChannel(ctx, u.UserID, u.ChannelID)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query (updated) user", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, user.JSON(true)))
}

// ListChannelMessages swaggerdoc
//
// @Summary     List messages of a channel
// @Description The next_page_token is an opaque token, the special value "@start" (or empty-string) is the beginning and "@end" is the end
// @Description Simply start the pagination without a next_page_token and get the next page by calling this endpoint with the returned next_page_token of the last query
// @Description If there are no more entries the token "@end" will be returned
// @Description By default we return long messages with a trimmed body, if trimmed=false is supplied we return full messages (this reduces the max page_size)
// @ID          api-channel-messages
//
// @Param       query_data query    handler.ListChannelMessages.query false " "
//
// @Success     200        {object} handler.ListChannelMessages.response
// @Failure     400        {object} ginresp.apiError
// @Failure     401        {object} ginresp.apiError
// @Failure     404        {object} ginresp.apiError
// @Failure     500        {object} ginresp.apiError
//
// @Router      /api-v2/users/{uid}/channels/{cid}/messages [GET]
func (h APIHandler) ListChannelMessages(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		ChannelUserID models.UserID    `uri:"uid"`
		ChannelID     models.ChannelID `uri:"cid"`
	}
	type query struct {
		PageSize      *int    `form:"page_size"`
		NextPageToken *string `form:"next_page_token"`
		Filter        *string `form:"filter"`
		Trimmed       *bool   `form:"trimmed"`
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

	if permResp := ctx.CheckPermissionRead(); permResp != nil {
		return *permResp
	}

	channel, err := h.database.GetChannel(ctx, u.ChannelUserID, u.ChannelID)
	if err == sql.ErrNoRows {
		return ginresp.APIError(g, 404, apierr.CHANNEL_NOT_FOUND, "Channel not found", err)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query channel", err)
	}

	userid := *ctx.GetPermissionUserID()

	sub, err := h.database.GetSubscriptionBySubscriber(ctx, userid, channel.ChannelID)
	if err == sql.ErrNoRows {
		return ginresp.APIError(g, 401, apierr.USER_AUTH_FAILED, "You are not authorized for this action", nil)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query subscription", err)
	}
	if !sub.Confirmed {
		return ginresp.APIError(g, 401, apierr.USER_AUTH_FAILED, "You are not authorized for this action", nil)
	}

	tok, err := cursortoken.Decode(langext.Coalesce(q.NextPageToken, ""))
	if err != nil {
		return ginresp.APIError(g, 500, apierr.PAGETOKEN_ERROR, "Failed to decode next_page_token", err)
	}

	messages, npt, err := h.database.ListChannelMessages(ctx, channel.ChannelID, pageSize, tok)
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
		UserID models.UserID `uri:"uid"`
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

	clients, err := h.database.ListSubscriptionsByOwner(ctx, u.UserID)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query channels", err)
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
		UserID    models.UserID    `uri:"uid"`
		ChannelID models.ChannelID `uri:"cid"`
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
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query channels", err)
	}

	clients, err := h.database.ListSubscriptionsByChannel(ctx, u.ChannelID)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query channels", err)
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
		UserID         models.UserID         `uri:"uid"`
		SubscriptionID models.SubscriptionID `uri:"sid"`
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

	if subscription.SubscriberUserID != u.UserID {
		return ginresp.APIError(g, 401, apierr.USER_AUTH_FAILED, "You are not authorized for this action", nil)
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
		UserID         models.UserID         `uri:"uid"`
		SubscriptionID models.SubscriptionID `uri:"sid"`
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
		return ginresp.APIError(g, 401, apierr.USER_AUTH_FAILED, "You are not authorized for this action", nil)
	}

	err = h.database.DeleteSubscription(ctx, u.SubscriptionID)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to delete subscription", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, subscription.JSON()))
}

// CreateSubscription swaggerdoc
//
// @Summary Creare/Request a subscription
// @ID      api-subscriptions-create
//
// @Param   uid        path     int                              true  "UserID"
// @Param   query_data query    handler.CreateSubscription.query false " "
// @Param   post_data  body     handler.CreateSubscription.body  false " "
//
// @Success 200        {object} models.SubscriptionJSON
// @Failure 400        {object} ginresp.apiError
// @Failure 401        {object} ginresp.apiError
// @Failure 404        {object} ginresp.apiError
// @Failure 500        {object} ginresp.apiError
//
// @Router  /api-v2/users/{uid}/subscriptions [POST]
func (h APIHandler) CreateSubscription(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID models.UserID `uri:"uid"`
	}
	type body struct {
		ChannelOwnerUserID models.UserID `form:"channel_owner_user_id" binding:"required"`
		Channel            string        `form:"channel_name" binding:"required"`
	}
	type query struct {
		ChanSubscribeKey *string `form:"chan_subscribe_key"`
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

	channel, err := h.database.GetChannelByName(ctx, b.ChannelOwnerUserID, h.app.NormalizeChannelName(b.Channel))
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query channel", err)
	}
	if channel == nil {
		return ginresp.APIError(g, 400, apierr.CHANNEL_NOT_FOUND, "Channel not found", err)
	}

	if channel.OwnerUserID != u.UserID && (q.ChanSubscribeKey == nil || *q.ChanSubscribeKey != channel.SubscribeKey) {
		ginresp.APIError(g, 401, apierr.USER_AUTH_FAILED, "You are not authorized for this action", nil)
	}

	sub, err := h.database.CreateSubscription(ctx, u.UserID, *channel, channel.OwnerUserID == u.UserID)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to create subscription", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, sub.JSON()))
}

// UpdateSubscription swaggerdoc
//
// @Summary Update a subscription (e.g. confirm)
// @ID      api-subscriptions-update
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
// @Router  /api-v2/users/{uid}/subscriptions/{sid} [PATCH]
func (h APIHandler) UpdateSubscription(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID         models.UserID         `uri:"uid"`
		SubscriptionID models.SubscriptionID `uri:"sid"`
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

	subscription, err := h.database.GetSubscription(ctx, u.SubscriptionID)
	if err == sql.ErrNoRows {
		return ginresp.APIError(g, 404, apierr.SUBSCRIPTION_NOT_FOUND, "Subscription not found", err)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query subscription", err)
	}

	if subscription.ChannelOwnerUserID != u.UserID {
		return ginresp.APIError(g, 401, apierr.USER_AUTH_FAILED, "You are not authorized for this action", nil)
	}

	if b.Confirmed != nil {
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
// @Summary     List all (subscribed) messages
// @Description The next_page_token is an opaque token, the special value "@start" (or empty-string) is the beginning and "@end" is the end
// @Description Simply start the pagination without a next_page_token and get the next page by calling this endpoint with the returned next_page_token of the last query
// @Description If there are no more entries the token "@end" will be returned
// @Description By default we return long messages with a trimmed body, if trimmed=false is supplied we return full messages (this reduces the max page_size)
// @ID          api-messages-list
//
// @Param       query_data query    handler.ListMessages.query false " "
//
// @Success     200        {object} handler.ListMessages.response
// @Failure     400        {object} ginresp.apiError
// @Failure     401        {object} ginresp.apiError
// @Failure     404        {object} ginresp.apiError
// @Failure     500        {object} ginresp.apiError
//
// @Router      /api-v2/messages [GET]
func (h APIHandler) ListMessages(g *gin.Context) ginresp.HTTPResponse {
	type query struct {
		PageSize      *int    `form:"page_size"`
		NextPageToken *string `form:"next_page_token"`
		Filter        *string `form:"filter"`
		Trimmed       *bool   `form:"trimmed"`
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

	if permResp := ctx.CheckPermissionRead(); permResp != nil {
		return *permResp
	}

	userid := *ctx.GetPermissionUserID()

	tok, err := cursortoken.Decode(langext.Coalesce(q.NextPageToken, ""))
	if err != nil {
		return ginresp.APIError(g, 500, apierr.PAGETOKEN_ERROR, "Failed to decode next_page_token", err)
	}

	err = h.database.UpdateUserLastRead(ctx, userid)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to update last-read", err)
	}

	messages, npt, err := h.database.ListMessages(ctx, userid, pageSize, tok)
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
// @Summary     Get a single message (untrimmed)
// @Description The user must either own the message and request the resource with the READ or ADMIN Key
// @Description Or the user must subscribe to the corresponding channel (and be confirmed) and request the resource with the READ or ADMIN Key
// @Description The returned message is never trimmed
// @ID          api-messages-get
//
// @Param       mid path     int true "SCNMessageID"
//
// @Success     200 {object} models.MessageJSON
// @Failure     400 {object} ginresp.apiError
// @Failure     401 {object} ginresp.apiError
// @Failure     404 {object} ginresp.apiError
// @Failure     500 {object} ginresp.apiError
//
// @Router      /api-v2/messages/{mid} [PATCH]
func (h APIHandler) GetMessage(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		MessageID models.SCNMessageID `uri:"mid"`
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

	msg, err := h.database.GetMessage(ctx, u.MessageID)
	if err == sql.ErrNoRows {
		return ginresp.APIError(g, 404, apierr.MESSAGE_NOT_FOUND, "message not found", err)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query message", err)
	}

	if !ctx.CheckPermissionMessageReadDirect(msg) {

		// either we have direct read permissions (it is our message + read/admin key)
		// or we subscribe (+confirmed) to the channel and have read/admin key

		if uid := ctx.GetPermissionUserID(); uid != nil && ctx.IsPermissionUserRead() {
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

		} else {
			// auth-key is not set or not a user:x variant
			return ginresp.APIError(g, 401, apierr.USER_AUTH_FAILED, "You are not authorized for this action", nil)
		}

	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, msg.FullJSON()))
}

// DeleteMessage swaggerdoc
//
// @Summary     Delete a single message
// @Description The user must own the message and request the resource with the ADMIN Key
// @ID          api-messages-delete
//
// @Param       mid path     int true "SCNMessageID"
//
// @Success     200 {object} models.MessageJSON
// @Failure     400 {object} ginresp.apiError
// @Failure     401 {object} ginresp.apiError
// @Failure     404 {object} ginresp.apiError
// @Failure     500 {object} ginresp.apiError
//
// @Router      /api-v2/messages/{mid} [PATCH]
func (h APIHandler) DeleteMessage(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		MessageID models.SCNMessageID `uri:"mid"`
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

	msg, err := h.database.GetMessage(ctx, u.MessageID)
	if err == sql.ErrNoRows {
		return ginresp.APIError(g, 404, apierr.MESSAGE_NOT_FOUND, "message not found", err)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query message", err)
	}

	if !ctx.CheckPermissionMessageReadDirect(msg) {
		return ginresp.APIError(g, 401, apierr.USER_AUTH_FAILED, "You are not authorized for this action", nil)
	}

	err = h.database.DeleteMessage(ctx, msg.SCNMessageID)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to delete message", err)
	}

	err = h.database.CancelPendingDeliveries(ctx, msg.SCNMessageID)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to cancel deliveries", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, msg.FullJSON()))
}

// CreateMessage swaggerdoc
//
// @Summary     Create a new message
// @Description This is similar to the main route `POST -> https://scn.blackfrestbytes.com/`
// @Description But this route can change in the future, for long-living scripts etc. it's better to use the normal POST route
// @ID          api-messages-create
//
// @Param       post_data query    handler.CreateMessage.body false " "
//
// @Success     200       {object} models.MessageJSON
// @Failure     400       {object} ginresp.apiError
// @Failure     401       {object} ginresp.apiError
// @Failure     404       {object} ginresp.apiError
// @Failure     500       {object} ginresp.apiError
//
// @Router      /api-v2/messages [POST]
func (h APIHandler) CreateMessage(g *gin.Context) ginresp.HTTPResponse {
	type body struct {
		Channel       *string  `json:"channel"`
		ChanKey       *string  `json:"chan_key"`
		Title         *string  `json:"title"`
		Content       *string  `json:"content"`
		Priority      *int     `json:"priority"`
		UserMessageID *string  `json:"msg_id"`
		SendTimestamp *float64 `json:"timestamp"`
	}

	var b body
	ctx, errResp := h.app.StartRequest(g, nil, nil, &b, nil)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	if permResp := ctx.CheckPermissionSend(); permResp != nil {
		return *permResp
	}

	userID := *ctx.GetPermissionUserID()

	if b.Title != nil {
		b.Title = langext.Ptr(strings.TrimSpace(*b.Title))
	}
	if b.UserMessageID != nil {
		b.UserMessageID = langext.Ptr(strings.TrimSpace(*b.UserMessageID))
	}

	if b.Title == nil {
		return ginresp.APIError(g, 400, apierr.MISSING_TITLE, "Missing parameter [[title]]", nil)
	}
	if b.SendTimestamp != nil && mathext.Abs(*b.SendTimestamp-float64(time.Now().Unix())) > (24*time.Hour).Seconds() {
		return ginresp.SendAPIError(g, 400, apierr.TIMESTAMP_OUT_OF_RANGE, -1, "The timestamp mus be within 24 hours of now()", nil)
	}
	if b.Priority != nil && (*b.Priority != 0 && *b.Priority != 1 && *b.Priority != 2) {
		return ginresp.SendAPIError(g, 400, apierr.INVALID_PRIO, 105, "Invalid priority", nil)
	}
	if len(*b.Title) == 0 {
		return ginresp.SendAPIError(g, 400, apierr.NO_TITLE, 103, "No title specified", nil)
	}
	if b.UserMessageID != nil && len(*b.UserMessageID) > 64 {
		return ginresp.SendAPIError(g, 400, apierr.USR_MSG_ID_TOO_LONG, -1, "MessageID too long (64 characters)", nil)
	}

	user, err := h.database.GetUser(ctx, userID)
	if err == sql.ErrNoRows {
		return ginresp.SendAPIError(g, 400, apierr.USER_NOT_FOUND, -1, "User not found", nil)
	}
	if err != nil {
		return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, -1, "Failed to query user", err)
	}

	channelName := user.DefaultChannel()
	if b.Channel != nil {
		channelName = h.app.NormalizeChannelName(*b.Channel)
	}

	if len(*b.Title) > user.MaxTitleLength() {
		return ginresp.SendAPIError(g, 400, apierr.TITLE_TOO_LONG, 103, fmt.Sprintf("Title too long (max %d characters)", user.MaxTitleLength()), nil)
	}
	if b.Content != nil && len(*b.Content) > user.MaxContentLength() {
		return ginresp.SendAPIError(g, 400, apierr.CONTENT_TOO_LONG, 104, fmt.Sprintf("Content too long (%d characters; max := %d characters)", len(*b.Content), user.MaxContentLength()), nil)
	}

	if b.UserMessageID != nil {
		msg, err := h.database.GetMessageByUserMessageID(ctx, *b.UserMessageID)
		if err != nil {
			return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, -1, "Failed to query existing message", err)
		}
		if msg != nil {
			return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, msg.FullJSON()))
		}
	}

	if user.QuotaRemainingToday() <= 0 {
		return ginresp.SendAPIError(g, 403, apierr.QUOTA_REACHED, -1, fmt.Sprintf("Daily quota reached (%d)", user.QuotaPerDay()), nil)
	}

	channel, err := h.app.GetOrCreateChannel(ctx, userID, channelName)
	if err != nil {
		return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, -1, "Failed to query/create channel", err)
	}

	var sendTimestamp *time.Time = nil
	if b.SendTimestamp != nil {
		sendTimestamp = langext.Ptr(timeext.UnixFloatSeconds(*b.SendTimestamp))
	}

	priority := langext.Coalesce(b.Priority, 1)

	msg, err := h.database.CreateMessage(ctx, userID, channel, sendTimestamp, *b.Title, b.Content, priority, b.UserMessageID)
	if err != nil {
		return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, -1, "Failed to create message in db", err)
	}

	subscriptions, err := h.database.ListSubscriptionsByChannel(ctx, channel.ChannelID)
	if err != nil {
		return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, -1, "Failed to query subscriptions", err)
	}

	err = h.database.IncUserMessageCounter(ctx, user)
	if err != nil {
		return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, -1, "Failed to inc user msg-counter", err)
	}

	err = h.database.IncChannelMessageCounter(ctx, channel)
	if err != nil {
		return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, -1, "Failed to inc channel msg-counter", err)
	}

	for _, sub := range subscriptions {
		clients, err := h.database.ListClients(ctx, sub.SubscriberUserID)
		if err != nil {
			return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, -1, "Failed to query clients", err)
		}

		if !sub.Confirmed {
			continue
		}

		for _, client := range clients {

			fcmDelivID, err := h.app.DeliverMessage(ctx, client, msg)
			if err != nil {
				_, err = h.database.CreateRetryDelivery(ctx, client, msg)
				if err != nil {
					return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, -1, "Failed to create delivery", err)
				}
			} else {
				_, err = h.database.CreateSuccessDelivery(ctx, client, msg, *fcmDelivID)
				if err != nil {
					return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, -1, "Failed to create delivery", err)
				}
			}

		}
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, msg.FullJSON()))
}
