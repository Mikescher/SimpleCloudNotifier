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

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/v2/users", gin.H{
		"no_client": true,
	})

	tt.AssertEqual(t, "len(clients)", 0, len(r0["clients"].([]any)))

	uid := fmt.Sprintf("%v", r0["user_id"])
	readtok := r0["read_key"].(string)
	sendtok := r0["send_key"].(string)

	tt.RequestAuthGetShouldFail(t, sendtok, baseUrl, "/api/v2/users/"+uid, 401, apierr.USER_AUTH_FAILED)
	tt.RequestAuthGetShouldFail(t, "", baseUrl, "/api/v2/users/"+uid, 401, apierr.USER_AUTH_FAILED)

	r1 := tt.RequestAuthGet[gin.H](t, readtok, baseUrl, "/api/v2/users/"+uid)

	tt.AssertEqual(t, "uid", uid, fmt.Sprintf("%v", r1["user_id"]))
}

func TestCreateUserDummyClient(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/v2/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := fmt.Sprintf("%v", r0["user_id"])

	tt.AssertEqual(t, "len(clients)", 1, len(r0["clients"].([]any)))

	admintok := r0["admin_key"].(string)

	r1 := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/v2/users/"+uid)

	tt.AssertEqual(t, "uid", uid, fmt.Sprintf("%v", r1["user_id"]))
	tt.AssertEqual(t, "username", nil, r1["username"])

	type rt2 struct {
		Clients []gin.H `json:"clients"`
	}

	r2 := tt.RequestAuthGet[rt2](t, admintok, baseUrl, "/api/v2/users/"+uid+"/clients")

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

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/v2/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
		"username":      "my_user",
	})

	tt.AssertEqual(t, "len(clients)", 1, len(r0["clients"].([]any)))

	uid := fmt.Sprintf("%v", r0["user_id"])

	admintok := r0["admin_key"].(string)

	r1 := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/v2/users/"+uid)

	tt.AssertEqual(t, "uid", uid, fmt.Sprintf("%v", r1["user_id"]))
	tt.AssertEqual(t, "username", "my_user", r1["username"])
}

func TestUpdateUsername(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/v2/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})
	tt.AssertEqual(t, "username", nil, r0["username"])

	uid := fmt.Sprintf("%v", r0["user_id"])
	admintok := r0["admin_key"].(string)

	r1 := tt.RequestAuthPatch[gin.H](t, admintok, baseUrl, "/api/v2/users/"+uid, gin.H{"username": "my_user_001"})
	tt.AssertEqual(t, "username", "my_user_001", r1["username"])

	r2 := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/v2/users/"+uid)
	tt.AssertEqual(t, "username", "my_user_001", r2["username"])

	r3 := tt.RequestAuthPatch[gin.H](t, admintok, baseUrl, "/api/v2/users/"+uid, gin.H{"username": "my_user_002"})
	tt.AssertEqual(t, "username", "my_user_002", r3["username"])

	r4 := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/v2/users/"+uid)
	tt.AssertEqual(t, "username", "my_user_002", r4["username"])

	r5 := tt.RequestAuthPatch[gin.H](t, admintok, baseUrl, "/api/v2/users/"+uid, gin.H{"username": ""})
	tt.AssertEqual(t, "username", nil, r5["username"])

	r6 := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/v2/users/"+uid)
	tt.AssertEqual(t, "username", nil, r6["username"])
}

func TestUgradeUserToPro(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/v2/users", gin.H{
		"no_client": true,
	})
	tt.AssertEqual(t, "is_pro", false, r0["is_pro"])

	uid0 := fmt.Sprintf("%v", r0["user_id"])
	admintok0 := r0["admin_key"].(string)

	r1 := tt.RequestAuthPatch[gin.H](t, admintok0, baseUrl, "/api/v2/users/"+uid0, gin.H{"pro_token": "ANDROID|v2|PURCHASED:000"})
	tt.AssertEqual(t, "is_pro", true, r1["is_pro"])

	r2 := tt.RequestAuthGet[gin.H](t, admintok0, baseUrl, "/api/v2/users/"+uid0)
	tt.AssertEqual(t, "is_pro", true, r2["is_pro"])
}

func TestDowngradeUserToNonPro(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/v2/users", gin.H{
		"no_client": true,
		"pro_token": "ANDROID|v2|PURCHASED:UNIQ_111",
	})
	tt.AssertEqual(t, "is_pro", true, r0["is_pro"])

	uid0 := fmt.Sprintf("%v", r0["user_id"])
	admintok0 := r0["admin_key"].(string)

	r1 := tt.RequestAuthPatch[gin.H](t, admintok0, baseUrl, "/api/v2/users/"+uid0, gin.H{"pro_token": ""})
	tt.AssertEqual(t, "is_pro", false, r1["is_pro"])

	r2 := tt.RequestAuthGet[gin.H](t, admintok0, baseUrl, "/api/v2/users/"+uid0)
	tt.AssertEqual(t, "is_pro", false, r2["is_pro"])
}

func TestFailedUgradeUserToPro(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/v2/users", gin.H{
		"no_client": true,
	})
	tt.AssertEqual(t, "is_pro", false, r0["is_pro"])

	uid0 := fmt.Sprintf("%v", r0["user_id"])
	admintok0 := r0["admin_key"].(string)

	tt.RequestAuthPatchShouldFail(t, admintok0, baseUrl, "/api/v2/users/"+uid0, gin.H{"pro_token": "ANDROID|v2|INVALID"}, 400, apierr.INVALID_PRO_TOKEN)

	tt.RequestAuthPatchShouldFail(t, admintok0, baseUrl, "/api/v2/users/"+uid0, gin.H{"pro_token": "ANDROID|v99|PURCHASED"}, 400, apierr.INVALID_PRO_TOKEN)

	tt.RequestAuthPatchShouldFail(t, admintok0, baseUrl, "/api/v2/users/"+uid0, gin.H{"pro_token": "@INVALID"}, 400, apierr.INVALID_PRO_TOKEN)
}

func TestDeleteUser(t *testing.T) {
	t.SkipNow() // TODO DeleteUser Not implemented

	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/v2/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := fmt.Sprintf("%v", r0["user_id"])
	admintok := r0["admin_key"].(string)

	tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/v2/users/"+uid)

	tt.RequestAuthDeleteShouldFail(t, admintok, baseUrl, "/api/v2/users/"+uid, nil, 401, apierr.USER_AUTH_FAILED)

	tt.RequestAuthDelete[tt.Void](t, admintok, baseUrl, "/api/v2/users/"+uid, nil)

	tt.RequestAuthGetShouldFail(t, admintok, baseUrl, "/api/v2/users/"+uid, 404, apierr.USER_NOT_FOUND)

}

func TestCreateProUser(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	{
		r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/v2/users", gin.H{
			"no_client": true,
		})

		tt.AssertEqual(t, "is_pro", false, r0["is_pro"])
	}

	{
		r1 := tt.RequestPost[gin.H](t, baseUrl, "/api/v2/users", gin.H{
			"no_client": true,
			"pro_token": "ANDROID|v2|PURCHASED:000",
		})

		tt.AssertEqual(t, "is_pro", true, r1["is_pro"])
	}

	{
		r2 := tt.RequestPost[gin.H](t, baseUrl, "/api/v2/users", gin.H{
			"agent_model":   "DUMMY_PHONE",
			"agent_version": "4X",
			"client_type":   "ANDROID",
			"fcm_token":     "DUMMY_FCM",
		})

		tt.AssertEqual(t, "is_pro", false, r2["is_pro"])
	}

	{
		r3 := tt.RequestPost[gin.H](t, baseUrl, "/api/v2/users", gin.H{
			"agent_model":   "DUMMY_PHONE",
			"agent_version": "4X",
			"client_type":   "ANDROID",
			"fcm_token":     "DUMMY_FCM",
			"pro_token":     "ANDROID|v2|PURCHASED:000",
		})

		tt.AssertEqual(t, "is_pro", true, r3["is_pro"])
	}

}

func TestFailToCreateProUser(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	tt.RequestPostShouldFail(t, baseUrl, "/api/v2/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
		"pro_token":     "ANDROID|v2|INVALID",
	}, 400, apierr.INVALID_PRO_TOKEN)

	tt.RequestPostShouldFail(t, baseUrl, "/api/v2/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
		"pro_token":     "_",
	}, 400, apierr.INVALID_PRO_TOKEN)

	tt.RequestPostShouldFail(t, baseUrl, "/api/v2/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
		"pro_token":     "ANDROID|v99|xxx",
	}, 400, apierr.INVALID_PRO_TOKEN)
}

func TestReuseProToken(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/v2/users", gin.H{
		"no_client": true,
	})
	tt.AssertEqual(t, "is_pro", false, r0["is_pro"])

	uid0 := fmt.Sprintf("%v", r0["user_id"])
	admintok0 := r0["admin_key"].(string)

	r1 := tt.RequestPost[gin.H](t, baseUrl, "/api/v2/users", gin.H{
		"no_client": true,
		"pro_token": "ANDROID|v2|PURCHASED:UNIQ_1",
	})
	tt.AssertEqual(t, "is_pro", true, r1["is_pro"])

	uid1 := fmt.Sprintf("%v", r1["user_id"])
	admintok1 := r1["admin_key"].(string)

	{
		rc := tt.RequestAuthGet[gin.H](t, admintok1, baseUrl, "/api/v2/users/"+uid1)
		tt.AssertEqual(t, "is_pro", true, rc["is_pro"])
	}

	r2 := tt.RequestPost[gin.H](t, baseUrl, "/api/v2/users", gin.H{
		"no_client": true,
		"pro_token": "ANDROID|v2|PURCHASED:UNIQ_1",
	})
	tt.AssertEqual(t, "is_pro", true, r2["is_pro"])

	uid2 := fmt.Sprintf("%v", r2["user_id"])
	admintok2 := r2["admin_key"].(string)

	{
		rc := tt.RequestAuthGet[gin.H](t, admintok0, baseUrl, "/api/v2/users/"+uid0)
		tt.AssertEqual(t, "is_pro", false, rc["is_pro"])
	}

	{
		rc := tt.RequestAuthGet[gin.H](t, admintok1, baseUrl, "/api/v2/users/"+uid1)
		tt.AssertEqual(t, "is_pro", false, rc["is_pro"])
	}

	{
		rc := tt.RequestAuthGet[gin.H](t, admintok2, baseUrl, "/api/v2/users/"+uid2)
		tt.AssertEqual(t, "is_pro", true, rc["is_pro"])
	}

	tt.RequestAuthPatch[gin.H](t, admintok0, baseUrl, "/api/v2/users/"+uid0, gin.H{"pro_token": "ANDROID|v2|PURCHASED:UNIQ_2"})

	{
		rc := tt.RequestAuthGet[gin.H](t, admintok0, baseUrl, "/api/v2/users/"+uid0)
		tt.AssertEqual(t, "is_pro", true, rc["is_pro"])
	}

	{
		rc := tt.RequestAuthGet[gin.H](t, admintok1, baseUrl, "/api/v2/users/"+uid1)
		tt.AssertEqual(t, "is_pro", false, rc["is_pro"])
	}

	{
		rc := tt.RequestAuthGet[gin.H](t, admintok2, baseUrl, "/api/v2/users/"+uid2)
		tt.AssertEqual(t, "is_pro", true, rc["is_pro"])
	}

	tt.RequestAuthPatch[gin.H](t, admintok0, baseUrl, "/api/v2/users/"+uid0, gin.H{"pro_token": "ANDROID|v2|PURCHASED:UNIQ_1"})

	{
		rc := tt.RequestAuthGet[gin.H](t, admintok0, baseUrl, "/api/v2/users/"+uid0)
		tt.AssertEqual(t, "is_pro", true, rc["is_pro"])
	}

	{
		rc := tt.RequestAuthGet[gin.H](t, admintok1, baseUrl, "/api/v2/users/"+uid1)
		tt.AssertEqual(t, "is_pro", false, rc["is_pro"])
	}

	{
		rc := tt.RequestAuthGet[gin.H](t, admintok2, baseUrl, "/api/v2/users/"+uid2)
		tt.AssertEqual(t, "is_pro", false, rc["is_pro"])
	}

}

func TestUserMessageCounter(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/v2/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := r0["user_id"].(string)
	admintok := r0["admin_key"].(string)

	assertCounter := func(c int) {
		r1 := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/v2/users/"+uid)
		tt.AssertStrRepEqual(t, "messages_sent", c, r1["messages_sent"])
		tt.AssertStrRepEqual(t, "quota_used", c, r1["quota_used"])
	}

	assertCounter(0)

	tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"key":     admintok,
		"user_id": uid,
		"title":   tt.ShortLipsum(1001, 1),
	})

	assertCounter(1)
	assertCounter(1)

	tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"key":     admintok,
		"user_id": uid,
		"title":   tt.ShortLipsum(1002, 1),
	})

	assertCounter(2)

	tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"key":     admintok,
		"user_id": uid,
		"title":   tt.ShortLipsum(1003, 1),
	})
	tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"key":     admintok,
		"user_id": uid,
		"title":   tt.ShortLipsum(1004, 1),
	})
	tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"key":     admintok,
		"user_id": uid,
		"title":   tt.ShortLipsum(1005, 1),
	})

	assertCounter(5)
}

func TestGetUserNoPro(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/v2/users", gin.H{
		"no_client": true,
	})

	uid := fmt.Sprintf("%v", r0["user_id"])
	readtok := r0["read_key"].(string)

	r1 := tt.RequestAuthGet[gin.H](t, readtok, baseUrl, "/api/v2/users/"+uid)

	tt.AssertEqual(t, "user_id", uid, fmt.Sprintf("%v", r1["user_id"]))
	tt.AssertEqual(t, "username", nil, r1["username"])
	tt.AssertNotEqual(t, "timestamp_created", nil, r1["timestamp_created"])
	tt.AssertEqual(t, "timestamp_lastread", nil, r1["timestamp_lastread"])
	tt.AssertEqual(t, "timestamp_lastsent", nil, r1["timestamp_lastsent"])
	tt.AssertEqual(t, "messages_sent", "0", fmt.Sprintf("%v", r1["messages_sent"]))
	tt.AssertEqual(t, "quota_used", "0", fmt.Sprintf("%v", r1["quota_used"]))
	tt.AssertEqual(t, "quota_remaining", "50", fmt.Sprintf("%v", r1["quota_remaining"]))
	tt.AssertEqual(t, "quota_max", "50", fmt.Sprintf("%v", r1["quota_max"]))
	tt.AssertEqual(t, "is_pro", "false", fmt.Sprintf("%v", r1["is_pro"]))
	tt.AssertEqual(t, "default_channel", "main", fmt.Sprintf("%v", r1["default_channel"]))
	tt.AssertEqual(t, "max_body_size", "2048", fmt.Sprintf("%v", r1["max_body_size"]))
	tt.AssertEqual(t, "max_title_length", "120", fmt.Sprintf("%v", r1["max_title_length"]))
	tt.AssertEqual(t, "default_priority", "1", fmt.Sprintf("%v", r1["default_priority"]))
	tt.AssertEqual(t, "max_channel_name_length", "120", fmt.Sprintf("%v", r1["max_channel_name_length"]))
	tt.AssertEqual(t, "max_channel_description_length", "300", fmt.Sprintf("%v", r1["max_channel_description_length"]))
	tt.AssertEqual(t, "max_sender_name_length", "120", fmt.Sprintf("%v", r1["max_sender_name_length"]))
	tt.AssertEqual(t, "max_user_message_id_length", "64", fmt.Sprintf("%v", r1["max_user_message_id_length"]))
}

func TestGetUserPro(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/v2/users", gin.H{
		"no_client": true,
		"pro_token": "ANDROID|v2|PURCHASED:DUMMY_TOK_XX",
	})

	uid := fmt.Sprintf("%v", r0["user_id"])
	readtok := r0["read_key"].(string)

	r1 := tt.RequestAuthGet[gin.H](t, readtok, baseUrl, "/api/v2/users/"+uid)

	tt.AssertEqual(t, "user_id", uid, fmt.Sprintf("%v", r1["user_id"]))
	tt.AssertEqual(t, "username", nil, r1["username"])
	tt.AssertNotEqual(t, "timestamp_created", nil, r1["timestamp_created"])
	tt.AssertEqual(t, "timestamp_lastread", nil, r1["timestamp_lastread"])
	tt.AssertEqual(t, "timestamp_lastsent", nil, r1["timestamp_lastsent"])
	tt.AssertEqual(t, "messages_sent", "0", fmt.Sprintf("%v", r1["messages_sent"]))
	tt.AssertEqual(t, "quota_used", "0", fmt.Sprintf("%v", r1["quota_used"]))
	tt.AssertEqual(t, "quota_remaining", "5000", fmt.Sprintf("%v", r1["quota_remaining"]))
	tt.AssertEqual(t, "quota_max", "5000", fmt.Sprintf("%v", r1["quota_max"]))
	tt.AssertEqual(t, "is_pro", "true", fmt.Sprintf("%v", r1["is_pro"]))
	tt.AssertEqual(t, "default_channel", "main", fmt.Sprintf("%v", r1["default_channel"]))
	tt.AssertEqual(t, "max_body_size", "2097152", fmt.Sprintf("%d", (int64)(r1["max_body_size"].(float64))))
	tt.AssertEqual(t, "max_title_length", "120", fmt.Sprintf("%v", r1["max_title_length"]))
	tt.AssertEqual(t, "default_priority", "1", fmt.Sprintf("%v", r1["default_priority"]))
	tt.AssertEqual(t, "max_channel_name_length", "120", fmt.Sprintf("%v", r1["max_channel_name_length"]))
	tt.AssertEqual(t, "max_channel_description_length", "300", fmt.Sprintf("%v", r1["max_channel_description_length"]))
	tt.AssertEqual(t, "max_sender_name_length", "120", fmt.Sprintf("%v", r1["max_sender_name_length"]))
	tt.AssertEqual(t, "max_user_message_id_length", "64", fmt.Sprintf("%v", r1["max_user_message_id_length"]))
}
