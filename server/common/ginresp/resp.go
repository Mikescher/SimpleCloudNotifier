package ginresp

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
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
	log.Error().Err(e).Msg("[InternalError] " + e.Error())

	return &jsonHTTPResponse{statusCode: http.StatusInternalServerError, data: apiError{Success: false, Error: int(apierr.INTERNAL_EXCEPTION), Message: e.Error()}}
}

func InternAPIError(status int, errorid apierr.APIError, msg string, e error) HTTPResponse {
	log.Error().Int("errorid", int(errorid)).Err(e).Msg("[InternAPIError] " + msg)

	if scn.Conf.ReturnRawErrors {
		return &jsonHTTPResponse{statusCode: status, data: apiError{Success: false, Error: int(errorid), Message: msg, RawError: fmt.Sprintf("%+v", e), Trace: string(debug.Stack())}}
	} else {
		return &jsonHTTPResponse{statusCode: status, data: apiError{Success: false, Error: int(errorid), Message: msg}}
	}
}

func CompatAPIError(errid int, msg string) HTTPResponse {
	log.Error().Int("errid", errid).Msg("[CompatAPIError] " + msg)

	return &jsonHTTPResponse{statusCode: 200, data: compatAPIError{Success: false, ErrorID: errid, Message: msg}}
}

func SendAPIError(status int, errorid apierr.APIError, highlight int, msg string, e error) HTTPResponse {
	log.Error().Int("errorid", int(errorid)).Int("highlight", highlight).Err(e).Msg("[SendAPIError] " + msg)

	if scn.Conf.ReturnRawErrors {
		return &jsonHTTPResponse{statusCode: status, data: apiError{Success: false, Error: int(errorid), ErrorHighlight: highlight, Message: msg, RawError: fmt.Sprintf("%+v", e), Trace: string(debug.Stack())}}
	} else {
		return &jsonHTTPResponse{statusCode: status, data: apiError{Success: false, Error: int(errorid), ErrorHighlight: highlight, Message: msg}}
	}
}

func NotImplemented() HTTPResponse {
	log.Error().Msg("[NotImplemented]")

	return &jsonHTTPResponse{statusCode: http.StatusInternalServerError, data: apiError{Success: false, Error: -1, ErrorHighlight: 0, Message: "Not Implemented"}}
}
