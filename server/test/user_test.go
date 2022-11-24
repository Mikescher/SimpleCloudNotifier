package test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"testing"
	"time"
)

func TestCreateUserNoClient(t *testing.T) {
	ws, stop := NewSimpleWebserver()
	defer stop()
	go func() { ws.Run() }()
	time.Sleep(100 * time.Millisecond)

	baseUrl := "http://127.0.0.1:" + ws.Port

	res := requestPost[gin.H](baseUrl, "/api/users", gin.H{
		"no_client": true,
	})

	uid := fmt.Sprintf("%v", res["user_id"])
	admintok := res["admin_key"].(string)

	fmt.Printf("uid       := %s\n", uid)
	fmt.Printf("admin_key := %s\n", admintok)

	requestAuthGet[Void](admintok, baseUrl, "/api/users/"+uid)
}
