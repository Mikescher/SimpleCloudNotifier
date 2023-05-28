package test

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	tt "blackforestbytes.com/simplecloudnotifier/test/util"
	"fmt"
	"github.com/gin-gonic/gin"
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

	tt.AssertEqual(t, "AllChannels", key7.KeytokenId, msg1["used_key_id"])

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
		"channels": []string{"main"},
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
