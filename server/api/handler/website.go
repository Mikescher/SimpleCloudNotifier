package handler

import (
	"blackforestbytes.com/simplecloudnotifier/common/ginresp"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"blackforestbytes.com/simplecloudnotifier/website"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type WebsiteHandler struct {
	app *logic.Application
}

func NewWebsiteHandler(app *logic.Application) WebsiteHandler {
	return WebsiteHandler{
		app: app,
	}
}

func (h WebsiteHandler) Index(g *gin.Context) ginresp.HTTPResponse {
	return h.serveAsset(g, "index.html")
}

func (h WebsiteHandler) APIDocs(g *gin.Context) ginresp.HTTPResponse {
	return h.serveAsset(g, "api.html")
}

func (h WebsiteHandler) APIDocsMore(g *gin.Context) ginresp.HTTPResponse {
	return h.serveAsset(g, "api_more.html")
}

func (h WebsiteHandler) MessageSent(g *gin.Context) ginresp.HTTPResponse {
	return h.serveAsset(g, "message_sent.html")
}

func (h WebsiteHandler) FaviconIco(g *gin.Context) ginresp.HTTPResponse {
	return h.serveAsset(g, "favicon.ico")
}

func (h WebsiteHandler) FaviconPNG(g *gin.Context) ginresp.HTTPResponse {
	return h.serveAsset(g, "favicon.png")
}

func (h WebsiteHandler) Javascript(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		Filename string `uri:"fn"`
	}

	var u uri
	if err := g.ShouldBindUri(&u); err != nil {
		return ginresp.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	return h.serveAsset(g, "js/"+u.Filename)
}

func (h WebsiteHandler) CSS(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		Filename string `uri:"fn"`
	}

	var u uri
	if err := g.ShouldBindUri(&u); err != nil {
		return ginresp.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	return h.serveAsset(g, "css/"+u.Filename)
}

func (h WebsiteHandler) serveAsset(g *gin.Context, fn string) ginresp.HTTPResponse {
	data, err := website.Assets.ReadFile(fn)
	if err != nil {
		return ginresp.Status(http.StatusNotFound)
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
