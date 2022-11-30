package test

import (
	"blackforestbytes.com/simplecloudnotifier/push"
	tt "blackforestbytes.com/simplecloudnotifier/test/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
	"testing"
)

func TestSendSimpleMessageJSON(t *testing.T) {
	ws, stop := tt.StartSimpleWebserver(t)
	defer stop()

	pusher := ws.Pusher.(*push.TestSink)

	baseUrl := "http://127.0.0.1:" + ws.Port

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := int(r0["user_id"].(float64))
	admintok := r0["admin_key"].(string)
	readtok := r0["read_key"].(string)
	sendtok := r0["send_key"].(string)

	msg1 := tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    "HelloWorld_001",
	})

	tt.RequestPostShouldFail(t, baseUrl, "/", gin.H{
		"user_key": readtok,
		"user_id":  uid,
		"title":    "HelloWorld_001",
	}, 401, 1311)

	tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.title", "HelloWorld_001", pusher.Last().Message.Title)
	tt.AssertStrRepEqual(t, "msg.content", nil, pusher.Last().Message.Content)
	tt.AssertStrRepEqual(t, "msg.scn_msg_id", msg1["scn_msg_id"], pusher.Last().Message.SCNMessageID)

	type mglist struct {
		Messages []gin.H `json:"messages"`
	}

	msgList1 := tt.RequestAuthGet[mglist](t, admintok, baseUrl, "/api/messages")
	tt.AssertEqual(t, "len(messages)", 1, len(msgList1.Messages))

	msg1Get := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/messages/"+fmt.Sprintf("%v", msg1["scn_msg_id"]))
	tt.AssertStrRepEqual(t, "msg.title", "HelloWorld_001", msg1Get["title"])
	tt.AssertStrRepEqual(t, "msg.channel_name", "main", msg1Get["channel_name"])
}

func TestSendSimpleMessageQuery(t *testing.T) {
	ws, stop := tt.StartSimpleWebserver(t)
	defer stop()

	pusher := ws.Pusher.(*push.TestSink)

	baseUrl := "http://127.0.0.1:" + ws.Port

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := int(r0["user_id"].(float64))
	admintok := r0["admin_key"].(string)
	sendtok := r0["send_key"].(string)

	msg1 := tt.RequestPost[gin.H](t, baseUrl, fmt.Sprintf("/?user_id=%d&user_key=%s&title=%s", uid, sendtok, url.QueryEscape("Hello World 2134")), nil)

	tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.title", "Hello World 2134", pusher.Last().Message.Title)
	tt.AssertStrRepEqual(t, "msg.content", nil, pusher.Last().Message.Content)
	tt.AssertStrRepEqual(t, "msg.scn_msg_id", msg1["scn_msg_id"], pusher.Last().Message.SCNMessageID)

	type mglist struct {
		Messages []gin.H `json:"messages"`
	}

	msgList1 := tt.RequestAuthGet[mglist](t, admintok, baseUrl, "/api/messages")
	tt.AssertEqual(t, "len(messages)", 1, len(msgList1.Messages))

	msg1Get := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/messages/"+fmt.Sprintf("%v", msg1["scn_msg_id"]))
	tt.AssertStrRepEqual(t, "msg.title", "Hello World 2134", msg1Get["title"])
	tt.AssertStrRepEqual(t, "msg.channel_name", "main", msg1Get["channel_name"])
}

func TestSendSimpleMessageForm(t *testing.T) {
	ws, stop := tt.StartSimpleWebserver(t)
	defer stop()

	pusher := ws.Pusher.(*push.TestSink)

	baseUrl := "http://127.0.0.1:" + ws.Port

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := int(r0["user_id"].(float64))
	admintok := r0["admin_key"].(string)
	sendtok := r0["send_key"].(string)

	msg1 := tt.RequestPost[gin.H](t, baseUrl, "/", tt.FormData{
		"user_key": sendtok,
		"user_id":  fmt.Sprintf("%d", uid),
		"title":    "Hello World 9999 [$$$]",
	})

	tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.title", "Hello World 9999 [$$$]", pusher.Last().Message.Title)
	tt.AssertStrRepEqual(t, "msg.content", nil, pusher.Last().Message.Content)
	tt.AssertStrRepEqual(t, "msg.scn_msg_id", msg1["scn_msg_id"], pusher.Last().Message.SCNMessageID)

	type mglist struct {
		Messages []gin.H `json:"messages"`
	}

	msgList1 := tt.RequestAuthGet[mglist](t, admintok, baseUrl, "/api/messages")
	tt.AssertEqual(t, "len(messages)", 1, len(msgList1.Messages))

	msg1Get := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/messages/"+fmt.Sprintf("%v", msg1["scn_msg_id"]))
	tt.AssertStrRepEqual(t, "msg.title", "Hello World 9999 [$$$]", msg1Get["title"])
	tt.AssertStrRepEqual(t, "msg.channel_name", "main", msg1Get["channel_name"])
}

func TestSendSimpleMessageFormAndQuery(t *testing.T) {
	ws, stop := tt.StartSimpleWebserver(t)
	defer stop()

	pusher := ws.Pusher.(*push.TestSink)

	baseUrl := "http://127.0.0.1:" + ws.Port

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := int(r0["user_id"].(float64))
	sendtok := r0["send_key"].(string)

	msg1 := tt.RequestPost[gin.H](t, baseUrl, fmt.Sprintf("/?user_id=%d&user_key=%s&title=%s", uid, sendtok, url.QueryEscape("1111111")), tt.FormData{
		"user_key": "ERR",
		"user_id":  "999999",
		"title":    "2222222",
	})

	tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.title", "1111111", pusher.Last().Message.Title)
	tt.AssertStrRepEqual(t, "msg.scn_msg_id", msg1["scn_msg_id"], pusher.Last().Message.SCNMessageID)
}

func TestSendSimpleMessageJSONAndQuery(t *testing.T) {
	ws, stop := tt.StartSimpleWebserver(t)
	defer stop()

	pusher := ws.Pusher.(*push.TestSink)

	baseUrl := "http://127.0.0.1:" + ws.Port

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := int(r0["user_id"].(float64))
	sendtok := r0["send_key"].(string)

	msg1 := tt.RequestPost[gin.H](t, baseUrl, fmt.Sprintf("/?user_id=%d&user_key=%s&title=%s", uid, sendtok, url.QueryEscape("1111111")), gin.H{
		"user_key": "ERR",
		"user_id":  999999,
		"title":    "2222222",
	})

	tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.title", "1111111", pusher.Last().Message.Title)
	tt.AssertStrRepEqual(t, "msg.scn_msg_id", msg1["scn_msg_id"], pusher.Last().Message.SCNMessageID)
}

func TestSendContentMessage(t *testing.T) {
	ws, stop := tt.StartSimpleWebserver(t)
	defer stop()

	pusher := ws.Pusher.(*push.TestSink)

	baseUrl := "http://127.0.0.1:" + ws.Port

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := int(r0["user_id"].(float64))
	sendtok := r0["send_key"].(string)

	msg1 := tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    "HelloWorld_042",
		"content":  "I am Content\nasdf",
	})

	tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.title", "HelloWorld_042", pusher.Last().Message.Title)
	tt.AssertStrRepEqual(t, "msg.content", "I am Content\nasdf", pusher.Last().Message.Content)
	tt.AssertStrRepEqual(t, "msg.scn_msg_id", msg1["scn_msg_id"], pusher.Last().Message.SCNMessageID)
}

func TestSendWithSendername(t *testing.T) {
	ws, stop := tt.StartSimpleWebserver(t)
	defer stop()

	pusher := ws.Pusher.(*push.TestSink)

	baseUrl := "http://127.0.0.1:" + ws.Port

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := int(r0["user_id"].(float64))
	sendtok := r0["send_key"].(string)

	msg1 := tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"user_key":    sendtok,
		"user_id":     uid,
		"title":       "HelloWorld_xyz",
		"content":     "Unicode: 日本 - yäy\000\n\t\x00...",
		"sender_name": "localhorst",
	})

	tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.title", "HelloWorld_xyz", pusher.Last().Message.Title)
	tt.AssertStrRepEqual(t, "msg.content", "Unicode: 日本 - yäy\000\n\t\x00...", pusher.Last().Message.Content)
	tt.AssertStrRepEqual(t, "msg.SenderName", "localhorst", pusher.Last().Message.SenderName)
	tt.AssertStrRepEqual(t, "msg.scn_msg_id", msg1["scn_msg_id"], pusher.Last().Message.SCNMessageID)
}

func TestSendLongContent(t *testing.T) {
	ws, stop := tt.StartSimpleWebserver(t)
	defer stop()

	pusher := ws.Pusher.(*push.TestSink)

	baseUrl := "http://127.0.0.1:" + ws.Port

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := int(r0["user_id"].(float64))
	sendtok := r0["send_key"].(string)

	longContent := ""
	for i := 0; i < 200; i++ {
		longContent += "123456789\n" // 10 * 200 = 2_000
	}

	msg1 := tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    "HelloWorld_042",
		"content":  longContent,
	})

	tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.title", "HelloWorld_042", pusher.Last().Message.Title)
	tt.AssertStrRepEqual(t, "msg.content", longContent, pusher.Last().Message.Content)
	tt.AssertStrRepEqual(t, "msg.scn_msg_id", msg1["scn_msg_id"], pusher.Last().Message.SCNMessageID)
}

func TestSendTooLongContent(t *testing.T) {
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

	longContent := ""
	for i := 0; i < 400; i++ {
		longContent += "123456789\n" // 10 * 400 = 4_000
	}

	tt.RequestPostShouldFail(t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    "HelloWorld_042",
		"content":  longContent,
	}, 400, 1203)
}

func TestSendTooLongTitle(t *testing.T) {
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

	tt.RequestPostShouldFail(t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
	}, 400, 1202)
}

func TestSendIdempotent(t *testing.T) {
	ws, stop := tt.StartSimpleWebserver(t)
	defer stop()

	pusher := ws.Pusher.(*push.TestSink)

	baseUrl := "http://127.0.0.1:" + ws.Port

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := int(r0["user_id"].(float64))
	sendtok := r0["send_key"].(string)

	msg1 := tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    "Hello SCN",
		"content":  "mamma mia",
		"msg_id":   "c0235a49-dabc-4cdc-a0ce-453966e0c2d5",
	})

	tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.scn_msg_id", msg1["scn_msg_id"], pusher.Last().Message.SCNMessageID)
	tt.AssertStrRepEqual(t, "msg.suppress_send", msg1["suppress_send"], false)
	tt.AssertStrRepEqual(t, "msg.msg_id", "c0235a49-dabc-4cdc-a0ce-453966e0c2d5", pusher.Last().Message.UserMessageID)
	tt.AssertStrRepEqual(t, "msg.title", "Hello SCN", pusher.Last().Message.Title)
	tt.AssertStrRepEqual(t, "msg.content", "mamma mia", pusher.Last().Message.Content)

	msg2 := tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    "Hello again",
		"content":  "mother mia",
		"msg_id":   "c0235a49-dabc-4cdc-a0ce-453966e0c2d5",
	})

	tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.scn_msg_id", msg1["scn_msg_id"], msg2["scn_msg_id"])
	tt.AssertStrRepEqual(t, "msg.scn_msg_id", msg2["scn_msg_id"], pusher.Last().Message.SCNMessageID)
	tt.AssertStrRepEqual(t, "msg.suppress_send", msg2["suppress_send"], true)
	tt.AssertStrRepEqual(t, "msg.msg_id", "c0235a49-dabc-4cdc-a0ce-453966e0c2d5", pusher.Last().Message.UserMessageID)
	tt.AssertStrRepEqual(t, "msg.title", "Hello SCN", pusher.Last().Message.Title)
	tt.AssertStrRepEqual(t, "msg.content", "mamma mia", pusher.Last().Message.Content)

	msg3 := tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    "Hello third",
		"content":  "let me go",
		"msg_id":   "3238e68e-c1ea-44ce-b21b-2576614082b5",
	})

	tt.AssertEqual(t, "messageCount", 2, len(pusher.Data))
	tt.AssertNotStrRepEqual(t, "msg.scn_msg_id", msg1["scn_msg_id"], msg3["scn_msg_id"])
	tt.AssertNotStrRepEqual(t, "msg.scn_msg_id", msg2["scn_msg_id"], msg3["scn_msg_id"])
	tt.AssertStrRepEqual(t, "msg.suppress_send", msg3["suppress_send"], false)
	tt.AssertStrRepEqual(t, "msg.msg_id", "3238e68e-c1ea-44ce-b21b-2576614082b5", pusher.Last().Message.UserMessageID)
	tt.AssertStrRepEqual(t, "msg.title", "Hello third", pusher.Last().Message.Title)
	tt.AssertStrRepEqual(t, "msg.content", "let me go", pusher.Last().Message.Content)

}

//TODO compat route

//TODO post to channel
//TODO post to newly-created-channel
//TODO post to foreign channel via send-key

//TODO quota exceed (+ quota counter)

//TODO invalid priority
//TODO chan_naem too long
//TODO chan_name normalization
//TODO custom_timestamp
//TODO invalid time

//TODO check message_counter + last_sent in channel
//TODO check message_counter + last_sent in user
