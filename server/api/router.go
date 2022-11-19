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

	commonHandler  handler.CommonHandler
	compatHandler  handler.CompatHandler
	websiteHandler handler.WebsiteHandler
	apiHandler     handler.APIHandler
	messageHandler handler.MessageHandler
}

func NewRouter(app *logic.Application) *Router {
	return &Router{
		app: app,

		commonHandler:  handler.NewCommonHandler(app),
		compatHandler:  handler.NewCompatHandler(app),
		websiteHandler: handler.NewWebsiteHandler(app),
		apiHandler:     handler.NewAPIHandler(app),
		messageHandler: handler.NewMessageHandler(app),
	}
}

// Init swaggerdocs
// @title       SimpleCloudNotifier API
// @version     2.0
// @description API for SCN
// @host        scn.blackforestbytes.com
// @BasePath    /
func (r *Router) Init(e *gin.Engine) {

	// ================ General ================

	e.Any("/api/common/ping", ginresp.Wrap(r.commonHandler.Ping))
	e.POST("/api/common/db-test", ginresp.Wrap(r.commonHandler.DatabaseTest))
	e.GET("/api/common/health", ginresp.Wrap(r.commonHandler.Health))

	// ================ Swagger ================

	e.GET("/documentation/swagger", ginext.RedirectTemporary("/documentation/swagger/"))
	e.GET("/documentation/swagger/", ginresp.Wrap(swagger.Handle))
	e.GET("/documentation/swagger/:fn", ginresp.Wrap(swagger.Handle))

	// ================ Website ================

	e.GET("/", ginresp.Wrap(r.websiteHandler.Index))
	e.GET("/index.php", ginresp.Wrap(r.websiteHandler.Index))
	e.GET("/index.html", ginresp.Wrap(r.websiteHandler.Index))
	e.GET("/index", ginresp.Wrap(r.websiteHandler.Index))

	e.GET("/api", ginresp.Wrap(r.websiteHandler.APIDocs))
	e.GET("/api.php", ginresp.Wrap(r.websiteHandler.APIDocs))
	e.GET("/api.html", ginresp.Wrap(r.websiteHandler.APIDocs))

	e.GET("/api_more", ginresp.Wrap(r.websiteHandler.APIDocsMore))
	e.GET("/api_more.php", ginresp.Wrap(r.websiteHandler.APIDocsMore))
	e.GET("/api_more.html", ginresp.Wrap(r.websiteHandler.APIDocsMore))

	e.GET("/message_sent", ginresp.Wrap(r.websiteHandler.MessageSent))
	e.GET("/message_sent.php", ginresp.Wrap(r.websiteHandler.MessageSent))
	e.GET("/message_sent.html", ginresp.Wrap(r.websiteHandler.MessageSent))

	e.GET("/favicon.ico", ginresp.Wrap(r.websiteHandler.FaviconIco))
	e.GET("/favicon.png", ginresp.Wrap(r.websiteHandler.FaviconPNG))

	e.GET("/js/:fn", ginresp.Wrap(r.websiteHandler.Javascript))
	e.GET("/css/:fn", ginresp.Wrap(r.websiteHandler.CSS))

	// ================ Compat (v1) ================

	compat := e.Group("/api/")
	{
		compat.GET("/register.php", ginresp.Wrap(r.compatHandler.Register))
		compat.GET("/info.php", ginresp.Wrap(r.compatHandler.Info))
		compat.GET("/ack.php", ginresp.Wrap(r.compatHandler.Ack))
		compat.GET("/requery.php", ginresp.Wrap(r.compatHandler.Requery))
		compat.GET("/update.php", ginresp.Wrap(r.compatHandler.Update))
		compat.GET("/expand.php", ginresp.Wrap(r.compatHandler.Expand))
		compat.GET("/upgrade.php", ginresp.Wrap(r.compatHandler.Upgrade))
	}

	// ================ Manage API ================

	apiv2 := e.Group("/api-v2/")
	{

		apiv2.POST("/users/", ginresp.Wrap(r.apiHandler.CreateUser))
		apiv2.GET("/users/:uid", ginresp.Wrap(r.apiHandler.GetUser))
		apiv2.PATCH("/users/:uid", ginresp.Wrap(r.apiHandler.UpdateUser))

		apiv2.GET("/users/:uid/clients", ginresp.Wrap(r.apiHandler.ListClients))
		apiv2.GET("/users/:uid/clients/:cid", ginresp.Wrap(r.apiHandler.GetClient))
		apiv2.POST("/users/:uid/clients", ginresp.Wrap(r.apiHandler.AddClient))
		apiv2.DELETE("/users/:uid/clients", ginresp.Wrap(r.apiHandler.DeleteClient))

		apiv2.GET("/users/:uid/channels", ginresp.Wrap(r.apiHandler.ListChannels))
		apiv2.GET("/users/:uid/channels/:cid", ginresp.Wrap(r.apiHandler.GetChannel))
		apiv2.GET("/users/:uid/channels/:cid/messages", ginresp.Wrap(r.apiHandler.ListChannelMessages))
		apiv2.GET("/users/:uid/channels/:cid/subscriptions", ginresp.Wrap(r.apiHandler.ListChannelSubscriptions))

		apiv2.GET("/users/:uid/subscriptions", ginresp.Wrap(r.apiHandler.ListUserSubscriptions))
		apiv2.GET("/users/:uid/subscriptions/:sid", ginresp.Wrap(r.apiHandler.GetSubscription))
		apiv2.DELETE("/users/:uid/subscriptions/:sid", ginresp.Wrap(r.apiHandler.CancelSubscription))
		apiv2.POST("/users/:uid/subscriptions", ginresp.Wrap(r.apiHandler.CreateSubscription))
		apiv2.PATCH("/users/:uid/subscriptions", ginresp.Wrap(r.apiHandler.UpdateSubscription))

		apiv2.GET("/messages", ginresp.Wrap(r.apiHandler.ListMessages))
		apiv2.GET("/messages/:mid", ginresp.Wrap(r.apiHandler.GetMessage))
		apiv2.DELETE("/messages/:mid", ginresp.Wrap(r.apiHandler.DeleteMessage))

		apiv2.POST("/messages", ginresp.Wrap(r.apiHandler.SendMessage))
	}

	// ================ Send API ================

	sendAPI := e.Group("")
	{
		sendAPI.POST("/", ginresp.Wrap(r.messageHandler.SendMessage))
		sendAPI.POST("/send", ginresp.Wrap(r.messageHandler.SendMessage))
		sendAPI.POST("/send.php", ginresp.Wrap(r.messageHandler.SendMessage))
	}

	if r.app.Config.ReturnRawErrors {
		e.NoRoute(ginresp.Wrap(r.commonHandler.NoRoute))
	}

}
