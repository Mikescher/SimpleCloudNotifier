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
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type Void = struct{}

func StartSimpleWebserver(t *testing.T) (*logic.Application, string, func()) {
	InitTests()

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

	TPrintln("DatabaseFile: " + dbfile)

	conf := scn.Config{
		Namespace:         "test",
		GinDebug:          true,
		ServerIP:          "0.0.0.0",
		ServerPort:        "0", // simply choose a free port
		DBFile:            dbfile,
		DBJournal:         "WAL",
		DBTimeout:         500 * time.Millisecond,
		DBMaxOpenConns:    2,
		DBMaxIdleConns:    2,
		DBConnMaxLifetime: 1 * time.Second,
		DBConnMaxIdleTime: 1 * time.Second,
		RequestTimeout:    30 * time.Second,
		ReturnRawErrors:   true,
		DummyFirebase:     true,
	}

	sqlite, err := db.NewDatabase(conf)
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

	stop := func() {
		app.Stop()
		_ = os.Remove(dbfile)
		_ = app.IsRunning.WaitWithTimeout(400*time.Millisecond, false)
	}

	go func() { app.Run() }()

	err = app.IsRunning.WaitWithTimeout(100*time.Millisecond, true)
	if err != nil {
		TestFailErr(t, err)
	}

	return app, "http://127.0.0.1:" + app.Port, stop
}
