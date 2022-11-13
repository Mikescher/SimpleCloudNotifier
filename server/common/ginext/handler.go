package ginext

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RedirectFound(newuri string) gin.HandlerFunc {
	return func(g *gin.Context) {
		g.Redirect(http.StatusFound, newuri)
	}
}

func RedirectTemporary(newuri string) gin.HandlerFunc {
	return func(g *gin.Context) {
		g.Redirect(http.StatusTemporaryRedirect, newuri)
	}
}

func RedirectPermanent(newuri string) gin.HandlerFunc {
	return func(g *gin.Context) {
		g.Redirect(http.StatusPermanentRedirect, newuri)
	}
}
