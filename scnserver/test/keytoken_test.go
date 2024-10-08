package test

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	tt "blackforestbytes.com/simplecloudnotifier/test/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"testing"
)

func TestTokenKeys(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	type keyobj struct {
		AllChannels  bool     `json:"all_channels"`
		Channels     []string `json:"channels"`
		KeytokenId   string   `json:"keytoken_id"`
		MessagesSent int      `json:"messages_sent"`
		Name         string   `json:"name"`
		OwnerUserId  string   `json:"owner_user_id"`
		Permissions  string   `json:"permissions"`
		Token        string   `json:"token"` // only in create
	}
	type keylist struct {
		Keys []keyobj `json:"keys"`
	}

	{
		klist := tt.RequestAuthGet[keylist](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", data.UID))
		tt.AssertEqual(t, "len(keys)", 3, len(klist.Keys))

		tt.AssertArrAny(t, "keys->any[Admin]", klist.Keys, func(s keyobj) bool {
			return s.AllChannels == true && s.Name == "AdminKey (default)" && s.Permissions == "A"
		})
		tt.AssertArrAny(t, "keys->any[Send]", klist.Keys, func(s keyobj) bool {
			return s.AllChannels == true && s.Name == "SendKey (default)" && s.Permissions == "CS"
		})
		tt.AssertArrAny(t, "keys->any[Read]", klist.Keys, func(s keyobj) bool {
			return s.AllChannels == true && s.Name == "ReadKey (default)" && s.Permissions == "UR;CR"
		})
	}

	key2 := tt.RequestAuthPost[keyobj](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", data.UID), gin.H{
		"all_channels": true,
		"channels":     []string{},
		"name":         "Admin2",
		"permissions":  "A",
	})

	tt.AssertEqual(t, "Name", "Admin2", key2.Name)
	tt.AssertEqual(t, "Permissions", "A", key2.Permissions)
	tt.AssertEqual(t, "AllChannels", true, key2.AllChannels)

	{
		klist := tt.RequestAuthGet[keylist](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", data.UID))
		tt.AssertEqual(t, "len(keys)", 4, len(klist.Keys))
	}

	key3 := tt.RequestAuthGet[keyobj](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys/%s", data.UID, key2.KeytokenId))

	tt.AssertEqual(t, "KeytokenId", key2.KeytokenId, key3.KeytokenId)
	tt.AssertEqual(t, "UserID", data.UID, key3.OwnerUserId)
	tt.AssertEqual(t, "Name", "Admin2", key3.Name)
	tt.AssertEqual(t, "Permissions", "A", key3.Permissions)
	tt.AssertEqual(t, "AllChannels", true, key3.AllChannels)

	tt.RequestAuthDelete[tt.Void](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys/%s", data.UID, key2.KeytokenId), gin.H{})

	tt.RequestAuthGetShouldFail(t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys/%s", data.UID, key2.KeytokenId), 404, apierr.KEY_NOT_FOUND)

	{
		klist := tt.RequestAuthGet[keylist](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", data.UID))
		tt.AssertEqual(t, "len(keys)", 3, len(klist.Keys))
	}

	chan0 := tt.RequestAuthPost[gin.H](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", data.UID), gin.H{
		"name": "testchan1",
	})
	chanid := fmt.Sprintf("%v", chan0["channel_id"])

	key4 := tt.RequestAuthPost[keyobj](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", data.UID), gin.H{
		"all_channels": false,
		"channels":     []string{chanid},
		"name":         "TKey1",
		"permissions":  "CS",
	})
	tt.AssertEqual(t, "Name", "TKey1", key4.Name)
	tt.AssertEqual(t, "Permissions", "CS", key4.Permissions)
	tt.AssertEqual(t, "AllChannels", false, key4.AllChannels)
	tt.AssertStrRepEqual(t, "Channels", []string{chanid}, key4.Channels)

	key5 := tt.RequestAuthPatch[keyobj](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys/%s", data.UID, key4.KeytokenId), gin.H{
		"all_channels": true,
		"channels":     []string{},
		"name":         "TKey2-A",
		"permissions":  "A",
	})
	tt.AssertEqual(t, "Name", "TKey2-A", key5.Name)
	tt.AssertEqual(t, "Permissions", "A", key5.Permissions)
	tt.AssertEqual(t, "AllChannels", true, key5.AllChannels)
	tt.AssertStrRepEqual(t, "Channels", []string{}, key5.Channels)

	key6 := tt.RequestAuthGet[keyobj](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys/%s", data.UID, key5.KeytokenId))
	tt.AssertEqual(t, "Name", "TKey2-A", key6.Name)
	tt.AssertEqual(t, "Permissions", "A", key6.Permissions)
	tt.AssertEqual(t, "AllChannels", true, key6.AllChannels)
	tt.AssertStrRepEqual(t, "Channels", []string{}, key6.Channels)

	key7 := tt.RequestAuthPost[keyobj](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", data.UID), gin.H{
		"all_channels": false,
		"channels":     []string{chanid},
		"name":         "TKey7",
		"permissions":  "CS",
	})

	{
		klist := tt.RequestAuthGet[keylist](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", data.UID))
		tt.AssertEqual(t, "len(keys)", 5, len(klist.Keys))
	}

	msg1s := tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"key":     key7.Token,
		"user_id": data.UID,
		"channel": "testchan1",
		"title":   "HelloWorld_001",
	})

	msg1 := tt.RequestAuthGet[gin.H](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages/%s", msg1s["scn_msg_id"]))

	tt.AssertEqual(t, "used_key_id", key7.KeytokenId, msg1["used_key_id"])

	tt.RequestPostShouldFail(t, baseUrl, "/", gin.H{
		"key":     key7.Token,
		"user_id": data.UID,
		"channel": "testchan2",
		"title":   "HelloWorld_001",
	}, 401, apierr.USER_AUTH_FAILED) // wrong channel

	tt.RequestPostShouldFail(t, baseUrl, "/", gin.H{
		"key":     key7.Token,
		"user_id": data.UID,
		"title":   "HelloWorld_001",
	}, 401, apierr.USER_AUTH_FAILED) // no channel (=main)

	tt.RequestAuthGetShouldFail(t, key7.Token, baseUrl, fmt.Sprintf("/api/v2/users/%s", data.UID), 401, apierr.USER_AUTH_FAILED) // no user read perm

	tt.RequestAuthPatchShouldFail(t, key7.Token, baseUrl, "/api/v2/users/"+data.UID, gin.H{"username": "my_user_001"}, 401, apierr.USER_AUTH_FAILED) // no user update perm

	key8 := tt.RequestAuthPost[keyobj](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", data.UID), gin.H{
		"all_channels": true,
		"channels":     []string{},
		"name":         "TKey7",
		"permissions":  "CR",
	})

	tt.RequestPostShouldFail(t, baseUrl, "/", gin.H{
		"key":     key8.Token,
		"user_id": data.UID,
		"title":   "HelloWorld_001",
	}, 401, apierr.USER_AUTH_FAILED) // no send perm

}

func TestTokenKeysInitial(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	type keyobj struct {
		AllChannels  bool     `json:"all_channels"`
		Channels     []string `json:"channels"`
		KeytokenId   string   `json:"keytoken_id"`
		MessagesSent int      `json:"messages_sent"`
		Name         string   `json:"name"`
		OwnerUserId  string   `json:"owner_user_id"`
		Permissions  string   `json:"permissions"`
		Token        string   `json:"token"` // only in create
	}
	type keylist struct {
		Keys []keyobj `json:"keys"`
	}

	klist := tt.RequestAuthGet[keylist](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", data.UID))
	tt.AssertEqual(t, "len(keys)", 3, len(klist.Keys))

	tt.AssertArrAny(t, "keys->any[Admin]", klist.Keys, func(s keyobj) bool {
		return s.AllChannels == true && s.Name == "AdminKey (default)" && s.Permissions == "A"
	})
	tt.AssertArrAny(t, "keys->any[Send]", klist.Keys, func(s keyobj) bool {
		return s.AllChannels == true && s.Name == "SendKey (default)" && s.Permissions == "CS"
	})
	tt.AssertArrAny(t, "keys->any[Read]", klist.Keys, func(s keyobj) bool {
		return s.AllChannels == true && s.Name == "ReadKey (default)" && s.Permissions == "UR;CR"
	})
}

func TestTokenKeysCreate(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	type keyobj struct {
		AllChannels  bool     `json:"all_channels"`
		Channels     []string `json:"channels"`
		KeytokenId   string   `json:"keytoken_id"`
		MessagesSent int      `json:"messages_sent"`
		Name         string   `json:"name"`
		OwnerUserId  string   `json:"owner_user_id"`
		Permissions  string   `json:"permissions"`
		Token        string   `json:"token"` // only in create
	}
	type keylist struct {
		Keys []keyobj `json:"keys"`
	}

	key2 := tt.RequestAuthPost[keyobj](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", data.UID), gin.H{
		"all_channels": true,
		"channels":     []string{},
		"name":         "Admin2",
		"permissions":  "A",
	})

	tt.AssertEqual(t, "Name", "Admin2", key2.Name)
	tt.AssertEqual(t, "Permissions", "A", key2.Permissions)
	tt.AssertEqual(t, "AllChannels", true, key2.AllChannels)

	{
		klist := tt.RequestAuthGet[keylist](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", data.UID))
		tt.AssertEqual(t, "len(keys)", 4, len(klist.Keys))
	}

	key3 := tt.RequestAuthGet[keyobj](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys/%s", data.UID, key2.KeytokenId))

	tt.AssertEqual(t, "KeytokenId", key2.KeytokenId, key3.KeytokenId)
	tt.AssertEqual(t, "UserID", data.UID, key3.OwnerUserId)
	tt.AssertEqual(t, "Name", "Admin2", key3.Name)
	tt.AssertEqual(t, "Permissions", "A", key3.Permissions)
	tt.AssertEqual(t, "AllChannels", true, key3.AllChannels)

}

func TestTokenKeysUpdate(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	type keyobj struct {
		AllChannels  bool     `json:"all_channels"`
		Channels     []string `json:"channels"`
		KeytokenId   string   `json:"keytoken_id"`
		MessagesSent int      `json:"messages_sent"`
		Name         string   `json:"name"`
		OwnerUserId  string   `json:"owner_user_id"`
		Permissions  string   `json:"permissions"`
		Token        string   `json:"token"` // only in create
	}
	type keylist struct {
		Keys []keyobj `json:"keys"`
	}

	key2 := tt.RequestAuthPost[keyobj](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", data.UID), gin.H{
		"all_channels": true,
		"channels":     []string{},
		"name":         "Admin2",
		"permissions":  "A",
	})

	tt.AssertEqual(t, "Name", "Admin2", key2.Name)
	tt.AssertEqual(t, "Permissions", "A", key2.Permissions)
	tt.AssertEqual(t, "AllChannels", true, key2.AllChannels)

	key3 := tt.RequestAuthGet[keyobj](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys/%s", data.UID, key2.KeytokenId))

	tt.AssertEqual(t, "KeytokenId", key2.KeytokenId, key3.KeytokenId)
	tt.AssertEqual(t, "UserID", data.UID, key3.OwnerUserId)
	tt.AssertEqual(t, "Name", "Admin2", key3.Name)
	tt.AssertEqual(t, "Permissions", "A", key3.Permissions)
	tt.AssertEqual(t, "AllChannels", true, key3.AllChannels)

	key5 := tt.RequestAuthPatch[keyobj](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys/%s", data.UID, key3.KeytokenId), gin.H{
		"name": "Hello",
	})
	tt.AssertEqual(t, "Name", "Hello", key5.Name)
	tt.AssertEqual(t, "Permissions", "A", key5.Permissions)
	tt.AssertEqual(t, "AllChannels", true, key5.AllChannels)
	tt.AssertStrRepEqual(t, "Channels", []string{}, key5.Channels)

	key6 := tt.RequestAuthGet[keyobj](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys/%s", data.UID, key5.KeytokenId))
	tt.AssertEqual(t, "Name", "Hello", key6.Name)
	tt.AssertEqual(t, "Permissions", "A", key6.Permissions)
	tt.AssertEqual(t, "AllChannels", true, key6.AllChannels)
	tt.AssertStrRepEqual(t, "Channels", []string{}, key6.Channels)

}

func TestTokenKeysDelete(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	type keyobj struct {
		AllChannels  bool     `json:"all_channels"`
		Channels     []string `json:"channels"`
		KeytokenId   string   `json:"keytoken_id"`
		MessagesSent int      `json:"messages_sent"`
		Name         string   `json:"name"`
		OwnerUserId  string   `json:"owner_user_id"`
		Permissions  string   `json:"permissions"`
		Token        string   `json:"token"` // only in create
	}
	type keylist struct {
		Keys []keyobj `json:"keys"`
	}

	key2 := tt.RequestAuthPost[keyobj](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", data.UID), gin.H{
		"all_channels": true,
		"channels":     []string{},
		"name":         "Admin2",
		"permissions":  "A",
	})

	tt.AssertEqual(t, "Name", "Admin2", key2.Name)
	tt.AssertEqual(t, "Permissions", "A", key2.Permissions)
	tt.AssertEqual(t, "AllChannels", true, key2.AllChannels)

	{
		klist := tt.RequestAuthGet[keylist](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", data.UID))
		tt.AssertEqual(t, "len(keys)", 4, len(klist.Keys))
	}

	tt.RequestAuthDelete[tt.Void](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys/%s", data.UID, key2.KeytokenId), gin.H{})

	{
		klist := tt.RequestAuthGet[keylist](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", data.UID))
		tt.AssertEqual(t, "len(keys)", 3, len(klist.Keys))
	}

}

func TestTokenKeysDeleteSelf(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	type keyobj struct {
		AllChannels  bool     `json:"all_channels"`
		Channels     []string `json:"channels"`
		KeytokenId   string   `json:"keytoken_id"`
		MessagesSent int      `json:"messages_sent"`
		Name         string   `json:"name"`
		OwnerUserId  string   `json:"owner_user_id"`
		Permissions  string   `json:"permissions"`
		Token        string   `json:"token"` // only in create
	}
	type keylist struct {
		Keys []keyobj `json:"keys"`
	}

	klist := tt.RequestAuthGet[keylist](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", data.UID))

	ak := ""
	for _, v := range klist.Keys {
		if v.Permissions == "A" {
			ak = v.KeytokenId
		}
	}

	tt.RequestAuthDeleteShouldFail(t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys/%s", data.UID, ak), gin.H{}, 400, apierr.CANNOT_SELFDELETE_KEY)
}

func TestTokenKeysDowngradeSelf(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	chan0 := tt.RequestAuthPost[gin.H](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", data.UID), gin.H{
		"name": "testchan1",
	})
	chanid := fmt.Sprintf("%v", chan0["channel_id"])

	type keyobj struct {
		AllChannels  bool     `json:"all_channels"`
		Channels     []string `json:"channels"`
		KeytokenId   string   `json:"keytoken_id"`
		MessagesSent int      `json:"messages_sent"`
		Name         string   `json:"name"`
		OwnerUserId  string   `json:"owner_user_id"`
		Permissions  string   `json:"permissions"`
		Token        string   `json:"token"` // only in create
	}
	type keylist struct {
		Keys []keyobj `json:"keys"`
	}

	klist := tt.RequestAuthGet[keylist](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", data.UID))

	ak := ""
	for _, v := range klist.Keys {
		if v.Permissions == "A" {
			ak = v.KeytokenId
		}
	}

	tt.RequestAuthPatchShouldFail(t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys/%s", data.UID, ak), gin.H{
		"permissions": "CR",
	}, 400, apierr.CANNOT_SELFUPDATE_KEY)

	tt.RequestAuthPatchShouldFail(t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys/%s", data.UID, ak), gin.H{
		"all_channels": false,
	}, 400, apierr.CANNOT_SELFUPDATE_KEY)

	tt.RequestAuthPatchShouldFail(t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys/%s", data.UID, ak), gin.H{
		"channels": []string{chanid},
	}, 400, apierr.CANNOT_SELFUPDATE_KEY)

	tt.RequestAuthPatch[tt.Void](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys/%s", data.UID, ak), gin.H{
		"name": "This-is-allowed",
	})

	keyOut := tt.RequestAuthGet[keyobj](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys/%s", data.UID, ak))

	tt.AssertEqual(t, "UserID", data.UID, keyOut.OwnerUserId)
	tt.AssertEqual(t, "Name", "This-is-allowed", keyOut.Name)
	tt.AssertEqual(t, "Permissions", "A", keyOut.Permissions)
	tt.AssertEqual(t, "AllChannels", true, keyOut.AllChannels)
	tt.AssertStrRepEqual(t, "Channels", []string{}, keyOut.Channels)

}

func TestTokenKeysPermissions(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	type keyobj struct {
		AllChannels  bool     `json:"all_channels"`
		Channels     []string `json:"channels"`
		KeytokenId   string   `json:"keytoken_id"`
		MessagesSent int      `json:"messages_sent"`
		Name         string   `json:"name"`
		OwnerUserId  string   `json:"owner_user_id"`
		Permissions  string   `json:"permissions"`
		Token        string   `json:"token"` // only in create
	}

	chan0 := tt.RequestAuthPost[gin.H](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", data.UID), gin.H{
		"name": "testchan1",
	})
	chanid := fmt.Sprintf("%v", chan0["channel_id"])

	key7 := tt.RequestAuthPost[keyobj](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", data.UID), gin.H{
		"all_channels": false,
		"channels":     []string{chanid},
		"name":         "TKey7",
		"permissions":  "CS",
	})

	tt.RequestPostShouldFail(t, baseUrl, "/", gin.H{
		"key":     key7.Token,
		"user_id": data.UID,
		"channel": "testchan2",
		"title":   "HelloWorld_001",
	}, 401, apierr.USER_AUTH_FAILED) // wrong channel

	tt.RequestPostShouldFail(t, baseUrl, "/", gin.H{
		"key":     key7.Token,
		"user_id": data.UID,
		"title":   "HelloWorld_001",
	}, 401, apierr.USER_AUTH_FAILED) // no channel (=main)

	tt.RequestAuthGetShouldFail(t, key7.Token, baseUrl, fmt.Sprintf("/api/v2/users/%s", data.UID), 401, apierr.USER_AUTH_FAILED) // no user read perm

	tt.RequestAuthPatchShouldFail(t, key7.Token, baseUrl, "/api/v2/users/"+data.UID, gin.H{"username": "my_user_001"}, 401, apierr.USER_AUTH_FAILED) // no user update perm

	key8 := tt.RequestAuthPost[keyobj](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", data.UID), gin.H{
		"all_channels": true,
		"channels":     []string{},
		"name":         "TKey7",
		"permissions":  "CR",
	})

	tt.RequestPostShouldFail(t, baseUrl, "/", gin.H{
		"key":     key8.Token,
		"user_id": data.UID,
		"title":   "HelloWorld_001",
	}, 401, apierr.USER_AUTH_FAILED) // no send perm

}

func TestTokenKeysMessageCounter(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/v2/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	type keyobj struct {
		AllChannels  bool     `json:"all_channels"`
		Channels     []string `json:"channels"`
		KeytokenId   string   `json:"keytoken_id"`
		MessagesSent int      `json:"messages_sent"`
		Name         string   `json:"name"`
		OwnerUserId  string   `json:"owner_user_id"`
		Permissions  string   `json:"permissions"`
		Token        string   `json:"token"` // only in create
	}
	type keylist struct {
		Keys []keyobj `json:"keys"`
	}

	uid := r0["user_id"].(string)
	admintok := r0["admin_key"].(string)
	sendtok := r0["send_key"].(string)

	klist := tt.RequestAuthGet[keylist](t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", uid))
	tt.AssertEqual(t, "len(keys)", 3, len(klist.Keys))

	admintokid := langext.ArrFirstOrNil(klist.Keys, func(v keyobj) bool { return v.Name == "AdminKey (default)" }).KeytokenId
	sendtokid := langext.ArrFirstOrNil(klist.Keys, func(v keyobj) bool { return v.Name == "SendKey (default)" }).KeytokenId
	readtokid := langext.ArrFirstOrNil(klist.Keys, func(v keyobj) bool { return v.Name == "ReadKey (default)" }).KeytokenId

	assertCounter := func(c0 int, c1 int, c2 int) {
		r1 := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/v2/users/"+uid+"/keys/"+admintokid)
		tt.AssertStrRepEqual(t, "c0.messages_sent", c0, r1["messages_sent"])

		r2 := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/v2/users/"+uid+"/keys/"+sendtokid)
		tt.AssertStrRepEqual(t, "c1.messages_sent", c1, r2["messages_sent"])

		r3 := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/v2/users/"+uid+"/keys/"+readtokid)
		tt.AssertStrRepEqual(t, "c2.messages_sent", c2, r3["messages_sent"])
	}

	assertCounter(0, 0, 0)

	tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"key":     admintok,
		"user_id": uid,
		"title":   tt.ShortLipsum(1001, 1),
	})

	assertCounter(1, 0, 0)

	tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"key":     admintok,
		"user_id": uid,
		"title":   tt.ShortLipsum(1002, 1),
	})

	assertCounter(2, 0, 0)

	tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"key":     sendtok,
		"user_id": uid,
		"title":   tt.ShortLipsum(1002, 1),
	})

	assertCounter(2, 1, 0)

	tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"key":     sendtok,
		"user_id": uid,
		"channel": "Chan1",
		"title":   tt.ShortLipsum(1003, 1),
	})
	tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"key":     sendtok,
		"user_id": uid,
		"channel": "Chan2",
		"title":   tt.ShortLipsum(1004, 1),
	})
	tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"key":     sendtok,
		"user_id": uid,
		"channel": "Chan2",
		"title":   tt.ShortLipsum(1005, 1),
	})

	assertCounter(2, 4, 0)
	assertCounter(2, 4, 0)

	tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"key":     admintok,
		"user_id": uid,
		"channel": "Chan2",
		"title":   tt.ShortLipsum(1004, 1),
	})

	assertCounter(3, 4, 0)

	tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"key":     admintok,
		"user_id": uid,
		"title":   tt.ShortLipsum(1002, 1),
	})

	assertCounter(4, 4, 0)

}

func TestTokenKeysCreateDefaultParam(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	type keyobj struct {
		AllChannels  bool     `json:"all_channels"`
		Channels     []string `json:"channels"`
		KeytokenId   string   `json:"keytoken_id"`
		MessagesSent int      `json:"messages_sent"`
		Name         string   `json:"name"`
		OwnerUserId  string   `json:"owner_user_id"`
		Permissions  string   `json:"permissions"`
		Token        string   `json:"token"` // only in create
	}

	chan0 := tt.RequestAuthPost[gin.H](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", data.UID), gin.H{
		"name": "testchan1",
	})
	chanid := fmt.Sprintf("%v", chan0["channel_id"])

	{
		key2 := tt.RequestAuthPost[keyobj](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", data.UID), gin.H{
			"name":        "K2",
			"permissions": "CS",
		})

		tt.AssertEqual(t, "Name", "K2", key2.Name)
		tt.AssertEqual(t, "Permissions", "CS", key2.Permissions)
		tt.AssertEqual(t, "AllChannels", true, key2.AllChannels)
		tt.AssertEqual(t, "Channels.Len", 0, len(key2.Channels))
	}

	{
		key2 := tt.RequestAuthPost[keyobj](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", data.UID), gin.H{
			"name":        "K3",
			"permissions": "CS",
			"channels":    []string{chanid},
		})

		tt.AssertEqual(t, "Name", "K3", key2.Name)
		tt.AssertEqual(t, "Permissions", "CS", key2.Permissions)
		tt.AssertEqual(t, "AllChannels", false, key2.AllChannels)
		tt.AssertEqual(t, "Channels.Len", 1, len(key2.Channels))
	}

	{
		key2 := tt.RequestAuthPost[keyobj](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", data.UID), gin.H{
			"name":         "K4",
			"permissions":  "CS",
			"all_channels": false,
		})

		tt.AssertEqual(t, "Name", "K4", key2.Name)
		tt.AssertEqual(t, "Permissions", "CS", key2.Permissions)
		tt.AssertEqual(t, "AllChannels", false, key2.AllChannels)
		tt.AssertEqual(t, "Channels.Len", 0, len(key2.Channels))
	}
}

func TestTokenKeysGetCurrent(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	type keyobj struct {
		AllChannels  bool     `json:"all_channels"`
		Channels     []string `json:"channels"`
		KeytokenId   string   `json:"keytoken_id"`
		MessagesSent int      `json:"messages_sent"`
		Name         string   `json:"name"`
		OwnerUserId  string   `json:"owner_user_id"`
		Permissions  string   `json:"permissions"`
		Token        string   `json:"token"` // only in create
	}
	type keylist struct {
		Keys []keyobj `json:"keys"`
	}

	klist := tt.RequestAuthGet[keylist](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", data.UID))
	tt.AssertEqual(t, "len(keys)", 3, len(klist.Keys))

	keyAdmin := *langext.ArrFirstOrNil(klist.Keys, func(v keyobj) bool { return v.Permissions == "A" })
	keySend := *langext.ArrFirstOrNil(klist.Keys, func(v keyobj) bool { return v.Permissions == "CS" })
	keyRead := *langext.ArrFirstOrNil(klist.Keys, func(v keyobj) bool { return v.Permissions == "UR;CR" })

	{
		currKey := tt.RequestAuthGet[keyobj](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys/current", data.UID))
		tt.AssertEqual(t, "currKey.KeytokenId", keyAdmin.KeytokenId, currKey.KeytokenId)
		tt.AssertEqual(t, "currKey.Permissions", "A", currKey.Permissions)
		tt.AssertEqual(t, "currKey.Token", data.AdminKey, currKey.Token)
	}

	{
		currKey := tt.RequestAuthGet[keyobj](t, data.SendKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys/current", data.UID))
		tt.AssertEqual(t, "currKey.KeytokenId", keySend.KeytokenId, currKey.KeytokenId)
		tt.AssertEqual(t, "currKey.Permissions", "CS", currKey.Permissions)
		tt.AssertEqual(t, "currKey.Token", data.SendKey, currKey.Token)
	}

	{
		currKey := tt.RequestAuthGet[keyobj](t, data.ReadKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys/current", data.UID))
		tt.AssertEqual(t, "currKey.KeytokenId", keyRead.KeytokenId, currKey.KeytokenId)
		tt.AssertEqual(t, "currKey.Permissions", "UR;CR", currKey.Permissions)
		tt.AssertEqual(t, "currKey.Token", data.ReadKey, currKey.Token)
	}

}
