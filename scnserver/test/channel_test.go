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
		tt.AssertMappedSet(t, "channels", []string{}, clist.Channels, "display_name")
		tt.AssertMappedSet(t, "channels", []string{}, clist.Channels, "internal_name")
	}

	tt.RequestAuthPost[gin.H](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid), gin.H{
		"name": "test",
	})

	{
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid))
		tt.AssertEqual(t, "chan.len", 1, len(clist.Channels))
		tt.AssertMappedSet(t, "channels", []string{"test"}, clist.Channels, "display_name")
		tt.AssertMappedSet(t, "channels", []string{"test"}, clist.Channels, "internal_name")
	}

	tt.RequestAuthPost[gin.H](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid), gin.H{
		"name": "asdf",
	})

	{
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid))
		tt.AssertMappedSet(t, "channels", []string{"asdf", "test"}, clist.Channels, "display_name")
		tt.AssertMappedSet(t, "channels", []string{"asdf", "test"}, clist.Channels, "internal_name")
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
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid))
		tt.AssertMappedSet(t, "channels", []string{}, clist.Channels, "display_name")
		tt.AssertMappedSet(t, "channels", []string{}, clist.Channels, "internal_name")
	}

	tt.RequestAuthPost[gin.H](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid), gin.H{
		"name": "tESt",
	})

	{
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid))
		tt.AssertMappedSet(t, "channels", []string{"tESt"}, clist.Channels, "display_name")
		tt.AssertMappedSet(t, "channels", []string{"test"}, clist.Channels, "internal_name")
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
		"name": "\rTeSt\n",
	}, 409, apierr.CHANNEL_ALREADY_EXISTS)

	{
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid))
		tt.AssertMappedSet(t, "channels", []string{"tESt"}, clist.Channels, "display_name")
		tt.AssertMappedSet(t, "channels", []string{"test"}, clist.Channels, "internal_name")
	}

	tt.RequestAuthPost[gin.H](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid), gin.H{
		"name": "  WeiRD_[\uF5FF]\\stUFf\r\n\t  ",
	})

	{
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid))
		tt.AssertMappedSet(t, "channels", []string{"tESt", "WeiRD_[\uF5FF]\\stUFf"}, clist.Channels, "display_name")
		tt.AssertMappedSet(t, "channels", []string{"test", "weird_[\uF5FF]\\stuff"}, clist.Channels, "internal_name")
	}

}

func TestListChannelsDefault(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	type chanlist struct {
		Channels []gin.H `json:"channels"`
	}

	testdata := map[int][]string{
		0:  {"main", "chatting chamber", "unicôdé häll \U0001f92a", "promotions", "reminders"},
		1:  {"main", "private"},
		2:  {"main", "ü", "ö", "ä"},
		3:  {"main", "\U0001f5ff", "innovations", "reminders"},
		4:  {"main"},
		5:  {"main", "test1", "test2", "test3", "test4", "test5"},
		6:  {"main", "security", "lipsum"},
		7:  {"main"},
		8:  {"main"},
		9:  {"main", "manual@chan"},
		10: {"main"},
		11: {"promotions"},
		12: {},
		13: {},
		14: {"main", "chan_self_subscribed", "chan_self_unsub"},
		15: {"main", "chan_other_nosub", "chan_other_request", "chan_other_accepted"},
	}

	for k, v := range testdata {
		r0 := tt.RequestAuthGet[chanlist](t, data.User[k].AdminKey, baseUrl, fmt.Sprintf("/api/users/%d/channels", data.User[k].UID))
		tt.AssertMappedSet(t, fmt.Sprintf("%d->chanlist", k), v, r0.Channels, "internal_name")
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
		0:  {"main", "chatting chamber", "unicôdé häll \U0001f92a", "promotions", "reminders"},
		1:  {"main", "private"},
		2:  {"main", "ü", "ö", "ä"},
		3:  {"main", "\U0001f5ff", "innovations", "reminders"},
		4:  {"main"},
		5:  {"main", "test1", "test2", "test3", "test4", "test5"},
		6:  {"main", "security", "lipsum"},
		7:  {"main"},
		8:  {"main"},
		9:  {"main", "manual@chan"},
		10: {"main"},
		11: {"promotions"},
		12: {},
		13: {},
		14: {"main", "chan_self_subscribed", "chan_self_unsub"},
		15: {"main", "chan_other_nosub", "chan_other_request", "chan_other_accepted"},
	}

	for k, v := range testdata {
		r0 := tt.RequestAuthGet[chanlist](t, data.User[k].AdminKey, baseUrl, fmt.Sprintf("/api/users/%d/channels?selector=%s", data.User[k].UID, "owned"))
		tt.AssertMappedSet(t, fmt.Sprintf("%d->chanlist", k), v, r0.Channels, "internal_name")
	}
}

func TestListChannelsSubscribedAny(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	type chanlist struct {
		Channels []gin.H `json:"channels"`
	}

	testdata := map[int][]string{
		0:  {"main", "chatting chamber", "unicôdé häll \U0001f92a", "promotions", "reminders"},
		1:  {"main", "private"},
		2:  {"main", "ü", "ö", "ä"},
		3:  {"main", "\U0001f5ff", "innovations", "reminders"},
		4:  {"main"},
		5:  {"main", "test1", "test2", "test3", "test4", "test5"},
		6:  {"main", "security", "lipsum"},
		7:  {"main"},
		8:  {"main"},
		9:  {"main", "manual@chan"},
		10: {"main"},
		11: {"promotions"},
		12: {},
		13: {},
		14: {"main", "chan_self_subscribed", "chan_other_request", "chan_other_accepted"},
		15: {"main", "chan_other_nosub", "chan_other_request", "chan_other_accepted"},
	}

	for k, v := range testdata {
		r0 := tt.RequestAuthGet[chanlist](t, data.User[k].AdminKey, baseUrl, fmt.Sprintf("/api/users/%d/channels?selector=%s", data.User[k].UID, "subscribed_any"))
		tt.AssertMappedSet(t, fmt.Sprintf("%d->chanlist", k), v, r0.Channels, "internal_name")
	}
}

func TestListChannelsAllAny(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	type chanlist struct {
		Channels []gin.H `json:"channels"`
	}

	testdata := map[int][]string{
		0:  {"main", "chatting chamber", "unicôdé häll \U0001f92a", "promotions", "reminders"},
		1:  {"main", "private"},
		2:  {"main", "ü", "ö", "ä"},
		3:  {"main", "\U0001f5ff", "innovations", "reminders"},
		4:  {"main"},
		5:  {"main", "test1", "test2", "test3", "test4", "test5"},
		6:  {"main", "security", "lipsum"},
		7:  {"main"},
		8:  {"main"},
		9:  {"main", "manual@chan"},
		10: {"main"},
		11: {"promotions"},
		12: {},
		13: {},
		14: {"main", "chan_self_subscribed", "chan_self_unsub", "chan_other_request", "chan_other_accepted"},
		15: {"main", "chan_other_nosub", "chan_other_request", "chan_other_accepted"},
	}

	for k, v := range testdata {
		r0 := tt.RequestAuthGet[chanlist](t, data.User[k].AdminKey, baseUrl, fmt.Sprintf("/api/users/%d/channels?selector=%s", data.User[k].UID, "all_any"))
		tt.AssertMappedSet(t, fmt.Sprintf("%d->chanlist", k), v, r0.Channels, "internal_name")
	}
}

func TestListChannelsSubscribed(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	type chanlist struct {
		Channels []gin.H `json:"channels"`
	}

	testdata := map[int][]string{
		0:  {"main", "chatting chamber", "unicôdé häll \U0001f92a", "promotions", "reminders"},
		1:  {"main", "private"},
		2:  {"main", "ü", "ö", "ä"},
		3:  {"main", "\U0001f5ff", "innovations", "reminders"},
		4:  {"main"},
		5:  {"main", "test1", "test2", "test3", "test4", "test5"},
		6:  {"main", "security", "lipsum"},
		7:  {"main"},
		8:  {"main"},
		9:  {"main", "manual@chan"},
		10: {"main"},
		11: {"promotions"},
		12: {},
		13: {},
		14: {"main", "chan_self_subscribed", "chan_other_accepted"},
		15: {"main", "chan_other_nosub", "chan_other_request", "chan_other_accepted"},
	}

	for k, v := range testdata {
		r0 := tt.RequestAuthGet[chanlist](t, data.User[k].AdminKey, baseUrl, fmt.Sprintf("/api/users/%d/channels?selector=%s", data.User[k].UID, "subscribed"))
		tt.AssertMappedSet(t, fmt.Sprintf("%d->chanlist", k), v, r0.Channels, "internal_name")
	}
}

func TestListChannelsAll(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	type chanlist struct {
		Channels []gin.H `json:"channels"`
	}

	testdata := map[int][]string{
		0:  {"main", "chatting chamber", "unicôdé häll \U0001f92a", "promotions", "reminders"},
		1:  {"main", "private"},
		2:  {"main", "ü", "ö", "ä"},
		3:  {"main", "\U0001f5ff", "innovations", "reminders"},
		4:  {"main"},
		5:  {"main", "test1", "test2", "test3", "test4", "test5"},
		6:  {"main", "security", "lipsum"},
		7:  {"main"},
		8:  {"main"},
		9:  {"main", "manual@chan"},
		10: {"main"},
		11: {"promotions"},
		12: {},
		13: {},
		14: {"main", "chan_self_subscribed", "chan_self_unsub", "chan_other_accepted"},
		15: {"main", "chan_other_nosub", "chan_other_request", "chan_other_accepted"},
	}

	for k, v := range testdata {
		r0 := tt.RequestAuthGet[chanlist](t, data.User[k].AdminKey, baseUrl, fmt.Sprintf("/api/users/%d/channels?selector=%s", data.User[k].UID, "all"))
		tt.AssertMappedSet(t, fmt.Sprintf("%d->chanlist", k), v, r0.Channels, "internal_name")
	}
}

//TODO test missing channel-xx methods