package handler

import (
	"blackforestbytes.com/simplecloudnotifier/api/ginresp"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"blackforestbytes.com/simplecloudnotifier/models"
	"blackforestbytes.com/simplecloudnotifier/website"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
	"gogs.mikescher.com/BlackForestBytes/goext/rext"
	"net/http"
	"regexp"
	"strings"
)

type WebsiteHandler struct {
	app         *logic.Application
	rexTemplate rext.Regex
	rexConfig   rext.Regex
}

func NewWebsiteHandler(app *logic.Application) WebsiteHandler {
	return WebsiteHandler{
		app:         app,
		rexTemplate: rext.W(regexp.MustCompile("{{template\\|[A-Za-z0-9_\\-\\[\\].]+}}")),
		rexConfig:   rext.W(regexp.MustCompile("{{config\\|[A-Za-z0-9_\\-.]+}}")),
	}
}

func (h WebsiteHandler) Index(pctx ginext.PreContext) ginext.HTTPResponse {
	ctx, g, errResp := pctx.Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	return h.app.DoRequest(ctx, g, models.TLockNone, func(ctx *logic.AppContext, finishSuccess func(r ginext.HTTPResponse) ginext.HTTPResponse) ginext.HTTPResponse {
		return h.serveAsset(g, "index.html", true)
	})
}

func (h WebsiteHandler) APIDocs(pctx ginext.PreContext) ginext.HTTPResponse {
	ctx, g, errResp := pctx.Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	return h.app.DoRequest(ctx, g, models.TLockNone, func(ctx *logic.AppContext, finishSuccess func(r ginext.HTTPResponse) ginext.HTTPResponse) ginext.HTTPResponse {
		return h.serveAsset(g, "api.html", true)
	})
}

func (h WebsiteHandler) APIDocsMore(pctx ginext.PreContext) ginext.HTTPResponse {
	ctx, g, errResp := pctx.Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	return h.app.DoRequest(ctx, g, models.TLockNone, func(ctx *logic.AppContext, finishSuccess func(r ginext.HTTPResponse) ginext.HTTPResponse) ginext.HTTPResponse {
		return h.serveAsset(g, "api_more.html", true)
	})
}

func (h WebsiteHandler) MessageSent(pctx ginext.PreContext) ginext.HTTPResponse {
	ctx, g, errResp := pctx.Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	return h.app.DoRequest(ctx, g, models.TLockNone, func(ctx *logic.AppContext, finishSuccess func(r ginext.HTTPResponse) ginext.HTTPResponse) ginext.HTTPResponse {
		return h.serveAsset(g, "message_sent.html", true)
	})
}

func (h WebsiteHandler) FaviconIco(pctx ginext.PreContext) ginext.HTTPResponse {
	ctx, g, errResp := pctx.Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	return h.app.DoRequest(ctx, g, models.TLockNone, func(ctx *logic.AppContext, finishSuccess func(r ginext.HTTPResponse) ginext.HTTPResponse) ginext.HTTPResponse {
		return h.serveAsset(g, "favicon.ico", false)
	})
}

func (h WebsiteHandler) FaviconPNG(pctx ginext.PreContext) ginext.HTTPResponse {
	ctx, g, errResp := pctx.Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	return h.app.DoRequest(ctx, g, models.TLockNone, func(ctx *logic.AppContext, finishSuccess func(r ginext.HTTPResponse) ginext.HTTPResponse) ginext.HTTPResponse {
		return h.serveAsset(g, "favicon.png", false)
	})
}

func (h WebsiteHandler) Javascript(pctx ginext.PreContext) ginext.HTTPResponse {
	ctx, g, errResp := pctx.Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	return h.app.DoRequest(ctx, g, models.TLockNone, func(ctx *logic.AppContext, finishSuccess func(r ginext.HTTPResponse) ginext.HTTPResponse) ginext.HTTPResponse {

		type uri struct {
			Filename string `uri:"fn"`
		}

		var u uri
		if err := g.ShouldBindUri(&u); err != nil {
			return ginext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		return h.serveAsset(g, "js/"+u.Filename, false)
	})
}

func (h WebsiteHandler) CSS(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
		Filename string `uri:"fn"`
	}

	var u uri
	ctx, g, errResp := pctx.URI(&u).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	return h.app.DoRequest(ctx, g, models.TLockNone, func(ctx *logic.AppContext, finishSuccess func(r ginext.HTTPResponse) ginext.HTTPResponse) ginext.HTTPResponse {
		return h.serveAsset(g, "css/"+u.Filename, false)
	})
}

func (h WebsiteHandler) serveAsset(g *gin.Context, fn string, repl bool) ginext.HTTPResponse {
	_data, err := website.Assets.ReadFile(fn)
	if err != nil {
		return ginext.Status(http.StatusNotFound)
	}

	data := string(_data)

	if repl {
		failed := false
		data = h.rexTemplate.ReplaceAllFunc(data, func(match string) string {
			prefix := len("{{template|")
			suffix := len("}}")
			fnSub := match[prefix : len(match)-suffix]

			fnSub = strings.ReplaceAll(fnSub, "[theme]", h.getTheme(g))

			subdata, err := website.Assets.ReadFile(fnSub)
			if err != nil {
				log.Error().Str("templ", string(match)).Str("fnSub", fnSub).Str("source", fn).Msg("Failed to replace template")
				failed = true
			}
			return string(subdata)
		})
		if failed {
			return ginresp.InternalError(errors.New("template replacement failed"))
		}

		data = h.rexConfig.ReplaceAllFunc(data, func(match string) string {
			prefix := len("{{config|")
			suffix := len("}}")
			cfgKey := match[prefix : len(match)-suffix]

			cval, ok := h.getReplConfig(cfgKey)
			if !ok {
				log.Error().Str("templ", match).Str("source", fn).Msg("Failed to replace config")
				failed = true
			}
			return cval
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

	return ginext.Data(http.StatusOK, mime, []byte(data))
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
