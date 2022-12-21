package test

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	tt "blackforestbytes.com/simplecloudnotifier/test/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"testing"
)

func TestCreateChannel(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := int(r0["user_id"].(float64))
	admintok := r0["admin_key"].(string)

	type chanlist struct {
		Channels []gin.H `json:"channels"`
	}

	{
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid))
		tt.AssertMappedSet(t, "channels", []string{}, clist.Channels, "name")
	}

	tt.RequestAuthPost[gin.H](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid), gin.H{
		"name": "test",
	})

	{
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid))
		tt.AssertEqual(t, "chan.len", 1, len(clist.Channels))
		tt.AssertMappedSet(t, "channels", []string{"test"}, clist.Channels, "name")
	}

	tt.RequestAuthPost[gin.H](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid), gin.H{
		"name": "asdf",
	})

	{
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid))
		tt.AssertMappedSet(t, "channels", []string{"asdf", "test"}, clist.Channels, "name")
	}
}

func TestCreateChannelNameTooLong(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := int(r0["user_id"].(float64))
	admintok := r0["admin_key"].(string)

	tt.RequestAuthPostShouldFail(t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid), gin.H{
		"name": langext.StrRepeat("X", 121),
	}, 400, apierr.CHANNEL_TOO_LONG)
}

func TestChannelNameNormalization(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := int(r0["user_id"].(float64))
	admintok := r0["admin_key"].(string)

	type chanlist struct {
		Channels []gin.H `json:"channels"`
	}

	{
		chan0 := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid))
		tt.AssertEqual(t, "chan-count", 0, len(chan0.Channels))
	}

	tt.RequestAuthPost[gin.H](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid), gin.H{
		"name": "tESt",
	})

	{
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid))
		tt.AssertEqual(t, "chan.len", 1, len(clist.Channels))
		tt.AssertEqual(t, "chan.name", "test", clist.Channels[0]["name"])
	}

	tt.RequestAuthPostShouldFail(t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid), gin.H{
		"name": "test",
	}, 409, apierr.CHANNEL_ALREADY_EXISTS)

	tt.RequestAuthPostShouldFail(t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid), gin.H{
		"name": "TEST",
	}, 409, apierr.CHANNEL_ALREADY_EXISTS)

	tt.RequestAuthPostShouldFail(t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid), gin.H{
		"name": "Test",
	}, 409, apierr.CHANNEL_ALREADY_EXISTS)

	tt.RequestAuthPostShouldFail(t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid), gin.H{
		"name": "Test ",
	}, 409, apierr.CHANNEL_ALREADY_EXISTS)

	tt.RequestAuthPostShouldFail(t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid), gin.H{
		"name": " Test",
	}, 409, apierr.CHANNEL_ALREADY_EXISTS)

	tt.RequestAuthPostShouldFail(t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid), gin.H{
		"name": " T e s t ",
	}, 409, apierr.CHANNEL_ALREADY_EXISTS)

	{
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid))
		tt.AssertEqual(t, "chan.len", 1, len(clist.Channels))
		tt.AssertEqual(t, "chan.name", "test", clist.Channels[0]["name"])
	}
}

func TestListChannelsOwned(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	type chanlist struct {
		Channels []gin.H `json:"channels"`
	}

	testdata := map[int][]string{
		0:  {"main", "chattingchamber", "unicdhll", "promotions", "reminders"},
		1:  {"promotions"},
		2:  {},
		3:  {},
		4:  {},
		5:  {},
		6:  {},
		7:  {},
		8:  {},
		9:  {},
		10: {},
		11: {},
		12: {},
		13: {},
	}

	for k, v := range testdata {
		r0 := tt.RequestAuthGet[chanlist](t, data.User[k].AdminKey, baseUrl, fmt.Sprintf("/api/users/%d/channels", data.User[k].UID))
		tt.AssertMappedSet(t, fmt.Sprintf("%d->chanlist", k), v, r0.Channels, "name")
	}
}

func TestListChannelsSubscribedAny(t *testing.T) {
	t.SkipNow() //TODO
}

func TestListChannelsAllAny(t *testing.T) {
	t.SkipNow() //TODO
}

func TestListChannelsSubscribed(t *testing.T) {
	t.SkipNow() //TODO
}

func TestListChannelsAll(t *testing.T) {
	t.SkipNow() //TODO
}

//TODO test missing channel-xx methods
