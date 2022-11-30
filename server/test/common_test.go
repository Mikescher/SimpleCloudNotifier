package test

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/api"
	"blackforestbytes.com/simplecloudnotifier/common/ginext"
	"blackforestbytes.com/simplecloudnotifier/db"
	"blackforestbytes.com/simplecloudnotifier/google"
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
	"runtime/debug"
	"strings"
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
		testFailErr(t, err)
	}

	f, err := os.Create(dbfile)
	if err != nil {
		testFailErr(t, err)
	}
	err = f.Close()
	if err != nil {
		testFailErr(t, err)
	}

	err = os.Chmod(dbfile, 0777)
	if err != nil {
		testFailErr(t, err)
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
		testFailErr(t, err)
	}

	app := logic.NewApp(sqlite)

	if err := app.Migrate(); err != nil {
		testFailErr(t, err)
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
	time.Sleep(100 * time.Millisecond)
	return app, stop
}

func requestGet[TResult any](t *testing.T, baseURL string, prefix string) TResult {
	return requestAny[TResult](t, "", "GET", baseURL, prefix, nil)
}

func requestAuthGet[TResult any](t *testing.T, akey string, baseURL string, prefix string) TResult {
	return requestAny[TResult](t, akey, "GET", baseURL, prefix, nil)
}

func requestPost[TResult any](t *testing.T, baseURL string, prefix string, body any) TResult {
	return requestAny[TResult](t, "", "POST", baseURL, prefix, body)
}

func requestAuthPost[TResult any](t *testing.T, akey string, baseURL string, prefix string, body any) TResult {
	return requestAny[TResult](t, akey, "POST", baseURL, prefix, body)
}

func requestPut[TResult any](t *testing.T, baseURL string, prefix string, body any) TResult {
	return requestAny[TResult](t, "", "PUT", baseURL, prefix, body)
}

func requestAuthPUT[TResult any](t *testing.T, akey string, baseURL string, prefix string, body any) TResult {
	return requestAny[TResult](t, akey, "PUT", baseURL, prefix, body)
}

func requestPatch[TResult any](t *testing.T, baseURL string, prefix string, body any) TResult {
	return requestAny[TResult](t, "", "PATCH", baseURL, prefix, body)
}

func requestAuthPatch[TResult any](t *testing.T, akey string, baseURL string, prefix string, body any) TResult {
	return requestAny[TResult](t, akey, "PATCH", baseURL, prefix, body)
}

func requestDelete[TResult any](t *testing.T, baseURL string, prefix string, body any) TResult {
	return requestAny[TResult](t, "", "DELETE", baseURL, prefix, body)
}

func requestAuthDelete[TResult any](t *testing.T, akey string, baseURL string, prefix string, body any) TResult {
	return requestAny[TResult](t, akey, "DELETE", baseURL, prefix, body)
}

func requestAny[TResult any](t *testing.T, akey string, method string, baseURL string, prefix string, body any) TResult {
	client := http.Client{}

	bytesbody := make([]byte, 0)
	if body != nil {
		bjson, err := json.Marshal(body)
		if err != nil {
			testFailErr(t, err)
		}
		bytesbody = bjson
	}

	req, err := http.NewRequest(method, baseURL+prefix, bytes.NewReader(bytesbody))
	if err != nil {
		testFailErr(t, err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if akey != "" {
		req.Header.Set("Authorization", "SCN "+akey)
	}

	resp, err := client.Do(req)
	if err != nil {
		testFailErr(t, err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBodyBin, err := io.ReadAll(resp.Body)
	if err != nil {
		testFailErr(t, err)
	}

	if resp.StatusCode != 200 {
		fmt.Println("Request: " + method + " :: " + baseURL + prefix)
		fmt.Println(string(respBodyBin))
		testFail(t, "Statuscode != 200")
	}

	var data TResult
	if err := json.Unmarshal(respBodyBin, &data); err != nil {
		testFailErr(t, err)
	}

	return data
}

func assertEqual(t *testing.T, key string, expected any, actual any) {
	if expected != actual {
		t.Errorf("Value [%s] differs (%T <-> %T):\n", key, expected, actual)

		str1 := fmt.Sprintf("%v", expected)
		str2 := fmt.Sprintf("%v", actual)

		if strings.Contains(str1, "\n") {
			t.Errorf("Actual  :\n~~~~~~~~~~~~~~~~\n%v\n~~~~~~~~~~~~~~~~\n\n", expected)
		} else {
			t.Errorf("Actual  : \"%v\"\n", expected)
		}

		if strings.Contains(str2, "\n") {
			t.Errorf("Expected  :\n~~~~~~~~~~~~~~~~\n%v\n~~~~~~~~~~~~~~~~\n\n", actual)
		} else {
			t.Errorf("Expected  : \"%v\"\n", actual)
		}

		t.FailNow()
	}
}

func testFail(t *testing.T, msg string) {
	t.Error(msg)
	t.FailNow()
}

func testFailFmt(t *testing.T, format string, args ...any) {
	t.Errorf(format, args...)
	t.FailNow()
}

func testFailErr(t *testing.T, e error) {
	t.Error(fmt.Sprintf("Failed with error:\n%s\n\nError:\n%+v\n\nTrace:\n%s", e.Error(), e, string(debug.Stack())))
	t.FailNow()
}
