package handler

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/api/ginresp"
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetUserPreview swaggerdoc
//
//	@Summary	Get a user (similar to api-user-get, but can be called from anyone and only returns a subset of fields)
//	@ID			api-user-get-preview
//	@Tags		API-v2
//
//	@Param		uid	path		string	true	"UserID"
//
//	@Success	200	{object}	models.UserPreviewJSON
//	@Failure	400	{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure	401	{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure	404	{object}	ginresp.apiError	"user not found"
//	@Failure	500	{object}	ginresp.apiError	"internal server error"
//
//	@Router		/api/v2/preview/users/{uid} [GET]
func (h APIHandler) GetUserPreview(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		UserID models.UserID `uri:"uid" binding:"entityid"`
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

	user, err := h.database.GetUser(ctx, u.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		return ginresp.APIError(g, 404, apierr.USER_NOT_FOUND, "User not found", err)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query user", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, user.JSONPreview()))
}

// GetChannelPreview swaggerdoc
//
//	@Summary	Get a single channel (similar to api-channels-get, but can be called from anyone and only returns a subset of fields)
//	@ID			api-channels-get-preview
//	@Tags		API-v2
//
//	@Param		cid	path		string	true	"ChannelID"
//
//	@Success	200	{object}	models.ChannelPreviewJSON
//	@Failure	400	{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure	401	{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure	404	{object}	ginresp.apiError	"channel not found"
//	@Failure	500	{object}	ginresp.apiError	"internal server error"
//
//	@Router		/api/v2/preview/channels/{cid} [GET]
func (h APIHandler) GetChannelPreview(g *gin.Context) ginresp.HTTPResponse {
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

	if permResp := ctx.CheckPermissionAny(); permResp != nil {
		return *permResp
	}

	channel, err := h.database.GetChannelByID(ctx, u.ChannelID)
	if errors.Is(err, sql.ErrNoRows) {
		return ginresp.APIError(g, 404, apierr.CHANNEL_NOT_FOUND, "Channel not found", err)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query channel", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, channel.JSONPreview()))
}

// GetUserKeyPreview swaggerdoc
//
//	@Summary	Get a single key (similar to api-tokenkeys-get, but can be called from anyone and only returns a subset of fields)
//	@ID			api-tokenkeys-get-preview
//	@Tags		API-v2
//
//	@Param		kid	path		string	true	"TokenKeyID"
//
//	@Success	200	{object}	models.KeyTokenPreviewJSON
//	@Failure	400	{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure	401	{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure	404	{object}	ginresp.apiError	"message not found"
//	@Failure	500	{object}	ginresp.apiError	"internal server error"
//
//	@Router		/api/v2/preview/keys/{kid} [GET]
func (h APIHandler) GetUserKeyPreview(g *gin.Context) ginresp.HTTPResponse {
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

	if permResp := ctx.CheckPermissionAny(); permResp != nil {
		return *permResp
	}

	keytoken, err := h.database.GetKeyToken(ctx, u.UserID, u.KeyID)
	if errors.Is(err, sql.ErrNoRows) {
		return ginresp.APIError(g, 404, apierr.KEY_NOT_FOUND, "Key not found", err)
	}
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query client", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, keytoken.JSONPreview()))
}
