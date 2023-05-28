package test

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/models"
	tt "blackforestbytes.com/simplecloudnotifier/test/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"gogs.mikescher.com/BlackForestBytes/goext/timeext"
	"net/url"
	"testing"
	"time"
)

func TestSearchMessageFTSSimple(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	type mglist struct {
		Messages []gin.H `json:"messages"`
	}

	msgList := tt.RequestAuthGet[mglist](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages?filter=%s", url.QueryEscape("Friday")))
	tt.AssertEqual(t, "msgList.len", 2, len(msgList.Messages))
	tt.AssertArrAny(t, "msgList.any<1>", msgList.Messages, func(msg gin.H) bool { return msg["title"].(string) == "Invitation" })
	tt.AssertArrAny(t, "msgList.any<2>", msgList.Messages, func(msg gin.H) bool { return msg["title"].(string) == "Important notice" })
}

func TestSearchMessageFTSMulti(t *testing.T) {
	t.SkipNow() //TODO search for messages by FTS
}

//TODO more search/list/filter message tests

//TODO list messages by chan_key

//TODO (fail to) list messages from channel that you cannot see

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
