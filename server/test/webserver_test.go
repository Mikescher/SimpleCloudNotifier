package test

import (
	tt "blackforestbytes.com/simplecloudnotifier/test/util"
	"fmt"
	"testing"
)

func TestWebserver(t *testing.T) {
	ws, stop := tt.StartSimpleWebserver(t)
	defer stop()

	fmt.Printf("Port       := %s\n", ws.Port)
}

func TestPing(t *testing.T) {
	ws, stop := tt.StartSimpleWebserver(t)
	defer stop()

	baseUrl := "http://127.0.0.1:" + ws.Port

	_ = tt.RequestGet[tt.Void](t, baseUrl, "/api/ping")
	_ = tt.RequestPut[tt.Void](t, baseUrl, "/api/ping", nil)
	_ = tt.RequestPost[tt.Void](t, baseUrl, "/api/ping", nil)
	_ = tt.RequestPatch[tt.Void](t, baseUrl, "/api/ping", nil)
	_ = tt.RequestDelete[tt.Void](t, baseUrl, "/api/ping", nil)
}

func TestMongo(t *testing.T) {
	ws, stop := tt.StartSimpleWebserver(t)
	defer stop()

	baseUrl := "http://127.0.0.1:" + ws.Port

	_ = tt.RequestPost[tt.Void](t, baseUrl, "/api/db-test", nil)
}

func TestHealth(t *testing.T) {
	ws, stop := tt.StartSimpleWebserver(t)
	defer stop()

	baseUrl := "http://127.0.0.1:" + ws.Port

	_ = tt.RequestGet[tt.Void](t, baseUrl, "/api/health")
}
