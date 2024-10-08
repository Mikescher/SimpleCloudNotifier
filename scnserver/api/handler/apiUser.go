package handler

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/api/ginresp"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
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
//	@Success	200			{object}	models.UserWithClientsAndKeys
//	@Failure	400			{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure	500			{object}	ginresp.apiError	"internal server error"
//
//	@Router		/api/v2/users [POST]
func (h APIHandler) CreateUser(pctx ginext.PreContext) ginext.HTTPResponse {
	type body struct {
		FCMToken     string            `json:"fcm_token"`
		ProToken     *string           `json:"pro_token"`
		Username     *string           `json:"username"`
		AgentModel   string            `json:"agent_model"`
		AgentVersion string            `json:"agent_version"`
		ClientName   *string           `json:"client_name"`
		ClientType   models.ClientType `json:"client_type"`
		NoClient     bool              `json:"no_client"`
	}

	var b body
	ctx, g, errResp := pctx.Body(&b).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	return h.app.DoRequest(ctx, g, models.TLockReadWrite, func(ctx *logic.AppContext, finishSuccess func(r ginext.HTTPResponse) ginext.HTTPResponse) ginext.HTTPResponse {

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
			if !b.ClientType.Valid() {
				return ginresp.APIError(g, 400, apierr.BINDFAIL_BODY_PARAM, "Invalid ClientType", nil)
			}
			clientType = b.ClientType
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

		err := h.database.DeleteClientsByFCM(ctx, b.FCMToken)
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

		log.Info().Msg(fmt.Sprintf("Sucessfully created new user %s (client: %v)", userobj.UserID, b.NoClient))

		if b.NoClient {
			return finishSuccess(ginext.JSON(http.StatusOK, userobj.PreMarshal().WithClients(make([]models.Client, 0), adminKey, sendKey, readKey)))
		} else {
			err := h.database.DeleteClientsByFCM(ctx, b.FCMToken)
			if err != nil {
				return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to delete existing clients in db", err)
			}

			client, err := h.database.CreateClient(ctx, userobj.UserID, clientType, b.FCMToken, b.AgentModel, b.AgentVersion, b.ClientName)
			if err != nil {
				return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to create client in db", err)
			}

			return finishSuccess(ginext.JSON(http.StatusOK, userobj.PreMarshal().WithClients([]models.Client{client}, adminKey, sendKey, readKey)))
		}
	})
}

// GetUser swaggerdoc
//
//	@Summary	Get a user
//	@ID			api-user-get
//	@Tags		API-v2
//
//	@Param		uid	path		string	true	"UserID"
//
//	@Success	200	{object}	models.User
//	@Failure	400	{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure	401	{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure	404	{object}	ginresp.apiError	"user not found"
//	@Failure	500	{object}	ginresp.apiError	"internal server error"
//
//	@Router		/api/v2/users/{uid} [GET]
func (h APIHandler) GetUser(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
		UserID models.UserID `uri:"uid" binding:"entityid"`
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

		user, err := h.database.GetUser(ctx, u.UserID)
		if errors.Is(err, sql.ErrNoRows) {
			return ginresp.APIError(g, 404, apierr.USER_NOT_FOUND, "User not found", err)
		}
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to query user", err)
		}

		return finishSuccess(ginext.JSON(http.StatusOK, user.PreMarshal()))

	})

}

// UpdateUser swaggerdoc
//
//	@Summary		(Partially) update a user
//	@Description	The body-values are optional, only send the ones you want to update
//	@ID				api-user-update
//	@Tags			API-v2
//
//	@Param			uid			path		string	true	"UserID"
//
//	@Param			username	body		string	false	"Change the username (send an empty string to clear it)"
//	@Param			pro_token	body		string	false	"Send a verification of premium purchase"
//
//	@Success		200			{object}	models.User
//	@Failure		400			{object}	ginresp.apiError	"supplied values/parameters cannot be parsed / are invalid"
//	@Failure		401			{object}	ginresp.apiError	"user is not authorized / has missing permissions"
//	@Failure		404			{object}	ginresp.apiError	"user not found"
//	@Failure		500			{object}	ginresp.apiError	"internal server error"
//
//	@Router			/api/v2/users/{uid} [PATCH]
func (h APIHandler) UpdateUser(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
		UserID models.UserID `uri:"uid" binding:"entityid"`
	}
	type body struct {
		Username *string `json:"username"`
		ProToken *string `json:"pro_token"`
	}

	var u uri
	var b body
	ctx, g, errResp := pctx.URI(&u).Body(&b).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	return h.app.DoRequest(ctx, g, models.TLockReadWrite, func(ctx *logic.AppContext, finishSuccess func(r ginext.HTTPResponse) ginext.HTTPResponse) ginext.HTTPResponse {

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

		return finishSuccess(ginext.JSON(http.StatusOK, user.PreMarshal()))
	})
}
