package logic

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/common/ginresp"
	"blackforestbytes.com/simplecloudnotifier/db"
	"blackforestbytes.com/simplecloudnotifier/firebase"
	"blackforestbytes.com/simplecloudnotifier/models"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
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
	Config   scn.Config
	Gin      *gin.Engine
	Database *db.Database
	Firebase *firebase.App
}

func NewApp(db *db.Database) *Application {
	return &Application{Database: db}
}

func (app *Application) Init(cfg scn.Config, g *gin.Engine, fb *firebase.App) {
	app.Config = cfg
	app.Gin = g
	app.Firebase = fb
}

func (app *Application) Run() {
	httpserver := &http.Server{
		Addr:    net.JoinHostPort(app.Config.ServerIP, app.Config.ServerPort),
		Handler: app.Gin,
	}

	errChan := make(chan error)

	go func() {
		log.Info().Str("address", httpserver.Addr).Msg("HTTP-Server started on http://localhost:" + app.Config.ServerPort)
		errChan <- httpserver.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	select {
	case <-stop:
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		log.Info().Msg("Stopping HTTP-Server")
		err := httpserver.Shutdown(ctx)
		if err != nil {
			log.Info().Err(err).Msg("Error while stopping the http-server")
			return
		}
		log.Info().Msg("Stopped HTTP-Server")

	case err := <-errChan:
		log.Error().Err(err).Msg("HTTP-Server failed")
	}

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

func (app *Application) VerifyProToken(token string) (bool, error) {
	return false, nil //TODO implement pro verification
}

func (app *Application) Migrate() error {
	ctx, cancel := context.WithTimeout(context.Background(), 24*time.Second)
	defer cancel()

	return app.Database.Migrate(ctx)
}

func (app *Application) StartRequest(g *gin.Context, uri any, query any, body any) (*AppContext, *ginresp.HTTPResponse) {

	if body != nil {
		if err := g.ShouldBindJSON(&body); err != nil {
			return nil, langext.Ptr(ginresp.InternAPIError(400, apierr.BINDFAIL_BODY_PARAM, "Failed to read body", err))
		}
	}

	if query != nil {
		if err := g.ShouldBindQuery(&query); err != nil {
			return nil, langext.Ptr(ginresp.InternAPIError(400, apierr.BINDFAIL_QUERY_PARAM, "Failed to read query", err))
		}
	}

	if uri != nil {
		if err := g.ShouldBindUri(&uri); err != nil {
			return nil, langext.Ptr(ginresp.InternAPIError(400, apierr.BINDFAIL_URI_PARAM, "Failed to read uri", err))
		}
	}

	ictx, cancel := context.WithTimeout(context.Background(), app.Config.RequestTimeout)
	actx := CreateAppContext(ictx, cancel)

	authheader := g.GetHeader("Authorization")

	perm, err := app.getPermissions(actx, authheader)
	if err != nil {
		cancel()
		return nil, langext.Ptr(ginresp.InternAPIError(400, apierr.PERM_QUERY_FAIL, "Failed to determine permissions", err))
	}

	actx.permissions = perm

	return actx, nil
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

func (app *Application) GetOrCreateChannel(ctx *AppContext, userid int64, chanName string) (models.Channel, error) {
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
