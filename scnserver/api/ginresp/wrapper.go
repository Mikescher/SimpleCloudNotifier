package ginresp

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/dataext"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"math/rand"
	"runtime/debug"
	"time"
)

type WHandlerFunc func(*gin.Context) HTTPResponse

type RequestLogAcceptor interface {
	InsertRequestLog(data models.RequestLog)
}

func Wrap(rlacc RequestLogAcceptor, fn WHandlerFunc) gin.HandlerFunc {

	maxRetry := scn.Conf.RequestMaxRetry
	retrySleep := scn.Conf.RequestRetrySleep

	return func(g *gin.Context) {

		reqctx := g.Request.Context()

		if g.Request.Body != nil {
			g.Request.Body = dataext.NewBufferedReadCloser(g.Request.Body)
		}

		t0 := time.Now()

		for ctr := 1; ; ctr++ {

			wrap, stackTrace, panicObj := callPanicSafe(fn, g)
			if panicObj != nil {
				log.Error().Interface("panicObj", panicObj).Msg("Panic occured (in gin handler)")
				log.Error().Msg(stackTrace)
				wrap = APIError(g, 500, apierr.PANIC, "A panic occured in the HTTP handler", errors.New(fmt.Sprintf("%+v\n\n@:\n%s", panicObj, stackTrace)))
			}

			if g.Writer.Written() {
				if scn.Conf.ReqLogEnabled {
					rlacc.InsertRequestLog(createRequestLog(g, t0, ctr, nil, langext.Ptr("Writing in WrapperFunc is not supported")))
				}
				panic("Writing in WrapperFunc is not supported")
			}

			if ctr < maxRetry && isSqlite3Busy(wrap) {
				log.Warn().Int("counter", ctr).Str("url", g.Request.URL.String()).Msg("Retry request (ErrBusy)")

				err := resetBody(g)
				if err != nil {
					panic(err)
				}

				time.Sleep(time.Duration(int64(float64(retrySleep) * (0.5 + rand.Float64()))))
				continue
			}

			if reqctx.Err() == nil {
				if scn.Conf.ReqLogEnabled {
					rlacc.InsertRequestLog(createRequestLog(g, t0, ctr, wrap, nil))
				}

				statuscode := wrap.Statuscode()
				if statuscode/100 != 2 {
					log.Warn().Str("url", g.Request.Method+"::"+g.Request.URL.String()).Msg(fmt.Sprintf("Request failed with statuscode %d", statuscode))
				}

				wrap.Write(g)
			}

			return
		}

	}

}

func createRequestLog(g *gin.Context, t0 time.Time, ctr int, resp HTTPResponse, panicstr *string) models.RequestLog {

	t1 := time.Now()

	ua := g.Request.UserAgent()
	auth := g.Request.Header.Get("Authorization")
	ct := g.Request.Header.Get("Content-Type")

	var reqbody []byte = nil
	if g.Request.Body != nil {
		brcbody, err := g.Request.Body.(dataext.BufferedReadCloser).BufferedAll()
		if err == nil {
			reqbody = brcbody
		}
	}
	var strreqbody *string = nil
	if len(reqbody) < scn.Conf.ReqLogMaxBodySize {
		strreqbody = langext.Ptr(string(reqbody))
	}

	var respbody *string = nil

	var strrespbody *string = nil
	if resp != nil {
		respbody = resp.BodyString()
		if respbody != nil && len(*respbody) < scn.Conf.ReqLogMaxBodySize {
			strrespbody = respbody
		}
	}

	permObj, hasPerm := g.Get("perm")

	hasTok := false
	if hasPerm {
		hasTok = permObj.(models.PermissionSet).Token != nil
	}

	return models.RequestLog{
		Method:              g.Request.Method,
		URI:                 g.Request.URL.String(),
		UserAgent:           langext.Conditional(ua == "", nil, &ua),
		Authentication:      langext.Conditional(auth == "", nil, &auth),
		RequestBody:         strreqbody,
		RequestBodySize:     int64(len(reqbody)),
		RequestContentType:  ct,
		RemoteIP:            g.RemoteIP(),
		KeyID:               langext.ConditionalFn10(hasTok, func() *models.KeyTokenID { return langext.Ptr(permObj.(models.PermissionSet).Token.KeyTokenID) }, nil),
		UserID:              langext.ConditionalFn10(hasTok, func() *models.UserID { return langext.Ptr(permObj.(models.PermissionSet).Token.OwnerUserID) }, nil),
		Permissions:         langext.ConditionalFn10(hasTok, func() *string { return langext.Ptr(permObj.(models.PermissionSet).Token.Permissions.String()) }, nil),
		ResponseStatuscode:  langext.ConditionalFn10(resp != nil, func() *int64 { return langext.Ptr(int64(resp.Statuscode())) }, nil),
		ResponseBodySize:    langext.ConditionalFn10(strrespbody != nil, func() *int64 { return langext.Ptr(int64(len(*respbody))) }, nil),
		ResponseBody:        strrespbody,
		ResponseContentType: langext.ConditionalFn10(resp != nil, func() string { return resp.ContentType() }, ""),
		RetryCount:          int64(ctr),
		Panicked:            panicstr != nil,
		PanicStr:            panicstr,
		ProcessingTime:      t1.Sub(t0),
		TimestampStart:      t0,
		TimestampFinish:     t1,
	}
}

func callPanicSafe(fn WHandlerFunc, g *gin.Context) (res HTTPResponse, stackTrace string, panicObj any) {
	defer func() {
		if rec := recover(); rec != nil {
			res = nil
			stackTrace = string(debug.Stack())
			panicObj = rec
		}
	}()

	res = fn(g)
	return res, "", nil
}

func resetBody(g *gin.Context) error {
	if g.Request.Body == nil {
		return nil
	}

	err := g.Request.Body.(dataext.BufferedReadCloser).Reset()
	if err != nil {
		return err
	}

	return nil
}

func isSqlite3Busy(r HTTPResponse) bool {
	if errwrap, ok := r.(*errorHTTPResponse); ok && errwrap != nil {

		if errors.Is(errwrap.error, sqlite3.ErrBusy) {
			return true
		}

		var s3err sqlite3.Error
		if errors.As(errwrap.error, &s3err) {
			if errors.Is(s3err.Code, sqlite3.ErrBusy) {
				return true
			}
		}
	}
	return false
}
