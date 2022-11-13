package ginext

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"github.com/gin-gonic/gin"
)

func NewEngine(cfg scn.Config) *gin.Engine {
	engine := gin.New()

	engine.RedirectFixedPath = false
	engine.RedirectTrailingSlash = false

	engine.Use(CorsMiddleware())

	if cfg.GinDebug {
		engine.Use(gin.Logger())
	}

	return engine
}
