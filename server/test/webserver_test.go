package test

import (
	"testing"
	"time"
)

func TestWebserver(t *testing.T) {
	ws, stop := NewSimpleWebserver()
	defer stop()
	go func() { ws.Run() }()
	time.Sleep(100 * time.Millisecond)
}

func TestPing(t *testing.T) {
	ws, stop := NewSimpleWebserver()
	defer stop()
	go func() { ws.Run() }()
	time.Sleep(100 * time.Millisecond)

	baseUrl := "http://127.0.0.1:" + ws.Port

	_ = requestGet[Void](baseUrl, "/api/ping")
	_ = requestPut[Void](baseUrl, "/api/ping", nil)
	_ = requestPost[Void](baseUrl, "/api/ping", nil)
	_ = requestPatch[Void](baseUrl, "/api/ping", nil)
	_ = requestDelete[Void](baseUrl, "/api/ping", nil)
}

func TestMongo(t *testing.T) {
	ws, stop := NewSimpleWebserver()
	defer stop()
	go func() { ws.Run() }()
	time.Sleep(100 * time.Millisecond)

	baseUrl := "http://127.0.0.1:" + ws.Port

	_ = requestPost[Void](baseUrl, "/api/db-test", nil)
}

func TestHealth(t *testing.T) {
	ws, stop := NewSimpleWebserver()
	defer stop()
	go func() { ws.Run() }()
	time.Sleep(100 * time.Millisecond)

	baseUrl := "http://127.0.0.1:" + ws.Port

	_ = requestGet[Void](baseUrl, "/api/health")
}
