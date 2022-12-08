package ginext

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"github.com/gin-gonic/gin"
)

var SuppressGinLogs = false

func NewEngine(cfg scn.Config) *gin.Engine {
	engine := gin.New()

	engine.RedirectFixedPath = false
	engine.RedirectTrailingSlash = false

	engine.Use(CorsMiddleware())

	if cfg.GinDebug {
		ginlogger := gin.Logger()
		engine.Use(func(context *gin.Context) {
			if SuppressGinLogs {
				return
			}
			ginlogger(context)
		})
	}

	return engine
}
