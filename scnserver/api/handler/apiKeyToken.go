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

// ListUserKeys swaggerdoc
//
//	@Summary		List keys of the user
//	@Description	The request must be done with an ADMIN key, the returned keys are without their token.
//	@ID				api-tokenkeys-list
//	@Tags			API-v2
//
//	@Param			uid	path		string	true	"UserID"
//
//	@Success		200	{object}	handler.ListUserKeys.response
//	@Failure		400	{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure		401	{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure		404	{object}	ginresp.apiError	"message not found"
//	@Failure		500	{object}	ginresp.apiError	"internal server error"
//
//	@Router			/api/v2/users/{uid}/keys [GET]
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

	toks, err := h.database.ListKeyTokens(ctx, u.UserID)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query keys", err)
	}

	res := langext.ArrMap(toks, func(v models.KeyToken) models.KeyTokenJSON { return v.JSON() })

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, response{Keys: res}))
}

// GetUserKey swaggerdoc
//
//	@Summary		Get a single key
//	@Description	The request must be done with an ADMIN key, the returned key does not include its token.
//	@ID				api-tokenkeys-get
//	@Tags			API-v2
//
//	@Param			uid	path		string	true	"UserID"
//	@Param			kid	path		string	true	"TokenKeyID"
//
//	@Success		200	{object}	models.KeyTokenJSON
//	@Failure		400	{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure		401	{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure		404	{object}	ginresp.apiError	"message not found"
//	@Failure		500	{object}	ginresp.apiError	"internal server error"
//
//	@Router			/api/v2/users/{uid}/keys/{kid} [GET]
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
	if errors.Is(err, sql.ErrNoRows) {
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
//	@Param		uid			path		string						true	"UserID"
//	@Param		kid			path		string						true	"TokenKeyID"
//
//	@Param		post_body	body		handler.UpdateUserKey.body	false	" "
//
//	@Success	200			{object}	models.KeyTokenJSON
//	@Failure	400			{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure	401			{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure	404			{object}	ginresp.apiError	"message not found"
//	@Failure	500			{object}	ginresp.apiError	"internal server error"
//
//	@Router		/api/v2/users/{uid}/keys/{kid} [PATCH]
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
	if errors.Is(err, sql.ErrNoRows) {
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
//	@Param		uid			path		string						true	"UserID"
//
//	@Param		post_body	body		handler.CreateUserKey.body	false	" "
//
//	@Success	200			{object}	models.KeyTokenJSON
//	@Failure	400			{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure	401			{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure	404			{object}	ginresp.apiError	"message not found"
//	@Failure	500			{object}	ginresp.apiError	"internal server error"
//
//	@Router		/api/v2/users/{uid}/keys [POST]
func (h APIHandler) CreateUserKey(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID models.UserID `uri:"uid" binding:"entityid"`
	}
	type body struct {
		Name        string              `json:"name"         binding:"required"`
		Permissions string              `json:"permissions"  binding:"required"`
		AllChannels *bool               `json:"all_channels"`
		Channels    *[]models.ChannelID `json:"channels"`
	}

	var u uri
	var b body
	ctx, errResp := h.app.StartRequest(g, &u, nil, &b, nil)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	channels := langext.Coalesce(b.Channels, make([]models.ChannelID, 0))

	var allChan bool
	if b.AllChannels == nil && b.Channels != nil {
		allChan = false
	} else if b.AllChannels == nil && b.Channels == nil {
		allChan = true
	} else {
		allChan = *b.AllChannels
	}

	for _, c := range channels {
		if err := c.Valid(); err != nil {
			return ginresp.APIError(g, 400, apierr.INVALID_BODY_PARAM, "Invalid ChannelID", err)
		}
	}

	if permResp := ctx.CheckPermissionUserAdmin(u.UserID); permResp != nil {
		return *permResp
	}

	token := h.app.GenerateRandomAuthKey()

	perms := models.ParseTokenPermissionList(b.Permissions)

	keytok, err := h.database.CreateKeyToken(ctx, b.Name, *ctx.GetPermissionUserID(), allChan, channels, perms, token)
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
//	@Param			uid	path		string	true	"UserID"
//	@Param			kid	path		string	true	"TokenKeyID"
//
//	@Success		200	{object}	models.KeyTokenJSON
//	@Failure		400	{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure		401	{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure		404	{object}	ginresp.apiError	"message not found"
//	@Failure		500	{object}	ginresp.apiError	"internal server error"
//
//	@Router			/api/v2/users/{uid}/keys/{kid} [DELETE]
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
	if errors.Is(err, sql.ErrNoRows) {
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
