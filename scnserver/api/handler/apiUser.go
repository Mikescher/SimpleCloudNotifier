package handler

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/api/ginresp"
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"github.com/gin-gonic/gin"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"net/http"
)

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
