package ginresp

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/api/apihighlight"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
	json "gogs.mikescher.com/BlackForestBytes/goext/gojson"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"runtime/debug"
	"strings"
)

type cookieval struct {
	name     string
	value    string
	maxAge   int
	path     string
	domain   string
	secure   bool
	httpOnly bool
}

type headerval struct {
	Key string
	Val string
}

type errorHTTPResponse struct {
	statusCode int
	data       any
	error      error
	headers    []headerval
	cookies    []cookieval
}

func (j errorHTTPResponse) Write(g *gin.Context) {
	for _, v := range j.headers {
		g.Header(v.Key, v.Val)
	}
	for _, v := range j.cookies {
		g.SetCookie(v.name, v.value, v.maxAge, v.path, v.domain, v.secure, v.httpOnly)
	}
	g.JSON(j.statusCode, j.data)
}

func (j errorHTTPResponse) Statuscode() int {
	return j.statusCode
}

func (j errorHTTPResponse) BodyString(g *gin.Context) *string {
	v, err := json.Marshal(j.data)
	if err != nil {
		return nil
	}
	return langext.Ptr(string(v))
}

func (j errorHTTPResponse) ContentType() string {
	return "application/json"
}

func (j errorHTTPResponse) WithHeader(k string, v string) ginext.HTTPResponse {
	j.headers = append(j.headers, headerval{k, v})
	return j
}

func (j errorHTTPResponse) WithCookie(name string, value string, maxAge int, path string, domain string, secure bool, httpOnly bool) ginext.HTTPResponse {
	j.cookies = append(j.cookies, cookieval{name, value, maxAge, path, domain, secure, httpOnly})
	return j
}

func (j errorHTTPResponse) IsSuccess() bool {
	return false
}

func (j errorHTTPResponse) Headers() []string {
	return langext.ArrMap(j.headers, func(v headerval) string { return v.Key + "=" + v.Val })
}

func InternalError(e error) ginext.HTTPResponse {
	return createApiError(nil, "InternalError", 500, apierr.INTERNAL_EXCEPTION, 0, e.Error(), e)
}

func APIError(g *gin.Context, status int, errorid apierr.APIError, msg string, e error) ginext.HTTPResponse {
	return createApiError(g, "APIError", status, errorid, 0, msg, e)
}

func SendAPIError(g *gin.Context, status int, errorid apierr.APIError, highlight apihighlight.ErrHighlight, msg string, e error) ginext.HTTPResponse {
	return createApiError(g, "SendAPIError", status, errorid, highlight, msg, e)
}

func createApiError(g *gin.Context, ident string, status int, errorid apierr.APIError, highlight apihighlight.ErrHighlight, msg string, e error) ginext.HTTPResponse {
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

func CompatAPIError(errid int, msg string) ginext.HTTPResponse {
	return ginext.JSON(200, compatAPIError{Success: false, ErrorID: errid, Message: msg})
}
