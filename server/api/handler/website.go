package handler

import (
	"blackforestbytes.com/simplecloudnotifier/common/ginresp"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"blackforestbytes.com/simplecloudnotifier/website"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strings"
)

type WebsiteHandler struct {
	app         *logic.Application
	rexTemplate *regexp.Regexp
}

func NewWebsiteHandler(app *logic.Application) WebsiteHandler {
	return WebsiteHandler{
		app:         app,
		rexTemplate: regexp.MustCompile("{{template\\|[A-Za-z0-9_\\-.]+}}"),
	}
}

func (h WebsiteHandler) Index(g *gin.Context) ginresp.HTTPResponse {
	return h.serveAsset(g, "index.html", true)
}

func (h WebsiteHandler) APIDocs(g *gin.Context) ginresp.HTTPResponse {
	return h.serveAsset(g, "api.html", true)
}

func (h WebsiteHandler) APIDocsMore(g *gin.Context) ginresp.HTTPResponse {
	return h.serveAsset(g, "api_more.html", true)
}

func (h WebsiteHandler) MessageSent(g *gin.Context) ginresp.HTTPResponse {
	return h.serveAsset(g, "message_sent.html", true)
}

func (h WebsiteHandler) FaviconIco(g *gin.Context) ginresp.HTTPResponse {
	return h.serveAsset(g, "favicon.ico", false)
}

func (h WebsiteHandler) FaviconPNG(g *gin.Context) ginresp.HTTPResponse {
	return h.serveAsset(g, "favicon.png", false)
}

func (h WebsiteHandler) Javascript(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		Filename string `uri:"fn"`
	}

	var u uri
	if err := g.ShouldBindUri(&u); err != nil {
		return ginresp.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	return h.serveAsset(g, "js/"+u.Filename, false)
}

func (h WebsiteHandler) CSS(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		Filename string `uri:"fn"`
	}

	var u uri
	if err := g.ShouldBindUri(&u); err != nil {
		return ginresp.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	return h.serveAsset(g, "css/"+u.Filename, false)
}

func (h WebsiteHandler) serveAsset(g *gin.Context, fn string, repl bool) ginresp.HTTPResponse {
	data, err := website.Assets.ReadFile(fn)
	if err != nil {
		return ginresp.Status(http.StatusNotFound)
	}

	if repl {
		failed := false
		data = h.rexTemplate.ReplaceAllFunc(data, func(match []byte) []byte {
			prefix := len("{{template|")
			suffix := len("}}")
			fnSub := match[prefix : len(match)-suffix]
			subdata, err := website.Assets.ReadFile(string(fnSub))
			if err != nil {
				failed = true
			}
			return subdata
		})
		if failed {
			return ginresp.InternalError(errors.New("template replacement failed"))
		}
	}

	mime := "text/plain"

	lowerFN := strings.ToLower(fn)
	if strings.HasSuffix(lowerFN, ".html") || strings.HasSuffix(lowerFN, ".htm") {
		mime = "text/html"
	} else if strings.HasSuffix(lowerFN, ".css") {
		mime = "text/css"
	} else if strings.HasSuffix(lowerFN, ".js") {
		mime = "text/javascript"
	} else if strings.HasSuffix(lowerFN, ".json") {
		mime = "application/json"
	} else if strings.HasSuffix(lowerFN, ".jpeg") || strings.HasSuffix(lowerFN, ".jpg") {
		mime = "image/jpeg"
	} else if strings.HasSuffix(lowerFN, ".png") {
		mime = "image/png"
	} else if strings.HasSuffix(lowerFN, ".svg") {
		mime = "image/svg+xml"
	}

	return ginresp.Data(http.StatusOK, mime, data)
}
