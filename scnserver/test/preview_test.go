package test

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	tt "blackforestbytes.com/simplecloudnotifier/test/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"testing"
)

func TestGetChannelPreview(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	type chanlist struct {
		Channels []gin.H `json:"channels"`
	}

	clist := tt.RequestAuthGet[chanlist](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", data.User[0].UID))

	chan1 := *langext.ArrFirstOrNil(clist.Channels, func(v gin.H) bool { return v["internal_name"] == "Reminders" })

	{
		chan1_rq := tt.RequestAuthGet[gin.H](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s", data.User[0].UID, chan1["channel_id"]))
		tt.AssertEqual(t, "channel_id", chan1["channel_id"], chan1_rq["channel_id"])
		tt.AssertEqual(t, "display_name", "Reminders", chan1_rq["display_name"])
		tt.AssertEqual(t, "internal_name", "Reminders", chan1_rq["internal_name"])
		tt.AssertEqual(t, "description_name", nil, chan1_rq["description_name"])
	}

	{
		chan1_rq := tt.RequestAuthGet[gin.H](t, data.User[0].ReadKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s", data.User[0].UID, chan1["channel_id"]))
		tt.AssertEqual(t, "channel_id", chan1["channel_id"], chan1_rq["channel_id"])
		tt.AssertEqual(t, "display_name", "Reminders", chan1_rq["display_name"])
		tt.AssertEqual(t, "internal_name", "Reminders", chan1_rq["internal_name"])
		tt.AssertEqual(t, "description_name", nil, chan1_rq["description_name"])
	}

	{
		chan1_rq := tt.RequestAuthGet[gin.H](t, data.User[0].ReadKey, baseUrl, fmt.Sprintf("/api/v2/preview/channels/%s", chan1["channel_id"]))
		tt.AssertEqual(t, "channel_id", chan1["channel_id"], chan1_rq["channel_id"])
		tt.AssertEqual(t, "display_name", "Reminders", chan1_rq["display_name"])
		tt.AssertEqual(t, "internal_name", "Reminders", chan1_rq["internal_name"])
		tt.AssertEqual(t, "description_name", nil, chan1_rq["description_name"])
	}

	{
		chan1_rq := tt.RequestAuthGet[gin.H](t, data.User[0].SendKey, baseUrl, fmt.Sprintf("/api/v2/preview/channels/%s", chan1["channel_id"]))
		tt.AssertEqual(t, "channel_id", chan1["channel_id"], chan1_rq["channel_id"])
		tt.AssertEqual(t, "display_name", "Reminders", chan1_rq["display_name"])
		tt.AssertEqual(t, "internal_name", "Reminders", chan1_rq["internal_name"])
		tt.AssertEqual(t, "description_name", nil, chan1_rq["description_name"])
	}

	{
		chan1_rq := tt.RequestAuthGet[gin.H](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/preview/channels/%s", chan1["channel_id"]))
		tt.AssertEqual(t, "channel_id", chan1["channel_id"], chan1_rq["channel_id"])
		tt.AssertEqual(t, "display_name", "Reminders", chan1_rq["display_name"])
		tt.AssertEqual(t, "internal_name", "Reminders", chan1_rq["internal_name"])
		tt.AssertEqual(t, "description_name", nil, chan1_rq["description_name"])
	}

	{
		tt.RequestAuthGetShouldFail(t, data.User[1].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s", data.User[0].UID, chan1["channel_id"]), 401, apierr.USER_AUTH_FAILED)
	}

	{
		tt.RequestAuthGetShouldFail(t, data.User[1].ReadKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s", data.User[0].UID, chan1["channel_id"]), 401, apierr.USER_AUTH_FAILED)
	}

	{
		chan1_rq := tt.RequestAuthGet[gin.H](t, data.User[1].AdminKey, baseUrl, fmt.Sprintf("/api/v2/preview/channels/%s", chan1["channel_id"]))
		tt.AssertEqual(t, "channel_id", chan1["channel_id"], chan1_rq["channel_id"])
		tt.AssertEqual(t, "display_name", "Reminders", chan1_rq["display_name"])
		tt.AssertEqual(t, "internal_name", "Reminders", chan1_rq["internal_name"])
		tt.AssertEqual(t, "description_name", nil, chan1_rq["description_name"])
	}

	{
		chan1_rq := tt.RequestAuthGet[gin.H](t, data.User[1].SendKey, baseUrl, fmt.Sprintf("/api/v2/preview/channels/%s", chan1["channel_id"]))
		tt.AssertEqual(t, "channel_id", chan1["channel_id"], chan1_rq["channel_id"])
		tt.AssertEqual(t, "display_name", "Reminders", chan1_rq["display_name"])
		tt.AssertEqual(t, "internal_name", "Reminders", chan1_rq["internal_name"])
		tt.AssertEqual(t, "description_name", nil, chan1_rq["description_name"])
	}

	{
		chan1_rq := tt.RequestAuthGet[gin.H](t, data.User[1].ReadKey, baseUrl, fmt.Sprintf("/api/v2/preview/channels/%s", chan1["channel_id"]))
		tt.AssertEqual(t, "channel_id", chan1["channel_id"], chan1_rq["channel_id"])
		tt.AssertEqual(t, "display_name", "Reminders", chan1_rq["display_name"])
		tt.AssertEqual(t, "internal_name", "Reminders", chan1_rq["internal_name"])
		tt.AssertEqual(t, "description_name", nil, chan1_rq["description_name"])
	}

}

func TestGetUserPreview(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	{
		user_rq_1 := tt.RequestAuthGet[gin.H](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s", data.User[0].UID))
		tt.AssertEqual(t, "user_id", data.User[0].UID, user_rq_1["user_id"])

		user_rq_2 := tt.RequestAuthGet[gin.H](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/preview/users/%s", data.User[0].UID))
		tt.AssertEqual(t, "user_id", user_rq_1["user_id"], user_rq_2["user_id"])
		tt.AssertEqual(t, "username", user_rq_1["username"], user_rq_2["username"])
	}

	{
		user_rq_1 := tt.RequestAuthGet[gin.H](t, data.User[0].ReadKey, baseUrl, fmt.Sprintf("/api/v2/users/%s", data.User[0].UID))
		tt.AssertEqual(t, "user_id", data.User[0].UID, user_rq_1["user_id"])

		user_rq_2 := tt.RequestAuthGet[gin.H](t, data.User[0].ReadKey, baseUrl, fmt.Sprintf("/api/v2/preview/users/%s", data.User[0].UID))
		tt.AssertEqual(t, "user_id", user_rq_1["user_id"], user_rq_2["user_id"])
		tt.AssertEqual(t, "username", user_rq_1["username"], user_rq_2["username"])
	}

	{
		tt.RequestAuthGetShouldFail(t, data.User[0].SendKey, baseUrl, fmt.Sprintf("/api/v2/users/%s", data.User[0].UID), 401, apierr.USER_AUTH_FAILED)
	}

	{
		tt.RequestAuthGetShouldFail(t, data.User[1].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s", data.User[0].UID), 401, apierr.USER_AUTH_FAILED)
		tt.RequestAuthGetShouldFail(t, data.User[1].ReadKey, baseUrl, fmt.Sprintf("/api/v2/users/%s", data.User[0].UID), 401, apierr.USER_AUTH_FAILED)
		tt.RequestAuthGetShouldFail(t, data.User[1].SendKey, baseUrl, fmt.Sprintf("/api/v2/users/%s", data.User[0].UID), 401, apierr.USER_AUTH_FAILED)
	}

	{
		tt.RequestAuthGet[gin.H](t, data.User[0].ReadKey, baseUrl, fmt.Sprintf("/api/v2/preview/users/%s", data.User[0].UID))
		tt.RequestAuthGet[gin.H](t, data.User[0].SendKey, baseUrl, fmt.Sprintf("/api/v2/preview/users/%s", data.User[0].UID))
		tt.RequestAuthGet[gin.H](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/preview/users/%s", data.User[0].UID))

		tt.RequestAuthGet[gin.H](t, data.User[1].ReadKey, baseUrl, fmt.Sprintf("/api/v2/preview/users/%s", data.User[0].UID))
		tt.RequestAuthGet[gin.H](t, data.User[1].SendKey, baseUrl, fmt.Sprintf("/api/v2/preview/users/%s", data.User[0].UID))
		tt.RequestAuthGet[gin.H](t, data.User[1].AdminKey, baseUrl, fmt.Sprintf("/api/v2/preview/users/%s", data.User[0].UID))
	}

	{
		user_rq_1 := tt.RequestAuthGet[gin.H](t, data.User[1].ReadKey, baseUrl, fmt.Sprintf("/api/v2/preview/users/%s", data.User[2].UID))
		tt.AssertEqual(t, "username", "Dreamer23", user_rq_1["username"])
	}

}

func TestGetKeyTokenPreview(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	type keyobj struct {
		AllChannels  bool     `json:"all_channels"`
		Channels     []string `json:"channels"`
		KeytokenId   string   `json:"keytoken_id"`
		MessagesSent int      `json:"messages_sent"`
		Name         string   `json:"name"`
		OwnerUserId  string   `json:"owner_user_id"`
		Permissions  string   `json:"permissions"`
	}
	type keylist struct {
		Keys []keyobj `json:"keys"`
	}

	klist := tt.RequestAuthGet[keylist](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", data.User[0].UID))

	{
		rq_1 := tt.RequestAuthGet[gin.H](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys/%s", data.User[0].UID, klist.Keys[0].KeytokenId))
		tt.AssertEqual(t, "all_channels", klist.Keys[0].AllChannels, rq_1["all_channels"])
		tt.AssertStrRepEqual(t, "channels", klist.Keys[0].Channels, rq_1["channels"])
		tt.AssertEqual(t, "keytoken_id", klist.Keys[0].KeytokenId, rq_1["keytoken_id"])
		tt.AssertEqual(t, "messages_sent", klist.Keys[0].MessagesSent, rq_1["messages_sent"])
		tt.AssertEqual(t, "name", klist.Keys[0].Name, rq_1["name"])
		tt.AssertEqual(t, "owner_user_id", klist.Keys[0].OwnerUserId, rq_1["owner_user_id"])
		tt.AssertEqual(t, "permissions", klist.Keys[0].Permissions, rq_1["permissions"])
	}

	{
		rq_2 := tt.RequestAuthGet[gin.H](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/preview/keys/%s", klist.Keys[0].KeytokenId))
		tt.AssertEqual(t, "all_channels", klist.Keys[0].AllChannels, rq_2["all_channels"])
		tt.AssertStrRepEqual(t, "channels", klist.Keys[0].Channels, rq_2["channels"])
		tt.AssertEqual(t, "keytoken_id", klist.Keys[0].KeytokenId, rq_2["keytoken_id"])
		tt.AssertEqual(t, "name", klist.Keys[0].Name, rq_2["name"])
		tt.AssertEqual(t, "owner_user_id", klist.Keys[0].OwnerUserId, rq_2["owner_user_id"])
		tt.AssertEqual(t, "permissions", klist.Keys[0].Permissions, rq_2["permissions"])
	}

	{
		tt.RequestAuthGetShouldFail(t, data.User[0].SendKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys/%s", data.User[0].UID, klist.Keys[0].KeytokenId), 401, apierr.USER_AUTH_FAILED)
		tt.RequestAuthGetShouldFail(t, data.User[1].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys/%s", data.User[0].UID, klist.Keys[0].KeytokenId), 401, apierr.USER_AUTH_FAILED)
		tt.RequestAuthGetShouldFail(t, data.User[1].ReadKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys/%s", data.User[0].UID, klist.Keys[0].KeytokenId), 401, apierr.USER_AUTH_FAILED)
		tt.RequestAuthGetShouldFail(t, data.User[1].SendKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys/%s", data.User[0].UID, klist.Keys[0].KeytokenId), 401, apierr.USER_AUTH_FAILED)
	}

	{
		rq_2 := tt.RequestAuthGet[gin.H](t, data.User[0].ReadKey, baseUrl, fmt.Sprintf("/api/v2/preview/keys/%s", klist.Keys[0].KeytokenId))
		tt.AssertEqual(t, "all_channels", klist.Keys[0].AllChannels, rq_2["all_channels"])
		tt.AssertStrRepEqual(t, "channels", klist.Keys[0].Channels, rq_2["channels"])
		tt.AssertEqual(t, "keytoken_id", klist.Keys[0].KeytokenId, rq_2["keytoken_id"])
		tt.AssertEqual(t, "name", klist.Keys[0].Name, rq_2["name"])
		tt.AssertEqual(t, "owner_user_id", klist.Keys[0].OwnerUserId, rq_2["owner_user_id"])
		tt.AssertEqual(t, "permissions", klist.Keys[0].Permissions, rq_2["permissions"])
	}

	{
		rq_2 := tt.RequestAuthGet[gin.H](t, data.User[0].SendKey, baseUrl, fmt.Sprintf("/api/v2/preview/keys/%s", klist.Keys[0].KeytokenId))
		tt.AssertEqual(t, "all_channels", klist.Keys[0].AllChannels, rq_2["all_channels"])
		tt.AssertStrRepEqual(t, "channels", klist.Keys[0].Channels, rq_2["channels"])
		tt.AssertEqual(t, "keytoken_id", klist.Keys[0].KeytokenId, rq_2["keytoken_id"])
		tt.AssertEqual(t, "name", klist.Keys[0].Name, rq_2["name"])
		tt.AssertEqual(t, "owner_user_id", klist.Keys[0].OwnerUserId, rq_2["owner_user_id"])
		tt.AssertEqual(t, "permissions", klist.Keys[0].Permissions, rq_2["permissions"])
	}

	{
		rq_2 := tt.RequestAuthGet[gin.H](t, data.User[1].AdminKey, baseUrl, fmt.Sprintf("/api/v2/preview/keys/%s", klist.Keys[0].KeytokenId))
		tt.AssertEqual(t, "all_channels", klist.Keys[0].AllChannels, rq_2["all_channels"])
		tt.AssertStrRepEqual(t, "channels", klist.Keys[0].Channels, rq_2["channels"])
		tt.AssertEqual(t, "keytoken_id", klist.Keys[0].KeytokenId, rq_2["keytoken_id"])
		tt.AssertEqual(t, "name", klist.Keys[0].Name, rq_2["name"])
		tt.AssertEqual(t, "owner_user_id", klist.Keys[0].OwnerUserId, rq_2["owner_user_id"])
		tt.AssertEqual(t, "permissions", klist.Keys[0].Permissions, rq_2["permissions"])
	}

	{
		rq_2 := tt.RequestAuthGet[gin.H](t, data.User[1].ReadKey, baseUrl, fmt.Sprintf("/api/v2/preview/keys/%s", klist.Keys[0].KeytokenId))
		tt.AssertEqual(t, "all_channels", klist.Keys[0].AllChannels, rq_2["all_channels"])
		tt.AssertStrRepEqual(t, "channels", klist.Keys[0].Channels, rq_2["channels"])
		tt.AssertEqual(t, "keytoken_id", klist.Keys[0].KeytokenId, rq_2["keytoken_id"])
		tt.AssertEqual(t, "name", klist.Keys[0].Name, rq_2["name"])
		tt.AssertEqual(t, "owner_user_id", klist.Keys[0].OwnerUserId, rq_2["owner_user_id"])
		tt.AssertEqual(t, "permissions", klist.Keys[0].Permissions, rq_2["permissions"])
	}

	{
		rq_2 := tt.RequestAuthGet[gin.H](t, data.User[1].SendKey, baseUrl, fmt.Sprintf("/api/v2/preview/keys/%s", klist.Keys[0].KeytokenId))
		tt.AssertEqual(t, "all_channels", klist.Keys[0].AllChannels, rq_2["all_channels"])
		tt.AssertStrRepEqual(t, "channels", klist.Keys[0].Channels, rq_2["channels"])
		tt.AssertEqual(t, "keytoken_id", klist.Keys[0].KeytokenId, rq_2["keytoken_id"])
		tt.AssertEqual(t, "name", klist.Keys[0].Name, rq_2["name"])
		tt.AssertEqual(t, "owner_user_id", klist.Keys[0].OwnerUserId, rq_2["owner_user_id"])
		tt.AssertEqual(t, "permissions", klist.Keys[0].Permissions, rq_2["permissions"])
	}

}
