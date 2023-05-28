package test

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	tt "blackforestbytes.com/simplecloudnotifier/test/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"strings"
	"testing"
)

func TestCreateChannel(t *testing.T) {
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

	type chanlist struct {
		Channels []gin.H `json:"channels"`
	}

	{
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", uid))
		tt.AssertMappedSet(t, "channels", []string{}, clist.Channels, "display_name")
		tt.AssertMappedSet(t, "channels", []string{}, clist.Channels, "internal_name")
	}

	tt.RequestAuthPost[gin.H](t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", uid), gin.H{
		"name": "test",
	})

	{
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", uid))
		tt.AssertEqual(t, "chan.len", 1, len(clist.Channels))
		tt.AssertMappedSet(t, "channels", []string{"test"}, clist.Channels, "display_name")
		tt.AssertMappedSet(t, "channels", []string{"test"}, clist.Channels, "internal_name")
	}

	tt.RequestAuthPost[gin.H](t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", uid), gin.H{
		"name": "asdf",
	})

	{
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", uid))
		tt.AssertMappedSet(t, "channels", []string{"asdf", "test"}, clist.Channels, "display_name")
		tt.AssertMappedSet(t, "channels", []string{"asdf", "test"}, clist.Channels, "internal_name")
	}
}

func TestCreateChannelNameTooLong(t *testing.T) {
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

	tt.RequestAuthPostShouldFail(t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", uid), gin.H{
		"name": langext.StrRepeat("X", 121),
	}, 400, apierr.CHANNEL_TOO_LONG)
}

func TestChannelNameNormalization(t *testing.T) {
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

	type chanlist struct {
		Channels []gin.H `json:"channels"`
	}

	{
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", uid))
		tt.AssertMappedSet(t, "channels", []string{}, clist.Channels, "display_name")
		tt.AssertMappedSet(t, "channels", []string{}, clist.Channels, "internal_name")
	}

	tt.RequestAuthPost[gin.H](t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", uid), gin.H{
		"name": "tESt",
	})

	{
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", uid))
		tt.AssertMappedSet(t, "channels", []string{"tESt"}, clist.Channels, "display_name")
		tt.AssertMappedSet(t, "channels", []string{"test"}, clist.Channels, "internal_name")
	}

	tt.RequestAuthPostShouldFail(t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", uid), gin.H{
		"name": "test",
	}, 409, apierr.CHANNEL_ALREADY_EXISTS)

	tt.RequestAuthPostShouldFail(t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", uid), gin.H{
		"name": "TEST",
	}, 409, apierr.CHANNEL_ALREADY_EXISTS)

	tt.RequestAuthPostShouldFail(t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", uid), gin.H{
		"name": "Test",
	}, 409, apierr.CHANNEL_ALREADY_EXISTS)

	tt.RequestAuthPostShouldFail(t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", uid), gin.H{
		"name": "Test ",
	}, 409, apierr.CHANNEL_ALREADY_EXISTS)

	tt.RequestAuthPostShouldFail(t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", uid), gin.H{
		"name": " Test",
	}, 409, apierr.CHANNEL_ALREADY_EXISTS)

	tt.RequestAuthPostShouldFail(t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", uid), gin.H{
		"name": "\rTeSt\n",
	}, 409, apierr.CHANNEL_ALREADY_EXISTS)

	{
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", uid))
		tt.AssertMappedSet(t, "channels", []string{"tESt"}, clist.Channels, "display_name")
		tt.AssertMappedSet(t, "channels", []string{"test"}, clist.Channels, "internal_name")
	}

	tt.RequestAuthPost[gin.H](t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", uid), gin.H{
		"name": "  WeiRD_[\uF5FF]\\stUFf\r\n\t  ",
	})

	{
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", uid))
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
		r0 := tt.RequestAuthGet[chanlist](t, data.User[k].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", data.User[k].UID))
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
		r0 := tt.RequestAuthGet[chanlist](t, data.User[k].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels?selector=%s", data.User[k].UID, "owned"))
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
		r0 := tt.RequestAuthGet[chanlist](t, data.User[k].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels?selector=%s", data.User[k].UID, "subscribed_any"))
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
		r0 := tt.RequestAuthGet[chanlist](t, data.User[k].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels?selector=%s", data.User[k].UID, "all_any"))
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
		r0 := tt.RequestAuthGet[chanlist](t, data.User[k].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels?selector=%s", data.User[k].UID, "subscribed"))
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
		r0 := tt.RequestAuthGet[chanlist](t, data.User[k].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels?selector=%s", data.User[k].UID, "all"))
		tt.AssertMappedSet(t, fmt.Sprintf("%d->chanlist", k), v, r0.Channels, "internal_name")
	}
}

func TestChannelUpdate(t *testing.T) {
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

	type chanlist struct {
		Channels []gin.H `json:"channels"`
	}

	{
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", uid))
		tt.AssertMappedSet(t, "channels", []string{}, clist.Channels, "display_name")
		tt.AssertMappedSet(t, "channels", []string{}, clist.Channels, "internal_name")
	}

	chan0 := tt.RequestAuthPost[gin.H](t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", uid), gin.H{
		"name": "server-alerts",
	})
	chanid := fmt.Sprintf("%v", chan0["channel_id"])

	{
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", uid))
		tt.AssertMappedSet(t, "channels", []string{"server-alerts"}, clist.Channels, "display_name")
		tt.AssertMappedSet(t, "channels", []string{"server-alerts"}, clist.Channels, "internal_name")
		tt.AssertEqual(t, "channels.descr", nil, clist.Channels[0]["description_name"])
	}

	{
		chan1 := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s", uid, chanid))
		tt.AssertEqual(t, "channels.display_name", "server-alerts", chan1["display_name"])
		tt.AssertEqual(t, "channels.internal_name", "server-alerts", chan1["internal_name"])
		tt.AssertEqual(t, "channels.description_name", nil, chan1["description_name"])
		tt.AssertEqual(t, "channels.subscribe_key", chan0["subscribe_key"], chan1["subscribe_key"])
		tt.AssertEqual(t, "channels.send_key", chan0["send_key"], chan1["send_key"])
	}

	// [1] update display_name

	tt.RequestAuthPatch[tt.Void](t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s", uid, chanid), gin.H{
		"display_name": "SERVER-ALERTS",
	})

	{
		chan1 := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s", uid, chanid))
		tt.AssertEqual(t, "channels.display_name", "SERVER-ALERTS", chan1["display_name"])
		tt.AssertEqual(t, "channels.internal_name", "server-alerts", chan1["internal_name"])
		tt.AssertEqual(t, "channels.description_name", nil, chan1["description_name"])
		tt.AssertEqual(t, "channels.subscribe_key", chan0["subscribe_key"], chan1["subscribe_key"])
		tt.AssertEqual(t, "channels.send_key", chan0["send_key"], chan1["send_key"])
	}

	// [2] fail to update display_name

	tt.RequestAuthPatchShouldFail(t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s", uid, chanid), gin.H{
		"display_name": "SERVER-ALERTS2",
	}, 400, apierr.CHANNEL_NAME_WOULD_CHANGE)

	// [3] renew subscribe_key

	tt.RequestAuthPatch[tt.Void](t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s", uid, chanid), gin.H{
		"subscribe_key": true,
	})

	{
		chan1 := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s", uid, chanid))
		tt.AssertNotEqual(t, "channels.subscribe_key", chan0["subscribe_key"], chan1["subscribe_key"])
		tt.AssertEqual(t, "channels.send_key", chan0["send_key"], chan1["send_key"])
	}

	// [5] update description_name

	tt.RequestAuthPatch[tt.Void](t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s", uid, chanid), gin.H{
		"description_name": "hello World",
	})

	{
		chan1 := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s", uid, chanid))
		tt.AssertEqual(t, "channels.description_name", "hello World", chan1["description_name"])
	}

	// [6] update description_name

	tt.RequestAuthPatch[tt.Void](t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s", uid, chanid), gin.H{
		"description_name": "  AXXhello World9  ",
	})

	{
		chan1 := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s", uid, chanid))
		tt.AssertEqual(t, "channels.description_name", "AXXhello World9", chan1["description_name"])
	}

	// [7] clear description_name

	tt.RequestAuthPatch[tt.Void](t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s", uid, chanid), gin.H{
		"description_name": "",
	})

	{
		chan1 := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s", uid, chanid))
		tt.AssertEqual(t, "channels.description_name", nil, chan1["description_name"])
	}

	// [8] fail to update description_name

	tt.RequestAuthPatchShouldFail(t, admintok, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s", uid, chanid), gin.H{
		"description_name": strings.Repeat("0123456789", 48),
	}, 400, apierr.CHANNEL_DESCRIPTION_TOO_LONG)

}

func TestListChannelMessages(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	type msg struct {
		ChannelId           string `json:"channel_id"`
		ChannelInternalName string `json:"channel_internal_name"`
		Content             string `json:"content"`
		MessageId           string `json:"message_id"`
		OwnerUserId         string `json:"owner_user_id"`
		Priority            int    `json:"priority"`
		SenderIp            string `json:"sender_ip"`
		SenderName          string `json:"sender_name"`
		SenderUserId        string `json:"sender_user_id"`
		Timestamp           string `json:"timestamp"`
		Title               string `json:"title"`
		Trimmed             bool   `json:"trimmed"`
		UsrMessageId        string `json:"usr_message_id"`
	}
	type mglist struct {
		Messages []msg  `json:"messages"`
		NPT      string `json:"next_page_token"`
		PageSize int    `json:"page_size"`
	}

	type chanobj struct {
		ChannelId       string `json:"channel_id"`
		DescriptionName string `json:"description_name"`
		DisplayName     string `json:"display_name"`
		InternalName    string `json:"internal_name"`
		MessagesSent    int    `json:"messages_sent"`
		OwnerUserId     string `json:"owner_user_id"`
		SubscribeKey    string `json:"subscribe_key"`
		Subscription    struct {
			ChannelId           string `json:"channel_id"`
			ChannelInternalName string `json:"channel_internal_name"`
			ChannelOwnerUserId  string `json:"channel_owner_user_id"`
			Confirmed           bool   `json:"confirmed"`
			SubscriberUserId    string `json:"subscriber_user_id"`
			SubscriptionId      string `json:"subscription_id"`
			TimestampCreated    string `json:"timestamp_created"`
		} `json:"subscription"`
		TimestampCreated  string `json:"timestamp_created"`
		TimestampLastsent string `json:"timestamp_lastsent"`
	}

	type chanlist struct {
		Channels []chanobj `json:"channels"`
	}

	clist := tt.RequestAuthGet[chanlist](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", data.User[16].UID))

	chan1 := langext.ArrFirstOrNil(clist.Channels, func(v chanobj) bool { return v.InternalName == "chan1" }).ChannelId
	chan2 := langext.ArrFirstOrNil(clist.Channels, func(v chanobj) bool { return v.InternalName == "chan2" }).ChannelId
	chan3 := langext.ArrFirstOrNil(clist.Channels, func(v chanobj) bool { return v.InternalName == "chan3" }).ChannelId

	{
		msgList0 := tt.RequestAuthGet[mglist](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/messages", data.User[16].UID, chan1))
		tt.AssertEqual(t, "msgList.len", 8, len(msgList0.Messages))
		tt.AssertEqual(t, "PageSize", 16, msgList0.PageSize)
		tt.AssertEqual(t, "msgList[0]", "Lorem Ipsum 11", msgList0.Messages[0].Title)
	}

	{
		msgList0 := tt.RequestAuthGet[mglist](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/messages", data.User[16].UID, chan2))
		tt.AssertEqual(t, "msgList.len", 10, len(msgList0.Messages))
		tt.AssertEqual(t, "PageSize", 16, msgList0.PageSize)
		tt.AssertEqual(t, "msgList[0]", "Lorem Ipsum 23", msgList0.Messages[0].Title)
	}

	{
		msgList0 := tt.RequestAuthGet[mglist](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/messages", data.User[16].UID, chan3))
		tt.AssertEqual(t, "msgList.len", 5, len(msgList0.Messages))
		tt.AssertEqual(t, "PageSize", 16, msgList0.PageSize)
		tt.AssertEqual(t, "msgList[0]", "Lorem Ipsum 20", msgList0.Messages[0].Title)
	}
}

func TestListSubscribedChannelMessages(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	type msg struct {
		ChannelId           string `json:"channel_id"`
		ChannelInternalName string `json:"channel_internal_name"`
		Content             string `json:"content"`
		MessageId           string `json:"message_id"`
		OwnerUserId         string `json:"owner_user_id"`
		Priority            int    `json:"priority"`
		SenderIp            string `json:"sender_ip"`
		SenderName          string `json:"sender_name"`
		SenderUserId        string `json:"sender_user_id"`
		Timestamp           string `json:"timestamp"`
		Title               string `json:"title"`
		Trimmed             bool   `json:"trimmed"`
		UsrMessageId        string `json:"usr_message_id"`
	}
	type mglist struct {
		Messages []msg  `json:"messages"`
		NPT      string `json:"next_page_token"`
		PageSize int    `json:"page_size"`
	}

	type chanobj struct {
		ChannelId       string `json:"channel_id"`
		DescriptionName string `json:"description_name"`
		DisplayName     string `json:"display_name"`
		InternalName    string `json:"internal_name"`
		MessagesSent    int    `json:"messages_sent"`
		OwnerUserId     string `json:"owner_user_id"`
		SubscribeKey    string `json:"subscribe_key"`
		Subscription    struct {
			ChannelId           string `json:"channel_id"`
			ChannelInternalName string `json:"channel_internal_name"`
			ChannelOwnerUserId  string `json:"channel_owner_user_id"`
			Confirmed           bool   `json:"confirmed"`
			SubscriberUserId    string `json:"subscriber_user_id"`
			SubscriptionId      string `json:"subscription_id"`
			TimestampCreated    string `json:"timestamp_created"`
		} `json:"subscription"`
		TimestampCreated  string `json:"timestamp_created"`
		TimestampLastsent string `json:"timestamp_lastsent"`
	}

	type chanlist struct {
		Channels []chanobj `json:"channels"`
	}

	clist := tt.RequestAuthGet[chanlist](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", data.User[16].UID))

	chan1 := langext.ArrFirstOrNil(clist.Channels, func(v chanobj) bool { return v.InternalName == "chan1" })
	chan2 := langext.ArrFirstOrNil(clist.Channels, func(v chanobj) bool { return v.InternalName == "chan2" })
	chan3 := langext.ArrFirstOrNil(clist.Channels, func(v chanobj) bool { return v.InternalName == "chan3" })

	{
		sub1 := tt.RequestAuthPost[gin.H](t, data.User[1].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions?chan_subscribe_key=%s", data.User[1].UID, chan1.SubscribeKey), gin.H{
			"channel_owner_user_id": data.User[16].UID,
			"channel_internal_name": "chan1",
		})
		sub2 := tt.RequestAuthPost[gin.H](t, data.User[1].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions?chan_subscribe_key=%s", data.User[1].UID, chan2.SubscribeKey), gin.H{
			"channel_owner_user_id": data.User[16].UID,
			"channel_internal_name": "chan2",
		})
		sub3 := tt.RequestAuthPost[gin.H](t, data.User[1].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions?chan_subscribe_key=%s", data.User[1].UID, chan3.SubscribeKey), gin.H{
			"channel_owner_user_id": data.User[16].UID,
			"channel_internal_name": "chan3",
		})

		tt.RequestAuthPatch[gin.H](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data.User[16].UID, sub1["subscription_id"]), gin.H{
			"confirmed": true,
		})
		tt.RequestAuthPatch[gin.H](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data.User[16].UID, sub2["subscription_id"]), gin.H{
			"confirmed": true,
		})
		tt.RequestAuthPatch[gin.H](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data.User[16].UID, sub3["subscription_id"]), gin.H{
			"confirmed": true,
		})

	}

	{
		msgList0 := tt.RequestAuthGet[mglist](t, data.User[1].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/messages", data.User[16].UID, chan1.ChannelId))
		tt.AssertEqual(t, "msgList.len", 8, len(msgList0.Messages))
		tt.AssertEqual(t, "PageSize", 16, msgList0.PageSize)
		tt.AssertEqual(t, "msgList[0]", "Lorem Ipsum 11", msgList0.Messages[0].Title)
	}

	{
		msgList0 := tt.RequestAuthGet[mglist](t, data.User[1].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/messages", data.User[16].UID, chan2.ChannelId))
		tt.AssertEqual(t, "msgList.len", 10, len(msgList0.Messages))
		tt.AssertEqual(t, "PageSize", 16, msgList0.PageSize)
		tt.AssertEqual(t, "msgList[0]", "Lorem Ipsum 23", msgList0.Messages[0].Title)
	}

	{
		msgList0 := tt.RequestAuthGet[mglist](t, data.User[1].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/messages", data.User[16].UID, chan3.ChannelId))
		tt.AssertEqual(t, "msgList.len", 5, len(msgList0.Messages))
		tt.AssertEqual(t, "PageSize", 16, msgList0.PageSize)
		tt.AssertEqual(t, "msgList[0]", "Lorem Ipsum 20", msgList0.Messages[0].Title)
	}
}

func TestListChannelSubscriptions(t *testing.T) {
	t.SkipNow() //TODO
}
