package ginresp

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HTTPResponse interface {
	Write(context *gin.Context)
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
	return &errHTTPResponse{statusCode: http.StatusInternalServerError, data: errBody{Success: false, Message: e.Error()}}
}

func NotImplemented(errid int) HTTPResponse {
	return &errHTTPResponse{statusCode: http.StatusInternalServerError, data: internAPIError{Success: false, ErrorID: errid, Message: "NotImplemented"}}
}

func SendAPIError(errorid apierr.APIError, highlight int, msg string) HTTPResponse {
	return &errHTTPResponse{statusCode: http.StatusInternalServerError, data: sendAPIError{Success: false, Error: int(errorid), ErrorHighlight: highlight, Message: msg}}
}
