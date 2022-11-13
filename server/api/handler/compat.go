package handler

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/api/models"
	"blackforestbytes.com/simplecloudnotifier/common/ginresp"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
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
// @Failure 500       {object} ginresp.internAPIError
// @Router  /register.php [get]
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var q query
	if err := g.ShouldBindQuery(&q); err != nil {
		return ginresp.InternAPIError(0, "Failed to read arguments")
	}

	if q.FCMToken == nil {
		return ginresp.InternAPIError(0, "Missing parameter [[fcm_token]]")
	}
	if q.Pro == nil {
		return ginresp.InternAPIError(0, "Missing parameter [[pro]]")
	}
	if q.ProToken == nil {
		return ginresp.InternAPIError(0, "Missing parameter [[pro_token]]")
	}

	isProInt := 0
	isProBool := false
	if *q.Pro == "true" {
		isProInt = 1
		isProBool = true
	} else {
		q.ProToken = nil
	}

	if isProBool {
		ptok, err := h.app.VerifyProToken(*q.ProToken)
		if err != nil {
			return ginresp.InternAPIError(0, fmt.Sprintf("Failed to query purchaste status: %v", err))
		}

		if !ptok {
			return ginresp.InternAPIError(0, "Purchase token could not be verified")
		}
	}

	userKey := h.app.GenerateRandomAuthKey()

	return h.app.RunTransaction(ctx, nil, func(tx *sql.Tx) (ginresp.HTTPResponse, bool) {

		res, err := tx.ExecContext(ctx, "INSERT INTO users (user_key, fcm_token, is_pro, pro_token, timestamp_accessed) VALUES (?, ?, ?, ?, NOW())", userKey, *q.FCMToken, isProInt, q.ProToken)
		if err != nil {
			return ginresp.InternAPIError(0, fmt.Sprintf("Failed to create user: %v", err)), false
		}

		userId, err := res.LastInsertId()
		if err != nil {
			return ginresp.InternAPIError(0, fmt.Sprintf("Failed to get user_id: %v", err)), false
		}

		_, err = tx.ExecContext(ctx, "UPDATE users SET fcm_token=NULL WHERE user_id <> ? AND fcm_token=?", userId, q.FCMToken)
		if err != nil {
			return ginresp.InternAPIError(0, fmt.Sprintf("Failed to update fcm: %v", err)), false
		}

		if isProInt == 1 {
			_, err := tx.ExecContext(ctx, "UPDATE users SET is_pro=0, pro_token=NULL WHERE user_id <> ? AND pro_token = ?", userId, q.ProToken)
			if err != nil {
				return ginresp.InternAPIError(0, fmt.Sprintf("Failed to update ispro: %v", err)), false
			}
		}

		return ginresp.JSON(http.StatusOK, response{
			Success:   true,
			Message:   "New user registered",
			UserID:    strconv.FormatInt(userId, 10),
			UserKey:   userKey,
			QuotaUsed: 0,
			QuotaMax:  h.app.QuotaMax(isProBool),
			IsPro:     isProInt,
		}), true

	})
}

// Info swaggerdoc
//
// @Summary Get information about the current user
// @ID      compat-info
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

	return ginresp.InternAPIError(0, "NotImplemented")
}

// Ack swaggerdoc
//
// @Summary Acknowledge that a message was received
// @ID      compat-ack
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

	return ginresp.InternAPIError(0, "NotImplemented")
}

// Requery swaggerdoc
//
// @Summary Return all not-acknowledged messages
// @ID      compat-requery
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

	return ginresp.InternAPIError(0, "NotImplemented")
}

// Update swaggerdoc
//
// @Summary Set the fcm-token (android)
// @ID      compat-update
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

	return ginresp.InternAPIError(0, "NotImplemented")
}

// Expand swaggerdoc
//
// @Summary Get a whole (potentially truncated) message
// @ID      compat-expand
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

	return ginresp.InternAPIError(0, "NotImplemented")
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

	return ginresp.InternAPIError(0, "NotImplemented")
}

// Send swaggerdoc
//
// @Summary     Send a message
// @Description (all arguments can either be supplied in the query or in the json body)
// @ID          compat-send
// @Accept      json
// @Produce     json
// @Param       _         query    handler.Send.query false " "
// @Param       post_body body     handler.Send.body  false " "
// @Success     200       {object} handler.Send.response
// @Failure     500       {object} ginresp.sendAPIError
// @Router      /send.php [post]
func (h CompatHandler) Send(g *gin.Context) ginresp.HTTPResponse {
	type query struct {
		UserID    string  `form:"user_id" required:"true"`
		UserKey   string  `form:"user_key" required:"true"`
		Title     string  `form:"title" required:"true"`
		Content   *string `form:"content"`
		Priority  *string `form:"priority"`
		MessageID *string `form:"msg_id"`
		Timestamp *string `form:"timestamp"`
	}
	type body struct {
		UserID    string  `json:"user_id" required:"true"`
		UserKey   string  `json:"user_key" required:"true"`
		Title     string  `json:"title" required:"true"`
		Content   *string `json:"content"`
		Priority  *string `json:"priority"`
		MessageID *string `json:"msg_id"`
		Timestamp *string `json:"timestamp"`
	}
	type response struct {
		Success string `json:"success"`
		Message string `json:"message"`
		//TODO
	}

	//TODO

	return ginresp.SendAPIError(apierr.INTERNAL_EXCEPTION, -1, "NotImplemented")
}
