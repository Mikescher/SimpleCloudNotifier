package logic

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/common/ginresp"
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Application struct {
	Config   scn.Config
	Gin      *gin.Engine
	Database *sql.DB
}

func NewApp(db *sql.DB) *Application {
	return &Application{Database: db}
}

func (app *Application) Init(cfg scn.Config, g *gin.Engine) {
	app.Config = cfg
	app.Gin = g
}

func (app *Application) Run() {
	httpserver := &http.Server{
		Addr:    net.JoinHostPort(app.Config.ServerIP, app.Config.ServerPort),
		Handler: app.Gin,
	}

	errChan := make(chan error)

	go func() {
		log.Info().Str("address", httpserver.Addr).Msg("HTTP-Server started")
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

func (app *Application) RunTransaction(ctx context.Context, opt *sql.TxOptions, fn func(tx *sql.Tx) (ginresp.HTTPResponse, bool)) ginresp.HTTPResponse {

	tx, err := app.Database.BeginTx(ctx, opt)
	if err != nil {
		return ginresp.InternAPIError(0, fmt.Sprintf("Failed to create transaction: %v", err))
	}

	res, commit := fn(tx)

	if commit {
		err = tx.Commit()
		if err != nil {
			return ginresp.InternAPIError(0, fmt.Sprintf("Failed to commit transaction: %v", err))
		}
	} else {
		err = tx.Rollback()
		if err != nil {
			return ginresp.InternAPIError(0, fmt.Sprintf("Failed to rollback transaction: %v", err))
		}
	}

	return res
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
