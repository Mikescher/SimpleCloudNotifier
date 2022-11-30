package test

import (
	"fmt"
	"testing"
)

func TestWebserver(t *testing.T) {
	ws, stop := StartSimpleWebserver(t)
	defer stop()

	fmt.Printf("Port       := %s\n", ws.Port)
}

func TestPing(t *testing.T) {
	ws, stop := StartSimpleWebserver(t)
	defer stop()

	baseUrl := "http://127.0.0.1:" + ws.Port

	_ = requestGet[Void](t, baseUrl, "/api/ping")
	_ = requestPut[Void](t, baseUrl, "/api/ping", nil)
	_ = requestPost[Void](t, baseUrl, "/api/ping", nil)
	_ = requestPatch[Void](t, baseUrl, "/api/ping", nil)
	_ = requestDelete[Void](t, baseUrl, "/api/ping", nil)
}

func TestMongo(t *testing.T) {
	ws, stop := StartSimpleWebserver(t)
	defer stop()

	baseUrl := "http://127.0.0.1:" + ws.Port

	_ = requestPost[Void](t, baseUrl, "/api/db-test", nil)
}

func TestHealth(t *testing.T) {
	ws, stop := StartSimpleWebserver(t)
	defer stop()

	baseUrl := "http://127.0.0.1:" + ws.Port

	_ = requestGet[Void](t, baseUrl, "/api/health")
}
