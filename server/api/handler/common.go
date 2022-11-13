package handler

import (
	"blackforestbytes.com/simplecloudnotifier/common/ginresp"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"bytes"
	"github.com/gin-gonic/gin"
	sqlite3 "github.com/mattn/go-sqlite3"
	"net/http"
)

type CommonHandler struct {
	app *logic.Application
}

func NewCommonHandler(app *logic.Application) CommonHandler {
	return CommonHandler{
		app: app,
	}
}

type pingResponse struct {
	Message string           `json:"message"`
	Info    pingResponseInfo `json:"info"`
}
type pingResponseInfo struct {
	Method  string              `json:"method"`
	Request string              `json:"request"`
	Headers map[string][]string `json:"headers"`
	URI     string              `json:"uri"`
	Address string              `json:"addr"`
}

// Ping swaggerdoc
//
// @Success 200 {object} pingResponse
// @Failure 500 {object} ginresp.errBody
// @Router  /ping [get]
// @Router  /ping [post]
// @Router  /ping [put]
// @Router  /ping [delete]
// @Router  /ping [patch]
func (h CommonHandler) Ping(g *gin.Context) ginresp.HTTPResponse {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(g.Request.Body)
	resuestBody := buf.String()

	return ginresp.JSON(http.StatusOK, pingResponse{
		Message: "Pong",
		Info: pingResponseInfo{
			Method:  g.Request.Method,
			Request: resuestBody,
			Headers: g.Request.Header,
			URI:     g.Request.RequestURI,
			Address: g.Request.RemoteAddr,
		},
	})
}

// DatabaseTest swaggerdoc
//
// @Success 200 {object} handler.DatabaseTest.response
// @Failure 500 {object} ginresp.errBody
// @Router  /db-test [get]
func (h CommonHandler) DatabaseTest(g *gin.Context) ginresp.HTTPResponse {
	type response struct {
		Success          bool   `json:"success"`
		LibVersion       string `json:"libVersion"`
		LibVersionNumber int    `json:"libVersionNumber"`
		SourceID         string `json:"sourceID"`
	}

	libVersion, libVersionNumber, sourceID := sqlite3.Version()

	err := h.app.Database.Ping()
	if err != nil {
		return ginresp.InternalError(err)
	}

	return ginresp.JSON(http.StatusOK, response{
		Success:          true,
		LibVersion:       libVersion,
		LibVersionNumber: libVersionNumber,
		SourceID:         sourceID,
	})
}

// Health swaggerdoc
//
// @Success 200 {object} handler.Health.response
// @Failure 500 {object} ginresp.errBody
// @Router  /health [get]
func (h CommonHandler) Health(*gin.Context) ginresp.HTTPResponse {
	type response struct {
		Status string `json:"status"`
	}
	return ginresp.JSON(http.StatusOK, response{Status: "ok"})
}
