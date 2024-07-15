package api

import (
	"blackforestbytes.com/simplecloudnotifier/api/ginresp"
	"blackforestbytes.com/simplecloudnotifier/api/handler"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"blackforestbytes.com/simplecloudnotifier/models"
	"blackforestbytes.com/simplecloudnotifier/swagger"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
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
func (r *Router) Init(e *ginext.GinWrapper) error {

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("entityid", models.ValidateEntityID, true)
		if err != nil {
			return err
		}
	} else {
		return errors.New("failed to add validators - wrong engine")
	}

	wrap := func(fn ginext.WHandlerFunc) ginext.WHandlerFunc{return Wrap(r.app, fn)}

	// ================ General (unversioned) ================

	commonAPI := e.Routes().Group("/api")
	{
		commonAPI.Any("/ping").Handle(wrap(r.commonHandler.Ping))
		commonAPI.POST("/db-test").Handle(wrap(r.commonHandler.DatabaseTest))
		commonAPI.GET("/health").Handle(wrap(r.commonHandler.Health))
		commonAPI.POST("/sleep/:secs").Handle(wrap(r.commonHandler.Sleep))
	}

	// ================ Swagger ================

	docs := e.Routes().Group("/documentation")
	{
		docs.GET("/swagger").Handle(wrap(ginext.RedirectTemporary("/documentation/swagger/")))
		docs.GET("/swagger/*sub").Handle(wrap(swagger.Handle))
	}

	// ================ Website ================

	frontend := e.Routes().Group("")
	{
		frontend.GET("/").Handle(wrap(r.websiteHandler.Index))
		frontend.GET("/index.php").Handle(wrap(r.websiteHandler.Index))
		frontend.GET("/index.html").Handle(wrap(r.websiteHandler.Index))
		frontend.GET("/index").Handle(wrap(r.websiteHandler.Index))

		frontend.GET("/api").Handle(wrap(r.websiteHandler.APIDocs))
		frontend.GET("/api.php").Handle(wrap(r.websiteHandler.APIDocs))
		frontend.GET("/api.html").Handle(wrap(r.websiteHandler.APIDocs))

		frontend.GET("/api_more").Handle(wrap(r.websiteHandler.APIDocsMore))
		frontend.GET("/api_more.php").Handle(wrap(r.websiteHandler.APIDocsMore))
		frontend.GET("/api_more.html").Handle(wrap(r.websiteHandler.APIDocsMore))

		frontend.GET("/message_sent").Handle(wrap(r.websiteHandler.MessageSent))
		frontend.GET("/message_sent.php").Handle(wrap(r.websiteHandler.MessageSent))
		frontend.GET("/message_sent.html").Handle(wrap(r.websiteHandler.MessageSent))

		frontend.GET("/favicon.ico").Handle(wrap(r.websiteHandler.FaviconIco))
		frontend.GET("/favicon.png").Handle(wrap(r.websiteHandler.FaviconPNG))

		frontend.GET("/js/:fn").Handle(wrap(r.websiteHandler.Javascript))
		frontend.GET("/css/:fn").Handle(wrap(r.websiteHandler.CSS))
	}

	// ================ Compat (v1) ================

	compat := e.Routes().Group("/api")
	{
		compat.GET("/register.php").Handle(wrap(r.compatHandler.Register))
		compat.GET("/info.php").Handle(wrap(r.compatHandler.Info))
		compat.GET("/ack.php").Handle(wrap(r.compatHandler.Ack))
		compat.GET("/requery.php").Handle(wrap(r.compatHandler.Requery))
		compat.GET("/update.php").Handle(wrap(r.compatHandler.Update))
		compat.GET("/expand.php").Handle(wrap(r.compatHandler.Expand))
		compat.GET("/upgrade.php").Handle(wrap(r.compatHandler.Upgrade))
	}

	// ================ Manage API (v2) ================

	apiv2 := e.Routes().Group("/api/v2/")
	{
		apiv2.POST("/users").Handle(wrap(r.apiHandler.CreateUser))
		apiv2.GET("/users/:uid").Handle(wrap(r.apiHandler.GetUser))
		apiv2.PATCH("/users/:uid").Handle(wrap(r.apiHandler.UpdateUser))

		apiv2.GET("/users/:uid/keys").Handle(wrap(r.apiHandler.ListUserKeys))
		apiv2.POST("/users/:uid/keys").Handle(wrap(r.apiHandler.CreateUserKey))
		apiv2.GET("/users/:uid/keys/current").Handle(wrap(r.apiHandler.GetCurrentUserKey))
		apiv2.GET("/users/:uid/keys/:kid").Handle(wrap(r.apiHandler.GetUserKey))
		apiv2.PATCH("/users/:uid/keys/:kid").Handle(wrap(r.apiHandler.UpdateUserKey))
		apiv2.DELETE("/users/:uid/keys/:kid").Handle(wrap(r.apiHandler.DeleteUserKey))

		apiv2.GET("/users/:uid/clients").Handle(wrap(r.apiHandler.ListClients))
		apiv2.GET("/users/:uid/clients/:cid").Handle(wrap(r.apiHandler.GetClient))
		apiv2.PATCH("/users/:uid/clients/:cid").Handle(wrap(r.apiHandler.UpdateClient))
		apiv2.POST("/users/:uid/clients").Handle(wrap(r.apiHandler.AddClient))
		apiv2.DELETE("/users/:uid/clients/:cid").Handle(wrap(r.apiHandler.DeleteClient))

		apiv2.GET("/users/:uid/channels").Handle(wrap(r.apiHandler.ListChannels))
		apiv2.POST("/users/:uid/channels").Handle(wrap(r.apiHandler.CreateChannel))
		apiv2.GET("/users/:uid/channels/:cid").Handle(wrap(r.apiHandler.GetChannel))
		apiv2.PATCH("/users/:uid/channels/:cid").Handle(wrap(r.apiHandler.UpdateChannel))
		apiv2.GET("/users/:uid/channels/:cid/messages").Handle(wrap(r.apiHandler.ListChannelMessages))
		apiv2.GET("/users/:uid/channels/:cid/subscriptions").Handle(wrap(r.apiHandler.ListChannelSubscriptions))

		apiv2.GET("/users/:uid/subscriptions").Handle(wrap(r.apiHandler.ListUserSubscriptions))
		apiv2.POST("/users/:uid/subscriptions").Handle(wrap(r.apiHandler.CreateSubscription))
		apiv2.GET("/users/:uid/subscriptions/:sid").Handle(wrap(r.apiHandler.GetSubscription))
		apiv2.DELETE("/users/:uid/subscriptions/:sid").Handle(wrap(r.apiHandler.CancelSubscription))
		apiv2.PATCH("/users/:uid/subscriptions/:sid").Handle(wrap(r.apiHandler.UpdateSubscription))

		apiv2.GET("/messages").Handle(wrap(r.apiHandler.ListMessages))
		apiv2.GET("/messages/:mid").Handle(wrap(r.apiHandler.GetMessage))
		apiv2.DELETE("/messages/:mid").Handle(wrap(r.apiHandler.DeleteMessage))

		apiv2.GET("/preview/users/:uid").Handle(wrap(r.apiHandler.GetUserPreview))
		apiv2.GET("/preview/keys/:kid").Handle(wrap(r.apiHandler.GetUserKeyPreview))
		apiv2.GET("/preview/channels/:cid").Handle(wrap(r.apiHandler.GetChannelPreview))
	}

	// ================ Send API (unversioned) ================

	sendAPI := e.Routes().Group("")
	{
		sendAPI.POST("/").Handle(wrap(r.messageHandler.SendMessage)
		sendAPI.POST("/send").Handle(wrap(r.messageHandler.SendMessage)
		sendAPI.POST("/send.php").Handle(wrap(r.compatHandler.SendMessage)

		sendAPI.POST("/external/v1/uptime-kuma").Handle(wrap(r.externalHandler.UptimeKuma)

	}

	// ================

	if r.app.Config.ReturnRawErrors {
		e.NoRoute(r.commonHandler.NoRoute)
	}

	// ================

	return nil
}
