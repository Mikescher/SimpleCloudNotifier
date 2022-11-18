package ginresp

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"github.com/gin-gonic/gin"
	"net/http"
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
	data       any
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

type errHTTPResponse struct {
	statusCode int
	data       any
}

func (j errHTTPResponse) Write(g *gin.Context) {
	g.JSON(j.statusCode, j.data)
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
	return &errHTTPResponse{statusCode: http.StatusInternalServerError, data: apiError{Success: false, Error: int(apierr.INTERNAL_EXCEPTION), Message: e.Error()}}
}

func InternAPIError(errorid apierr.APIError, msg string, e error) HTTPResponse {
	if scn.Conf.ReturnRawErrors {
		return &errHTTPResponse{statusCode: http.StatusInternalServerError, data: apiError{Success: false, Error: int(errorid), Message: msg, RawError: e}}
	} else {
		return &errHTTPResponse{statusCode: http.StatusInternalServerError, data: apiError{Success: false, Error: int(errorid), Message: msg}}
	}
}

func SendAPIError(errorid apierr.APIError, highlight int, msg string) HTTPResponse {
	return &errHTTPResponse{statusCode: http.StatusInternalServerError, data: apiError{Success: false, Error: int(errorid), ErrorHighlight: highlight, Message: msg}}
}

func NotImplemented() HTTPResponse {
	return &errHTTPResponse{statusCode: http.StatusInternalServerError, data: apiError{Success: false, Error: -1, ErrorHighlight: 0, Message: "Not Implemented"}}
}
