package logic

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/db"
	"blackforestbytes.com/simplecloudnotifier/db/simplectx"
	"blackforestbytes.com/simplecloudnotifier/google"
	"blackforestbytes.com/simplecloudnotifier/models"
	"blackforestbytes.com/simplecloudnotifier/push"
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
	"gogs.mikescher.com/BlackForestBytes/goext/rext"
	"gogs.mikescher.com/BlackForestBytes/goext/syncext"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"
)

var rexWhitespaceStart = rext.W(regexp.MustCompile("^\\s+"))
var rexWhitespaceEnd = rext.W(regexp.MustCompile("\\s+$"))
var rexNormalizeUsername = rext.W(regexp.MustCompile("[^[:alnum:]\\-_ ]"))
var rexCompatTitleChannel = rext.W(regexp.MustCompile("^\\[(?P<channel>[A-Za-z\\-0-9_ ]+)] (?P<title>(.|\\r|\\n)+)$"))

type Application struct {
	Config           scn.Config
	Gin              *ginext.GinWrapper
	Database         *DBPool
	Pusher           push.NotificationClient
	AndroidPublisher google.AndroidPublisherClient
	Jobs             []Job
	stopChan         chan bool
	Port             string
	IsRunning        *syncext.AtomicBool
	RequestLogQueue  chan models.RequestLog
}

func NewApp(db *DBPool) *Application {
	return &Application{
		Database:        db,
		stopChan:        make(chan bool),
		IsRunning:       syncext.NewAtomicBool(false),
		RequestLogQueue: make(chan models.RequestLog, 1024),
	}
}

func (app *Application) Init(cfg scn.Config, g *ginext.GinWrapper, fb push.NotificationClient, apc google.AndroidPublisherClient, jobs []Job) {
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

	// ================== START HTTP ==================

	addr := net.JoinHostPort(app.Config.ServerIP, app.Config.ServerPort)

	errChan, httpserver := app.Gin.ListenAndServeHTTP(addr, func(port string) {
		app.Port = port
		app.IsRunning.Set(true)
	})

	// ================== START JOBS ==================

	for _, job := range app.Jobs {
		err := job.Start()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to start job")
		}
	}

	// ================== LISTEN FOR SIGNALS ==================

	sigstop := make(chan os.Signal, 1)
	signal.Notify(sigstop, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sigstop:
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

	case <-app.stopChan:
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

	// ================== STOP JOBS ==================

	for _, job := range app.Jobs {
		job.Stop()
	}

	// ================== STOP DB ==================

	{
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		err := app.Database.Stop(ctx)
		if err != nil {
			log.Err(err).Msg("Failed to stop database")
		}
	}
	log.Info().Msg("Stopped Databases")

	// ================== FINISH ==================

	app.IsRunning.Set(false)
}

func (app *Application) GenerateRandomAuthKey() string {
	return scn.RandomAuthKey()
}

func (app *Application) VerifyProToken(ctx *AppContext, token string) (bool, error) {

	if strings.HasPrefix(token, "ANDROID|v1|") {
		subToken := token[len("ANDROID|v1|"):]
		return app.VerifyAndroidProToken(ctx, subToken)
	}

	if strings.HasPrefix(token, "ANDROID|v2|") {
		subToken := token[len("ANDROID|v2|"):]
		return app.VerifyAndroidProToken(ctx, subToken)
	}

	if strings.HasPrefix(token, "IOS|v1|") {
		return false, errors.New("invalid token-version: ios-v1")
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

func (app *Application) NewSimpleTransactionContext(timeout time.Duration) *simplectx.SimpleContext {
	ictx, cancel := context.WithTimeout(context.Background(), timeout)
	return simplectx.CreateSimpleContext(ictx, cancel)
}

func (app *Application) getPermissions(ctx db.TxContext, hdr string) (models.PermissionSet, error) {
	if hdr == "" {
		return models.NewEmptyPermissions(), nil
	}

	if !strings.HasPrefix(hdr, "SCN ") {
		return models.NewEmptyPermissions(), nil
	}

	key := strings.TrimSpace(hdr[4:])

	tok, err := app.Database.Primary.GetKeyTokenByToken(ctx, key)
	if err != nil {
		return models.PermissionSet{}, err
	}

	if tok != nil {

		err = app.Database.Primary.UpdateKeyTokenLastUsed(ctx, tok.KeyTokenID)
		if err != nil {
			return models.PermissionSet{}, err
		}

		return models.PermissionSet{Token: tok}, nil
	}

	return models.NewEmptyPermissions(), nil
}

func (app *Application) GetOrCreateChannel(ctx *AppContext, userid models.UserID, displayChanName string, intChanName string) (models.Channel, error) {
	existingChan, err := app.Database.Primary.GetChannelByName(ctx, userid, intChanName)
	if err != nil {
		return models.Channel{}, err
	}

	if existingChan != nil {
		return *existingChan, nil
	}

	subscribeKey := app.GenerateRandomAuthKey()

	newChan, err := app.Database.Primary.CreateChannel(ctx, userid, displayChanName, intChanName, subscribeKey, nil)
	if err != nil {
		return models.Channel{}, err
	}

	_, err = app.Database.Primary.CreateSubscription(ctx, userid, newChan, true)
	if err != nil {
		return models.Channel{}, err
	}

	return newChan, nil
}

func (app *Application) NormalizeChannelDisplayName(v string) string {
	return strings.TrimSpace(v)
}

func (app *Application) NormalizeChannelInternalName(v string) string {
	return strings.TrimSpace(v)
}

func (app *Application) NormalizeUsername(v string) string {
	return strings.TrimSpace(v)
}

func (app *Application) DeliverMessage(ctx context.Context, user models.User, client models.Client, channel models.Channel, msg models.Message) (string, error) {
	fcmDelivID, err := app.Pusher.SendNotification(ctx, user, client, channel, msg)
	if err != nil {
		log.Warn().Str("MessageID", msg.MessageID.String()).Str("ClientID", client.ClientID.String()).Err(err).Msg("FCM Delivery failed")
		return "", err
	}
	return fcmDelivID, nil
}

func (app *Application) InsertRequestLog(data models.RequestLog) {
	ok := syncext.WriteNonBlocking(app.RequestLogQueue, data)
	if !ok {
		log.Error().Msg("failed to insert request-log (queue full)")
	}
}
