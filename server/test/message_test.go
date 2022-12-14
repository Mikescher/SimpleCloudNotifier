package test

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
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
	ws, stop := tt.StartSimpleWebserver(t)
	defer stop()

	baseUrl := "http://127.0.0.1:" + ws.Port

	data := tt.InitDefaultData(t, ws)

	type mglist struct {
		Messages []gin.H `json:"messages"`
	}

	msgList := tt.RequestAuthGet[mglist](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/messages?filter=%s", url.QueryEscape("Friday")))
	tt.AssertEqual(t, "msgList.len", 2, len(msgList.Messages))
	tt.AssertTrue(t, "msgList.any<1>", langext.ArrAny(msgList.Messages, func(msg gin.H) bool { return msg["title"].(string) == "Invitation" }))
	tt.AssertTrue(t, "msgList.any<2>", langext.ArrAny(msgList.Messages, func(msg gin.H) bool { return msg["title"].(string) == "Important notice" }))
}

func TestSearchMessageFTSMulti(t *testing.T) {
	//TODO search for messages by FTS
}

//TODO more search/list/filter message tests

//TODO list messages by chan_key

//TODO list messages from channel that you cannot see

func TestDeleteMessage(t *testing.T) {
	ws, stop := tt.StartSimpleWebserver(t)
	defer stop()

	baseUrl := "http://127.0.0.1:" + ws.Port

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := int(r0["user_id"].(float64))
	sendtok := r0["send_key"].(string)
	admintok := r0["admin_key"].(string)

	msg1 := tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    "Message_1",
	})

	tt.RequestAuthGet[tt.Void](t, admintok, baseUrl, "/api/messages/"+fmt.Sprintf("%v", msg1["scn_msg_id"]))

	tt.RequestAuthDelete[tt.Void](t, admintok, baseUrl, "/api/messages/"+fmt.Sprintf("%v", msg1["scn_msg_id"]), gin.H{})

	tt.RequestAuthGetShouldFail(t, admintok, baseUrl, "/api/messages/"+fmt.Sprintf("%v", msg1["scn_msg_id"]), 404, apierr.MESSAGE_NOT_FOUND)
}

func TestDeleteMessageAndResendUsrMsgId(t *testing.T) {
	ws, stop := tt.StartSimpleWebserver(t)
	defer stop()

	baseUrl := "http://127.0.0.1:" + ws.Port

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := int(r0["user_id"].(float64))
	sendtok := r0["send_key"].(string)
	admintok := r0["admin_key"].(string)

	msg1 := tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    "Message_1",
		"msg_id":   "bef8dd3d-078e-4f89-abf4-5258ad22a2e4",
	})

	tt.AssertEqual(t, "suppress_send", false, msg1["suppress_send"])

	tt.RequestAuthGet[tt.Void](t, admintok, baseUrl, "/api/messages/"+fmt.Sprintf("%v", msg1["scn_msg_id"]))

	msg2 := tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    "Message_1",
		"msg_id":   "bef8dd3d-078e-4f89-abf4-5258ad22a2e4",
	})

	tt.AssertEqual(t, "suppress_send", true, msg2["suppress_send"])

	tt.RequestAuthDelete[tt.Void](t, admintok, baseUrl, "/api/messages/"+fmt.Sprintf("%v", msg1["scn_msg_id"]), gin.H{})

	// even though message is deleted, we still get a `suppress_send` on send_message

	msg3 := tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    "Message_1",
		"msg_id":   "bef8dd3d-078e-4f89-abf4-5258ad22a2e4",
	})

	tt.AssertEqual(t, "suppress_send", true, msg3["suppress_send"])

}

func TestGetMessageSimple(t *testing.T) {
	ws, stop := tt.StartSimpleWebserver(t)
	defer stop()

	baseUrl := "http://127.0.0.1:" + ws.Port

	data := tt.InitDefaultData(t, ws)

	msgOut := tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"user_key": data.User[0].SendKey,
		"user_id":  data.User[0].UID,
		"title":    "Message_1",
	})

	msgIn := tt.RequestAuthGet[gin.H](t, data.User[0].AdminKey, baseUrl, "/api/messages/"+fmt.Sprintf("%v", msgOut["scn_msg_id"]))

	tt.AssertEqual(t, "msg.title", "Message_1", msgIn["title"])
}

func TestGetMessageNotFound(t *testing.T) {
	ws, stop := tt.StartSimpleWebserver(t)
	defer stop()

	baseUrl := "http://127.0.0.1:" + ws.Port

	data := tt.InitDefaultData(t, ws)

	tt.RequestAuthGetShouldFail(t, data.User[0].AdminKey, baseUrl, "/api/messages/8963586", 404, apierr.MESSAGE_NOT_FOUND)
}

func TestGetMessageFull(t *testing.T) {
	ws, stop := tt.StartSimpleWebserver(t)
	defer stop()

	baseUrl := "http://127.0.0.1:" + ws.Port

	data := tt.InitDefaultData(t, ws)

	ts := time.Now().Unix() - 735
	content := tt.Lipsum0(2)

	msgOut := tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"user_key":    data.User[0].SendKey,
		"user_id":     data.User[0].UID,
		"title":       "Message_1",
		"content":     content,
		"channel":     "demo-channel-007",
		"msg_id":      "580b5055-a9b5-4cee-b53c-28cf304d25b0",
		"priority":    0,
		"sender_name": "unit-test-[TestGetMessageFull]",
		"timestamp":   ts,
	})

	msgIn := tt.RequestAuthGet[gin.H](t, data.User[0].AdminKey, baseUrl, "/api/messages/"+fmt.Sprintf("%v", msgOut["scn_msg_id"]))

	tt.AssertEqual(t, "msg.title", "Message_1", msgIn["title"])
	tt.AssertEqual(t, "msg.content", content, msgIn["content"])
	tt.AssertEqual(t, "msg.channel", "demo-channel-007", msgIn["channel_name"])
	tt.AssertEqual(t, "msg.msg_id", "580b5055-a9b5-4cee-b53c-28cf304d25b0", msgIn["usr_message_id"])
	tt.AssertStrRepEqual(t, "msg.priority", 0, msgIn["priority"])
	tt.AssertEqual(t, "msg.sender_name", "unit-test-[TestGetMessageFull]", msgIn["sender_name"])
	tt.AssertEqual(t, "msg.timestamp", time.Unix(ts, 0).In(timeext.TimezoneBerlin).Format(time.RFC3339Nano), msgIn["timestamp"])
}
