package handler

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	hl "blackforestbytes.com/simplecloudnotifier/api/apihighlight"
	"blackforestbytes.com/simplecloudnotifier/api/ginresp"
	ct "blackforestbytes.com/simplecloudnotifier/db/cursortoken"
	primarydb "blackforestbytes.com/simplecloudnotifier/db/impl/primary"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"gogs.mikescher.com/BlackForestBytes/goext/dataext"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"net/http"
)

type CompatHandler struct {
	app      *logic.Application
	database *primarydb.Database
}

func NewCompatHandler(app *logic.Application) CompatHandler {
	return CompatHandler{
		app:      app,
		database: app.Database.Primary,
	}
}

// SendMessageCompat swaggerdoc
//
//	@Deprecated
//
//	@Summary		Send a new message (compatibility)
//	@Description	All parameter can be set via query-parameter or form-data body. Only UserID, UserKey and Title are required
//	@Tags			External
//
//	@Param			query_data	query		handler.SendMessageCompat.combined	false	" "
//	@Param			form_data	formData	handler.SendMessageCompat.combined	false	" "
//
//	@Success		200			{object}	handler.SendMessageCompat.response
//	@Failure		400			{object}	ginresp.apiError
//	@Failure		401			{object}	ginresp.apiError
//	@Failure		403			{object}	ginresp.apiError
//	@Failure		500			{object}	ginresp.apiError
//
//	@Router			/send.php [POST]
func (h MessageHandler) SendMessageCompat(g *gin.Context) ginresp.HTTPResponse {
	type combined struct {
		UserID        *int64   `json:"user_id"   form:"user_id"`
		UserKey       *string  `json:"user_key"  form:"user_key"`
		Title         *string  `json:"title"     form:"title"`
		Content       *string  `json:"content"   form:"content"`
		Priority      *int     `json:"priority"  form:"priority"`
		UserMessageID *string  `json:"msg_id"    form:"msg_id"`
		SendTimestamp *float64 `json:"timestamp" form:"timestamp"`
	}
	type response struct {
		Success        bool            `json:"success"`
		ErrorID        apierr.APIError `json:"error"`
		ErrorHighlight int             `json:"errhighlight"`
		Message        string          `json:"message"`
		SuppressSend   bool            `json:"suppress_send"`
		MessageCount   int             `json:"messagecount"`
		Quota          int             `json:"quota"`
		IsPro          bool            `json:"is_pro"`
		QuotaMax       int             `json:"quota_max"`
		SCNMessageID   int64           `json:"scn_msg_id"`
	}

	var f combined
	var q combined
	ctx, errResp := h.app.StartRequest(g, nil, &q, nil, &f)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	data := dataext.ObjectMerge(f, q)

	newid, err := h.database.ConvertCompatID(ctx, langext.Coalesce(data.UserID, -1), "userid")
	if err != nil {
		return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query userid<old>", err)
	}
	if newid == nil {
		return ginresp.SendAPIError(g, 400, apierr.USER_NOT_FOUND, hl.USER_ID, "User not found (compat)", nil)
	}

	okResp, errResp := h.sendMessageInternal(g, ctx, langext.Ptr(models.UserID(*newid)), data.UserKey, nil, data.Title, data.Content, data.Priority, data.UserMessageID, data.SendTimestamp, nil)
	if errResp != nil {
		return *errResp
	} else {
		if okResp.MessageIsOld {
			return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, response{
				Success:        true,
				ErrorID:        apierr.NO_ERROR,
				ErrorHighlight: -1,
				Message:        "Message already sent",
				SuppressSend:   true,
				MessageCount:   okResp.User.MessagesSent,
				Quota:          okResp.User.QuotaUsedToday(),
				IsPro:          okResp.User.IsPro,
				QuotaMax:       okResp.User.QuotaPerDay(),
				SCNMessageID:   okResp.CompatMessageID,
			}))
		} else {
			return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, response{
				Success:        true,
				ErrorID:        apierr.NO_ERROR,
				ErrorHighlight: -1,
				Message:        "Message sent",
				SuppressSend:   false,
				MessageCount:   okResp.User.MessagesSent + 1,
				Quota:          okResp.User.QuotaUsedToday() + 1,
				IsPro:          okResp.User.IsPro,
				QuotaMax:       okResp.User.QuotaPerDay(),
				SCNMessageID:   okResp.CompatMessageID,
			}))
		}
	}
}

// Register swaggerdoc
//
//	@Summary	Register a new account
//	@ID			compat-register
//	@Tags		API-v1
//
//	@Deprecated
//
//	@Param		fcm_token	query		string	true	"the (android) fcm token"
//	@Param		pro			query		string	true	"if the user is a paid account"	Enums(true, false)
//	@Param		pro_token	query		string	true	"the (android) IAP token"
//
//	@Param		fcm_token	formData	string	true	"the (android) fcm token"
//	@Param		pro			formData	string	true	"if the user is a paid account"	Enums(true, false)
//	@Param		pro_token	formData	string	true	"the (android) IAP token"
//
//	@Success	200			{object}	handler.Register.response
//	@Failure	default		{object}	ginresp.compatAPIError
//
//	@Router		/api/register.php [get]
func (h CompatHandler) Register(g *gin.Context) ginresp.HTTPResponse {
	type query struct {
		FCMToken *string `json:"fcm_token" form:"fcm_token"`
		Pro      *string `json:"pro"       form:"pro"`
		ProToken *string `json:"pro_token" form:"pro_token"`
	}
	type response struct {
		Success   bool   `json:"success"`
		Message   string `json:"message"`
		UserID    int64  `json:"user_id"`
		UserKey   string `json:"user_key"`
		QuotaUsed int    `json:"quota"`
		QuotaMax  int    `json:"quota_max"`
		IsPro     bool   `json:"is_pro"`
	}

	var datq query
	var datb query
	ctx, errResp := h.app.StartRequest(g, nil, &datq, nil, &datb)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	data := dataext.ObjectMerge(datb, datq)

	if data.FCMToken == nil {
		return ginresp.CompatAPIError(0, "Missing parameter [[fcm_token]]")
	}
	if data.Pro == nil {
		return ginresp.CompatAPIError(0, "Missing parameter [[pro]]")
	}
	if data.ProToken == nil {
		return ginresp.CompatAPIError(0, "Missing parameter [[pro_token]]")
	}

	if data.ProToken != nil {
		data.ProToken = langext.Ptr("ANDROID|v1|" + *data.ProToken)
	}

	if *data.Pro != "true" {
		data.ProToken = nil
	}

	if data.ProToken != nil {
		ptok, err := h.app.VerifyProToken(ctx, *data.ProToken)
		if err != nil {
			return ginresp.CompatAPIError(0, "Failed to query purchase status")
		}

		if !ptok {
			return ginresp.CompatAPIError(0, "Purchase token could not be verified")
		}
	}

	adminKey := h.app.GenerateRandomAuthKey()

	err := h.database.ClearFCMTokens(ctx, *data.FCMToken)
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to clear existing fcm tokens")
	}

	if data.ProToken != nil {
		err := h.database.ClearProTokens(ctx, *data.ProToken)
		if err != nil {
			return ginresp.CompatAPIError(0, "Failed to clear existing pro tokens")
		}
	}

	user, err := h.database.CreateUser(ctx, data.ProToken, nil)
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to create user in db")
	}

	_, err = h.database.CreateKeyToken(ctx, "CompatKey", user.UserID, true, make([]models.ChannelID, 0), models.TokenPermissionList{models.PermAdmin}, adminKey)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to create admin-key in db", err)
	}

	_, err = h.database.CreateClient(ctx, user.UserID, models.ClientTypeAndroid, *data.FCMToken, "compat", "compat")
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to create client in db")
	}

	oldid, err := h.database.CreateCompatID(ctx, "userid", user.UserID.String())
	if err != nil {
		return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to create userid<old>", err)
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, response{
		Success:   true,
		Message:   "New user registered",
		UserID:    oldid,
		UserKey:   adminKey,
		QuotaUsed: user.QuotaUsedToday(),
		QuotaMax:  user.QuotaPerDay(),
		IsPro:     user.IsPro,
	}))
}

// Info swaggerdoc
//
//	@Summary	Get information about the current user
//	@ID			compat-info
//	@Tags		API-v1
//
//	@Deprecated
//
//	@Param		user_id		query		string	true	"the user_id"
//	@Param		user_key	query		string	true	"the user_key"
//
//	@Param		user_id		formData	string	true	"the user_id"
//	@Param		user_key	formData	string	true	"the user_key"
//
//	@Success	200			{object}	handler.Info.response
//	@Failure	default		{object}	ginresp.compatAPIError
//
//	@Router		/api/info.php [get]
func (h CompatHandler) Info(g *gin.Context) ginresp.HTTPResponse {
	type query struct {
		UserID  *int64  `json:"user_id"  form:"user_id"`
		UserKey *string `json:"user_key" form:"user_key"`
	}
	type response struct {
		Success    bool   `json:"success"`
		Message    string `json:"message"`
		UserID     int64  `json:"user_id"`
		UserKey    string `json:"user_key"`
		QuotaUsed  int    `json:"quota"`
		QuotaMax   int    `json:"quota_max"`
		IsPro      int    `json:"is_pro"`
		FCMSet     bool   `json:"fcm_token_set"`
		UnackCount int    `json:"unack_count"`
	}

	var datq query
	var datb query
	ctx, errResp := h.app.StartRequest(g, nil, &datq, nil, &datb)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	data := dataext.ObjectMerge(datb, datq)

	if data.UserID == nil {
		return ginresp.CompatAPIError(101, "Missing parameter [[user_id]]")
	}
	if data.UserKey == nil {
		return ginresp.CompatAPIError(102, "Missing parameter [[user_key]]")
	}

	useridCompNew, err := h.database.ConvertCompatID(ctx, *data.UserID, "userid")
	if err != nil {
		return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query userid<old>", err)
	}
	if useridCompNew == nil {
		return ginresp.SendAPIError(g, 400, apierr.USER_NOT_FOUND, hl.USER_ID, "User not found (compat)", nil)
	}

	user, err := h.database.GetUser(ctx, models.UserID(*useridCompNew))
	if err == sql.ErrNoRows {
		return ginresp.CompatAPIError(201, "User not found")
	}
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to query user")
	}

	keytok, err := h.database.GetKeyTokenByToken(ctx, *data.UserKey)
	if err == sql.ErrNoRows {
		return ginresp.CompatAPIError(204, "Authentification failed")
	}
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to query token")
	}
	if !keytok.IsAdmin(user.UserID) {
		return ginresp.CompatAPIError(204, "Authentification failed")
	}

	clients, err := h.database.ListClients(ctx, user.UserID)
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to query clients")
	}

	fcmSet := langext.ArrAny(clients, func(c models.Client) bool { return c.FCMToken != nil })

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, response{
		Success:    true,
		Message:    "ok",
		UserID:     *data.UserID,
		UserKey:    keytok.Token,
		QuotaUsed:  user.QuotaUsedToday(),
		QuotaMax:   user.QuotaPerDay(),
		IsPro:      langext.Conditional(user.IsPro, 1, 0),
		FCMSet:     fcmSet,
		UnackCount: 0,
	}))
}

// Ack swaggerdoc
//
//	@Summary	Acknowledge that a message was received
//	@ID			compat-ack
//	@Tags		API-v1
//
//	@Deprecated
//
//	@Param		user_id		query		string	true	"the user_id"
//	@Param		user_key	query		string	true	"the user_key"
//	@Param		scn_msg_id	query		string	true	"the message id"
//
//	@Param		user_id		formData	string	true	"the user_id"
//	@Param		user_key	formData	string	true	"the user_key"
//	@Param		scn_msg_id	formData	string	true	"the message id"
//
//	@Success	200			{object}	handler.Ack.response
//	@Failure	default		{object}	ginresp.compatAPIError
//
//	@Router		/api/ack.php [get]
func (h CompatHandler) Ack(g *gin.Context) ginresp.HTTPResponse {
	type query struct {
		UserID    *int64  `json:"user_id"    form:"user_id"`
		UserKey   *string `json:"user_key"   form:"user_key"`
		MessageID *int64  `json:"scn_msg_id" form:"scn_msg_id"`
	}
	type response struct {
		Success      bool   `json:"success"`
		Message      string `json:"message"`
		PrevAckValue int    `json:"prev_ack"`
		NewAckValue  int    `json:"new_ack"`
	}

	var datq query
	var datb query
	ctx, errResp := h.app.StartRequest(g, nil, &datq, nil, &datb)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	data := dataext.ObjectMerge(datb, datq)

	if data.UserID == nil {
		return ginresp.CompatAPIError(101, "Missing parameter [[user_id]]")
	}
	if data.UserKey == nil {
		return ginresp.CompatAPIError(102, "Missing parameter [[user_key]]")
	}
	if data.MessageID == nil {
		return ginresp.CompatAPIError(103, "Missing parameter [[scn_msg_id]]")
	}

	useridCompNew, err := h.database.ConvertCompatID(ctx, *data.UserID, "userid")
	if err != nil {
		return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query userid<old>", err)
	}
	if useridCompNew == nil {
		return ginresp.SendAPIError(g, 400, apierr.USER_NOT_FOUND, hl.USER_ID, "User not found (compat)", nil)
	}

	user, err := h.database.GetUser(ctx, models.UserID(*useridCompNew))
	if err == sql.ErrNoRows {
		return ginresp.CompatAPIError(201, "User not found")
	}
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to query user")
	}

	keytok, err := h.database.GetKeyTokenByToken(ctx, *data.UserKey)
	if err == sql.ErrNoRows {
		return ginresp.CompatAPIError(204, "Authentification failed")
	}
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to query token")
	}
	if !keytok.IsAdmin(user.UserID) {
		return ginresp.CompatAPIError(204, "Authentification failed")
	}

	messageIdComp, err := h.database.ConvertCompatID(ctx, *data.MessageID, "messageid")
	if err != nil {
		return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query messageid<old>", err)
	}
	if useridCompNew == nil {
		return ginresp.SendAPIError(g, 400, apierr.MESSAGE_NOT_FOUND, hl.USER_ID, "Message not found (compat)", nil)
	}

	ackBefore, err := h.database.GetAck(ctx, models.MessageID(*messageIdComp))
	if err != nil {
		return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query ack", err)
	}

	if !ackBefore {
		err = h.database.SetAck(ctx, user.UserID, models.MessageID(*messageIdComp))
		if err != nil {
			return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to set ack", err)
		}
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, response{
		Success:      true,
		Message:      "ok",
		PrevAckValue: langext.Conditional(ackBefore, 1, 0),
		NewAckValue:  1,
	}))
}

// Requery swaggerdoc
//
//	@Summary	Return all not-acknowledged messages
//	@ID			compat-requery
//	@Tags		API-v1
//
//	@Deprecated
//
//	@Param		user_id		query		string	true	"the user_id"
//	@Param		user_key	query		string	true	"the user_key"
//
//	@Param		user_id		formData	string	true	"the user_id"
//	@Param		user_key	formData	string	true	"the user_key"
//
//	@Success	200			{object}	handler.Requery.response
//	@Failure	default		{object}	ginresp.compatAPIError
//
//	@Router		/api/requery.php [get]
func (h CompatHandler) Requery(g *gin.Context) ginresp.HTTPResponse {
	type query struct {
		UserID  *int64  `json:"user_id"  form:"user_id"`
		UserKey *string `json:"user_key" form:"user_key"`
	}
	type response struct {
		Success bool                   `json:"success"`
		Message string                 `json:"message"`
		Count   int                    `json:"count"`
		Data    []models.CompatMessage `json:"data"`
	}

	var datq query
	var datb query
	ctx, errResp := h.app.StartRequest(g, nil, &datq, nil, &datb)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	data := dataext.ObjectMerge(datb, datq)

	if data.UserID == nil {
		return ginresp.CompatAPIError(101, "Missing parameter [[user_id]]")
	}
	if data.UserKey == nil {
		return ginresp.CompatAPIError(102, "Missing parameter [[user_key]]")
	}

	useridCompNew, err := h.database.ConvertCompatID(ctx, *data.UserID, "userid")
	if err != nil {
		return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query userid<old>", err)
	}
	if useridCompNew == nil {
		return ginresp.SendAPIError(g, 400, apierr.USER_NOT_FOUND, hl.USER_ID, "User not found (compat)", nil)
	}

	user, err := h.database.GetUser(ctx, models.UserID(*useridCompNew))
	if err == sql.ErrNoRows {
		return ginresp.CompatAPIError(201, "User not found")
	}
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to query user")
	}

	keytok, err := h.database.GetKeyTokenByToken(ctx, *data.UserKey)
	if err == sql.ErrNoRows {
		return ginresp.CompatAPIError(204, "Authentification failed")
	}
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to query token")
	}
	if !keytok.IsAdmin(user.UserID) {
		return ginresp.CompatAPIError(204, "Authentification failed")
	}

	filter := models.MessageFilter{
		Owner:              langext.Ptr([]models.UserID{user.UserID}),
		CompatAcknowledged: langext.Ptr(false),
	}

	msgs, _, err := h.database.ListMessages(ctx, filter, langext.Ptr(16), ct.Start())
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to query user")
	}

	compMsgs := make([]models.CompatMessage, 0, len(msgs))
	for _, v := range msgs {

		messageIdComp, err := h.database.ConvertToCompatIDOrCreate(ctx, v.MessageID.String(), "messageid")
		if err != nil {
			return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query/create messageid<old>", err)
		}

		compMsgs = append(compMsgs, models.CompatMessage{
			Title:         compatizeMessageTitle(ctx, h.app, v),
			Body:          v.Content,
			Priority:      v.Priority,
			Timestamp:     v.Timestamp().Unix(),
			UserMessageID: v.UserMessageID,
			SCNMessageID:  messageIdComp,
			Trimmed:       nil,
		})
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, response{
		Success: true,
		Message: "ok",
		Count:   len(compMsgs),
		Data:    compMsgs,
	}))
}

// Update swaggerdoc
//
//	@Summary	Set the fcm-token (android)
//	@ID			compat-update
//	@Tags		API-v1
//
//	@Deprecated
//
//	@Param		user_id		query		string	true	"the user_id"
//	@Param		user_key	query		string	true	"the user_key"
//	@Param		fcm_token	query		string	true	"the (android) fcm token"
//
//	@Param		user_id		formData	string	true	"the user_id"
//	@Param		user_key	formData	string	true	"the user_key"
//	@Param		fcm_token	formData	string	true	"the (android) fcm token"
//
//	@Success	200			{object}	handler.Update.response
//	@Failure	default		{object}	ginresp.compatAPIError
//
//	@Router		/api/update.php [get]
func (h CompatHandler) Update(g *gin.Context) ginresp.HTTPResponse {
	type query struct {
		UserID   *int64  `json:"user_id"   form:"user_id"`
		UserKey  *string `json:"user_key"  form:"user_key"`
		FCMToken *string `json:"fcm_token" form:"fcm_token"`
	}
	type response struct {
		Success   bool   `json:"success"`
		Message   string `json:"message"`
		UserID    int64  `json:"user_id"`
		UserKey   string `json:"user_key"`
		QuotaUsed int    `json:"quota"`
		QuotaMax  int    `json:"quota_max"`
		IsPro     int    `json:"is_pro"`
	}

	var datq query
	var datb query
	ctx, errResp := h.app.StartRequest(g, nil, &datq, nil, &datb)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	data := dataext.ObjectMerge(datb, datq)

	if data.UserID == nil {
		return ginresp.CompatAPIError(101, "Missing parameter [[user_id]]")
	}
	if data.UserKey == nil {
		return ginresp.CompatAPIError(102, "Missing parameter [[user_key]]")
	}

	useridCompNew, err := h.database.ConvertCompatID(ctx, *data.UserID, "userid")
	if err != nil {
		return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query userid<old>", err)
	}
	if useridCompNew == nil {
		return ginresp.SendAPIError(g, 400, apierr.USER_NOT_FOUND, hl.USER_ID, "User not found (compat)", nil)
	}

	user, err := h.database.GetUser(ctx, models.UserID(*useridCompNew))
	if err == sql.ErrNoRows {
		return ginresp.CompatAPIError(201, "User not found")
	}
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to query user")
	}

	keytok, err := h.database.GetKeyTokenByToken(ctx, *data.UserKey)
	if err == sql.ErrNoRows {
		return ginresp.CompatAPIError(204, "Authentification failed")
	}
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to query token")
	}
	if !keytok.IsAdmin(user.UserID) {
		return ginresp.CompatAPIError(204, "Authentification failed")
	}

	clients, err := h.database.ListClients(ctx, user.UserID)
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to list clients")
	}

	newAdminKey := h.app.GenerateRandomAuthKey()

	_, err = h.database.CreateKeyToken(ctx, "CompatKey", user.UserID, true, make([]models.ChannelID, 0), models.TokenPermissionList{models.PermAdmin}, newAdminKey)
	if err != nil {
		return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to create admin-key in db", err)
	}

	err = h.database.DeleteKeyToken(ctx, keytok.KeyTokenID)
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to update keys")
	}

	if data.FCMToken != nil {

		for _, client := range clients {

			err = h.database.DeleteClient(ctx, client.ClientID)
			if err != nil {
				return ginresp.CompatAPIError(0, "Failed to delete client")
			}

		}

		_, err = h.database.CreateClient(ctx, user.UserID, models.ClientTypeAndroid, *data.FCMToken, "compat", "compat")
		if err != nil {
			return ginresp.CompatAPIError(0, "Failed to delete client")
		}

	}

	user, err = h.database.GetUser(ctx, user.UserID)
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to query user")
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, response{
		Success:   true,
		Message:   "user updated",
		UserID:    *data.UserID,
		UserKey:   newAdminKey,
		QuotaUsed: user.QuotaUsedToday(),
		QuotaMax:  user.QuotaPerDay(),
		IsPro:     langext.Conditional(user.IsPro, 1, 0),
	}))
}

// Expand swaggerdoc
//
//	@Summary	Get a whole (potentially truncated) message
//	@ID			compat-expand
//	@Tags		API-v1
//
//	@Deprecated
//
//	@Param		user_id		query		string	true	"The user_id"
//	@Param		user_key	query		string	true	"The user_key"
//	@Param		scn_msg_id	query		string	true	"The message-id"
//
//	@Param		user_id		formData	string	true	"The user_id"
//	@Param		user_key	formData	string	true	"The user_key"
//	@Param		scn_msg_id	formData	string	true	"The message-id"
//
//	@Success	200			{object}	handler.Expand.response
//	@Failure	default		{object}	ginresp.compatAPIError
//
//	@Router		/api/expand.php [get]
func (h CompatHandler) Expand(g *gin.Context) ginresp.HTTPResponse {
	type query struct {
		UserID    *int64  `json:"user_id"    form:"user_id"`
		UserKey   *string `json:"user_key"   form:"user_key"`
		MessageID *int64  `json:"scn_msg_id" form:"scn_msg_id"`
	}
	type response struct {
		Success bool                 `json:"success"`
		Message string               `json:"message"`
		Data    models.CompatMessage `json:"data"`
	}

	var datq query
	var datb query
	ctx, errResp := h.app.StartRequest(g, nil, &datq, nil, &datb)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	data := dataext.ObjectMerge(datb, datq)

	if data.UserID == nil {
		return ginresp.CompatAPIError(101, "Missing parameter [[user_id]]")
	}
	if data.UserKey == nil {
		return ginresp.CompatAPIError(102, "Missing parameter [[user_key]]")
	}
	if data.MessageID == nil {
		return ginresp.CompatAPIError(103, "Missing parameter [[scn_msg_id]]")
	}

	useridCompNew, err := h.database.ConvertCompatID(ctx, *data.UserID, "userid")
	if err != nil {
		return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query userid<old>", err)
	}
	if useridCompNew == nil {
		return ginresp.SendAPIError(g, 400, apierr.USER_NOT_FOUND, hl.USER_ID, "User not found (compat)", nil)
	}

	user, err := h.database.GetUser(ctx, models.UserID(*useridCompNew))
	if err == sql.ErrNoRows {
		return ginresp.CompatAPIError(201, "User not found")
	}
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to query user")
	}

	keytok, err := h.database.GetKeyTokenByToken(ctx, *data.UserKey)
	if err == sql.ErrNoRows {
		return ginresp.CompatAPIError(204, "Authentification failed")
	}
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to query token")
	}
	if !keytok.IsAdmin(user.UserID) {
		return ginresp.CompatAPIError(204, "Authentification failed")
	}

	messageCompNew, err := h.database.ConvertCompatID(ctx, *data.MessageID, "messageid")
	if err != nil {
		return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query messagid<old>", err)
	}
	if messageCompNew == nil {
		return ginresp.CompatAPIError(301, "Message not found")
	}

	msg, err := h.database.GetMessage(ctx, models.MessageID(*messageCompNew), false)
	if err == sql.ErrNoRows {
		return ginresp.CompatAPIError(301, "Message not found")
	}
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to query message")
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, response{
		Success: true,
		Message: "ok",
		Data: models.CompatMessage{
			Title:         compatizeMessageTitle(ctx, h.app, msg),
			Body:          msg.Content,
			Trimmed:       langext.Ptr(false),
			Priority:      msg.Priority,
			Timestamp:     msg.Timestamp().Unix(),
			UserMessageID: msg.UserMessageID,
			SCNMessageID:  *data.MessageID,
		},
	}))
}

// Upgrade swaggerdoc
//
//	@Summary	Upgrade a free account to a paid account
//	@ID			compat-upgrade
//	@Tags		API-v1
//
//	@Deprecated
//
//	@Param		user_id		query		string	true	"the user_id"
//	@Param		user_key	query		string	true	"the user_key"
//	@Param		pro			query		string	true	"if the user is a paid account"	Enums(true, false)
//	@Param		pro_token	query		string	true	"the (android) IAP token"
//
//	@Param		user_id		formData	string	true	"the user_id"
//	@Param		user_key	formData	string	true	"the user_key"
//	@Param		pro			formData	string	true	"if the user is a paid account"	Enums(true, false)
//	@Param		pro_token	formData	string	true	"the (android) IAP token"
//
//	@Success	200			{object}	handler.Upgrade.response
//	@Failure	default		{object}	ginresp.compatAPIError
//
//	@Router		/api/upgrade.php [get]
func (h CompatHandler) Upgrade(g *gin.Context) ginresp.HTTPResponse {
	type query struct {
		UserID   *int64  `json:"user_id"   form:"user_id"`
		UserKey  *string `json:"user_key"  form:"user_key"`
		Pro      *string `json:"pro"       form:"pro"`
		ProToken *string `json:"pro_token" form:"pro_token"`
	}
	type response struct {
		Success   bool   `json:"success"`
		Message   string `json:"message"`
		UserID    int64  `json:"user_id"`
		QuotaUsed int    `json:"quota"`
		QuotaMax  int    `json:"quota_max"`
		IsPro     bool   `json:"is_pro"`
	}

	var datq query
	var datb query
	ctx, errResp := h.app.StartRequest(g, nil, &datq, nil, &datb)
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	data := dataext.ObjectMerge(datb, datq)

	if data.UserID == nil {
		return ginresp.CompatAPIError(101, "Missing parameter [[user_id]]")
	}
	if data.UserKey == nil {
		return ginresp.CompatAPIError(102, "Missing parameter [[user_key]]")
	}
	if data.Pro == nil {
		return ginresp.CompatAPIError(103, "Missing parameter [[pro]]")
	}
	if data.ProToken == nil {
		return ginresp.CompatAPIError(104, "Missing parameter [[pro_token]]")
	}

	useridCompNew, err := h.database.ConvertCompatID(ctx, *data.UserID, "userid")
	if err != nil {
		return ginresp.SendAPIError(g, 500, apierr.DATABASE_ERROR, hl.NONE, "Failed to query userid<old>", err)
	}
	if useridCompNew == nil {
		return ginresp.SendAPIError(g, 400, apierr.USER_NOT_FOUND, hl.USER_ID, "User not found (compat)", nil)
	}

	user, err := h.database.GetUser(ctx, models.UserID(*useridCompNew))
	if err == sql.ErrNoRows {
		return ginresp.CompatAPIError(201, "User not found")
	}
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to query user")
	}

	keytok, err := h.database.GetKeyTokenByToken(ctx, *data.UserKey)
	if err == sql.ErrNoRows {
		return ginresp.CompatAPIError(204, "Authentification failed")
	}
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to query token")
	}
	if !keytok.IsAdmin(user.UserID) {
		return ginresp.CompatAPIError(204, "Authentification failed")
	}

	if data.ProToken != nil {
		data.ProToken = langext.Ptr("ANDROID|v1|" + *data.ProToken)
	}

	if *data.Pro != "true" {
		data.ProToken = nil
	}

	if data.ProToken != nil {
		ptok, err := h.app.VerifyProToken(ctx, *data.ProToken)
		if err != nil {
			return ginresp.CompatAPIError(0, "Failed to query purchase status")
		}

		if !ptok {
			return ginresp.CompatAPIError(0, "Purchase token could not be verified")
		}

		err = h.database.ClearProTokens(ctx, *data.ProToken)
		if err != nil {
			return ginresp.APIError(g, 500, apierr.DATABASE_ERROR, "Failed to clear existing fcm tokens", err)
		}

		err = h.database.UpdateUserProToken(ctx, user.UserID, langext.Ptr(*data.ProToken))
		if err != nil {
			return ginresp.CompatAPIError(0, "Failed to update user")
		}
	} else {
		err = h.database.UpdateUserProToken(ctx, user.UserID, nil)
		if err != nil {
			return ginresp.CompatAPIError(0, "Failed to update user")
		}
	}

	user, err = h.database.GetUser(ctx, user.UserID)
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to query user")
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, response{
		Success:   true,
		Message:   "user updated",
		UserID:    *data.UserID,
		QuotaUsed: user.QuotaUsedToday(),
		QuotaMax:  user.QuotaPerDay(),
		IsPro:     user.IsPro,
	}))
}

func compatizeMessageTitle(ctx *logic.AppContext, app *logic.Application, msg models.Message) string {
	if msg.ChannelInternalName == "main" {
		return msg.Title
	}

	channel, err := app.Database.Primary.GetChannelByID(ctx, msg.ChannelID)
	if err != nil {
		return fmt.Sprintf("[%s] %s", "%SCN-ERR%", msg.Title)
	}

	return fmt.Sprintf("[%s] %s", channel.DisplayName, msg.Title)
}
