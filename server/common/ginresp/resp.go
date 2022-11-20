package ginresp

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"runtime/debug"
)

type HTTPResponse interface {
	Write(g *gin.Context)
}

type jsonHTTPResponse struct {
	statusCode int
	data       any
}

func (j jsonHTTPResponse) Write(g *gin.Context) {
	g.JSON(j.statusCode, j.data)
}

type emptyHTTPResponse struct {
	statusCode int
}

func (j emptyHTTPResponse) Write(g *gin.Context) {
	g.Status(j.statusCode)
}

type textHTTPResponse struct {
	statusCode int
	data       string
}

func (j textHTTPResponse) Write(g *gin.Context) {
	g.String(j.statusCode, "%s", j.data)
}

type dataHTTPResponse struct {
	statusCode  int
	data        []byte
	contentType string
}

func (j dataHTTPResponse) Write(g *gin.Context) {
	g.Data(j.statusCode, j.contentType, j.data)
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

func InternAPIError(g *gin.Context, status int, errorid apierr.APIError, msg string, e error) HTTPResponse {
	return createApiError(g, "InternAPIError", status, errorid, 0, msg, e)
}

func SendAPIError(g *gin.Context, status int, errorid apierr.APIError, highlight int, msg string, e error) HTTPResponse {
	return createApiError(g, "SendAPIError", status, errorid, highlight, msg, e)
}

func NotImplemented(g *gin.Context) HTTPResponse {
	return createApiError(g, "NotImplemented", 500, apierr.UNDEFINED, 0, "Not Implemented", nil)
}

func createApiError(g *gin.Context, ident string, status int, errorid apierr.APIError, highlight int, msg string, e error) HTTPResponse {
	reqUri := ""
	if g != nil && g.Request != nil {
		reqUri = g.Request.Method + " :: " + g.Request.RequestURI
	}

	log.Error().
		Int("errorid", int(errorid)).
		Int("highlight", highlight).
		Str("uri", reqUri).
		AnErr("err", e).
		Stack().
		Msg(fmt.Sprintf("[%s] %s", ident, msg))

	if scn.Conf.ReturnRawErrors {
		return &jsonHTTPResponse{
			statusCode: status,
			data: apiError{
				Success:        false,
				Error:          int(errorid),
				ErrorHighlight: highlight,
				Message:        msg,
				RawError:       fmt.Sprintf("%+v", e),
				Trace:          string(debug.Stack()),
			},
		}
	} else {
		return &jsonHTTPResponse{
			statusCode: status,
			data: apiError{
				Success:        false,
				Error:          int(errorid),
				ErrorHighlight: highlight,
				Message:        msg,
			},
		}
	}
}

func CompatAPIError(errid int, msg string) HTTPResponse {
	return &jsonHTTPResponse{statusCode: 200, data: compatAPIError{Success: false, ErrorID: errid, Message: msg}}
}
