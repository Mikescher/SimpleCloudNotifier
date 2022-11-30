package util

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/api"
	"blackforestbytes.com/simplecloudnotifier/common/ginext"
	"blackforestbytes.com/simplecloudnotifier/db"
	"blackforestbytes.com/simplecloudnotifier/google"
	"blackforestbytes.com/simplecloudnotifier/jobs"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"blackforestbytes.com/simplecloudnotifier/push"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type Void = struct{}

func StartSimpleWebserver(t *testing.T) (*logic.Application, func()) {
	cw := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05 Z07:00",
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	multi := zerolog.MultiLevelWriter(cw)
	logger := zerolog.New(multi).With().
		Timestamp().
		Caller().
		Logger()

	log.Logger = logger

	gin.SetMode(gin.TestMode)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	uuid2, _ := langext.NewHexUUID()
	dbdir := t.TempDir()
	dbfile := filepath.Join(dbdir, uuid2+".sqlite3")

	err := os.MkdirAll(dbdir, os.ModePerm)
	if err != nil {
		TestFailErr(t, err)
	}

	f, err := os.Create(dbfile)
	if err != nil {
		TestFailErr(t, err)
	}
	err = f.Close()
	if err != nil {
		TestFailErr(t, err)
	}

	err = os.Chmod(dbfile, 0777)
	if err != nil {
		TestFailErr(t, err)
	}

	fmt.Println("DatabaseFile: " + dbfile)

	conf := scn.Config{
		Namespace:       "test",
		GinDebug:        true,
		ServerIP:        "0.0.0.0",
		ServerPort:      "0", // simply choose a free port
		DBFile:          dbfile,
		RequestTimeout:  30 * time.Second,
		ReturnRawErrors: true,
		DummyFirebase:   true,
	}

	sqlite, err := db.NewDatabase(dbfile)
	if err != nil {
		TestFailErr(t, err)
	}

	app := logic.NewApp(sqlite)

	if err := app.Migrate(); err != nil {
		TestFailErr(t, err)
	}

	ginengine := ginext.NewEngine(conf)

	router := api.NewRouter(app)

	nc := push.NewTestSink()

	apc := google.NewDummy()

	jobRetry := jobs.NewDeliveryRetryJob(app)
	app.Init(conf, ginengine, nc, apc, []logic.Job{jobRetry})

	router.Init(ginengine)

	stop := func() { app.Stop(); _ = os.Remove(dbfile) }
	go func() { app.Run() }()

	time.Sleep(100 * time.Millisecond) // wait until http server is up

	return app, stop
}
