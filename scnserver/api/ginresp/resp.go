package ginresp

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/api/apihighlight"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	json "gogs.mikescher.com/BlackForestBytes/goext/gojson"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"runtime/debug"
	"strings"
)

type HTTPResponse interface {
	Write(g *gin.Context)
	Statuscode() int
	BodyString() *string
	ContentType() string
}

type jsonHTTPResponse struct {
	statusCode int
	data       any
}

func (j jsonHTTPResponse) Write(g *gin.Context) {
	g.Render(j.statusCode, json.GoJsonRender{Data: j.data, NilSafeSlices: true, NilSafeMaps: true})
}

func (j jsonHTTPResponse) Statuscode() int {
	return j.statusCode
}

func (j jsonHTTPResponse) BodyString() *string {
	v, err := json.Marshal(j.data)
	if err != nil {
		return nil
	}
	return langext.Ptr(string(v))
}

func (j jsonHTTPResponse) ContentType() string {
	return "application/json"
}

type emptyHTTPResponse struct {
	statusCode int
}

func (j emptyHTTPResponse) Write(g *gin.Context) {
	g.Status(j.statusCode)
}

func (j emptyHTTPResponse) Statuscode() int {
	return j.statusCode
}

func (j emptyHTTPResponse) BodyString() *string {
	return nil
}

func (j emptyHTTPResponse) ContentType() string {
	return ""
}

type textHTTPResponse struct {
	statusCode int
	data       string
}

func (j textHTTPResponse) Write(g *gin.Context) {
	g.String(j.statusCode, "%s", j.data)
}

func (j textHTTPResponse) Statuscode() int {
	return j.statusCode
}

func (j textHTTPResponse) BodyString() *string {
	return langext.Ptr(j.data)
}

func (j textHTTPResponse) ContentType() string {
	return "text/plain"
}

type dataHTTPResponse struct {
	statusCode  int
	data        []byte
	contentType string
}

func (j dataHTTPResponse) Write(g *gin.Context) {
	g.Data(j.statusCode, j.contentType, j.data)
}

func (j dataHTTPResponse) Statuscode() int {
	return j.statusCode
}

func (j dataHTTPResponse) BodyString() *string {
	return langext.Ptr(string(j.data))
}

func (j dataHTTPResponse) ContentType() string {
	return j.contentType
}

type errorHTTPResponse struct {
	statusCode int
	data       any
	error      error
}

func (j errorHTTPResponse) Write(g *gin.Context) {
	g.JSON(j.statusCode, j.data)
}

func (j errorHTTPResponse) Statuscode() int {
	return j.statusCode
}

func (j errorHTTPResponse) BodyString() *string {
	v, err := json.Marshal(j.data)
	if err != nil {
		return nil
	}
	return langext.Ptr(string(v))
}

func (j errorHTTPResponse) ContentType() string {
	return "application/json"
}

func Status(sc int) HTTPResponse {
	return &emptyHTTPResponse{statusCode: sc}
}

func JSON(sc int, data any) HTTPResponse {
	return &jsonHTTPResponse{statusCode: sc, data: data}
}

func Data(sc int, contentType string, data []byte) HTTPResponse {
	return &dataHTTPResponse{statusCode: sc, contentType: contentType, data: data}
}

func Text(sc int, data string) HTTPResponse {
	return &textHTTPResponse{statusCode: sc, data: data}
}

func InternalError(e error) HTTPResponse {
	return createApiError(nil, "InternalError", 500, apierr.INTERNAL_EXCEPTION, 0, e.Error(), e)
}

func APIError(g *gin.Context, status int, errorid apierr.APIError, msg string, e error) HTTPResponse {
	return createApiError(g, "APIError", status, errorid, 0, msg, e)
}

func SendAPIError(g *gin.Context, status int, errorid apierr.APIError, highlight apihighlight.ErrHighlight, msg string, e error) HTTPResponse {
	return createApiError(g, "SendAPIError", status, errorid, highlight, msg, e)
}

func NotImplemented(g *gin.Context) HTTPResponse {
	return createApiError(g, "NotImplemented", 500, apierr.NOT_IMPLEMENTED, 0, "Not Implemented", nil)
}

func createApiError(g *gin.Context, ident string, status int, errorid apierr.APIError, highlight apihighlight.ErrHighlight, msg string, e error) HTTPResponse {
	reqUri := ""
	if g != nil && g.Request != nil {
		reqUri = g.Request.Method + " :: " + g.Request.RequestURI
	}

	log.Error().
		Int("errorid", int(errorid)).
		Int("highlight", int(highlight)).
		Str("uri", reqUri).
		AnErr("err", e).
		Stack().
		Msg(fmt.Sprintf("[%s] %s", ident, msg))

	if scn.Conf.ReturnRawErrors {
		return &errorHTTPResponse{
			statusCode: status,
			data: extendedAPIError{
				Success:        false,
				Error:          int(errorid),
				ErrorHighlight: int(highlight),
				Message:        msg,
				RawError:       langext.Ptr(langext.Conditional(e == nil, "", fmt.Sprintf("%+v", e))),
				Trace:          strings.Split(string(debug.Stack()), "\n"),
			},
			error: e,
		}
	} else {
		return &errorHTTPResponse{
			statusCode: status,
			data: apiError{
				Success:        false,
				Error:          int(errorid),
				ErrorHighlight: int(highlight),
				Message:        msg,
			},
			error: e,
		}
	}
}

func CompatAPIError(errid int, msg string) HTTPResponse {
	return &jsonHTTPResponse{statusCode: 200, data: compatAPIError{Success: false, ErrorID: errid, Message: msg}}
}
