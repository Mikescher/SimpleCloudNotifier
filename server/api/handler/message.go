package handler

import (
	"blackforestbytes.com/simplecloudnotifier/common/ginresp"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	app *logic.Application
}

func (h MessageHandler) SendMessage(g *gin.Context) ginresp.HTTPResponse {
	return ginresp.NotImplemented()
}

func NewMessageHandler(app *logic.Application) MessageHandler {
	return MessageHandler{
		app: app,
	}
}
