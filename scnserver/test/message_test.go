package test

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/models"
	tt "blackforestbytes.com/simplecloudnotifier/test/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/timeext"
	"net/url"
	"testing"
	"time"
)

func TestSearchMessageFTSSimple(t *testing.T) {
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
		Messages []msg `json:"messages"`
	}

	msgList := tt.RequestAuthGet[mglist](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?search=%s", url.QueryEscape("Friday")))
	tt.AssertEqual(t, "msgList.len", 2, len(msgList.Messages))
	tt.AssertArrAny(t, "msgList.any<1>", msgList.Messages, func(msg msg) bool { return msg.Title == "Invitation" })
	tt.AssertArrAny(t, "msgList.any<2>", msgList.Messages, func(msg msg) bool { return msg.Title == "Important notice" })
}

func TestSearchMessageFTSMulti(t *testing.T) {
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
		Messages []msg `json:"messages"`
	}

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

	klist := tt.RequestAuthGet[keylist](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", data.User[0].UID))

	akey := langext.ArrFirstOrNil(klist.Keys, func(v keyobj) bool { return v.Permissions == "A" }).KeytokenId
	skey := langext.ArrFirstOrNil(klist.Keys, func(v keyobj) bool { return v.Permissions == "CS" }).KeytokenId

	{
		msgList := tt.RequestAuthGet[mglist](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?sender=%s&priority=0", url.QueryEscape("Mobile Mate")))
		tt.AssertEqual(t, "msgList.len", 1, len(msgList.Messages))
		tt.AssertEqual(t, "msg.Title", "System update", msgList.Messages[0].Title)
	}

	{
		msgList := tt.RequestAuthGet[mglist](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?used_key=%s&priority=%d&channel=%s", skey, 2, "Reminders"))
		tt.AssertEqual(t, "msgList.len", 1, len(msgList.Messages))
		tt.AssertEqual(t, "msg.Content", "Don't forget to clock in before starting your shift", msgList.Messages[0].Content)
	}

	{
		msgList := tt.RequestAuthGet[mglist](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?used_key=%s&used_key=%s&channel=%s&search=%s", skey, akey, "main", url.QueryEscape("tomorrow")))
		tt.AssertEqual(t, "msgList.len", 1, len(msgList.Messages))
		tt.AssertEqual(t, "msg.Title", "Notice", msgList.Messages[0].Title)
	}
}

func TestListMessagesFilteredChannels(t *testing.T) {
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
		Messages []msg `json:"messages"`
	}

	{
		msgList := tt.RequestAuthGet[mglist](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?channel=%s", url.QueryEscape("Chatting Chamber")))
		tt.AssertEqual(t, "msgList.len", 3, len(msgList.Messages))
	}

	{
		msgList := tt.RequestAuthGet[mglist](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?channel=%s&channel=%s", url.QueryEscape("Chatting Chamber"), url.QueryEscape("Reminders")))
		tt.AssertEqual(t, "msgList.len", 3+6, len(msgList.Messages))
	}
}

func TestListMessagesFilteredChannelIDs(t *testing.T) {
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
		Messages []msg `json:"messages"`
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

	clist := tt.RequestAuthGet[chanlist](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", data.User[0].UID))

	chan1 := langext.ArrFirstOrNil(clist.Channels, func(v chanobj) bool { return v.InternalName == "Chatting Chamber" }).ChannelId
	chan2 := langext.ArrFirstOrNil(clist.Channels, func(v chanobj) bool { return v.InternalName == "Reminders" }).ChannelId

	{
		msgList := tt.RequestAuthGet[mglist](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?channel_id=%s", chan1))
		tt.AssertEqual(t, "msgList.len", 3, len(msgList.Messages))
	}

	{
		msgList := tt.RequestAuthGet[mglist](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?channel_id=%s&channel_id=%s", chan1, chan2))
		tt.AssertEqual(t, "msgList.len", 3+6, len(msgList.Messages))
	}
}

func TestListMessagesFilteredSenders(t *testing.T) {
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
		Messages []msg `json:"messages"`
	}

	{
		msgList := tt.RequestAuthGet[mglist](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?sender=%s", url.QueryEscape("Mobile Mate")))
		tt.AssertEqual(t, "msgList.len", 4, len(msgList.Messages))
	}

	{
		msgList := tt.RequestAuthGet[mglist](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?sender=%s&sender=%s", url.QueryEscape("Mobile Mate"), url.QueryEscape("Pocket Pal")))
		tt.AssertEqual(t, "msgList.len", 4+3, len(msgList.Messages))
	}
}

func TestListMessagesFilteredTime(t *testing.T) {
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
		Messages []msg `json:"messages"`
	}

	{
		msgList := tt.RequestAuthGet[mglist](t, data.User[4].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages"))
		tt.AssertEqual(t, "msgList.len", 9, len(msgList.Messages))
	}
	{
		msgList := tt.RequestAuthGet[mglist](t, data.User[4].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?before=%s", url.QueryEscape(time.Now().Add(time.Minute).Format(time.RFC3339))))
		tt.AssertEqual(t, "msgList.len", 7, len(msgList.Messages))
	}
	{
		msgList := tt.RequestAuthGet[mglist](t, data.User[4].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?before=%s", url.QueryEscape(time.Now().Add(-1*time.Hour).Format(time.RFC3339))))
		tt.AssertEqual(t, "msgList.len", 1, len(msgList.Messages))
	}
	{
		msgList := tt.RequestAuthGet[mglist](t, data.User[4].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?after=%s", url.QueryEscape(time.Now().Add(+1*time.Minute).Format(time.RFC3339))))
		tt.AssertEqual(t, "msgList.len", 2, len(msgList.Messages))
	}
	{
		msgList := tt.RequestAuthGet[mglist](t, data.User[4].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?after=%s&before=%s", url.QueryEscape(time.Now().Add(-1*time.Minute).Format(time.RFC3339)), url.QueryEscape(time.Now().Add(+1*time.Minute).Format(time.RFC3339))))
		tt.AssertEqual(t, "msgList.len", 6, len(msgList.Messages))
	}
	{
		msgList := tt.RequestAuthGet[mglist](t, data.User[4].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?after=%s&before=%s", url.QueryEscape(time.Now().Add(-5*time.Hour).Format(time.RFC3339)), url.QueryEscape(time.Now().Add(-4*time.Hour).Format(time.RFC3339))))
		tt.AssertEqual(t, "msgList.len", 1, len(msgList.Messages))
	}
	{
		msgList := tt.RequestAuthGet[mglist](t, data.User[4].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?after=%s&before=%s", url.QueryEscape(time.Now().Add(6*time.Hour).Format(time.RFC3339)), url.QueryEscape(time.Now().Add(7*time.Hour).Format(time.RFC3339))))
		tt.AssertEqual(t, "msgList.len", 2, len(msgList.Messages))
	}
}

func TestListMessagesFilteredPriority(t *testing.T) {
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
		Messages []msg `json:"messages"`
	}

	{
		msgList := tt.RequestAuthGet[mglist](t, data.User[4].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?priority=%d", 0))
		tt.AssertEqual(t, "msgList.len", 6, len(msgList.Messages))
	}

	{
		msgList := tt.RequestAuthGet[mglist](t, data.User[4].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?priority=%d", 1))
		tt.AssertEqual(t, "msgList.len", 1+1, len(msgList.Messages))
	}

	{
		msgList := tt.RequestAuthGet[mglist](t, data.User[4].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?priority=%d", 2))
		tt.AssertEqual(t, "msgList.len", 1, len(msgList.Messages))
	}

	{
		msgList := tt.RequestAuthGet[mglist](t, data.User[4].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?priority=%d&priority=%d", 1, 2))
		tt.AssertEqual(t, "msgList.len", 2+1, len(msgList.Messages))
	}

}

func TestListMessagesFilteredKeyTokens(t *testing.T) {
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
		Messages []msg `json:"messages"`
	}

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

	klist := tt.RequestAuthGet[keylist](t, data.User[4].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", data.User[4].UID))

	akey := langext.ArrFirstOrNil(klist.Keys, func(v keyobj) bool { return v.Permissions == "A" }).KeytokenId
	skey := langext.ArrFirstOrNil(klist.Keys, func(v keyobj) bool { return v.Permissions == "CS" }).KeytokenId

	{
		msgList := tt.RequestAuthGet[mglist](t, data.User[4].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?used_key=%s", skey))
		tt.AssertEqual(t, "msgList.len", 7, len(msgList.Messages))
	}

	{
		msgList := tt.RequestAuthGet[mglist](t, data.User[4].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?used_key=%s", akey))
		tt.AssertEqual(t, "msgList.len", 2, len(msgList.Messages))
	}

	{
		msgList := tt.RequestAuthGet[mglist](t, data.User[4].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?used_key=%s&used_key=%s", skey, akey))
		tt.AssertEqual(t, "msgList.len", 9, len(msgList.Messages))
	}
}

func TestDeleteMessage(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/v2/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := r0["user_id"].(string)
	sendtok := r0["send_key"].(string)
	admintok := r0["admin_key"].(string)

	msg1 := tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"key":     sendtok,
		"user_id": uid,
		"title":   "Message_1",
	})

	tt.RequestAuthGet[tt.Void](t, admintok, baseUrl, "/api/v2/messages/"+fmt.Sprintf("%v", msg1["scn_msg_id"]))

	tt.RequestAuthDelete[tt.Void](t, admintok, baseUrl, "/api/v2/messages/"+fmt.Sprintf("%v", msg1["scn_msg_id"]), gin.H{})

	tt.RequestAuthGetShouldFail(t, admintok, baseUrl, "/api/v2/messages/"+fmt.Sprintf("%v", msg1["scn_msg_id"]), 404, apierr.MESSAGE_NOT_FOUND)
}

func TestDeleteMessageAndResendUsrMsgId(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/v2/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := r0["user_id"].(string)
	sendtok := r0["send_key"].(string)
	admintok := r0["admin_key"].(string)

	msg1 := tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"key":     sendtok,
		"user_id": uid,
		"title":   "Message_1",
		"msg_id":  "bef8dd3d-078e-4f89-abf4-5258ad22a2e4",
	})

	tt.AssertEqual(t, "suppress_send", false, msg1["suppress_send"])

	tt.RequestAuthGet[tt.Void](t, admintok, baseUrl, "/api/v2/messages/"+fmt.Sprintf("%v", msg1["scn_msg_id"]))

	msg2 := tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"key":     sendtok,
		"user_id": uid,
		"title":   "Message_1",
		"msg_id":  "bef8dd3d-078e-4f89-abf4-5258ad22a2e4",
	})

	tt.AssertEqual(t, "suppress_send", true, msg2["suppress_send"])

	tt.RequestAuthDelete[tt.Void](t, admintok, baseUrl, "/api/v2/messages/"+fmt.Sprintf("%v", msg1["scn_msg_id"]), gin.H{})

	// even though message is deleted, we still get a `suppress_send` on send_message

	msg3 := tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"key":     sendtok,
		"user_id": uid,
		"title":   "Message_1",
		"msg_id":  "bef8dd3d-078e-4f89-abf4-5258ad22a2e4",
	})

	tt.AssertEqual(t, "suppress_send", true, msg3["suppress_send"])

}

func TestGetMessageSimple(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	msgOut := tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"key":     data.User[0].SendKey,
		"user_id": data.User[0].UID,
		"title":   "Message_1",
	})

	msgIn := tt.RequestAuthGet[gin.H](t, data.User[0].AdminKey, baseUrl, "/api/v2/messages/"+fmt.Sprintf("%v", msgOut["scn_msg_id"]))

	tt.AssertEqual(t, "msg.title", "Message_1", msgIn["title"])
}

func TestGetMessageNotFound(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	tt.RequestAuthGetShouldFail(t, data.User[0].AdminKey, baseUrl, "/api/v2/messages/"+models.NewMessageID().String(), 404, apierr.MESSAGE_NOT_FOUND)
}

func TestGetMessageInvalidID(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	tt.RequestAuthGetShouldFail(t, data.User[0].AdminKey, baseUrl, "/api/v2/messages/"+models.NewUserID().String(), 400, apierr.BINDFAIL_URI_PARAM)

	tt.RequestAuthGetShouldFail(t, data.User[0].AdminKey, baseUrl, "/api/v2/messages/"+"asdfxxx", 400, apierr.BINDFAIL_URI_PARAM)
}

func TestGetMessageFull(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	ts := time.Now().Unix() - 735
	content := tt.ShortLipsum0(2)

	msgOut := tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"key":         data.User[0].SendKey,
		"user_id":     data.User[0].UID,
		"title":       "Message_1",
		"content":     content,
		"channel":     "demo-channel-007",
		"msg_id":      "580b5055-a9b5-4cee-b53c-28cf304d25b0",
		"priority":    0,
		"sender_name": "unit-test-[TestGetMessageFull]",
		"timestamp":   ts,
	})

	msgIn := tt.RequestAuthGet[gin.H](t, data.User[0].AdminKey, baseUrl, "/api/v2/messages/"+fmt.Sprintf("%v", msgOut["scn_msg_id"]))

	tt.AssertEqual(t, "msg.title", "Message_1", msgIn["title"])
	tt.AssertEqual(t, "msg.content", content, msgIn["content"])
	tt.AssertEqual(t, "msg.channel", "demo-channel-007", msgIn["channel_internal_name"])
	tt.AssertEqual(t, "msg.msg_id", "580b5055-a9b5-4cee-b53c-28cf304d25b0", msgIn["usr_message_id"])
	tt.AssertStrRepEqual(t, "msg.priority", 0, msgIn["priority"])
	tt.AssertEqual(t, "msg.sender_name", "unit-test-[TestGetMessageFull]", msgIn["sender_name"])
	tt.AssertEqual(t, "msg.timestamp", time.Unix(ts, 0).In(timeext.TimezoneBerlin).Format(time.RFC3339Nano), msgIn["timestamp"])
}

func TestListMessages(t *testing.T) {
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
		Messages []msg `json:"messages"`
	}

	msgList := tt.RequestAuthGet[mglist](t, data.User[7].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages"))
	tt.AssertEqual(t, "msgList.len", 6, len(msgList.Messages))
	tt.AssertEqual(t, "msgList[0]", "Server outage status", msgList.Messages[0].Title)
	tt.AssertEqual(t, "msgList[1]", "Server maintenance reminder", msgList.Messages[1].Title)
	tt.AssertEqual(t, "msgList[2]", "Server security alert", msgList.Messages[2].Title)
	tt.AssertEqual(t, "msgList[3]", "Server traffic warning", msgList.Messages[3].Title)
	tt.AssertEqual(t, "msgList[4]", "New server release update", msgList.Messages[4].Title)
	tt.AssertEqual(t, "msgList[5]", "Server outage resolution update", msgList.Messages[5].Title)
}

func TestListMessagesPaginated(t *testing.T) {
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
	{
		msgList0 := tt.RequestAuthGet[mglist](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages"))
		tt.AssertEqual(t, "msgList.len", 16, len(msgList0.Messages))
		tt.AssertEqual(t, "msgList.len", 16, msgList0.PageSize)
	}
	npt := ""
	{
		msgList1 := tt.RequestAuthGet[mglist](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?page_size=%d&next_page_token=%s", 10, "@start"))
		tt.AssertEqual(t, "msgList.len", 10, len(msgList1.Messages))
		tt.AssertEqual(t, "msgList.PageSize", 10, msgList1.PageSize)
		tt.AssertEqual(t, "msgList[0]", "Lorem Ipsum 23", msgList1.Messages[0].Title)
	}
	{
		msgList1 := tt.RequestAuthGet[mglist](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?page_size=%d&next_page_token=%s", 10, "@START"))
		tt.AssertEqual(t, "msgList.len", 10, len(msgList1.Messages))
		tt.AssertEqual(t, "msgList.PageSize", 10, msgList1.PageSize)
		tt.AssertEqual(t, "msgList[0]", "Lorem Ipsum 23", msgList1.Messages[0].Title)
	}
	{
		msgList1 := tt.RequestAuthGet[mglist](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?page_size=%d&next_page_token=%s", 10, ""))
		tt.AssertEqual(t, "msgList.len", 10, len(msgList1.Messages))
		tt.AssertEqual(t, "msgList.PageSize", 10, msgList1.PageSize)
		tt.AssertEqual(t, "msgList[0]", "Lorem Ipsum 23", msgList1.Messages[0].Title)
		npt = msgList1.NPT
	}
	{
		msgList2 := tt.RequestAuthGet[mglist](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?page_size=%d&next_page_token=%s", 10, npt))
		tt.AssertEqual(t, "msgList.len", 10, len(msgList2.Messages))
		tt.AssertEqual(t, "msgList.PageSize", 10, msgList2.PageSize)
		tt.AssertEqual(t, "msgList[0]", "Lorem Ipsum 13", msgList2.Messages[0].Title)
		npt = msgList2.NPT
	}
	{
		msgList3 := tt.RequestAuthGet[mglist](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?page_size=%d&next_page_token=%s", 10, npt))
		tt.AssertEqual(t, "msgList.len", 3, len(msgList3.Messages))
		tt.AssertEqual(t, "msgList.PageSize", 10, msgList3.PageSize)
		tt.AssertEqual(t, "msgList[0]", "Lorem Ipsum 03", msgList3.Messages[0].Title)
		tt.AssertEqual(t, "msgList[0]", "@end", msgList3.NPT)
		npt = msgList3.NPT
	}
}

func TestListMessagesPaginatedInvalid(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	tt.RequestAuthGetShouldFail(t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?page_size=%d&next_page_token=%s", 10, "INVALID"), 400, apierr.PAGETOKEN_ERROR)
}

func TestListMessagesZeroPagesize(t *testing.T) {
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

	msgList1 := tt.RequestAuthGet[mglist](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?page_size=%d&next_page_token=%s", 0, "@start"))
	tt.AssertEqual(t, "msgList.len", 1, len(msgList1.Messages))
	tt.AssertEqual(t, "msgList.PageSize", 1, msgList1.PageSize)
	tt.AssertEqual(t, "msgList[0]", "Lorem Ipsum 23", msgList1.Messages[0].Title)
}

func TestListMessagesFiltered(t *testing.T) {
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
		Messages   []msg `json:"messages"`
		TotalCount int   `json:"total_count"`
	}

	cid1 := ""
	for _, channel := range data.User[0].Channels {
		if channel.InternalName == "Reminders" {
			cid1 = channel.ChannelID
		}
	}

	cid2 := ""
	for _, channel := range data.User[0].Channels {
		if channel.InternalName == "Chatting Chamber" {
			cid2 = channel.ChannelID
		}
	}

	skey := ""
	for _, key := range data.User[0].Keys {
		if key.Name == "SendKey (default)" {
			skey = key.KeyID
		}
	}

	akey := ""
	for _, key := range data.User[0].Keys {
		if key.Name == "AdminKey (default)" {
			akey = key.KeyID
		}
	}

	filterTests := []struct {
		Name  string
		Count int
		Query string
	}{
		{"all", 22, fmt.Sprintf("/api/v2/messages")},
		{"channel=Reminders|Promotions", 9, fmt.Sprintf("/api/v2/messages?channel=%s&channel=%s", "Reminders", "Promotions")},
		{"channel=Reminders", 6, fmt.Sprintf("/api/v2/messages?channel=%s", "Reminders")},
		{"channel_id=1", 6, fmt.Sprintf("/api/v2/messages?channel_id=%s", cid1)},
		{"channel_id=1|2", 9, fmt.Sprintf("/api/v2/messages?channel_id=%s&channel_id=%s", cid1, cid2)},
		{"search=unusual", 1, fmt.Sprintf("/api/v2/messages?search=%s", "unusual")},
		{"search=your", 6, fmt.Sprintf("/api/v2/messages?search=%s", "your")},
		{"prio=0", 5, fmt.Sprintf("/api/v2/messages?priority=%s", "0")},
		{"prio=1", 4 + 7, fmt.Sprintf("/api/v2/messages?priority=%s", "1")},
		{"prio=2", 6, fmt.Sprintf("/api/v2/messages?priority=%s", "2")},
		{"prio=0|2", 5 + 6, fmt.Sprintf("/api/v2/messages?priority=%s&priority=%s", "0", "2")},
		{"key=a", 11, fmt.Sprintf("/api/v2/messages?used_key=%s", akey)},
		{"key=s", 11, fmt.Sprintf("/api/v2/messages?used_key=%s", skey)},
		{"key=a|s", 11 + 11, fmt.Sprintf("/api/v2/messages?used_key=%s&used_key=%s", akey, skey)},
		{"key=a&prio=0", 0, fmt.Sprintf("/api/v2/messages?used_key=%s&priority=%d", akey, 0)},
		{"key=s&prio=0", 5, fmt.Sprintf("/api/v2/messages?used_key=%s&priority=%d", skey, 0)},
		{"key=s&prio=2", 6, fmt.Sprintf("/api/v2/messages?used_key=%s&priority=%d", skey, 2)},
		{"sender=MobileMate", 4, fmt.Sprintf("/api/v2/messages?sender=%s", url.QueryEscape("Mobile Mate"))},
		{"sender=PocketPal", 3, fmt.Sprintf("/api/v2/messages?sender=%s", url.QueryEscape("Pocket Pal"))},
		{"sender=MobileMate|PocketPal", 3 + 4, fmt.Sprintf("/api/v2/messages?sender=%s&sender=%s", url.QueryEscape("Pocket Pal"), url.QueryEscape("Mobile Mate"))},
		{"sender=empty", 12, fmt.Sprintf("/api/v2/messages?has_sender=%s", "false")},
		{"sender=any", 10, fmt.Sprintf("/api/v2/messages?has_sender=%s", "true")},
		{"before=-1H", 2, fmt.Sprintf("/api/v2/messages?before=%s", url.QueryEscape(time.Now().Add(-time.Hour).Format(time.RFC3339Nano)))},
		{"after=-1H", 20, fmt.Sprintf("/api/v2/messages?after=%s", url.QueryEscape(time.Now().Add(-time.Hour).Format(time.RFC3339Nano)))},
		{"after=+5min", 3, fmt.Sprintf("/api/v2/messages?after=%s", url.QueryEscape(time.Now().Add(5*time.Minute).Format(time.RFC3339Nano)))},
	}

	for _, testdata := range filterTests {
		msgList := tt.RequestAuthGet[mglist](t, data.User[0].AdminKey, baseUrl, testdata.Query)
		tt.AssertEqual(t, "msgList.filter["+testdata.Name+"].len", testdata.Count, msgList.TotalCount)
	}
}

func TestListMessagesSearch(t *testing.T) {
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
		Messages   []msg `json:"messages"`
		TotalCount int   `json:"total_count"`
	}

	filterTests := []struct {
		Name  string
		Count int
		Query string
	}{
		{"all", 22, fmt.Sprintf("/api/v2/messages")},

		{"search=Promotions", 3, fmt.Sprintf("/api/v2/messages?search=%s", url.QueryEscape("Promotions"))},
		{"search=Important(1)", 3, fmt.Sprintf("/api/v2/messages?search=%s", url.QueryEscape("Important"))},
		{"search=Important(2)", 3, fmt.Sprintf("/api/v2/messages?search=%s", url.QueryEscape("important"))},
		{"search=Important(3)", 3, fmt.Sprintf("/api/v2/messages?search=%s", url.QueryEscape("IMPORTANT"))},
		{"search=safetyTraining(1)", 1, fmt.Sprintf("/api/v2/messages?search=%s", url.QueryEscape("safety training"))},
		{"search=safetyTraining(2)", 1, fmt.Sprintf("/api/v2/messages?search=%s", url.QueryEscape("\"safety training\""))},
		{"search=staffMeeting(1)", 2, fmt.Sprintf("/api/v2/messages?search=%s", url.QueryEscape("staff meeting"))},
		{"search=staffMeeting(2)", 1, fmt.Sprintf("/api/v2/messages?search=%s", url.QueryEscape("\"staff meeting\""))},
		{"search=?", 0, fmt.Sprintf("/api/v2/messages?search=%s", url.QueryEscape("\"?\""))},   // fails because FTS searches for full words
		{"search=Prom", 0, fmt.Sprintf("/api/v2/messages?search=%s", url.QueryEscape("Prom"))}, // fails because FTS searches for full words
		{"search=the(1)", 17, fmt.Sprintf("/api/v2/messages?search=%s", url.QueryEscape("the"))},
		{"search=the(2)", 17, fmt.Sprintf("/api/v2/messages?search=%s", url.QueryEscape("THE"))},
		{"search=please", 9, fmt.Sprintf("/api/v2/messages?search=%s", url.QueryEscape("please"))},
		{"search=11pm", 2, fmt.Sprintf("/api/v2/messages?search=%s", url.QueryEscape("\"11:00pm\""))},
	}

	for _, testdata := range filterTests {
		msgList := tt.RequestAuthGet[mglist](t, data.User[0].AdminKey, baseUrl, testdata.Query)
		tt.AssertEqual(t, "msgList.filter["+testdata.Name+"].len", testdata.Count, msgList.TotalCount)
	}
}

func TestListMessagesStringSearch(t *testing.T) {
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
		Messages   []msg `json:"messages"`
		TotalCount int   `json:"total_count"`
	}

	filterTests := []struct {
		Name  string
		Count int
		Query string
	}{
		{"all", 22, fmt.Sprintf("/api/v2/messages")},

		{"search=Promotions", 3, fmt.Sprintf("/api/v2/messages?string_search=%s", url.QueryEscape("Promotions"))},
		{"search=Important(1)", 3, fmt.Sprintf("/api/v2/messages?string_search=%s", url.QueryEscape("Important"))},
		{"search=Important(2)", 3, fmt.Sprintf("/api/v2/messages?string_search=%s", url.QueryEscape("important"))},
		{"search=Important(3)", 3, fmt.Sprintf("/api/v2/messages?string_search=%s", url.QueryEscape("IMPORTANT"))},
		{"search=safetyTraining", 1, fmt.Sprintf("/api/v2/messages?string_search=%s", url.QueryEscape("safety training"))},
		{"search=?", 1, fmt.Sprintf("/api/v2/messages?string_search=%s", url.QueryEscape("?"))},
		{"search=the(1)", 17, fmt.Sprintf("/api/v2/messages?string_search=%s", url.QueryEscape("the"))},
		{"search=the(2)", 17, fmt.Sprintf("/api/v2/messages?string_search=%s", url.QueryEscape("THE"))},
		{"search=please", 9, fmt.Sprintf("/api/v2/messages?string_search=%s", url.QueryEscape("please"))},
		{"search=there", 3, fmt.Sprintf("/api/v2/messages?string_search=%s", url.QueryEscape("tHERe"))},
		{"search=11pm", 2, fmt.Sprintf("/api/v2/messages?string_search=%s", url.QueryEscape("11:00pm"))},

		{"search=Prom", 3, fmt.Sprintf("/api/v2/messages?string_search=%s", url.QueryEscape("Prom"))},
		{"search=run", 1, fmt.Sprintf("/api/v2/messages?string_search=%s", url.QueryEscape("run"))},

		{"search=please+there", 10, fmt.Sprintf("/api/v2/messages?string_search=%s&string_search=%s", url.QueryEscape("please"), url.QueryEscape("THERE"))},
	}

	for _, testdata := range filterTests {
		msgList := tt.RequestAuthGet[mglist](t, data.User[0].AdminKey, baseUrl, testdata.Query)
		tt.AssertEqual(t, "msgList.filter["+testdata.Name+"].len", testdata.Count, msgList.TotalCount)
	}
}
