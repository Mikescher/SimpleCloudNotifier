package api

import (
	"blackforestbytes.com/simplecloudnotifier/api/handler"
	"blackforestbytes.com/simplecloudnotifier/common/ginext"
	"blackforestbytes.com/simplecloudnotifier/common/ginresp"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"blackforestbytes.com/simplecloudnotifier/swagger"
	"github.com/gin-gonic/gin"
)

type Router struct {
	app *logic.Application

	commonHandler handler.CommonHandler
	compatHandler handler.CompatHandler
}

func NewRouter(app *logic.Application) *Router {
	return &Router{
		app: app,

		commonHandler: handler.NewCommonHandler(app),
		compatHandler: handler.NewCompatHandler(app),
	}
}

// Init swaggerdocs
// @title       SimpleCloudNotifier API
// @version     2.0
// @description API for SCN
// @host        scn.blackforestbytes.com
// @BasePath    /api/
func (r *Router) Init(e *gin.Engine) {

	e.Any("/ping", ginresp.Wrap(r.commonHandler.Ping))
	e.POST("/db-test", ginresp.Wrap(r.commonHandler.DatabaseTest))
	e.GET("/health", ginresp.Wrap(r.commonHandler.Health))

	e.GET("documentation/swagger", ginext.RedirectTemporary("/documentation/swagger/"))
	e.GET("documentation/swagger/", ginresp.Wrap(swagger.Handle))
	e.GET("documentation/swagger/:fn", ginresp.Wrap(swagger.Handle))

	e.POST("/send.php", ginresp.Wrap(r.compatHandler.Send))
	e.GET("/register.php", ginresp.Wrap(r.compatHandler.Register))
	e.GET("/info.php", ginresp.Wrap(r.compatHandler.Info))
	e.GET("/ack.php", ginresp.Wrap(r.compatHandler.Ack))
	e.GET("/requery.php", ginresp.Wrap(r.compatHandler.Requery))
	e.GET("/update.php", ginresp.Wrap(r.compatHandler.Update))
	e.GET("/expand.php", ginresp.Wrap(r.compatHandler.Expand))
	e.GET("/upgrade.php", ginresp.Wrap(r.compatHandler.Upgrade))
}
