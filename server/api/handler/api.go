package handler

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	hl "blackforestbytes.com/simplecloudnotifier/api/apihighlight"
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
	"net/http"
	"strings"
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
// @Tags    API-v2
//
// @Param   post_body body     handler.CreateUser.body false " "
//
// @Success 200       {object} models.UserJSONWithClients
// @Failure 400       {object} ginresp.apiError "supplied values/parameters cannot be parsed / are invalid"
// @Failure 500       {object} ginresp.apiError "internal server error"
//
// @Router  /api/users [POST]
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

	if b.NoClient {
		return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, userobj.JSONWithClients(make([]models.Client, 0))))
	} else {
		err := h.database.DeleteClientsByFCM(ctx, b.FCMToken)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to delete existing clients in db", err)
		}

		client, err := h.database.CreateClient(ctx, userobj.UserID, clientType, b.FCMToken, b.AgentModel, b.AgentVersion)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to create client in db", err)
		}

		return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, userobj.JSONWithClients([]models.Client{client})))
	}

}

// GetUser swaggerdoc
//
// @Summary Get a user
// @ID      api-user-get
// @Tags    API-v2
//
// @Param   uid path     int true "UserID"
//
// @Success 200 {object} models.UserJSON
// @Failure 400 {object} ginresp.apiError "supplied values/parameters cannot be parsed / are invalid"
// @Failure 401 {object} ginresp.apiError "user is not authorized / has missing permissions"
// @Failure 404 {object} ginresp.apiError "user not found"
// @Failure 500 {object} ginresp.apiError "internal server error"
//
// @Router  /api/users/{uid} [GET]
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
// @Tags        API-v2
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
// @Failure     400       {object} ginresp.apiError "supplied values/parameters cannot be parsed / are invalid"
// @Failure     401       {object} ginresp.apiError "user is not authorized / has missing permissions"
// @Failure     404       {object} ginresp.apiError "user not found"
// @Failure     500       {object} ginresp.apiError "internal server error"
//
// @Router      /api/users/{uid} [PATCH]
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

		err := h.database.UpdateUserUsername(ctx, u.UserID, username)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to update user", err)
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
// @Tags    API-v2
//
// @Param   uid path     int true "UserID"
//
// @Success 200 {object} handler.ListClients.response
// @Failure 400 {object} ginresp.apiError "supplied values/parameters cannot be parsed / are invalid"
// @Failure 401 {object} ginresp.apiError "user is not authorized / has missing permissions"
// @Failure 500 {object} ginresp.apiError "internal server error"
//
// @Router  /api/users/{uid}/clients [GET]
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
// @Summary Get a single client
// @ID      api-clients-get
// @Tags    API-v2
//
// @Param   uid path     int true "UserID"
// @Param   cid path     int true "ClientID"
//
// @Success 200 {object} models.ClientJSON
// @Failure 400 {object} ginresp.apiError "supplied values/parameters cannot be parsed / are invalid"
// @Failure 401 {object} ginresp.apiError "user is not authorized / has missing permissions"
// @Failure 404 {object} ginresp.apiError "client not found"
// @Failure 500 {object} ginresp.apiError "internal server error"
//
// @Router  /api/users/{uid}/clients/{cid} [GET]
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
// @Tags    API-v2
//
// @Param   uid       path     int                    true  "UserID"
//
// @Param   post_body body     handler.AddClient.body false " "
//
// @Success 200       {object} models.ClientJSON
// @Failure 400       {object} ginresp.apiError "supplied values/parameters cannot be parsed / are invalid"
// @Failure 401       {object} ginresp.apiError "user is not authorized / has missing permissions"
// @Failure 500       {object} ginresp.apiError "internal server error"
//
// @Router  /api/users/{uid}/clients [POST]
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
// @Summary Delete a client
// @ID      api-clients-delete
// @Tags    API-v2
//
// @Param   uid path     int true "UserID"
// @Param   cid path     int true "ClientID"
//
// @Success 200 {object} models.ClientJSON
// @Failure 400 {object} ginresp.apiError "supplied values/parameters cannot be parsed / are invalid"
// @Failure 401 {object} ginresp.apiError "user is not authorized / has missing permissions"
// @Failure 404 {object} ginresp.apiError "client not found"
// @Failure 500 {object} ginresp.apiError "internal server error"
//
// @Router  /api/users/{uid}/clients/{cid} [DELETE]
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
// @Tags        API-v2
//
// @Param       uid      path     int    true "UserID"
// @Param       selector query    string true "Filter channels (default: owned)" Enums(owned, subscribed, all, subscribed_any, all_any)
//
// @Success     200      {object} handler.ListChannels.response
// @Failure     400      {object} ginresp.apiError "supplied values/parameters cannot be parsed / are invalid"
// @Failure     401      {object} ginresp.apiError "user is not authorized / has missing permissions"
// @Failure     500      {object} ginresp.apiError "internal server error"
//
// @Router      /api/users/{uid}/channels [GET]
func (h APIHandler) ListChannels(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID models.UserID `uri:"uid"`
	}
	type query struct {
		Selector *string `json:"selector" form:"selector"  enums:"owned,subscribed_any,all_any,subscribed,all"`
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
// @Tags    API-v2
//
// @Param   uid path     int true "UserID"
// @Param   cid path     int true "ChannelID"
//
// @Success 200 {object} models.ChannelJSON
// @Failure 400 {object} ginresp.apiError "supplied values/parameters cannot be parsed / are invalid"
// @Failure 401 {object} ginresp.apiError "user is not authorized / has missing permissions"
// @Failure 404 {object} ginresp.apiError "channel not found"
// @Failure 500 {object} ginresp.apiError "internal server error"
//
// @Router  /api/users/{uid}/channels/{cid} [GET]
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

// CreateChannel swaggerdoc
//
// @Summary Create a new (empty) channel
// @ID      api-channels-create
// @Tags    API-v2
//
// @Param   uid       path     int                        true  "UserID"
// @Param   post_body body     handler.CreateChannel.body false " "
//
// @Success 200       {object} models.ChannelJSON
// @Failure 400       {object} ginresp.apiError "supplied values/parameters cannot be parsed / are invalid"
// @Failure 401       {object} ginresp.apiError "user is not authorized / has missing permissions"
// @Failure 409       {object} ginresp.apiError "channel already exists"
// @Failure 500       {object} ginresp.apiError "internal server error"
//
// @Router  /api/users/{uid}/channels/ [POST]
func (h APIHandler) CreateChannel(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID models.UserID `uri:"uid"`
	}
	type body struct {
		Name string `json:"name"`
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

	channelName := h.app.NormalizeChannelName(b.Name)

	channelExisting, err := h.database.GetChannelByName(ctx, u.UserID, channelName)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query channel", err)
	}

	user, err := h.database.GetUser(ctx, u.UserID)
	if err == sql.ErrNoRows {
		return ginresp.SendAPIError(g, 400, apierr.USER_NOT_FOUND, hl.USER_ID, "User not found", nil)
	}
	if err != nil {
		return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query user", err)
	}

	if len(channelName) > user.MaxChannelNameLength() {
		return ginresp.SendAPIError(g, 400, apierr.CHANNEL_TOO_LONG, hl.CHANNEL, fmt.Sprintf("Channel too long (max %d characters)", user.MaxChannelNameLength()), nil)
	}

	if channelExisting != nil {
		return ginresp.APIError(g, 409, apierr.CHANNEL_ALREADY_EXISTS, "Channel with this name already exists", nil)
	}

	subscribeKey := h.app.GenerateRandomAuthKey()
	sendKey := h.app.GenerateRandomAuthKey()

	channel, err := h.database.CreateChannel(ctx, u.UserID, channelName, subscribeKey, sendKey)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to create channel", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, channel.JSON(true)))
}

// UpdateChannel swaggerdoc
//
// @Summary (Partially) update a channel
// @ID      api-channels-update
// @Tags    API-v2
//
// @Param   uid           path     int    true  "UserID"
// @Param   cid           path     int    true  "ChannelID"
//
// @Param   subscribe_key body     string false "Send `true` to create a new subscribe_key"
// @Param   send_key      body     string false "Send `true` to create a new send_key"
//
// @Success 200           {object} models.ChannelJSON
// @Failure 400           {object} ginresp.apiError "supplied values/parameters cannot be parsed / are invalid"
// @Failure 401           {object} ginresp.apiError "user is not authorized / has missing permissions"
// @Failure 404           {object} ginresp.apiError "channel not found"
// @Failure 500           {object} ginresp.apiError "internal server error"
//
// @Router  /api/users/{uid}/channels/{cid} [PATCH]
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
// @Tags        API-v2
//
// @Param       query_data query    handler.ListChannelMessages.query false " "
// @Param       uid        path     int                               true  "UserID"
// @Param       cid        path     int                               true  "ChannelID"
//
// @Success     200        {object} handler.ListChannelMessages.response
// @Failure     400        {object} ginresp.apiError "supplied values/parameters cannot be parsed / are invalid"
// @Failure     401        {object} ginresp.apiError "user is not authorized / has missing permissions"
// @Failure     404        {object} ginresp.apiError "channel not found"
// @Failure     500        {object} ginresp.apiError "internal server error"
//
// @Router      /api/users/{uid}/channels/{cid}/messages [GET]
func (h APIHandler) ListChannelMessages(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		ChannelUserID models.UserID    `uri:"uid"`
		ChannelID     models.ChannelID `uri:"cid"`
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

	filter := models.MessageFilter{
		ChannelID: langext.Ptr([]models.ChannelID{channel.ChannelID}),
	}

	messages, npt, err := h.database.ListMessages(ctx, filter, pageSize, tok)
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
// @Tags    API-v2
//
// @Param   uid path     int true "UserID"
//
// @Success 200 {object} handler.ListUserSubscriptions.response
// @Failure 400 {object} ginresp.apiError "supplied values/parameters cannot be parsed / are invalid"
// @Failure 401 {object} ginresp.apiError "user is not authorized / has missing permissions"
// @Failure 500 {object} ginresp.apiError "internal server error"
//
// @Router  /api/users/{uid}/subscriptions [GET]
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
// @Tags    API-v2
//
// @Param   uid path     int true "UserID"
// @Param   cid path     int true "ChannelID"
//
// @Success 200 {object} handler.ListChannelSubscriptions.response
// @Failure 400 {object} ginresp.apiError "supplied values/parameters cannot be parsed / are invalid"
// @Failure 401 {object} ginresp.apiError "user is not authorized / has missing permissions"
// @Failure 404 {object} ginresp.apiError "channel not found"
// @Failure 500 {object} ginresp.apiError "internal server error"
//
// @Router  /api/users/{uid}/channels/{cid}/subscriptions [GET]
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
// @Tags    API-v2
//
// @Param   uid path     int true "UserID"
// @Param   sid path     int true "SubscriptionID"
//
// @Success 200 {object} models.SubscriptionJSON
// @Failure 400 {object} ginresp.apiError "supplied values/parameters cannot be parsed / are invalid"
// @Failure 401 {object} ginresp.apiError "user is not authorized / has missing permissions"
// @Failure 404 {object} ginresp.apiError "subscription not found"
// @Failure 500 {object} ginresp.apiError "internal server error"
//
// @Router  /api/users/{uid}/subscriptions/{sid} [GET]
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
// @Tags    API-v2
//
// @Param   uid path     int true "UserID"
// @Param   sid path     int true "SubscriptionID"
//
// @Success 200 {object} models.SubscriptionJSON
// @Failure 400 {object} ginresp.apiError "supplied values/parameters cannot be parsed / are invalid"
// @Failure 401 {object} ginresp.apiError "user is not authorized / has missing permissions"
// @Failure 404 {object} ginresp.apiError "subscription not found"
// @Failure 500 {object} ginresp.apiError "internal server error"
//
// @Router  /api/users/{uid}/subscriptions/{sid} [DELETE]
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
// @Tags    API-v2
//
// @Param   uid        path     int                              true  "UserID"
// @Param   query_data query    handler.CreateSubscription.query false " "
// @Param   post_data  body     handler.CreateSubscription.body  false " "
//
// @Success 200        {object} models.SubscriptionJSON
// @Failure 400        {object} ginresp.apiError "supplied values/parameters cannot be parsed / are invalid"
// @Failure 401        {object} ginresp.apiError "user is not authorized / has missing permissions"
// @Failure 500        {object} ginresp.apiError "internal server error"
//
// @Router  /api/users/{uid}/subscriptions [POST]
func (h APIHandler) CreateSubscription(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID models.UserID `uri:"uid"`
	}
	type body struct {
		ChannelOwnerUserID models.UserID `form:"channel_owner_user_id" binding:"required"`
		Channel            string        `form:"channel_name" binding:"required"`
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
// @Tags    API-v2
//
// @Param   uid path     int true "UserID"
// @Param   sid path     int true "SubscriptionID"
//
// @Success 200 {object} models.SubscriptionJSON
// @Failure 400 {object} ginresp.apiError "supplied values/parameters cannot be parsed / are invalid"
// @Failure 401 {object} ginresp.apiError "user is not authorized / has missing permissions"
// @Failure 404 {object} ginresp.apiError "subscription not found"
// @Failure 500 {object} ginresp.apiError "internal server error"
//
// @Router  /api/users/{uid}/subscriptions/{sid} [PATCH]
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
// @Tags        API-v2
//
// @Param       query_data query    handler.ListMessages.query false " "
//
// @Success     200        {object} handler.ListMessages.response
// @Failure     400        {object} ginresp.apiError "supplied values/parameters cannot be parsed / are invalid"
// @Failure     401        {object} ginresp.apiError "user is not authorized / has missing permissions"
// @Failure     500        {object} ginresp.apiError "internal server error"
//
// @Router      /api/messages [GET]
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

	filter := models.MessageFilter{
		ConfirmedSubscriptionBy: langext.Ptr(userid),
	}

	if q.Filter != nil && strings.TrimSpace(*q.Filter) != "" {
		filter.SearchString = langext.Ptr([]string{strings.TrimSpace(*q.Filter)})
	}

	messages, npt, err := h.database.ListMessages(ctx, filter, pageSize, tok)
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
// @Tags        API-v2
//
// @Param       mid path     int true "SCNMessageID"
//
// @Success     200 {object} models.MessageJSON
// @Failure     400 {object} ginresp.apiError "supplied values/parameters cannot be parsed / are invalid"
// @Failure     401 {object} ginresp.apiError "user is not authorized / has missing permissions"
// @Failure     404 {object} ginresp.apiError "message not found"
// @Failure     500 {object} ginresp.apiError "internal server error"
//
// @Router      /api/messages/{mid} [PATCH]
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

	msg, err := h.database.GetMessage(ctx, u.MessageID, false)
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
// @Tags        API-v2
//
// @Param       mid path     int true "SCNMessageID"
//
// @Success     200 {object} models.MessageJSON
// @Failure     400 {object} ginresp.apiError "supplied values/parameters cannot be parsed / are invalid"
// @Failure     401 {object} ginresp.apiError "user is not authorized / has missing permissions"
// @Failure     404 {object} ginresp.apiError "message not found"
// @Failure     500 {object} ginresp.apiError "internal server error"
//
// @Router      /api/messages/{mid} [DELETE]
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

	msg, err := h.database.GetMessage(ctx, u.MessageID, false)
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
