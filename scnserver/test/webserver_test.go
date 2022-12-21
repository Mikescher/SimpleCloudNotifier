package test

import (
	tt "blackforestbytes.com/simplecloudnotifier/test/util"
	"fmt"
	"testing"
)

func TestWebserver(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	fmt.Printf("URL       := %s\n", baseUrl)
}

func TestPing(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	_ = tt.RequestGet[tt.Void](t, baseUrl, "/api/ping")
	_ = tt.RequestPut[tt.Void](t, baseUrl, "/api/ping", nil)
	_ = tt.RequestPost[tt.Void](t, baseUrl, "/api/ping", nil)
	_ = tt.RequestPatch[tt.Void](t, baseUrl, "/api/ping", nil)
	_ = tt.RequestDelete[tt.Void](t, baseUrl, "/api/ping", nil)
}

func TestMongo(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	_ = tt.RequestPost[tt.Void](t, baseUrl, "/api/db-test", nil)
}

func TestHealth(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	_ = tt.RequestGet[tt.Void](t, baseUrl, "/api/health")
}
