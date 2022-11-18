package handler

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/api/models"
	"blackforestbytes.com/simplecloudnotifier/common/ginresp"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

type APIHandler struct {
	app *logic.Application
}

// CreateUser swaggerdoc
//
// @Summary Create a new user
// @ID      api-user-create
//
// @Param       post_body body     handler.CreateUser.body  false " "
//
// @Success 200 {object} models.UserJSON
// @Failure 400 {object} ginresp.apiError
// @Failure 500 {object} ginresp.apiError
//
// @Router  /api-v2/user/ [POST]
func (h APIHandler) CreateUser(g *gin.Context) ginresp.HTTPResponse {
	type body struct {
		FCMToken     string  `form:"fcm_token"`
		ProToken     *string `form:"pro_token"`
		Username     *string `form:"username"`
		AgentModel   string  `form:"agent_model"`
		AgentVersion string  `form:"agent_version"`
		ClientType   string  `form:"client_type"`
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
		return ginresp.InternAPIError(apierr.INVALID_CLIENTTYPE, "Invalid ClientType", nil)
	}

	if b.ProToken != nil {
		ptok, err := h.app.VerifyProToken(*b.ProToken)
		if err != nil {
			return ginresp.InternAPIError(apierr.FAILED_VERIFY_PRO_TOKEN, "Failed to query purchase status", err)
		}

		if !ptok {
			return ginresp.InternAPIError(apierr.INVALID_PRO_TOKEN, "Purchase token could not be verified", nil)
		}
	}

	readKey := h.app.GenerateRandomAuthKey()
	sendKey := h.app.GenerateRandomAuthKey()
	adminKey := h.app.GenerateRandomAuthKey()

	err := h.app.Database.ClearFCMTokens(ctx, b.FCMToken)
	if err != nil {
		return ginresp.InternAPIError(apierr.DATABASE_ERROR, "Failed to clear existing fcm tokens", err)
	}

	if b.ProToken != nil {
		err := h.app.Database.ClearProTokens(ctx, b.FCMToken)
		if err != nil {
			return ginresp.InternAPIError(apierr.DATABASE_ERROR, "Failed to clear existing fcm tokens", err)
		}
	}

	userobj, err := h.app.Database.CreateUser(ctx, readKey, sendKey, adminKey, b.ProToken, b.Username)
	if err != nil {
		return ginresp.InternAPIError(apierr.DATABASE_ERROR, "Failed to create user in db", err)
	}

	_, err = h.app.Database.CreateClient(ctx, userobj.UserID, clientType, b.FCMToken, b.AgentModel, b.AgentVersion)
	if err != nil {
		return ginresp.InternAPIError(apierr.DATABASE_ERROR, "Failed to create user in db", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, userobj.JSON()))
}

// GetUser swaggerdoc
//
// @Summary Create a new user
// @ID      api-user-create
//
// @Param       post_body body     handler.CreateUser.body  false " "
// @Param       uid       path     int                      true  "UserID"
//
// @Success 200 {object} models.UserJSON
// @Failure 400 {object} ginresp.apiError
// @Failure 401 {object} ginresp.apiError
// @Failure 404 {object} ginresp.apiError
// @Failure 500 {object} ginresp.apiError
//
// @Router  /api-v2/user/{uid} [GET]
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

	user, err := h.app.Database.GetUser(ctx, u.UserID)
	if err == sql.ErrNoRows {
		return ginresp.InternAPIError(apierr.USER_NOT_FOUND, "User not found", err)
	}
	if err != nil {
		return ginresp.InternAPIError(apierr.DATABASE_ERROR, "Failed to query user", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, user.JSON()))
}

func (h APIHandler) UpdateUser(g *gin.Context) ginresp.HTTPResponse {
	return ginresp.NotImplemented()
}

func (h APIHandler) ListClients(g *gin.Context) ginresp.HTTPResponse {
	return ginresp.NotImplemented()
}

func (h APIHandler) GetClient(g *gin.Context) ginresp.HTTPResponse {
	return ginresp.NotImplemented()
}

func (h APIHandler) AddClient(g *gin.Context) ginresp.HTTPResponse {
	return ginresp.NotImplemented()
}

func (h APIHandler) DeleteClient(g *gin.Context) ginresp.HTTPResponse {
	return ginresp.NotImplemented()
}

func (h APIHandler) ListChannels(g *gin.Context) ginresp.HTTPResponse {
	return ginresp.NotImplemented()
}

func (h APIHandler) GetChannel(g *gin.Context) ginresp.HTTPResponse {
	return ginresp.NotImplemented()
}

func (h APIHandler) GetChannelMessages(g *gin.Context) ginresp.HTTPResponse {
	return ginresp.NotImplemented()
}

func (h APIHandler) ListUserSubscriptions(g *gin.Context) ginresp.HTTPResponse {
	return ginresp.NotImplemented()
}

func (h APIHandler) ListChannelSubscriptions(g *gin.Context) ginresp.HTTPResponse {
	return ginresp.NotImplemented()
}

func (h APIHandler) GetSubscription(g *gin.Context) ginresp.HTTPResponse {
	return ginresp.NotImplemented()
}

func (h APIHandler) CancelSubscription(g *gin.Context) ginresp.HTTPResponse {
	return ginresp.NotImplemented()
}

func (h APIHandler) CreateSubscription(g *gin.Context) ginresp.HTTPResponse {
	return ginresp.NotImplemented()
}

func (h APIHandler) ListMessages(g *gin.Context) ginresp.HTTPResponse {
	return ginresp.NotImplemented()
}

func (h APIHandler) GetMessage(g *gin.Context) ginresp.HTTPResponse {
	return ginresp.NotImplemented()
}

func (h APIHandler) DeleteMessage(g *gin.Context) ginresp.HTTPResponse {
	return ginresp.NotImplemented()
}

func NewAPIHandler(app *logic.Application) APIHandler {
	return APIHandler{
		app: app,
	}
}
