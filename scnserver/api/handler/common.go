package handler

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/api/ginresp"
	"blackforestbytes.com/simplecloudnotifier/db/simplectx"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-sqlite3"
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/timeext"
	"net/http"
	"time"
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
//	@Summary	Simple endpoint to test connection (any http method)
//	@Tags		Common
//
//	@Success	200	{object}	pingResponse
//	@Failure	500	{object}	ginresp.apiError
//
//	@Router		/api/ping [get]
//	@Router		/api/ping [post]
//	@Router		/api/ping [put]
//	@Router		/api/ping [delete]
//	@Router		/api/ping [patch]
func (h CommonHandler) Ping(pctx ginext.PreContext) ginext.HTTPResponse {
	ctx, g, errResp := pctx.Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(g.Request.Body)
	resuestBody := buf.String()

	return ginext.JSON(http.StatusOK, pingResponse{
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
//	@Summary	Check for a working database connection
//	@ID			api-common-dbtest
//	@Tags		Common
//
//	@Success	200	{object}	handler.DatabaseTest.response
//	@Failure	500	{object}	ginresp.apiError
//
//	@Router		/api/db-test [post]
func (h CommonHandler) DatabaseTest(pctx ginext.PreContext) ginext.HTTPResponse {
	type response struct {
		Success          bool   `json:"success"`
		LibVersion       string `json:"libVersion"`
		LibVersionNumber int    `json:"libVersionNumber"`
		SourceID         string `json:"sourceID"`
	}

	ctx, _, errResp := pctx.Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	libVersion, libVersionNumber, sourceID := sqlite3.Version()

	err := h.app.Database.Ping(ctx)
	if err != nil {
		return ginresp.InternalError(err)
	}

	return ginext.JSON(http.StatusOK, response{
		Success:          true,
		LibVersion:       libVersion,
		LibVersionNumber: libVersionNumber,
		SourceID:         sourceID,
	})
}

// Health swaggerdoc
//
//	@Summary	Server Health-checks
//	@ID			api-common-health
//	@Tags		Common
//
//	@Success	200	{object}	handler.Health.response
//	@Failure	500	{object}	ginresp.apiError
//
//	@Router		/api/health [get]
func (h CommonHandler) Health(pctx ginext.PreContext) ginext.HTTPResponse {
	type response struct {
		Status string `json:"status"`
	}

	ctx, _, errResp := pctx.Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	_, libVersionNumber, _ := sqlite3.Version()

	if libVersionNumber < 3039000 {
		return ginresp.InternalError(errors.New("sqlite version too low"))
	}

	tctx := simplectx.CreateSimpleContext(ctx, nil)

	err := h.app.Database.Ping(tctx)
	if err != nil {
		return ginresp.InternalError(err)
	}

	for _, subdb := range h.app.Database.List() {

		uuidKey, _ := langext.NewHexUUID()
		uuidWrite, _ := langext.NewHexUUID()

		err = subdb.WriteMetaString(tctx, uuidKey, uuidWrite)
		if err != nil {
			return ginresp.InternalError(err)
		}

		uuidRead, err := subdb.ReadMetaString(tctx, uuidKey)
		if err != nil {
			return ginresp.InternalError(err)
		}

		if uuidRead == nil || uuidWrite != *uuidRead {
			return ginresp.InternalError(errors.New("writing into DB was not consistent"))
		}

		err = subdb.DeleteMeta(tctx, uuidKey)
		if err != nil {
			return ginresp.InternalError(err)
		}

	}

	return ginext.JSON(http.StatusOK, response{Status: "ok"})
}

// Sleep swaggerdoc
//
//	@Summary	Return 200 after x seconds
//	@ID			api-common-sleep
//	@Tags		Common
//
//	@Param		secs	path		number	true	"sleep delay (in seconds)"
//
//	@Success	200		{object}	handler.Sleep.response
//	@Failure	400		{object}	ginresp.apiError
//	@Failure	500		{object}	ginresp.apiError
//
//	@Router		/api/sleep/{secs} [post]
func (h CommonHandler) Sleep(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
		Seconds float64 `uri:"secs"`
	}
	type response struct {
		Start    string  `json:"start"`
		End      string  `json:"end"`
		Duration float64 `json:"duration"`
	}

	ctx, g, errResp := pctx.Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	t0 := time.Now().Format(time.RFC3339Nano)

	var u uri
	if err := g.ShouldBindUri(&u); err != nil {
		return ginresp.APIError(g, 400, apierr.BINDFAIL_URI_PARAM, "Failed to read uri", err)
	}

	time.Sleep(timeext.FromSeconds(u.Seconds))

	t1 := time.Now().Format(time.RFC3339Nano)

	return ginext.JSON(http.StatusOK, response{
		Start:    t0,
		End:      t1,
		Duration: u.Seconds,
	})
}

func (h CommonHandler) NoRoute(pctx ginext.PreContext) ginext.HTTPResponse {
	ctx, g, errResp := pctx.Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	return ginext.JSON(http.StatusNotFound, gin.H{
		"":           "================ ROUTE NOT FOUND ================",
		"FullPath":   g.FullPath(),
		"Method":     g.Request.Method,
		"URL":        g.Request.URL.String(),
		"RequestURI": g.Request.RequestURI,
		"Proto":      g.Request.Proto,
		"Header":     g.Request.Header,
		"~":          "================ ROUTE NOT FOUND ================",
	})
}
