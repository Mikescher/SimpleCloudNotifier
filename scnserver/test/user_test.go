package test

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	tt "blackforestbytes.com/simplecloudnotifier/test/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"testing"
)

func TestCreateUserNoClient(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"no_client": true,
	})

	tt.AssertEqual(t, "len(clients)", 0, len(r0["clients"].([]any)))

	uid := fmt.Sprintf("%v", r0["user_id"])
	admintok := r0["admin_key"].(string)
	readtok := r0["read_key"].(string)
	sendtok := r0["send_key"].(string)

	tt.RequestAuthGetShouldFail(t, sendtok, baseUrl, "/api/users/"+uid, 401, apierr.USER_AUTH_FAILED)
	tt.RequestAuthGetShouldFail(t, "", baseUrl, "/api/users/"+uid, 401, apierr.USER_AUTH_FAILED)

	r1 := tt.RequestAuthGet[gin.H](t, readtok, baseUrl, "/api/users/"+uid)

	tt.AssertEqual(t, "uid", uid, fmt.Sprintf("%v", r1["user_id"]))
	tt.AssertEqual(t, "admin_key", admintok, r1["admin_key"])
}

func TestCreateUserDummyClient(t *testing.T) {
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
}

func TestCreateUserWithUsername(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
		"username":      "my_user",
	})

	tt.AssertEqual(t, "len(clients)", 1, len(r0["clients"].([]any)))

	uid := fmt.Sprintf("%v", r0["user_id"])

	admintok := r0["admin_key"].(string)

	r1 := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/users/"+uid)

	tt.AssertEqual(t, "uid", uid, fmt.Sprintf("%v", r1["user_id"]))
	tt.AssertEqual(t, "admin_key", admintok, r1["admin_key"])
	tt.AssertEqual(t, "username", "my_user", r1["username"])
}

func TestUpdateUsername(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})
	tt.AssertEqual(t, "username", nil, r0["username"])

	uid := fmt.Sprintf("%v", r0["user_id"])
	admintok := r0["admin_key"].(string)

	r1 := tt.RequestAuthPatch[gin.H](t, admintok, baseUrl, "/api/users/"+uid, gin.H{"username": "my_user_001"})
	tt.AssertEqual(t, "username", "my_user_001", r1["username"])

	r2 := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/users/"+uid)
	tt.AssertEqual(t, "username", "my_user_001", r2["username"])

	r3 := tt.RequestAuthPatch[gin.H](t, admintok, baseUrl, "/api/users/"+uid, gin.H{"username": "my_user_002"})
	tt.AssertEqual(t, "username", "my_user_002", r3["username"])

	r4 := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/users/"+uid)
	tt.AssertEqual(t, "username", "my_user_002", r4["username"])

	r5 := tt.RequestAuthPatch[gin.H](t, admintok, baseUrl, "/api/users/"+uid, gin.H{"username": ""})
	tt.AssertEqual(t, "username", nil, r5["username"])

	r6 := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/users/"+uid)
	tt.AssertEqual(t, "username", nil, r6["username"])
}

func TestRecreateKeys(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})
	tt.AssertEqual(t, "username", nil, r0["username"])

	uid := fmt.Sprintf("%v", r0["user_id"])

	admintok := r0["admin_key"].(string)
	readtok := r0["read_key"].(string)
	sendtok := r0["send_key"].(string)

	tt.RequestAuthPatchShouldFail(t, readtok, baseUrl, "/api/users/"+uid, gin.H{"read_key": true}, 401, apierr.USER_AUTH_FAILED)

	tt.RequestAuthPatchShouldFail(t, sendtok, baseUrl, "/api/users/"+uid, gin.H{"read_key": true}, 401, apierr.USER_AUTH_FAILED)

	r1 := tt.RequestAuthPatch[gin.H](t, admintok, baseUrl, "/api/users/"+uid, gin.H{})
	tt.AssertEqual(t, "admin_key", admintok, r1["admin_key"])
	tt.AssertEqual(t, "read_key", readtok, r1["read_key"])
	tt.AssertEqual(t, "send_key", sendtok, r1["send_key"])

	r2 := tt.RequestAuthPatch[gin.H](t, admintok, baseUrl, "/api/users/"+uid, gin.H{"read_key": true})
	tt.AssertEqual(t, "admin_key", admintok, r2["admin_key"])
	tt.AssertNotEqual(t, "read_key", readtok, r2["read_key"])
	tt.AssertEqual(t, "send_key", sendtok, r2["send_key"])
	readtok = r2["read_key"].(string)

	r3 := tt.RequestAuthPatch[gin.H](t, admintok, baseUrl, "/api/users/"+uid, gin.H{"read_key": true, "send_key": true})
	tt.AssertEqual(t, "admin_key", admintok, r3["admin_key"])
	tt.AssertNotEqual(t, "read_key", readtok, r3["read_key"])
	tt.AssertNotEqual(t, "send_key", sendtok, r3["send_key"])
	readtok = r3["read_key"].(string)
	sendtok = r3["send_key"].(string)

	r4 := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/users/"+uid)
	tt.AssertEqual(t, "admin_key", admintok, r4["admin_key"])
	tt.AssertEqual(t, "read_key", readtok, r4["read_key"])
	tt.AssertEqual(t, "send_key", sendtok, r4["send_key"])

	r5 := tt.RequestAuthPatch[gin.H](t, admintok, baseUrl, "/api/users/"+uid, gin.H{"admin_key": true})
	tt.AssertNotEqual(t, "admin_key", admintok, r5["admin_key"])
	tt.AssertEqual(t, "read_key", readtok, r5["read_key"])
	tt.AssertEqual(t, "send_key", sendtok, r5["send_key"])
	admintokNew := r5["admin_key"].(string)

	tt.RequestAuthGetShouldFail(t, admintok, baseUrl, "/api/users/"+uid, 401, apierr.USER_AUTH_FAILED)

	r6 := tt.RequestAuthGet[gin.H](t, admintokNew, baseUrl, "/api/users/"+uid)
	tt.AssertEqual(t, "admin_key", admintokNew, r6["admin_key"])
	tt.AssertEqual(t, "read_key", readtok, r6["read_key"])
	tt.AssertEqual(t, "send_key", sendtok, r6["send_key"])
}

func TestDeleteUser(t *testing.T) {
	t.SkipNow() // TODO DeleteUser Not implemented

	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := fmt.Sprintf("%v", r0["user_id"])
	admintok := r0["admin_key"].(string)

	tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/users/"+uid)

	tt.RequestAuthDeleteShouldFail(t, admintok, baseUrl, "/api/users/"+uid, nil, 401, apierr.USER_AUTH_FAILED)

	tt.RequestAuthDelete[tt.Void](t, admintok, baseUrl, "/api/users/"+uid, nil)

	tt.RequestAuthGetShouldFail(t, admintok, baseUrl, "/api/users/"+uid, 404, apierr.USER_NOT_FOUND)

}

func TestCreateProUser(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	{
		r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
			"no_client": true,
		})

		tt.AssertEqual(t, "is_pro", false, r0["is_pro"])
	}

	{
		r1 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
			"no_client": true,
			"pro_token": "ANDROID|v2|PURCHASED:000",
		})

		tt.AssertEqual(t, "is_pro", true, r1["is_pro"])
	}

	{
		r2 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
			"agent_model":   "DUMMY_PHONE",
			"agent_version": "4X",
			"client_type":   "ANDROID",
			"fcm_token":     "DUMMY_FCM",
		})

		tt.AssertEqual(t, "is_pro", false, r2["is_pro"])
	}

	{
		r3 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
			"agent_model":   "DUMMY_PHONE",
			"agent_version": "4X",
			"client_type":   "ANDROID",
			"fcm_token":     "DUMMY_FCM",
			"pro_token":     "ANDROID|v2|PURCHASED:000",
		})

		tt.AssertEqual(t, "is_pro", true, r3["is_pro"])
	}

}
