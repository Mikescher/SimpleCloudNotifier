package handler

import (
	"blackforestbytes.com/simplecloudnotifier/common/ginresp"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"blackforestbytes.com/simplecloudnotifier/website"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"regexp"
	"strings"
)

type WebsiteHandler struct {
	app         *logic.Application
	rexTemplate *regexp.Regexp
	rexConfig   *regexp.Regexp
}

func NewWebsiteHandler(app *logic.Application) WebsiteHandler {
	return WebsiteHandler{
		app:         app,
		rexTemplate: regexp.MustCompile("{{template\\|[A-Za-z0-9_\\-\\[\\].]+}}"),
		rexConfig:   regexp.MustCompile("{{config\\|[A-Za-z0-9_\\-.]+}}"),
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
			fnSub := string(match[prefix : len(match)-suffix])

			fnSub = strings.ReplaceAll(fnSub, "[theme]", h.getTheme(g))

			subdata, err := website.Assets.ReadFile(fnSub)
			if err != nil {
				log.Error().Str("templ", string(match)).Str("fnSub", fnSub).Str("source", fn).Msg("Failed to replace template")
				failed = true
			}
			return subdata
		})
		if failed {
			return ginresp.InternalError(errors.New("template replacement failed"))
		}

		data = h.rexConfig.ReplaceAllFunc(data, func(match []byte) []byte {
			prefix := len("{{config|")
			suffix := len("}}")
			cfgKey := match[prefix : len(match)-suffix]

			cval, ok := h.getReplConfig(string(cfgKey))
			if !ok {
				log.Error().Str("templ", string(match)).Str("source", fn).Msg("Failed to replace config")
				failed = true
			}
			return []byte(cval)
		})
		if failed {
			return ginresp.InternalError(errors.New("config replacement failed"))
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

func (h WebsiteHandler) getReplConfig(key string) (string, bool) {
	key = strings.TrimSpace(strings.ToLower(key))

	if key == "baseurl" {
		return h.app.Config.BaseURL, true
	}
	if key == "ip" {
		return h.app.Config.ServerIP, true
	}
	if key == "port" {
		return h.app.Config.ServerPort, true
	}
	if key == "namespace" {
		return h.app.Config.Namespace, true
	}

	return "", false

}

func (h WebsiteHandler) getTheme(g *gin.Context) string {
	if c, err := g.Cookie("theme"); err != nil {
		return "light"
	} else if c == "light" {
		return "light"
	} else if c == "dark" {
		return "dark"
	} else {
		return "light"
	}
}
