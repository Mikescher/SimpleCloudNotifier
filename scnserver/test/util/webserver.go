package util

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/api"
	"blackforestbytes.com/simplecloudnotifier/google"
	"blackforestbytes.com/simplecloudnotifier/jobs"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"blackforestbytes.com/simplecloudnotifier/push"
	"github.com/rs/zerolog"
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type Void = struct{}

func StartSimpleWebserver(t *testing.T) (*logic.Application, string, func()) {
	InitTests()

	uuid1, _ := langext.NewHexUUID()
	uuid2, _ := langext.NewHexUUID()
	uuid3, _ := langext.NewHexUUID()

	dbdir := t.TempDir()
	dbfile1 := filepath.Join(dbdir, uuid1+".sqlite3")
	dbfile2 := filepath.Join(dbdir, uuid2+".sqlite3")
	dbfile3 := filepath.Join(dbdir, uuid3+".sqlite3")

	err := os.MkdirAll(dbdir, os.ModePerm)
	if err != nil {
		TestFailErr(t, err)
	}

	f1, err := os.Create(dbfile1)
	if err != nil {
		TestFailErr(t, err)
	}
	err = f1.Close()
	if err != nil {
		TestFailErr(t, err)
	}
	err = os.Chmod(dbfile1, 0777)
	if err != nil {
		TestFailErr(t, err)
	}
	f2, err := os.Create(dbfile2)
	if err != nil {
		TestFailErr(t, err)
	}
	err = f2.Close()
	if err != nil {
		TestFailErr(t, err)
	}
	err = os.Chmod(dbfile2, 0777)
	if err != nil {
		TestFailErr(t, err)
	}
	f3, err := os.Create(dbfile3)
	if err != nil {
		TestFailErr(t, err)
	}
	err = f3.Close()
	if err != nil {
		TestFailErr(t, err)
	}
	err = os.Chmod(dbfile3, 0777)
	if err != nil {
		TestFailErr(t, err)
	}

	TPrintln(zerolog.InfoLevel, "DatabaseFile<main>:      "+dbfile1)
	TPrintln(zerolog.InfoLevel, "DatabaseFile<requests>:  "+dbfile2)
	TPrintln(zerolog.InfoLevel, "DatabaseFile<logs>:      "+dbfile3)

	scn.Conf = CreateTestConfig(t, dbfile1, dbfile2, dbfile3)

	sqlite, err := logic.NewDBPool(scn.Conf)
	if err != nil {
		TestFailErr(t, err)
	}

	app := logic.NewApp(sqlite)

	if err := app.Migrate(); err != nil {
		TestFailErr(t, err)
	}

	ginengine := ginext.NewEngine(ginext.Options{
		AllowCors:             &scn.Conf.Cors,
		GinDebug:              &scn.Conf.GinDebug,
		BufferBody:            langext.PTrue,
		Timeout:               langext.Ptr(time.Duration(int64(scn.Conf.RequestTimeout) * int64(scn.Conf.RequestMaxRetry))),
		BuildRequestBindError: logic.BuildGinRequestError,
	})

	router := api.NewRouter(app)

	nc := push.NewTestSink()

	apc := google.NewDummy()

	app.Init(scn.Conf, ginengine, nc, apc, []logic.Job{
		jobs.NewDeliveryRetryJob(app),
		jobs.NewRequestLogCollectorJob(app),
	})

	err = router.Init(ginengine)
	if err != nil {
		panic(err)
	}

	stop := func() {
		TPrintln(zerolog.InfoLevel, "Stopping App")
		app.Stop()
		_ = app.IsRunning.WaitWithTimeout(5*time.Second, false)
		TPrintln(zerolog.InfoLevel, "Stopped App")
		_ = os.Remove(dbfile1)
		_ = os.Remove(dbfile2)
		_ = os.Remove(dbfile3)
	}

	go func() { app.Run() }()

	err = app.IsRunning.WaitWithTimeout(100*time.Millisecond, true)
	if err != nil {
		TestFailErr(t, err)
	}

	return app, "http://127.0.0.1:" + app.Port, stop
}

func StartSimpleTestspace(t *testing.T) (string, string, string, scn.Config, func()) {
	InitTests()

	uuid1, _ := langext.NewHexUUID()
	uuid2, _ := langext.NewHexUUID()
	uuid3, _ := langext.NewHexUUID()

	dbdir := t.TempDir()
	dbfile1 := filepath.Join(dbdir, uuid1+".sqlite3")
	dbfile2 := filepath.Join(dbdir, uuid2+".sqlite3")
	dbfile3 := filepath.Join(dbdir, uuid3+".sqlite3")

	err := os.MkdirAll(dbdir, os.ModePerm)
	if err != nil {
		TestFailErr(t, err)
	}

	f1, err := os.Create(dbfile1)
	if err != nil {
		TestFailErr(t, err)
	}
	err = f1.Close()
	if err != nil {
		TestFailErr(t, err)
	}
	err = os.Chmod(dbfile1, 0777)
	if err != nil {
		TestFailErr(t, err)
	}
	f2, err := os.Create(dbfile2)
	if err != nil {
		TestFailErr(t, err)
	}
	err = f2.Close()
	if err != nil {
		TestFailErr(t, err)
	}
	err = os.Chmod(dbfile2, 0777)
	if err != nil {
		TestFailErr(t, err)
	}
	f3, err := os.Create(dbfile3)
	if err != nil {
		TestFailErr(t, err)
	}
	err = f3.Close()
	if err != nil {
		TestFailErr(t, err)
	}
	err = os.Chmod(dbfile3, 0777)
	if err != nil {
		TestFailErr(t, err)
	}

	TPrintln(zerolog.InfoLevel, "DatabaseFile<main>:      "+dbfile1)
	TPrintln(zerolog.InfoLevel, "DatabaseFile<requests>:  "+dbfile2)
	TPrintln(zerolog.InfoLevel, "DatabaseFile<logs>:      "+dbfile3)

	scn.Conf = CreateTestConfig(t, dbfile1, dbfile2, dbfile3)

	stop := func() {
		_ = os.Remove(dbfile1)
		_ = os.Remove(dbfile2)
		_ = os.Remove(dbfile3)
	}

	return dbfile1, dbfile2, dbfile3, scn.Conf, stop
}

func CreateTestConfig(t *testing.T, dbfile1 string, dbfile2 string, dbfile3 string) scn.Config {
	conf, ok := scn.GetConfig("local-host")
	if !ok {
		TestFail(t, "conf not found")
	}

	conf.ServerPort = "0" // simply choose a free port
	conf.DBMain.File = dbfile1
	conf.DBLogs.File = dbfile2
	conf.DBRequests.File = dbfile3
	conf.DBMain.Timeout = 500 * time.Millisecond
	conf.DBLogs.Timeout = 500 * time.Millisecond
	conf.DBRequests.Timeout = 500 * time.Millisecond
	conf.DBMain.ConnMaxLifetime = 1 * time.Second
	conf.DBLogs.ConnMaxLifetime = 1 * time.Second
	conf.DBRequests.ConnMaxLifetime = 1 * time.Second
	conf.DBMain.ConnMaxIdleTime = 1 * time.Second
	conf.DBLogs.ConnMaxIdleTime = 1 * time.Second
	conf.DBRequests.ConnMaxIdleTime = 1 * time.Second
	conf.RequestMaxRetry = 32
	conf.RequestRetrySleep = 100 * time.Millisecond
	conf.ReturnRawErrors = true
	conf.DummyFirebase = true

	return conf
}
