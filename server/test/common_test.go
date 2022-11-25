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
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Void = struct{}

func NewSimpleWebserver() (*logic.Application, func()) {
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

	uuid1, _ := langext.NewHexUUID()
	uuid2, _ := langext.NewHexUUID()
	dbdir := filepath.Join(os.TempDir(), uuid1)
	dbfile := filepath.Join(dbdir, uuid2+".sqlite3")

	err := os.MkdirAll(dbdir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(dbfile)
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}

	err = os.Chmod(dbfile, 0777)
	if err != nil {
		panic(err)
	}

	fmt.Println("DatabaseFile: " + dbfile)

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

	apc := google.NewDummy()

	jobRetry := jobs.NewDeliveryRetryJob(app)
	app.Init(conf, ginengine, nc, apc, []logic.Job{jobRetry})

	router.Init(ginengine)

	return app, func() { app.Stop(); _ = os.Remove(dbfile) }
}

func requestGet[TResult any](baseURL string, prefix string) TResult {
	return requestAny[TResult]("", "GET", baseURL, prefix, nil)
}

func requestAuthGet[TResult any](akey string, baseURL string, prefix string) TResult {
	return requestAny[TResult](akey, "GET", baseURL, prefix, nil)
}

func requestPost[TResult any](baseURL string, prefix string, body any) TResult {
	return requestAny[TResult]("", "POST", baseURL, prefix, body)
}

func requestAuthPost[TResult any](akey string, baseURL string, prefix string, body any) TResult {
	return requestAny[TResult](akey, "POST", baseURL, prefix, body)
}

func requestPut[TResult any](baseURL string, prefix string, body any) TResult {
	return requestAny[TResult]("", "PUT", baseURL, prefix, body)
}

func requestAuthPUT[TResult any](akey string, baseURL string, prefix string, body any) TResult {
	return requestAny[TResult](akey, "PUT", baseURL, prefix, body)
}

func requestPatch[TResult any](baseURL string, prefix string, body any) TResult {
	return requestAny[TResult]("", "PATCH", baseURL, prefix, body)
}

func requestAuthPatch[TResult any](akey string, baseURL string, prefix string, body any) TResult {
	return requestAny[TResult](akey, "PATCH", baseURL, prefix, body)
}

func requestDelete[TResult any](baseURL string, prefix string, body any) TResult {
	return requestAny[TResult]("", "DELETE", baseURL, prefix, body)
}

func requestAuthDelete[TResult any](akey string, baseURL string, prefix string, body any) TResult {
	return requestAny[TResult](akey, "DELETE", baseURL, prefix, body)
}

func requestAny[TResult any](akey string, method string, baseURL string, prefix string, body any) TResult {
	client := http.Client{}

	bytesbody := make([]byte, 0)
	if body != nil {
		bjson, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}
		bytesbody = bjson
	}

	req, err := http.NewRequest(method, baseURL+prefix, bytes.NewReader(bytesbody))
	if err != nil {
		panic(err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if akey != "" {
		req.Header.Set("Authorization", "SCN "+akey)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBodyBin, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		fmt.Println("Request: " + method + " :: " + baseURL + prefix)
		fmt.Println(string(respBodyBin))
		panic("Statuscode != 200")
	}

	var data TResult
	if err := json.Unmarshal(respBodyBin, &data); err != nil {
		panic(err)
	}

	return data
}
