package handler

import (
	"blackforestbytes.com/simplecloudnotifier/api/models"
	"blackforestbytes.com/simplecloudnotifier/common/ginresp"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"github.com/gin-gonic/gin"
)

type CompatHandler struct {
	app *logic.Application
}

func NewCompatHandler(app *logic.Application) CompatHandler {
	return CompatHandler{
		app: app,
	}
}

// Register swaggerdoc
//
// @Summary Register a new account
// @ID      compat-register
// @Param   fcm_token query    string true "the (android) fcm token"
// @Param   pro       query    string true "if the user is a paid account" Enums(true, false)
// @Param   pro_token query    string true "the (android) IAP token"
// @Success 200       {object} handler.Register.response
// @Failure 500       {object} ginresp.apiError
// @Router  /api/register.php [get]
func (h CompatHandler) Register(g *gin.Context) ginresp.HTTPResponse {
	type query struct {
		FCMToken *string `form:"fcm_token"`
		Pro      *string `form:"pro"`
		ProToken *string `form:"pro_token"`
	}
	type response struct {
		Success   bool   `json:"success"`
		Message   string `json:"message"`
		UserID    string `json:"user_id"`
		UserKey   string `json:"user_key"`
		QuotaUsed int    `json:"quota"`
		QuotaMax  int    `json:"quota_max"`
		IsPro     int    `json:"is_pro"`
	}

	//TODO

	return ginresp.NotImplemented()
}

// Info swaggerdoc
//
// @Summary Get information about the current user
// @ID      compat-info
// @Param   user_id  query    string true "the user_id"
// @Param   user_key query    string true "the user_key"
// @Success 200      {object} handler.Info.response
// @Failure 500      {object} ginresp.apiError
// @Router  /api/info.php [get]
func (h CompatHandler) Info(g *gin.Context) ginresp.HTTPResponse {
	type query struct {
		UserID  string `form:"user_id"`
		UserKey string `form:"user_key"`
	}
	type response struct {
		Success    string `json:"success"`
		Message    string `json:"message"`
		UserID     string `json:"user_id"`
		UserKey    string `json:"user_key"`
		QuotaUsed  string `json:"quota"`
		QuotaMax   string `json:"quota_max"`
		IsPro      string `json:"is_pro"`
		FCMSet     bool   `json:"fcm_token_set"`
		UnackCount int    `json:"unack_count"`
	}

	//TODO

	return ginresp.NotImplemented()
}

// Ack swaggerdoc
//
// @Summary Acknowledge that a message was received
// @ID      compat-ack
// @Param   user_id    query    string true "the user_id"
// @Param   user_key   query    string true "the user_key"
// @Param   scn_msg_id query    string true "the message id"
// @Success 200        {object} handler.Ack.response
// @Failure 500        {object} ginresp.apiError
// @Router  /api/ack.php [get]
func (h CompatHandler) Ack(g *gin.Context) ginresp.HTTPResponse {
	type query struct {
		UserID    string `form:"user_id"`
		UserKey   string `form:"user_key"`
		MessageID string `form:"scn_msg_id"`
	}
	type response struct {
		Success      string `json:"success"`
		Message      string `json:"message"`
		PrevAckValue int    `json:"prev_ack"`
		NewAckValue  int    `json:"new_ack"`
	}

	//TODO

	return ginresp.NotImplemented()
}

// Requery swaggerdoc
//
// @Summary Return all not-acknowledged messages
// @ID      compat-requery
// @Param   user_id  query    string true "the user_id"
// @Param   user_key query    string true "the user_key"
// @Success 200      {object} handler.Requery.response
// @Failure 500      {object} ginresp.apiError
// @Router  /api/requery.php [get]
func (h CompatHandler) Requery(g *gin.Context) ginresp.HTTPResponse {
	type query struct {
		UserID  string `form:"user_id"`
		UserKey string `form:"user_key"`
	}
	type response struct {
		Success string                 `json:"success"`
		Message string                 `json:"message"`
		Count   int                    `json:"count"`
		Data    []models.CompatMessage `json:"data"`
	}

	//TODO

	return ginresp.NotImplemented()
}

// Update swaggerdoc
//
// @Summary Set the fcm-token (android)
// @ID      compat-update
// @Param   user_id   query    string true "the user_id"
// @Param   user_key  query    string true "the user_key"
// @Param   fcm_token query    string true "the (android) fcm token"
// @Success 200       {object} handler.Update.response
// @Failure 500       {object} ginresp.apiError
// @Router  /api/update.php [get]
func (h CompatHandler) Update(g *gin.Context) ginresp.HTTPResponse {
	type query struct {
		UserID   string `form:"user_id"`
		UserKey  string `form:"user_key"`
		FCMToken string `form:"fcm_token"`
	}
	type response struct {
		Success   string `json:"success"`
		Message   string `json:"message"`
		UserID    string `json:"user_id"`
		UserKey   string `json:"user_key"`
		QuotaUsed string `json:"quota"`
		QuotaMax  string `json:"quota_max"`
		IsPro     string `json:"is_pro"`
	}

	//TODO

	return ginresp.NotImplemented()
}

// Expand swaggerdoc
//
// @Summary Get a whole (potentially truncated) message
// @ID      compat-expand
// @Success 200 {object} handler.Expand.response
// @Failure 500 {object} ginresp.apiError
// @Router  /api/expand.php [get]
func (h CompatHandler) Expand(g *gin.Context) ginresp.HTTPResponse {
	type query struct {
		UserID    string `form:"user_id"`
		UserKey   string `form:"user_key"`
		MessageID string `form:"scn_msg_id"`
	}
	type response struct {
		Success string                    `json:"success"`
		Message string                    `json:"message"`
		Data    models.ShortCompatMessage `json:"data"`
	}

	//TODO

	return ginresp.NotImplemented()
}

// Upgrade swaggerdoc
//
// @Summary Upgrade a free account to a paid account
// @ID      compat-upgrade
// @Param   user_id   query    string true "the user_id"
// @Param   user_key  query    string true "the user_key"
// @Param   pro       query    string true "if the user is a paid account" Enums(true, false)
// @Param   pro_token query    string true "the (android) IAP token"
// @Success 200       {object} handler.Upgrade.response
// @Failure 500       {object} ginresp.apiError
// @Router  /api/upgrade.php [get]
func (h CompatHandler) Upgrade(g *gin.Context) ginresp.HTTPResponse {
	type query struct {
		UserID   string `form:"user_id"`
		UserKey  string `form:"user_key"`
		Pro      string `form:"pro"`
		ProToken string `form:"pro_token"`
	}
	type response struct {
		Success string                    `json:"success"`
		Message string                    `json:"message"`
		Data    models.ShortCompatMessage `json:"data"`
	}

	//TODO

	return ginresp.NotImplemented()
}
