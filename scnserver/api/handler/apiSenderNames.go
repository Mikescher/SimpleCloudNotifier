package handler

import (
	"blackforestbytes.com/simplecloudnotifier/logic"
	"blackforestbytes.com/simplecloudnotifier/models"
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
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

		return nil //TODO

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

		userid := *ctx.GetPermissionUserID()

		if permResp := ctx.CheckPermissionUserRead(userid); permResp != nil {
			return *permResp
		}

		return nil //TODO

	})
}
