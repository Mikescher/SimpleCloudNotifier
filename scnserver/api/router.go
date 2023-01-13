package api

import (
	"blackforestbytes.com/simplecloudnotifier/api/ginext"
	"blackforestbytes.com/simplecloudnotifier/api/ginresp"
	"blackforestbytes.com/simplecloudnotifier/api/handler"
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
//
// @title       SimpleCloudNotifier API
// @version     2.0
// @description API for SCN
// @host        scn.blackforestbytes.com
//
// @tag.name    External
// @tag.name    API-v1
// @tag.name    API-v2
// @tag.name    Common
//
// @BasePath    /
func (r *Router) Init(e *gin.Engine) {

	// ================ General ================

	commonAPI := e.Group("/api")
	{
		commonAPI.Any("/ping", r.Wrap(r.commonHandler.Ping))
		commonAPI.POST("/db-test", r.Wrap(r.commonHandler.DatabaseTest))
		commonAPI.GET("/health", r.Wrap(r.commonHandler.Health))
		commonAPI.POST("/sleep/:secs", r.Wrap(r.commonHandler.Sleep))
	}

	// ================ Swagger ================

	docs := e.Group("/documentation")
	{
		docs.GET("/swagger", ginext.RedirectTemporary("/documentation/swagger/"))
		docs.GET("/swagger/*sub", r.Wrap(swagger.Handle))
	}

	// ================ Website ================

	frontend := e.Group("")
	{
		frontend.GET("/", r.Wrap(r.websiteHandler.Index))
		frontend.GET("/index.php", r.Wrap(r.websiteHandler.Index))
		frontend.GET("/index.html", r.Wrap(r.websiteHandler.Index))
		frontend.GET("/index", r.Wrap(r.websiteHandler.Index))

		frontend.GET("/api", r.Wrap(r.websiteHandler.APIDocs))
		frontend.GET("/api.php", r.Wrap(r.websiteHandler.APIDocs))
		frontend.GET("/api.html", r.Wrap(r.websiteHandler.APIDocs))

		frontend.GET("/api_more", r.Wrap(r.websiteHandler.APIDocsMore))
		frontend.GET("/api_more.php", r.Wrap(r.websiteHandler.APIDocsMore))
		frontend.GET("/api_more.html", r.Wrap(r.websiteHandler.APIDocsMore))

		frontend.GET("/message_sent", r.Wrap(r.websiteHandler.MessageSent))
		frontend.GET("/message_sent.php", r.Wrap(r.websiteHandler.MessageSent))
		frontend.GET("/message_sent.html", r.Wrap(r.websiteHandler.MessageSent))

		frontend.GET("/favicon.ico", r.Wrap(r.websiteHandler.FaviconIco))
		frontend.GET("/favicon.png", r.Wrap(r.websiteHandler.FaviconPNG))

		frontend.GET("/js/:fn", r.Wrap(r.websiteHandler.Javascript))
		frontend.GET("/css/:fn", r.Wrap(r.websiteHandler.CSS))
	}

	// ================ Compat (v1) ================

	compat := e.Group("/api/")
	{
		compat.GET("/register.php", r.Wrap(r.compatHandler.Register))
		compat.GET("/info.php", r.Wrap(r.compatHandler.Info))
		compat.GET("/ack.php", r.Wrap(r.compatHandler.Ack))
		compat.GET("/requery.php", r.Wrap(r.compatHandler.Requery))
		compat.GET("/update.php", r.Wrap(r.compatHandler.Update))
		compat.GET("/expand.php", r.Wrap(r.compatHandler.Expand))
		compat.GET("/upgrade.php", r.Wrap(r.compatHandler.Upgrade))
	}

	// ================ Manage API ================

	apiv2 := e.Group("/api/")
	{

		apiv2.POST("/users", r.Wrap(r.apiHandler.CreateUser))
		apiv2.GET("/users/:uid", r.Wrap(r.apiHandler.GetUser))
		apiv2.PATCH("/users/:uid", r.Wrap(r.apiHandler.UpdateUser))

		apiv2.GET("/users/:uid/clients", r.Wrap(r.apiHandler.ListClients))
		apiv2.GET("/users/:uid/clients/:cid", r.Wrap(r.apiHandler.GetClient))
		apiv2.POST("/users/:uid/clients", r.Wrap(r.apiHandler.AddClient))
		apiv2.DELETE("/users/:uid/clients/:cid", r.Wrap(r.apiHandler.DeleteClient))

		apiv2.GET("/users/:uid/channels", r.Wrap(r.apiHandler.ListChannels))
		apiv2.POST("/users/:uid/channels", r.Wrap(r.apiHandler.CreateChannel))
		apiv2.GET("/users/:uid/channels/:cid", r.Wrap(r.apiHandler.GetChannel))
		apiv2.PATCH("/users/:uid/channels/:cid", r.Wrap(r.apiHandler.UpdateChannel))
		apiv2.GET("/users/:uid/channels/:cid/messages", r.Wrap(r.apiHandler.ListChannelMessages))
		apiv2.GET("/users/:uid/channels/:cid/subscriptions", r.Wrap(r.apiHandler.ListChannelSubscriptions))

		apiv2.GET("/users/:uid/subscriptions", r.Wrap(r.apiHandler.ListUserSubscriptions))
		apiv2.POST("/users/:uid/subscriptions", r.Wrap(r.apiHandler.CreateSubscription))
		apiv2.GET("/users/:uid/subscriptions/:sid", r.Wrap(r.apiHandler.GetSubscription))
		apiv2.DELETE("/users/:uid/subscriptions/:sid", r.Wrap(r.apiHandler.CancelSubscription))
		apiv2.PATCH("/users/:uid/subscriptions/:sid", r.Wrap(r.apiHandler.UpdateSubscription))

		apiv2.GET("/messages", r.Wrap(r.apiHandler.ListMessages))
		apiv2.GET("/messages/:mid", r.Wrap(r.apiHandler.GetMessage))
		apiv2.DELETE("/messages/:mid", r.Wrap(r.apiHandler.DeleteMessage))
	}

	// ================ Send API ================

	sendAPI := e.Group("")
	{
		sendAPI.POST("/", r.Wrap(r.messageHandler.SendMessage))
		sendAPI.POST("/send", r.Wrap(r.messageHandler.SendMessage))
		sendAPI.POST("/send.php", r.Wrap(r.messageHandler.SendMessageCompat))
	}

	if r.app.Config.ReturnRawErrors {
		e.NoRoute(r.Wrap(r.commonHandler.NoRoute))
	}

}

func (r *Router) Wrap(fn ginresp.WHandlerFunc) gin.HandlerFunc {
	return ginresp.Wrap(r.app, fn)
}
