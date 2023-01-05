package handler

import (
	"blackforestbytes.com/simplecloudnotifier/api/ginresp"
	primarydb "blackforestbytes.com/simplecloudnotifier/db/impl/primary"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
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

// Register swaggerdoc
//
// @Summary Register a new account
// @ID      compat-register
// @Tags    API-v1
//
// @Deprecated
//
// @Param   fcm_token query    string true "the (android) fcm token"
// @Param   pro       query    string true "if the user is a paid account" Enums(true, false)
// @Param   pro_token query    string true "the (android) IAP token"
//
// @Param   fcm_token formData string true "the (android) fcm token"
// @Param   pro       formData string true "if the user is a paid account" Enums(true, false)
// @Param   pro_token formData string true "the (android) IAP token"
//
// @Success 200       {object} handler.Register.response
// @Failure 200       {object} ginresp.compatAPIError
//
// @Router  /api/register.php [get]
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

	if data.FCMToken == nil {
		return ginresp.CompatAPIError(0, "Missing parameter [[fcm_token]]")
	}
	if data.Pro == nil {
		return ginresp.CompatAPIError(0, "Missing parameter [[pro]]")
	}
	if data.ProToken == nil {
		return ginresp.CompatAPIError(0, "Missing parameter [[pro_token]]")
	}

	if *data.Pro != "true" {
		data.ProToken = nil
	}

	if data.ProToken != nil {
		ptok, err := h.app.VerifyProToken(ctx, "ANDROID|v2|"+*data.ProToken)
		if err != nil {
			return ginresp.CompatAPIError(0, "Failed to query purchase status")
		}

		if !ptok {
			return ginresp.CompatAPIError(0, "Purchase token could not be verified")
		}
	}

	readKey := h.app.GenerateRandomAuthKey()
	sendKey := h.app.GenerateRandomAuthKey()
	adminKey := h.app.GenerateRandomAuthKey()

	err := h.database.ClearFCMTokens(ctx, *data.FCMToken)
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to clear existing fcm tokens")
	}

	if data.ProToken != nil {
		err := h.database.ClearProTokens(ctx, *data.ProToken)
		if err != nil {
			return ginresp.CompatAPIError(0, "Failed to clear existing fcm tokens")
		}
	}

	user, err := h.database.CreateUser(ctx, readKey, sendKey, adminKey, data.ProToken, nil)
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to create user in db")
	}

	_, err = h.database.CreateClient(ctx, user.UserID, models.ClientTypeAndroid, *data.FCMToken, "compat", "compat")
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to create client in db")
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, response{
		Success:   true,
		Message:   "New user registered",
		UserID:    user.UserID.IntID(),
		UserKey:   user.AdminKey,
		QuotaUsed: user.QuotaUsedToday(),
		QuotaMax:  user.QuotaPerDay(),
		IsPro:     langext.Conditional(user.IsPro, 1, 0),
	}))
}

// Info swaggerdoc
//
// @Summary Get information about the current user
// @ID      compat-info
// @Tags    API-v1
//
// @Deprecated
//
// @Param   user_id  query    string true "the user_id"
// @Param   user_key query    string true "the user_key"
//
// @Param   user_id  formData string true "the user_id"
// @Param   user_key formData string true "the user_key"
//
// @Success 200      {object} handler.Info.response
// @Failure 200      {object} ginresp.compatAPIError
//
// @Router  /api/info.php [get]
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

	user, err := h.database.GetUser(ctx, models.UserID(*data.UserID))
	if err == sql.ErrNoRows {
		return ginresp.CompatAPIError(201, "User not found")
	}
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to query user")
	}

	if user.AdminKey != *data.UserKey {
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
		UserID:     user.UserID.IntID(),
		UserKey:    user.AdminKey,
		QuotaUsed:  user.QuotaUsedToday(),
		QuotaMax:   user.QuotaPerDay(),
		IsPro:      langext.Conditional(user.IsPro, 1, 0),
		FCMSet:     fcmSet,
		UnackCount: 0,
	}))
}

// Ack swaggerdoc
//
// @Summary Acknowledge that a message was received
// @ID      compat-ack
// @Tags    API-v1
//
// @Deprecated
//
// @Param   user_id    query    string true "the user_id"
// @Param   user_key   query    string true "the user_key"
// @Param   scn_msg_id query    string true "the message id"
//
// @Param   user_id    formData string true "the user_id"
// @Param   user_key   formData string true "the user_key"
// @Param   scn_msg_id formData string true "the message id"
//
// @Success 200        {object} handler.Ack.response
// @Failure 200        {object} ginresp.compatAPIError
//
// @Router  /api/ack.php [get]
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

	user, err := h.database.GetUser(ctx, models.UserID(*data.UserID))
	if err == sql.ErrNoRows {
		return ginresp.CompatAPIError(201, "User not found")
	}
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to query user")
	}

	if user.AdminKey != *data.UserKey {
		return ginresp.CompatAPIError(204, "Authentification failed")
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, response{
		Success:      true,
		Message:      "ok",
		PrevAckValue: 0,
		NewAckValue:  1,
	}))
}

// Requery swaggerdoc
//
// @Summary Return all not-acknowledged messages
// @ID      compat-requery
// @Tags    API-v1
//
// @Deprecated
//
// @Param   user_id  query    string true "the user_id"
// @Param   user_key query    string true "the user_key"
//
// @Param   user_id  formData string true "the user_id"
// @Param   user_key formData string true "the user_key"
//
// @Success 200      {object} handler.Requery.response
// @Failure 200      {object} ginresp.compatAPIError
//
// @Router  /api/requery.php [get]
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

	user, err := h.database.GetUser(ctx, models.UserID(*data.UserID))
	if err == sql.ErrNoRows {
		return ginresp.CompatAPIError(201, "User not found")
	}
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to query user")
	}

	if user.AdminKey != *data.UserKey {
		return ginresp.CompatAPIError(204, "Authentification failed")
	}

	return ctx.FinishSuccess(ginresp.JSON(http.StatusOK, response{
		Success: true,
		Message: "ok",
		Count:   0,
		Data:    make([]models.CompatMessage, 0),
	}))
}

// Update swaggerdoc
//
// @Summary Set the fcm-token (android)
// @ID      compat-update
// @Tags    API-v1
//
// @Deprecated
//
// @Param   user_id   query    string true "the user_id"
// @Param   user_key  query    string true "the user_key"
// @Param   fcm_token query    string true "the (android) fcm token"
//
// @Param   user_id   formData string true "the user_id"
// @Param   user_key  formData string true "the user_key"
// @Param   fcm_token formData string true "the (android) fcm token"
//
// @Success 200       {object} handler.Update.response
// @Failure 200       {object} ginresp.compatAPIError
//
// @Router  /api/update.php [get]
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

	user, err := h.database.GetUser(ctx, models.UserID(*data.UserID))
	if err == sql.ErrNoRows {
		return ginresp.CompatAPIError(201, "User not found")
	}
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to query user")
	}

	if user.AdminKey != *data.UserKey {
		return ginresp.CompatAPIError(204, "Authentification failed")
	}

	clients, err := h.database.ListClients(ctx, user.UserID)
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to list clients")
	}

	newAdminKey := h.app.GenerateRandomAuthKey()
	newReadKey := h.app.GenerateRandomAuthKey()
	newSendKey := h.app.GenerateRandomAuthKey()

	err = h.database.UpdateUserKeys(ctx, user.UserID, newSendKey, newReadKey, newAdminKey)
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
		UserID:    user.UserID.IntID(),
		UserKey:   user.AdminKey,
		QuotaUsed: user.QuotaUsedToday(),
		QuotaMax:  user.QuotaPerDay(),
		IsPro:     langext.Conditional(user.IsPro, 1, 0),
	}))
}

// Expand swaggerdoc
//
// @Summary Get a whole (potentially truncated) message
// @ID      compat-expand
// @Tags    API-v1
//
// @Deprecated
//
// @Param   user_id    query    string true "The user_id"
// @Param   user_key   query    string true "The user_key"
// @Param   scn_msg_id query    string true "The message-id"
//
// @Param   user_id    formData string true "The user_id"
// @Param   user_key   formData string true "The user_key"
// @Param   scn_msg_id formData string true "The message-id"
//
// @Success 200        {object} handler.Expand.response
// @Failure 200        {object} ginresp.compatAPIError
//
// @Router  /api/expand.php [get]
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

	user, err := h.database.GetUser(ctx, models.UserID(*data.UserID))
	if err == sql.ErrNoRows {
		return ginresp.CompatAPIError(201, "User not found")
	}
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to query user")
	}

	if user.AdminKey != *data.UserKey {
		return ginresp.CompatAPIError(204, "Authentification failed")
	}

	msg, err := h.database.GetMessage(ctx, models.SCNMessageID(*data.MessageID), false)
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
			Title:         msg.Title,
			Body:          langext.Coalesce(msg.Content, ""),
			Trimmed:       langext.Ptr(false),
			Priority:      msg.Priority,
			Timestamp:     msg.Timestamp().Unix(),
			UserMessageID: msg.UserMessageID,
			SCNMessageID:  msg.SCNMessageID.IntID(),
		},
	}))
}

// Upgrade swaggerdoc
//
// @Summary Upgrade a free account to a paid account
// @ID      compat-upgrade
// @Tags    API-v1
//
// @Deprecated
//
// @Param   user_id   query    string true "the user_id"
// @Param   user_key  query    string true "the user_key"
// @Param   pro       query    string true "if the user is a paid account" Enums(true, false)
// @Param   pro_token query    string true "the (android) IAP token"
//
// @Param   user_id   formData string true "the user_id"
// @Param   user_key  formData string true "the user_key"
// @Param   pro       formData string true "if the user is a paid account" Enums(true, false)
// @Param   pro_token formData string true "the (android) IAP token"
//
// @Success 200       {object} handler.Upgrade.response
// @Failure 200       {object} ginresp.compatAPIError
//
// @Router  /api/upgrade.php [get]
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

	user, err := h.database.GetUser(ctx, models.UserID(*data.UserID))
	if err == sql.ErrNoRows {
		return ginresp.CompatAPIError(201, "User not found")
	}
	if err != nil {
		return ginresp.CompatAPIError(0, "Failed to query user")
	}

	if user.AdminKey != *data.UserKey {
		return ginresp.CompatAPIError(204, "Authentification failed")
	}

	if *data.Pro != "true" {
		data.ProToken = nil
	}

	if data.ProToken != nil {
		ptok, err := h.app.VerifyProToken(ctx, "ANDROID|v2|"+*data.ProToken)
		if err != nil {
			return ginresp.CompatAPIError(0, "Failed to query purchase status")
		}

		if !ptok {
			return ginresp.CompatAPIError(0, "Purchase token could not be verified")
		}

		err = h.database.UpdateUserProToken(ctx, user.UserID, data.ProToken)
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
		UserID:    user.UserID.IntID(),
		QuotaUsed: user.QuotaUsedToday(),
		QuotaMax:  user.QuotaPerDay(),
		IsPro:     user.IsPro,
	}))
}
