package test

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/api"
	"blackforestbytes.com/simplecloudnotifier/common/ginext"
	"blackforestbytes.com/simplecloudnotifier/db"
	"blackforestbytes.com/simplecloudnotifier/jobs"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"blackforestbytes.com/simplecloudnotifier/push"
	"bytes"
	"encoding/json"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func NewSimpleWebserver(t *testing.T) *logic.Application {

	uuid, err := langext.NewHexUUID()
	if err != nil {
		panic(err)
	}

	dbfile := filepath.Join(os.TempDir(), uuid+"sqlite3")
	defer func() {
		_ = os.Remove(dbfile)
	}()

	conf := scn.Config{
		Namespace:       "test",
		GinDebug:        true,
		ServerIP:        "0.0.0.0",
		ServerPort:      "0", // simply choose a free port
		DBFile:          dbfile,
		RequestTimeout:  500 * time.Millisecond,
		ReturnRawErrors: true,
		DummyFirebase:   true,
	}

	sqlite, err := db.NewDatabase(dbfile)
	if err != nil {
		panic(err)
	}

	app := logic.NewApp(sqlite)

	if err := app.Migrate(); err != nil {
		panic(err)
	}

	ginengine := ginext.NewEngine(conf)

	router := api.NewRouter(app)

	nc := push.NewTestSink()

	jobRetry := jobs.NewDeliveryRetryJob(app)
	app.Init(conf, ginengine, nc, []logic.Job{jobRetry})

	router.Init(ginengine)

	return app
}

func requestGet[T any](t *testing.T, baseURL string, prefix string) T {
	client := http.Client{}

	req, err := http.NewRequest("GET", baseURL+prefix, bytes.NewReader([]byte{}))
	if err != nil {
		t.Error(err)
		return *new(T)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
		return *new(T)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		t.Error("Statuscode != 200")
	}

	respBodyBin, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return *new(T)
	}

	var data T
	if err := json.Unmarshal(respBodyBin, &data); err != nil {
		t.Error(err)
		return *new(T)
	}

	return data
}
