package test

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
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
	}, 401, apierr.USER_AUTH_FAILED)

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
	admintok := r0["admin_key"].(string)
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

	type mglist struct {
		Messages []gin.H `json:"messages"`
	}

	msgList1 := tt.RequestAuthGet[mglist](t, admintok, baseUrl, "/api/messages")
	tt.AssertEqual(t, "len(messages)", 1, len(msgList1.Messages))

	msg1Get := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/messages/"+fmt.Sprintf("%v", msg1["scn_msg_id"]))
	tt.AssertStrRepEqual(t, "msg.title", "HelloWorld_042", msg1Get["title"])
	tt.AssertStrRepEqual(t, "msg.content", "I am Content\nasdf", msg1Get["content"])
	tt.AssertStrRepEqual(t, "msg.channel_name", "main", msg1Get["channel_name"])
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
	admintok := r0["admin_key"].(string)
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

	type mglist struct {
		Messages []gin.H `json:"messages"`
	}

	msgList1 := tt.RequestAuthGet[mglist](t, admintok, baseUrl, "/api/messages")
	tt.AssertEqual(t, "len(messages)", 1, len(msgList1.Messages))
	tt.AssertStrRepEqual(t, "msg.title", "HelloWorld_042", msgList1.Messages[0]["title"])
	tt.AssertNotStrRepEqual(t, "msg.content", longContent, msgList1.Messages[0]["content"])
	tt.AssertStrRepEqual(t, "msg.channel_name", "main", msgList1.Messages[0]["channel_name"])
	tt.AssertStrRepEqual(t, "msg.trimmmed", true, msgList1.Messages[0]["trimmed"])

	msgList2 := tt.RequestAuthGet[mglist](t, admintok, baseUrl, "/api/messages?trimmed=false")
	tt.AssertEqual(t, "len(messages)", 1, len(msgList2.Messages))
	tt.AssertStrRepEqual(t, "msg.title", "HelloWorld_042", msgList2.Messages[0]["title"])
	tt.AssertStrRepEqual(t, "msg.content", longContent, msgList2.Messages[0]["content"])
	tt.AssertStrRepEqual(t, "msg.channel_name", "main", msgList2.Messages[0]["channel_name"])
	tt.AssertStrRepEqual(t, "msg.trimmmed", false, msgList2.Messages[0]["trimmed"])

	msg1Get := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/messages/"+fmt.Sprintf("%v", msg1["scn_msg_id"]))
	tt.AssertStrRepEqual(t, "msg.title", "HelloWorld_042", msg1Get["title"])
	tt.AssertStrRepEqual(t, "msg.titcontentle", longContent, msg1Get["content"])
	tt.AssertStrRepEqual(t, "msg.channel_name", "main", msg1Get["channel_name"])
	tt.AssertStrRepEqual(t, "msg.trimmmed", false, msg1Get["trimmed"])
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
	}, 400, apierr.CONTENT_TOO_LONG)
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
	}, 400, apierr.TITLE_TOO_LONG)
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
	readtok := r0["admin_key"].(string)
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

	type mglist struct {
		Messages []gin.H `json:"messages"`
	}

	msgList1 := tt.RequestAuthGet[mglist](t, readtok, baseUrl, "/api/messages")
	tt.AssertEqual(t, "len(messages)", 1, len(msgList1.Messages))

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

	msgList2 := tt.RequestAuthGet[mglist](t, readtok, baseUrl, "/api/messages")
	tt.AssertEqual(t, "len(messages)", 1, len(msgList2.Messages))

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

	msgList3 := tt.RequestAuthGet[mglist](t, readtok, baseUrl, "/api/messages")
	tt.AssertEqual(t, "len(messages)", 2, len(msgList3.Messages))
}

func TestSendWithPriority(t *testing.T) {
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
	admintok := r0["admin_key"].(string)

	{
		msg1 := tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
			"user_key": sendtok,
			"user_id":  uid,
			"title":    "M_001",
			"content":  "TestSendWithPriority#001",
		})

		tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))

		tt.AssertStrRepEqual(t, "msg.prio", 1, pusher.Last().Message.Priority)

		msg1Get := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/messages/"+fmt.Sprintf("%v", msg1["scn_msg_id"]))
		tt.AssertStrRepEqual(t, "msg.title", "M_001", msg1Get["title"])
		tt.AssertStrRepEqual(t, "msg.content", "TestSendWithPriority#001", msg1Get["content"])
		tt.AssertStrRepEqual(t, "msg.content", 1, msg1Get["priority"])
	}

	{
		msg2 := tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
			"user_key": sendtok,
			"user_id":  uid,
			"title":    "M_002",
			"content":  "TestSendWithPriority#002",
			"priority": 0,
		})

		tt.AssertEqual(t, "messageCount", 2, len(pusher.Data))

		tt.AssertStrRepEqual(t, "msg.prio", 0, pusher.Last().Message.Priority)

		msg2Get := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/messages/"+fmt.Sprintf("%v", msg2["scn_msg_id"]))
		tt.AssertStrRepEqual(t, "msg.title", "M_002", msg2Get["title"])
		tt.AssertStrRepEqual(t, "msg.content", "TestSendWithPriority#002", msg2Get["content"])
		tt.AssertStrRepEqual(t, "msg.content", 0, msg2Get["priority"])
	}

	{
		msg3 := tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
			"user_key": sendtok,
			"user_id":  uid,
			"title":    "M_003",
			"content":  "TestSendWithPriority#003",
			"priority": 1,
		})

		tt.AssertEqual(t, "messageCount", 3, len(pusher.Data))

		tt.AssertStrRepEqual(t, "msg.prio", 1, pusher.Last().Message.Priority)

		msg3Get := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/messages/"+fmt.Sprintf("%v", msg3["scn_msg_id"]))
		tt.AssertStrRepEqual(t, "msg.title", "M_003", msg3Get["title"])
		tt.AssertStrRepEqual(t, "msg.content", "TestSendWithPriority#003", msg3Get["content"])
		tt.AssertStrRepEqual(t, "msg.content", 1, msg3Get["priority"])
	}

	{
		msg4 := tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
			"user_key": sendtok,
			"user_id":  uid,
			"title":    "M_004",
			"content":  "TestSendWithPriority#004",
			"priority": 2,
		})

		tt.AssertEqual(t, "messageCount", 4, len(pusher.Data))

		tt.AssertStrRepEqual(t, "msg.prio", 2, pusher.Last().Message.Priority)

		msg4Get := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/messages/"+fmt.Sprintf("%v", msg4["scn_msg_id"]))
		tt.AssertStrRepEqual(t, "msg.title", "M_004", msg4Get["title"])
		tt.AssertStrRepEqual(t, "msg.content", "TestSendWithPriority#004", msg4Get["content"])
		tt.AssertStrRepEqual(t, "msg.content", 2, msg4Get["priority"])
	}
}

func TestSendInvalidPriority(t *testing.T) {
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
	admintok := r0["admin_key"].(string)

	tt.RequestPostShouldFail(t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    "(title)",
		"content":  "(content)",
		"priority": -1,
	}, 400, apierr.INVALID_PRIO)

	tt.RequestPostShouldFail(t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    "(title)",
		"content":  "(content)",
		"priority": 4,
	}, 400, apierr.INVALID_PRIO)

	tt.RequestPostShouldFail(t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    "(title)",
		"content":  "(content)",
		"priority": 9999,
	}, 400, apierr.INVALID_PRIO)

	tt.RequestPostShouldFail(t, baseUrl, "/", gin.H{
		"user_key": admintok,
		"user_id":  uid,
		"title":    "(title)",
		"content":  "(content)",
		"priority": -1,
	}, 400, apierr.INVALID_PRIO)

	tt.RequestPostShouldFail(t, baseUrl, "/", gin.H{
		"user_key": admintok,
		"user_id":  uid,
		"title":    "(title)",
		"content":  "(content)",
		"priority": 4,
	}, 400, apierr.INVALID_PRIO)

	tt.RequestPostShouldFail(t, baseUrl, "/", gin.H{
		"user_key": admintok,
		"user_id":  uid,
		"title":    "(title)",
		"content":  "(content)",
		"priority": 9999,
	}, 400, apierr.INVALID_PRIO)

	struid := fmt.Sprintf("%d", uid)

	tt.RequestPostShouldFail(t, baseUrl, "/", tt.FormData{
		"user_key": sendtok,
		"user_id":  struid,
		"title":    "(title)",
		"content":  "(content)",
		"priority": "-1",
	}, 400, apierr.INVALID_PRIO)

	tt.RequestPostShouldFail(t, baseUrl, "/", tt.FormData{
		"user_key": sendtok,
		"user_id":  struid,
		"title":    "(title)",
		"content":  "(content)",
		"priority": "4",
	}, 400, apierr.INVALID_PRIO)

	tt.RequestPostShouldFail(t, baseUrl, "/", tt.FormData{
		"user_key": sendtok,
		"user_id":  struid,
		"title":    "(title)",
		"content":  "(content)",
		"priority": "9999",
	}, 400, apierr.INVALID_PRIO)

	tt.RequestPostShouldFail(t, baseUrl, "/", tt.FormData{
		"user_key": admintok,
		"user_id":  struid,
		"title":    "(title)",
		"content":  "(content)",
		"priority": "-1",
	}, 400, apierr.INVALID_PRIO)

	tt.RequestPostShouldFail(t, baseUrl, "/", tt.FormData{
		"user_key": admintok,
		"user_id":  struid,
		"title":    "(title)",
		"content":  "(content)",
		"priority": "4",
	}, 400, apierr.INVALID_PRIO)

	tt.RequestPostShouldFail(t, baseUrl, "/", tt.FormData{
		"user_key": admintok,
		"user_id":  struid,
		"title":    "(title)",
		"content":  "(content)",
		"priority": "9999",
	}, 400, apierr.INVALID_PRIO)

	tt.AssertEqual(t, "messageCount", 0, len(pusher.Data))
}

//TODO compat route

//TODO post to channel

//TODO post to newly-created-channel

//TODO post to foreign channel via send-key

//TODO quota exceed (+ quota counter)

//TODO chan_naem too long

//TODO chan_name normalization

//TODO custom_timestamp

//TODO invalid time

//TODO check message_counter + last_sent in channel

//TODO check message_counter + last_sent in user

//todo test pagination
