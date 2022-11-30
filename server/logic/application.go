package logic

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/common/ginresp"
	"blackforestbytes.com/simplecloudnotifier/db"
	"blackforestbytes.com/simplecloudnotifier/google"
	"blackforestbytes.com/simplecloudnotifier/models"
	"blackforestbytes.com/simplecloudnotifier/push"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/syncext"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"
)

type Application struct {
	Config           scn.Config
	Gin              *gin.Engine
	Database         *db.Database
	Pusher           push.NotificationClient
	AndroidPublisher google.AndroidPublisherClient
	Jobs             []Job
	stopChan         chan bool
	Port             string
	IsRunning        *syncext.AtomicBool
}

func NewApp(db *db.Database) *Application {
	return &Application{
		Database:  db,
		stopChan:  make(chan bool),
		IsRunning: syncext.NewAtomicBool(false),
	}
}

func (app *Application) Init(cfg scn.Config, g *gin.Engine, fb push.NotificationClient, apc google.AndroidPublisherClient, jobs []Job) {
	app.Config = cfg
	app.Gin = g
	app.Pusher = fb
	app.AndroidPublisher = apc
	app.Jobs = jobs
}

func (app *Application) Stop() {
	// non-blocking send
	select {
	case app.stopChan <- true:
	}
}

func (app *Application) Run() {
	httpserver := &http.Server{
		Addr:    net.JoinHostPort(app.Config.ServerIP, app.Config.ServerPort),
		Handler: app.Gin,
	}

	errChan := make(chan error)

	go func() {

		ln, err := net.Listen("tcp", httpserver.Addr)
		if err != nil {
			errChan <- err
			return
		}

		_, port, err := net.SplitHostPort(ln.Addr().String())
		if err != nil {
			errChan <- err
			return
		}

		log.Info().Str("address", httpserver.Addr).Msg("HTTP-Server started on http://localhost:" + port)

		app.Port = port

		app.IsRunning.Set(true) // the net.Listener a few lines above is at this point actually already buffering requests

		errChan <- httpserver.Serve(ln)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	for _, job := range app.Jobs {
		job.Start()
	}

	select {
	case <-stop:
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		log.Info().Msg("Stopping HTTP-Server")

		err := httpserver.Shutdown(ctx)

		if err != nil {
			log.Info().Err(err).Msg("Error while stopping the http-server")
		} else {
			log.Info().Msg("Stopped HTTP-Server")
		}

	case err := <-errChan:
		log.Error().Err(err).Msg("HTTP-Server failed")

	case _ = <-app.stopChan:
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		log.Info().Msg("Manually stopping HTTP-Server")

		err := httpserver.Shutdown(ctx)

		if err != nil {
			log.Info().Err(err).Msg("Error while stopping the http-server")
		} else {
			log.Info().Msg("Manually stopped HTTP-Server")
		}
	}

	for _, job := range app.Jobs {
		job.Stop()
	}

	app.IsRunning.Set(false)
}

func (app *Application) GenerateRandomAuthKey() string {
	charset := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	k := ""
	for i := 0; i < 64; i++ {
		k += string(charset[rand.Int()%len(charset)])
	}
	return k
}

func (app *Application) QuotaMax(ispro bool) int {
	if ispro {
		return 1000
	} else {
		return 50
	}
}

func (app *Application) VerifyProToken(ctx *AppContext, token string) (bool, error) {
	if strings.HasPrefix(token, "ANDROID|v2|") {
		subToken := token[len("ANDROID|v2|"):]
		return app.VerifyAndroidProToken(ctx, subToken)
	}
	if strings.HasPrefix(token, "IOS|v2|") {
		subToken := token[len("IOS|v2|"):]
		return app.VerifyIOSProToken(ctx, subToken)
	}

	return false, nil
}

func (app *Application) VerifyAndroidProToken(ctx *AppContext, token string) (bool, error) {

	purchase, err := app.AndroidPublisher.GetProductPurchase(ctx, app.Config.GooglePackageName, app.Config.GoogleProProductID, token)
	if err != nil {
		return false, err
	}

	if purchase == nil {
		return false, nil
	}
	if purchase.PurchaseState == nil {
		return false, nil
	}
	if *purchase.PurchaseState != google.PurchaseStatePurchased {
		return false, nil
	}

	return true, nil
}

func (app *Application) VerifyIOSProToken(ctx *AppContext, token string) (bool, error) {
	return false, nil //TODO IOS
}

func (app *Application) Migrate() error {
	ctx, cancel := context.WithTimeout(context.Background(), 24*time.Second)
	defer cancel()

	return app.Database.Migrate(ctx)
}

func (app *Application) StartRequest(g *gin.Context, uri any, query any, body any, form any) (*AppContext, *ginresp.HTTPResponse) {

	if uri != nil {
		if err := g.ShouldBindUri(uri); err != nil {
			return nil, langext.Ptr(ginresp.APIError(g, 400, apierr.BINDFAIL_URI_PARAM, "Failed to read uri", err))
		}
	}

	if query != nil {
		if err := g.ShouldBindQuery(query); err != nil {
			return nil, langext.Ptr(ginresp.APIError(g, 400, apierr.BINDFAIL_QUERY_PARAM, "Failed to read query", err))
		}
	}

	if body != nil && g.ContentType() == "application/json" {
		if err := g.ShouldBindJSON(body); err != nil {
			return nil, langext.Ptr(ginresp.APIError(g, 400, apierr.BINDFAIL_BODY_PARAM, "Failed to read body", err))
		}
	}

	if form != nil && g.ContentType() == "multipart/form-data" {
		if err := g.ShouldBindWith(form, binding.Form); err != nil {
			return nil, langext.Ptr(ginresp.APIError(g, 400, apierr.BINDFAIL_BODY_PARAM, "Failed to read multipart-form", err))
		}
	}

	ictx, cancel := context.WithTimeout(context.Background(), app.Config.RequestTimeout)
	actx := CreateAppContext(g, ictx, cancel)

	authheader := g.GetHeader("Authorization")

	perm, err := app.getPermissions(actx, authheader)
	if err != nil {
		cancel()
		return nil, langext.Ptr(ginresp.APIError(g, 400, apierr.PERM_QUERY_FAIL, "Failed to determine permissions", err))
	}

	actx.permissions = perm

	return actx, nil
}

func (app *Application) NewSimpleTransactionContext(timeout time.Duration) *SimpleContext {
	ictx, cancel := context.WithTimeout(context.Background(), timeout)
	return CreateSimpleContext(ictx, cancel)
}

func (app *Application) getPermissions(ctx *AppContext, hdr string) (PermissionSet, error) {
	if hdr == "" {
		return NewEmptyPermissions(), nil
	}

	if !strings.HasPrefix(hdr, "SCN ") {
		return NewEmptyPermissions(), nil
	}

	key := strings.TrimSpace(hdr[4:])

	user, err := app.Database.GetUserByKey(ctx, key)
	if err != nil {
		return PermissionSet{}, err
	}

	if user != nil && user.SendKey == key {
		return PermissionSet{UserID: langext.Ptr(user.UserID), KeyType: PermKeyTypeUserSend}, nil
	}
	if user != nil && user.ReadKey == key {
		return PermissionSet{UserID: langext.Ptr(user.UserID), KeyType: PermKeyTypeUserRead}, nil
	}
	if user != nil && user.AdminKey == key {
		return PermissionSet{UserID: langext.Ptr(user.UserID), KeyType: PermKeyTypeUserAdmin}, nil
	}

	return NewEmptyPermissions(), nil
}

func (app *Application) GetOrCreateChannel(ctx *AppContext, userid models.UserID, chanName string) (models.Channel, error) {
	chanName = app.NormalizeChannelName(chanName)

	existingChan, err := app.Database.GetChannelByName(ctx, userid, chanName)
	if err != nil {
		return models.Channel{}, err
	}

	if existingChan != nil {
		return *existingChan, nil
	}

	subscribeKey := app.GenerateRandomAuthKey()
	sendKey := app.GenerateRandomAuthKey()

	newChan, err := app.Database.CreateChannel(ctx, userid, chanName, subscribeKey, sendKey)
	if err != nil {
		return models.Channel{}, err
	}

	_, err = app.Database.CreateSubscription(ctx, userid, newChan, true)
	if err != nil {
		return models.Channel{}, err
	}

	return newChan, nil
}

func (app *Application) NormalizeChannelName(v string) string {
	rex := regexp.MustCompile("[^[:alnum:]\\-_]")

	v = strings.TrimSpace(v)
	v = strings.ToLower(v)
	v = rex.ReplaceAllString(v, "")

	return v
}

func (app *Application) NormalizeUsername(v string) string {
	rex := regexp.MustCompile("[^[:alnum:]\\-_ ]")

	v = strings.TrimSpace(v)
	v = rex.ReplaceAllString(v, "")

	return v
}

func (app *Application) DeliverMessage(ctx context.Context, client models.Client, msg models.Message) (*string, error) {
	if client.FCMToken != nil {
		fcmDelivID, err := app.Pusher.SendNotification(ctx, client, msg)
		if err != nil {
			log.Warn().Int64("SCNMessageID", msg.SCNMessageID.IntID()).Int64("ClientID", client.ClientID.IntID()).Err(err).Msg("FCM Delivery failed")
			return nil, err
		}
		return langext.Ptr(fcmDelivID), nil
	} else {
		return langext.Ptr(""), nil
	}
}
