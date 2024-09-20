package handler

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/api/ginresp"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"blackforestbytes.com/simplecloudnotifier/models"
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
	"net/http"
)

// ListUserSenderNames swaggerdoc
//
//	@Summary		List sender-names (of allthe messages of this user)
//	@ID				api-usersendernames-list
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
func (h APIHandler) ListUserSenderNames(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
		UserID models.UserID `uri:"uid" binding:"entityid"`
	}
	type response struct {
		SenderNames []models.SenderNameStatistics `json:"sender_names"`
	}

	var u uri
	ctx, g, errResp := pctx.URI(&u).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	return h.app.DoRequest(ctx, g, models.TLockRead, func(ctx *logic.AppContext, finishSuccess func(r ginext.HTTPResponse) ginext.HTTPResponse) ginext.HTTPResponse {

		if permResp := ctx.CheckPermissionUserRead(u.UserID); permResp != nil {
			return *permResp
		}

		names, err := h.database.ListSenderNames(ctx, u.UserID, false)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query messages", err)
		}

		return finishSuccess(ginext.JSON(http.StatusOK, response{SenderNames: names}))

	})
}

// ListSenderNames swaggerdoc
//
//	@Summary		List sender-names (of all messages this user can view, eitehr own or foreign-subscribed)
//	@ID				api-sendernames-list
//	@Tags			API-v2
//
//	@Success		200	{object}	handler.ListSenderNames.response
//	@Failure		400	{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure		401	{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure		404	{object}	ginresp.apiError	"message not found"
//	@Failure		500	{object}	ginresp.apiError	"internal server error"
//
//	@Router			/api/v2/sender-names [GET]
func (h APIHandler) ListSenderNames(pctx ginext.PreContext) ginext.HTTPResponse {
	type response struct {
		SenderNames []models.SenderNameStatistics `json:"sender_names"`
	}

	ctx, g, errResp := pctx.Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	return h.app.DoRequest(ctx, g, models.TLockRead, func(ctx *logic.AppContext, finishSuccess func(r ginext.HTTPResponse) ginext.HTTPResponse) ginext.HTTPResponse {

		if permResp := ctx.CheckPermissionAny(); permResp != nil {
			return *permResp
		}

		userID := *ctx.GetPermissionUserID()

		if permResp := ctx.CheckPermissionUserRead(userID); permResp != nil {
			return *permResp
		}

		names, err := h.database.ListSenderNames(ctx, userID, true)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query messages", err)
		}

		return finishSuccess(ginext.JSON(http.StatusOK, response{SenderNames: names}))

	})
}
