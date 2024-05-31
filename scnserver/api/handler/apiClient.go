package handler

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/api/ginresp"
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"net/http"
)

// ListClients swaggerdoc
//
//	@Summary	List all clients
//	@ID			api-clients-list
//	@Tags		API-v2
//
//	@Param		uid	path		string	true	"UserID"
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
//	@Param		uid	path		string	true	"UserID"
//	@Param		cid	path		string	true	"ClientID"
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
	if errors.Is(err, sql.ErrNoRows) {
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
//	@Param		uid			path		string					true	"UserID"
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
		FCMToken     string            `json:"fcm_token" binding:"required"`
		AgentModel   string            `json:"agent_model" binding:"required"`
		AgentVersion string            `json:"agent_version" binding:"required"`
		ClientType   models.ClientType `json:"client_type" binding:"required"`
	}

	var u uri
	var b body
	ctx, errResp := h.app.StartRequest(g, &u, nil, &b, nil)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	if !b.ClientType.Valid() {
		return ginresp.APIError(g, 400, apierr.INVALID_CLIENTTYPE, "Invalid ClientType", nil)
	}
	clientType := b.ClientType

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
//	@Param		uid	path		string	true	"UserID"
//	@Param		cid	path		string	true	"ClientID"
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
	if errors.Is(err, sql.ErrNoRows) {
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

// UpdateClient swaggerdoc
//
//	@Summary		(Partially) update a client
//	@Description	The body-values are optional, only send the ones you want to update
//	@ID				api-client-update
//	@Tags			API-v2
//
//	@Param			uid			path		string	true	"UserID"
//	@Param			cid			path		string	true	"ClientID"
//
//	@Param			clientname	body		string	false	"Change the clientname (send an empty string to clear it)"
//	@Param			pro_token	body		string	false	"Send a verification of premium purchase"
//
//	@Success		200			{object}	models.ClientJSON
//	@Failure		400			{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure		401			{object}	ginresp.apiError	"client is not authorized / has missing permissions"
//	@Failure		404			{object}	ginresp.apiError	"client not found"
//	@Failure		500			{object}	ginresp.apiError	"internal server error"
//
//	@Router			/api/v2/users/{uid}/clients/{cid} [PATCH]
func (h APIHandler) UpdateClient(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID   models.UserID   `uri:"uid" binding:"entityid"`
		ClientID models.ClientID `uri:"cid" binding:"entityid"`
	}
	type body struct {
		FCMToken     *string `json:"fcm_token"`
		AgentModel   *string `json:"agent_model"`
		AgentVersion *string `json:"agent_version"`
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

	client, err := h.database.GetClient(ctx, u.UserID, u.ClientID)
	if errors.Is(err, sql.ErrNoRows) {
		return ginresp.APIError(g, 404, apierr.CLIENT_NOT_FOUND, "Client not found", err)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query client", err)
	}

	if b.FCMToken != nil && *b.FCMToken != client.FCMToken {

		err = h.database.DeleteClientsByFCM(ctx, *b.FCMToken)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to delete existing clients in db", err)
		}

		err = h.database.UpdateClientFCMToken(ctx, u.ClientID, *b.FCMToken)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to update client", err)
		}
	}

	if b.AgentModel != nil {
		err = h.database.UpdateClientAgentModel(ctx, u.ClientID, *b.AgentModel)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to update client", err)
		}
	}

	if b.AgentVersion != nil {
		err = h.database.UpdateClientAgentVersion(ctx, u.ClientID, *b.AgentVersion)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to update client", err)
		}
	}

	client, err = h.database.GetClient(ctx, u.UserID, u.ClientID)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query (updated) client", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, client.JSON()))
}
