package test

import (
	"testing"
	"time"
)

func TestWebserver(t *testing.T) {
	ws := NewSimpleWebserver(t)
	defer ws.Stop()
	go func() { ws.Run() }()
	time.Sleep(100 * time.Millisecond)
}

func TestPing(t *testing.T) {
	ws := NewSimpleWebserver(t)
	defer ws.Stop()
	go func() { ws.Run() }()
	time.Sleep(100 * time.Millisecond)

	baseUrl := "http://127.0.0.1:" + ws.Port

	_ = requestGet[struct{}](t, baseUrl, "/api/ping")

}
