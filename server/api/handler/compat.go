package handler

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
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
// @Param   fcm_token query    string true "the (android) fcm token"
// @Param   pro       query    string true "if the user is a paid account" Enums(true, false)
// @Param   pro_token query    string true "the (android) IAP token"
// @Success 200       {object} handler.Register.response
// @Failure 500       {object} ginresp.internAPIError
// @Router  /register.php [get]
func (h CompatHandler) Register(g *gin.Context) ginresp.HTTPResponse {
	type query struct {
		FCMToken string `form:"fcm_token"`
		Pro      string `form:"pro"`
		ProToken string `form:"pro_token"`
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

	return ginresp.NotImplemented(0)
}

// Info swaggerdoc
//
// @Summary Get information about the current user
// @Param   user_id  query    string true "the user_id"
// @Param   user_key query    string true "the user_key"
// @Success 200      {object} handler.Info.response
// @Failure 500      {object} ginresp.internAPIError
// @Router  /info.php [get]
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

	return ginresp.NotImplemented(0)
}

// Ack swaggerdoc
//
// @Summary Acknowledge that a message was received
// @Param   user_id    query    string true "the user_id"
// @Param   user_key   query    string true "the user_key"
// @Param   scn_msg_id query    string true "the message id"
// @Success 200        {object} handler.Ack.response
// @Failure 500        {object} ginresp.internAPIError
// @Router  /ack.php [get]
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

	return ginresp.NotImplemented(0)
}

// Requery swaggerdoc
//
// @Summary Return all not-acknowledged messages
// @Param   user_id  query    string true "the user_id"
// @Param   user_key query    string true "the user_key"
// @Success 200      {object} handler.Requery.response
// @Failure 500      {object} ginresp.internAPIError
// @Router  /requery.php [get]
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

	return ginresp.NotImplemented(0)
}

// Update swaggerdoc
//
// @Summary Set the fcm-token (android)
// @Param   user_id   query    string true "the user_id"
// @Param   user_key  query    string true "the user_key"
// @Param   fcm_token query    string true "the (android) fcm token"
// @Success 200       {object} handler.Update.response
// @Failure 500       {object} ginresp.internAPIError
// @Router  /update.php [get]
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

	return ginresp.NotImplemented(0)
}

// Expand swaggerdoc
//
// @Summary Get a whole (potentially truncated) message
// @Success 200 {object} handler.Expand.response
// @Failure 500 {object} ginresp.internAPIError
// @Router  /expand.php [get]
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

	return ginresp.NotImplemented(0)
}

// Upgrade swaggerdoc
//
// @Summary Upgrade a free account to a paid account
// @Param   user_id   query    string true "the user_id"
// @Param   user_key  query    string true "the user_key"
// @Param   pro       query    string true "if the user is a paid account" Enums(true, false)
// @Param   pro_token query    string true "the (android) IAP token"
// @Success 200       {object} handler.Upgrade.response
// @Failure 500       {object} ginresp.internAPIError
// @Router  /upgrade.php [get]
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

	return ginresp.NotImplemented(0)
}

// Send swaggerdoc
//
// @Summary Send a message
// @Description all aeguments can either be supplied in the query or in the json body
// @Param   user_id   query    string true "the user_id"
// @Param   user_key  query    string true "the user_key"
// @Param   title     query    string true "The message title"
// @Param   content   query    string false "The message content"
// @Param   priority  query    string false "The message priority" Enum(0, 1, 2)
// @Param   msg_id    query    string false "The message idempotency id"
// @Param   timestamp query    string false "The message timestamp"
// @Param   user_id   body     string true "the user_id"
// @Param   user_key  body     string true "the user_key"
// @Param   title     body     string true "The message title"
// @Param   content   body     string false "The message content"
// @Param   priority  body     string false "The message priority" Enum(0, 1, 2)
// @Param   msg_id    body     string false "The message idempotency id"
// @Param   timestamp body     string false "The message timestamp"
// @Success 200       {object} handler.Send.response
// @Failure 500       {object} ginresp.sendAPIError
// @Router  /send.php [post]
func (h CompatHandler) Send(g *gin.Context) ginresp.HTTPResponse {
	type query struct {
		//TODO
	}
	type response struct {
		Success string `json:"success"`
		Message string `json:"message"`
		//TODO
	}

	//TODO

	return ginresp.SendAPIError(apierr.INTERNAL_EXCEPTION, -1, "NotImplemented")
}
