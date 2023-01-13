package test

import (
	tt "blackforestbytes.com/simplecloudnotifier/test/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"testing"
)

func TestGetClient(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := fmt.Sprintf("%v", r0["user_id"])

	tt.AssertEqual(t, "len(clients)", 1, len(r0["clients"].([]any)))

	admintok := r0["admin_key"].(string)

	fmt.Printf("uid       := %s\n", uid)
	fmt.Printf("admin_key := %s\n", admintok)

	r1 := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/users/"+uid)

	tt.AssertEqual(t, "uid", uid, fmt.Sprintf("%v", r1["user_id"]))
	tt.AssertEqual(t, "admin_key", admintok, r1["admin_key"])
	tt.AssertEqual(t, "username", nil, r1["username"])

	type rt2 struct {
		Clients []gin.H `json:"clients"`
	}

	r2 := tt.RequestAuthGet[rt2](t, admintok, baseUrl, "/api/users/"+uid+"/clients")

	tt.AssertEqual(t, "len(clients)", 1, len(r2.Clients))

	c0 := r2.Clients[0]

	tt.AssertEqual(t, "agent_model", "DUMMY_PHONE", c0["agent_model"])
	tt.AssertEqual(t, "agent_version", "4X", c0["agent_version"])
	tt.AssertEqual(t, "fcm_token", "DUMMY_FCM", c0["fcm_token"])
	tt.AssertEqual(t, "client_type", "ANDROID", c0["type"])
	tt.AssertEqual(t, "user_id", uid, fmt.Sprintf("%v", c0["user_id"]))

	cid := fmt.Sprintf("%v", c0["client_id"])

	r3 := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/users/"+uid+"/clients/"+cid)

	tt.AssertJsonMapEqual(t, "client", r3, c0)
}

func TestCreateAndDeleteClient(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := fmt.Sprintf("%v", r0["user_id"])

	tt.AssertEqual(t, "len(clients)", 1, len(r0["clients"].([]any)))

	admintok := r0["admin_key"].(string)

	fmt.Printf("uid       := %s\n", uid)
	fmt.Printf("admin_key := %s\n", admintok)

	r2 := tt.RequestAuthPost[gin.H](t, admintok, baseUrl, "/api/users/"+uid+"/clients", gin.H{
		"agent_model":   "DUMMY_PHONE_2",
		"agent_version": "99X",
		"client_type":   "IOS",
		"fcm_token":     "DUMMY_FCM_2",
	})

	cid2 := fmt.Sprintf("%v", r2["client_id"])

	type rt3 struct {
		Clients []gin.H `json:"clients"`
	}

	r3 := tt.RequestAuthGet[rt3](t, admintok, baseUrl, "/api/users/"+uid+"/clients")
	tt.AssertEqual(t, "len(clients)", 2, len(r3.Clients))

	r4 := tt.RequestAuthDelete[gin.H](t, admintok, baseUrl, "/api/users/"+uid+"/clients/"+cid2, nil)
	tt.AssertEqual(t, "client_id", cid2, fmt.Sprintf("%v", r4["client_id"]))

	r5 := tt.RequestAuthGet[rt3](t, admintok, baseUrl, "/api/users/"+uid+"/clients")
	tt.AssertEqual(t, "len(clients)", 1, len(r5.Clients))
}

func TestReuseFCM(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM_001",
	})

	uid := fmt.Sprintf("%v", r0["user_id"])

	tt.AssertEqual(t, "len(clients)", 1, len(r0["clients"].([]any)))

	admintok := r0["admin_key"].(string)

	fmt.Printf("uid       := %s\n", uid)
	fmt.Printf("admin_key := %s\n", admintok)

	type rt2 struct {
		Clients []gin.H `json:"clients"`
	}

	r1 := tt.RequestAuthGet[rt2](t, admintok, baseUrl, "/api/users/"+uid+"/clients")

	tt.AssertEqual(t, "len(clients)", 1, len(r1.Clients))

	r2 := tt.RequestAuthPost[gin.H](t, admintok, baseUrl, "/api/users/"+uid+"/clients", gin.H{
		"agent_model":   "DUMMY_PHONE_2",
		"agent_version": "99X",
		"client_type":   "IOS",
		"fcm_token":     "DUMMY_FCM_001",
	})

	cid2 := fmt.Sprintf("%v", r2["client_id"])

	type rt3 struct {
		Clients []gin.H `json:"clients"`
	}

	r3 := tt.RequestAuthGet[rt3](t, admintok, baseUrl, "/api/users/"+uid+"/clients")
	tt.AssertEqual(t, "len(clients)", 1, len(r3.Clients))

	tt.AssertEqual(t, "clients->client_id", cid2, fmt.Sprintf("%v", r3.Clients[0]["client_id"]))
}

func TestListClients(t *testing.T) {
	t.SkipNow() //TODO
}
