package api

import (
	"blackforestbytes.com/simplecloudnotifier/api/ginext"
	"blackforestbytes.com/simplecloudnotifier/api/ginresp"
	"blackforestbytes.com/simplecloudnotifier/api/handler"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"blackforestbytes.com/simplecloudnotifier/models"
	"blackforestbytes.com/simplecloudnotifier/swagger"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Router struct {
	app *logic.Application

	commonHandler   handler.CommonHandler
	compatHandler   handler.CompatHandler
	websiteHandler  handler.WebsiteHandler
	apiHandler      handler.APIHandler
	messageHandler  handler.MessageHandler
	externalHandler handler.ExternalHandler
}

func NewRouter(app *logic.Application) *Router {
	return &Router{
		app: app,

		commonHandler:   handler.NewCommonHandler(app),
		compatHandler:   handler.NewCompatHandler(app),
		websiteHandler:  handler.NewWebsiteHandler(app),
		apiHandler:      handler.NewAPIHandler(app),
		messageHandler:  handler.NewMessageHandler(app),
		externalHandler: handler.NewExternalHandler(app),
	}
}

// Init swaggerdocs
//
//	@title			SimpleCloudNotifier API
//	@version		2.0
//	@description	API for SCN
//	@host			simplecloudnotifier.de
//
//	@tag.name		External
//	@tag.name		API-v1
//	@tag.name		API-v2
//	@tag.name		Common
//
//	@BasePath		/
func (r *Router) Init(e *gin.Engine) error {

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("entityid", models.ValidateEntityID, true)
		if err != nil {
			return err
		}
	} else {
		return errors.New("failed to add validators - wrong engine")
	}

	// ================ General (unversioned) ================

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

	compat := e.Group("/api")
	{
		compat.GET("/register.php", r.Wrap(r.compatHandler.Register))
		compat.GET("/info.php", r.Wrap(r.compatHandler.Info))
		compat.GET("/ack.php", r.Wrap(r.compatHandler.Ack))
		compat.GET("/requery.php", r.Wrap(r.compatHandler.Requery))
		compat.GET("/update.php", r.Wrap(r.compatHandler.Update))
		compat.GET("/expand.php", r.Wrap(r.compatHandler.Expand))
		compat.GET("/upgrade.php", r.Wrap(r.compatHandler.Upgrade))
	}

	// ================ Manage API (v2) ================

	apiv2 := e.Group("/api/v2/")
	{
		apiv2.POST("/users", r.Wrap(r.apiHandler.CreateUser))
		apiv2.GET("/users/:uid", r.Wrap(r.apiHandler.GetUser))
		apiv2.PATCH("/users/:uid", r.Wrap(r.apiHandler.UpdateUser))

		apiv2.GET("/users/:uid/keys", r.Wrap(r.apiHandler.ListUserKeys))
		apiv2.POST("/users/:uid/keys", r.Wrap(r.apiHandler.CreateUserKey))
		apiv2.GET("/users/:uid/keys/current", r.Wrap(r.apiHandler.GetCurrentUserKey))
		apiv2.GET("/users/:uid/keys/:kid", r.Wrap(r.apiHandler.GetUserKey))
		apiv2.PATCH("/users/:uid/keys/:kid", r.Wrap(r.apiHandler.UpdateUserKey))
		apiv2.DELETE("/users/:uid/keys/:kid", r.Wrap(r.apiHandler.DeleteUserKey))

		apiv2.GET("/users/:uid/clients", r.Wrap(r.apiHandler.ListClients))
		apiv2.GET("/users/:uid/clients/:cid", r.Wrap(r.apiHandler.GetClient))
		apiv2.PATCH("/users/:uid/clients/:cid", r.Wrap(r.apiHandler.UpdateClient))
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

		apiv2.GET("/preview/users/:uid", r.Wrap(r.apiHandler.GetUserPreview))
		apiv2.GET("/preview/keys/:kid", r.Wrap(r.apiHandler.GetUserKeyPreview))
		apiv2.GET("/preview/channels/:cid", r.Wrap(r.apiHandler.GetChannelPreview))
	}

	// ================ Send API (unversioned) ================

	sendAPI := e.Group("")
	{
		sendAPI.POST("/", r.Wrap(r.messageHandler.SendMessage))
		sendAPI.POST("/send", r.Wrap(r.messageHandler.SendMessage))
		sendAPI.POST("/send.php", r.Wrap(r.compatHandler.SendMessage))

		sendAPI.POST("/external/v1/uptime-kuma", r.Wrap(r.externalHandler.UptimeKuma))

	}

	// ================

	if r.app.Config.ReturnRawErrors {
		e.NoRoute(r.Wrap(r.commonHandler.NoRoute))
	}

	// ================

	return nil
}

func (r *Router) Wrap(fn ginresp.WHandlerFunc) gin.HandlerFunc {
	return ginresp.Wrap(r.app, fn)
}
