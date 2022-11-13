package logic

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
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
