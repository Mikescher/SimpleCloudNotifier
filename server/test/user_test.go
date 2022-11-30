package test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"testing"
)

func TestCreateUserNoClient(t *testing.T) {
	ws, stop := StartSimpleWebserver(t)
	defer stop()

	baseUrl := "http://127.0.0.1:" + ws.Port

	r0 := requestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"no_client": true,
	})

	assertEqual(t, "len(clients)", 0, len(r0["clients"].([]any)))

	uid := fmt.Sprintf("%v", r0["user_id"])
	admintok := r0["admin_key"].(string)

	fmt.Printf("uid       := %s\n", uid)
	fmt.Printf("admin_key := %s\n", admintok)

	r1 := requestAuthGet[gin.H](t, admintok, baseUrl, "/api/users/"+uid)

	assertEqual(t, "uid", uid, fmt.Sprintf("%v", r1["user_id"]))
	assertEqual(t, "admin_key", admintok, r1["admin_key"])
}

func TestCreateUserDummyClient(t *testing.T) {
	ws, stop := StartSimpleWebserver(t)
	defer stop()

	baseUrl := "http://127.0.0.1:" + ws.Port

	r0 := requestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := fmt.Sprintf("%v", r0["user_id"])

	assertEqual(t, "len(clients)", 1, len(r0["clients"].([]any)))

	admintok := r0["admin_key"].(string)

	fmt.Printf("uid       := %s\n", uid)
	fmt.Printf("admin_key := %s\n", admintok)

	r1 := requestAuthGet[gin.H](t, admintok, baseUrl, "/api/users/"+uid)

	assertEqual(t, "uid", uid, fmt.Sprintf("%v", r1["user_id"]))
	assertEqual(t, "admin_key", admintok, r1["admin_key"])
	assertEqual(t, "username", nil, r1["username"])

	type rt2 struct {
		Clients []gin.H `json:"clients"`
	}

	r2 := requestAuthGet[rt2](t, admintok, baseUrl, "/api/users/"+uid+"/clients")

	assertEqual(t, "len(clients)", 1, len(r2.Clients))

	c0 := r2.Clients[0]

	assertEqual(t, "agent_model", "DUMMY_PHONE", c0["agent_model"])
	assertEqual(t, "agent_version", "4X", c0["agent_version"])
	assertEqual(t, "fcm_token", "DUMMY_FCM", c0["fcm_token"])
	assertEqual(t, "client_type", "ANDROID", c0["type"])
}

func TestCreateUserWithUsername(t *testing.T) {
	ws, stop := StartSimpleWebserver(t)
	defer stop()

	baseUrl := "http://127.0.0.1:" + ws.Port

	r0 := requestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
		"username":      "my_user",
	})

	assertEqual(t, "len(clients)", 1, len(r0["clients"].([]any)))

	uid := fmt.Sprintf("%v", r0["user_id"])

	admintok := r0["admin_key"].(string)

	fmt.Printf("uid       := %s\n", uid)
	fmt.Printf("admin_key := %s\n", admintok)

	r1 := requestAuthGet[gin.H](t, admintok, baseUrl, "/api/users/"+uid)

	assertEqual(t, "uid", uid, fmt.Sprintf("%v", r1["user_id"]))
	assertEqual(t, "admin_key", admintok, r1["admin_key"])
	assertEqual(t, "username", "my_user", r1["username"])
}
