package test

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/push"
	tt "blackforestbytes.com/simplecloudnotifier/test/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
	"testing"
	"time"
)

func TestSendSimpleMessageJSON(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	pusher := ws.Pusher.(*push.TestSink)

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

	tt.RequestPostShouldFail(t, baseUrl, "/", gin.H{
		"user_key": "asdf",
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
	tt.AssertStrRepEqual(t, "msg.channel_internal_name", "main", msg1Get["channel_internal_name"])
}

func TestSendSimpleMessageQuery(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	pusher := ws.Pusher.(*push.TestSink)

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
	tt.AssertStrRepEqual(t, "msg.channel_internal_name", "main", msg1Get["channel_internal_name"])
}

func TestSendSimpleMessageForm(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	pusher := ws.Pusher.(*push.TestSink)

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
	tt.AssertStrRepEqual(t, "msg.channel_internal_name", "main", msg1Get["channel_internal_name"])
}

func TestSendSimpleMessageFormAndQuery(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	pusher := ws.Pusher.(*push.TestSink)

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
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	pusher := ws.Pusher.(*push.TestSink)

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

func TestSendSimpleMessageAlt1(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	pusher := ws.Pusher.(*push.TestSink)

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

	msg1 := tt.RequestPost[gin.H](t, baseUrl, "/send", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    "HelloWorld_001",
	})

	tt.RequestPostShouldFail(t, baseUrl, "/send", gin.H{
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
	tt.AssertStrRepEqual(t, "msg.channel_internal_name", "main", msg1Get["channel_internal_name"])
}

func TestSendContentMessage(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	pusher := ws.Pusher.(*push.TestSink)

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
	tt.AssertStrRepEqual(t, "msg.title", "HelloWorld_042", msgList1.Messages[0]["title"])
	tt.AssertStrRepEqual(t, "msg.content", "I am Content\nasdf", msgList1.Messages[0]["content"])
	tt.AssertStrRepEqual(t, "msg.channel_internal_name", "main", msgList1.Messages[0]["channel_internal_name"])

	msg1Get := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/messages/"+fmt.Sprintf("%v", msg1["scn_msg_id"]))
	tt.AssertStrRepEqual(t, "msg.title", "HelloWorld_042", msg1Get["title"])
	tt.AssertStrRepEqual(t, "msg.content", "I am Content\nasdf", msg1Get["content"])
	tt.AssertStrRepEqual(t, "msg.channel_internal_name", "main", msg1Get["channel_internal_name"])
}

func TestSendWithSendername(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	pusher := ws.Pusher.(*push.TestSink)

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

	type mglist struct {
		Messages []gin.H `json:"messages"`
	}

	msgList1 := tt.RequestAuthGet[mglist](t, admintok, baseUrl, "/api/messages")
	tt.AssertEqual(t, "len(messages)", 1, len(msgList1.Messages))
	tt.AssertStrRepEqual(t, "msg.title", "HelloWorld_xyz", msgList1.Messages[0]["title"])
	tt.AssertStrRepEqual(t, "msg.content", "Unicode: 日本 - yäy\000\n\t\x00...", msgList1.Messages[0]["content"])
	tt.AssertStrRepEqual(t, "msg.sender_name", "localhorst", msgList1.Messages[0]["sender_name"])
	tt.AssertStrRepEqual(t, "msg.channel_internal_name", "main", msgList1.Messages[0]["channel_internal_name"])

	msg1Get := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/messages/"+fmt.Sprintf("%v", msg1["scn_msg_id"]))
	tt.AssertStrRepEqual(t, "msg.title", "HelloWorld_xyz", msg1Get["title"])
	tt.AssertStrRepEqual(t, "msg.content", "Unicode: 日本 - yäy\000\n\t\x00...", msg1Get["content"])
	tt.AssertStrRepEqual(t, "msg.sender_name", "localhorst", msg1Get["sender_name"])
	tt.AssertStrRepEqual(t, "msg.channel_internal_name", "main", msg1Get["channel_internal_name"])
}

func TestSendLongContent(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	pusher := ws.Pusher.(*push.TestSink)

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
	tt.AssertStrRepEqual(t, "msg.channel_internal_name", "main", msgList1.Messages[0]["channel_internal_name"])
	tt.AssertStrRepEqual(t, "msg.trimmmed", true, msgList1.Messages[0]["trimmed"])

	msgList2 := tt.RequestAuthGet[mglist](t, admintok, baseUrl, "/api/messages?trimmed=false")
	tt.AssertEqual(t, "len(messages)", 1, len(msgList2.Messages))
	tt.AssertStrRepEqual(t, "msg.title", "HelloWorld_042", msgList2.Messages[0]["title"])
	tt.AssertStrRepEqual(t, "msg.content", longContent, msgList2.Messages[0]["content"])
	tt.AssertStrRepEqual(t, "msg.channel_internal_name", "main", msgList2.Messages[0]["channel_internal_name"])
	tt.AssertStrRepEqual(t, "msg.trimmmed", false, msgList2.Messages[0]["trimmed"])

	msg1Get := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/messages/"+fmt.Sprintf("%v", msg1["scn_msg_id"]))
	tt.AssertStrRepEqual(t, "msg.title", "HelloWorld_042", msg1Get["title"])
	tt.AssertStrRepEqual(t, "msg.titcontentle", longContent, msg1Get["content"])
	tt.AssertStrRepEqual(t, "msg.channel_internal_name", "main", msg1Get["channel_internal_name"])
	tt.AssertStrRepEqual(t, "msg.trimmmed", false, msg1Get["trimmed"])
}

func TestSendTooLongContent(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

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

func TestSendLongContentPro(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
		"pro_token":     "ANDROID|v2|PURCHASED:DUMMY_TOK_XX",
	})

	uid := int(r0["user_id"].(float64))
	sendtok := r0["send_key"].(string)

	{
		longContent := ""
		for i := 0; i < 400; i++ {
			longContent += "123456789\n" // 10 * 400 = 4_000 (max = 16_384)
		}

		tt.RequestPost[tt.Void](t, baseUrl, "/", gin.H{
			"user_key": sendtok,
			"user_id":  uid,
			"title":    "HelloWorld_042",
			"content":  longContent,
		})
	}

	{
		longContent := ""
		for i := 0; i < 800; i++ {
			longContent += "123456789\n" // 10 * 800 = 8_000 (max = 16_384)
		}

		tt.RequestPost[tt.Void](t, baseUrl, "/", gin.H{
			"user_key": sendtok,
			"user_id":  uid,
			"title":    "HelloWorld_042",
			"content":  longContent,
		})

	}

	{
		longContent := ""
		for i := 0; i < 1600; i++ {
			longContent += "123456789\n" // 10 * 1600 = 16_000 (max = 16_384)
		}

		tt.RequestPost[tt.Void](t, baseUrl, "/", gin.H{
			"user_key": sendtok,
			"user_id":  uid,
			"title":    "HelloWorld_042",
			"content":  longContent,
		})
	}

	{
		longContent := ""
		for i := 0; i < 1630; i++ {
			longContent += "123456789\n" // 10 * 1630 = 163_000 (max = 16_384)
		}

		tt.RequestPost[tt.Void](t, baseUrl, "/", gin.H{
			"user_key": sendtok,
			"user_id":  uid,
			"title":    "HelloWorld_042",
			"content":  longContent,
		})
	}

	{
		longContent := ""
		for i := 0; i < 1640; i++ {
			longContent += "123456789\n" // 10 * 1640 = 164_000 (max = 16_384)
		}

		tt.RequestPostShouldFail(t, baseUrl, "/", gin.H{
			"user_key": sendtok,
			"user_id":  uid,
			"title":    "HelloWorld_042",
			"content":  longContent,
		}, 400, apierr.CONTENT_TOO_LONG)
	}
}

func TestSendTooLongTitle(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

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
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	pusher := ws.Pusher.(*push.TestSink)

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
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	pusher := ws.Pusher.(*push.TestSink)

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
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	pusher := ws.Pusher.(*push.TestSink)

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

func TestSendWithTimestamp(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	pusher := ws.Pusher.(*push.TestSink)

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := int(r0["user_id"].(float64))
	sendtok := r0["send_key"].(string)
	admintok := r0["admin_key"].(string)

	ts := time.Now().Unix() - int64(time.Hour.Seconds())

	msg1 := tt.RequestPost[gin.H](t, baseUrl, "/", tt.FormData{
		"user_key":  sendtok,
		"user_id":   fmt.Sprintf("%d", uid),
		"title":     "TTT",
		"timestamp": fmt.Sprintf("%d", ts),
	})

	tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.title", "TTT", pusher.Last().Message.Title)
	tt.AssertStrRepEqual(t, "msg.TimestampClient", ts, pusher.Last().Message.TimestampClient.Unix())
	tt.AssertStrRepEqual(t, "msg.Timestamp", ts, pusher.Last().Message.Timestamp().Unix())
	tt.AssertNotStrRepEqual(t, "msg.ts", pusher.Last().Message.TimestampClient, pusher.Last().Message.TimestampReal)
	tt.AssertStrRepEqual(t, "msg.scn_msg_id", msg1["scn_msg_id"], pusher.Last().Message.SCNMessageID)

	type mglist struct {
		Messages []gin.H `json:"messages"`
	}

	msgList1 := tt.RequestAuthGet[mglist](t, admintok, baseUrl, "/api/messages")
	tt.AssertEqual(t, "len(messages)", 1, len(msgList1.Messages))
	tt.AssertStrRepEqual(t, "msg.title", "TTT", msgList1.Messages[0]["title"])
	tt.AssertStrRepEqual(t, "msg.content", nil, msgList1.Messages[0]["sender_name"])
	tt.AssertStrRepEqual(t, "msg.channel_internal_name", "main", msgList1.Messages[0]["channel_internal_name"])

	tm1, err := time.Parse(time.RFC3339Nano, msgList1.Messages[0]["timestamp"].(string))
	tt.TestFailIfErr(t, err)
	tt.AssertStrRepEqual(t, "msg.timestamp", ts, tm1.Unix())

	msg1Get := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/messages/"+fmt.Sprintf("%v", msg1["scn_msg_id"]))
	tt.AssertStrRepEqual(t, "msg.title", "TTT", msg1Get["title"])
	tt.AssertStrRepEqual(t, "msg.content", nil, msg1Get["sender_name"])
	tt.AssertStrRepEqual(t, "msg.channel_internal_name", "main", msg1Get["channel_internal_name"])

	tmg1, err := time.Parse(time.RFC3339Nano, msg1Get["timestamp"].(string))
	tt.TestFailIfErr(t, err)
	tt.AssertStrRepEqual(t, "msg.timestamp", ts, tmg1.Unix())
}

func TestSendInvalidTimestamp(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	pusher := ws.Pusher.(*push.TestSink)

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := int(r0["user_id"].(float64))
	sendtok := r0["send_key"].(string)

	tt.RequestPostShouldFail(t, baseUrl, "/", tt.FormData{
		"user_key":  sendtok,
		"user_id":   fmt.Sprintf("%d", uid),
		"title":     "TTT",
		"timestamp": "-10000",
	}, 400, apierr.TIMESTAMP_OUT_OF_RANGE)

	tt.RequestPostShouldFail(t, baseUrl, "/", tt.FormData{
		"user_key":  sendtok,
		"user_id":   fmt.Sprintf("%d", uid),
		"title":     "TTT",
		"timestamp": "0",
	}, 400, apierr.TIMESTAMP_OUT_OF_RANGE)

	tt.RequestPostShouldFail(t, baseUrl, "/", tt.FormData{
		"user_key":  sendtok,
		"user_id":   fmt.Sprintf("%d", uid),
		"title":     "TTT",
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()-int64(25*time.Hour.Seconds())),
	}, 400, apierr.TIMESTAMP_OUT_OF_RANGE)

	tt.RequestPostShouldFail(t, baseUrl, "/", tt.FormData{
		"user_key":  sendtok,
		"user_id":   fmt.Sprintf("%d", uid),
		"title":     "TTT",
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()+int64(25*time.Hour.Seconds())),
	}, 400, apierr.TIMESTAMP_OUT_OF_RANGE)

	tt.RequestPostShouldFail(t, baseUrl, "/", gin.H{
		"user_key":  sendtok,
		"user_id":   uid,
		"title":     "TTT",
		"timestamp": -10000,
	}, 400, apierr.TIMESTAMP_OUT_OF_RANGE)

	tt.RequestPostShouldFail(t, baseUrl, "/", gin.H{
		"user_key":  sendtok,
		"user_id":   uid,
		"title":     "TTT",
		"timestamp": 0,
	}, 400, apierr.TIMESTAMP_OUT_OF_RANGE)

	tt.RequestPostShouldFail(t, baseUrl, "/", gin.H{
		"user_key":  sendtok,
		"user_id":   uid,
		"title":     "TTT",
		"timestamp": time.Now().Unix() - int64(25*time.Hour.Seconds()),
	}, 400, apierr.TIMESTAMP_OUT_OF_RANGE)

	tt.RequestPostShouldFail(t, baseUrl, "/", gin.H{
		"user_key":  sendtok,
		"user_id":   uid,
		"title":     "TTT",
		"timestamp": time.Now().Unix() + int64(25*time.Hour.Seconds()),
	}, 400, apierr.TIMESTAMP_OUT_OF_RANGE)

	tt.RequestPostShouldFail(t, baseUrl, fmt.Sprintf("/?user_key=%s&user_id=%d&title=%s&timestamp=%d",
		sendtok,
		uid,
		"TTT",
		-10000,
	), nil, 400, apierr.TIMESTAMP_OUT_OF_RANGE)

	tt.RequestPostShouldFail(t, baseUrl, fmt.Sprintf("/?user_key=%s&user_id=%d&title=%s&timestamp=%d",
		sendtok,
		uid,
		"TTT",
		0,
	), nil, 400, apierr.TIMESTAMP_OUT_OF_RANGE)

	tt.RequestPostShouldFail(t, baseUrl, fmt.Sprintf("/?user_key=%s&user_id=%d&title=%s&timestamp=%d",
		sendtok,
		uid,
		"TTT",
		time.Now().Unix()-int64(25*time.Hour.Seconds()),
	), nil, 400, apierr.TIMESTAMP_OUT_OF_RANGE)

	tt.RequestPostShouldFail(t, baseUrl, fmt.Sprintf("/?user_key=%s&user_id=%d&title=%s&timestamp=%d",
		sendtok,
		uid,
		"TTT",
		time.Now().Unix()+int64(25*time.Hour.Seconds()),
	), nil, 400, apierr.TIMESTAMP_OUT_OF_RANGE)

	tt.AssertEqual(t, "messageCount", 0, len(pusher.Data))
}

func TestSendCompat(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	pusher := ws.Pusher.(*push.TestSink)

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

	msg1 := tt.RequestPost[gin.H](t, baseUrl, "/send.php", tt.FormData{
		"user_key": sendtok,
		"user_id":  fmt.Sprintf("%d", uid),
		"title":    "HelloWorld_001",
	})

	tt.RequestPostShouldFail(t, baseUrl, "/send.php", gin.H{
		"user_key": readtok,
		"user_id":  uid,
		"title":    "HelloWorld_001",
	}, 0, 0)

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
	tt.AssertStrRepEqual(t, "msg.channel_internal_name", "main", msg1Get["channel_internal_name"])

	msg2 := tt.RequestPost[gin.H](t, baseUrl, fmt.Sprintf("/send.php?user_key=%s&user_id=%d&title=%s", sendtok, uid, "HelloWorld_002"), nil)

	tt.AssertEqual(t, "messageCount", 2, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.title", "HelloWorld_002", pusher.Last().Message.Title)
	tt.AssertStrRepEqual(t, "msg.content", nil, pusher.Last().Message.Content)
	tt.AssertStrRepEqual(t, "msg.scn_msg_id", msg2["scn_msg_id"], pusher.Last().Message.SCNMessageID)

	tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/messages/"+fmt.Sprintf("%v", msg2["scn_msg_id"]))

	content3 := "039c1817-76ee-44ab-972a-4cec0a15a791\n" +
		"046f59ea-9a49-4060-93e6-8a4e14134faf\n" +
		"ab566fbe-9020-41b6-afa6-94f3d8d7c7b4\n" +
		"d52e5f7d-26a8-45b9-befc-da44a3f112da\n" +
		"d19fae55-d52a-4753-b9f1-66a935d68b1e\n" +
		"99a4099d-44d5-497a-a69b-18e277400d6e\n" +
		"a55757aa-afaa-420e-afaf-f3951e9e2434\n" +
		"ee58f5fc-b384-49f4-bc2c-c5b3c7bd54b7\n" +
		"5a7008d9-dd15-406a-83d1-fd6209c56141\n"
	ts3 := time.Now().Unix() - int64(time.Hour.Seconds())

	msg3 := tt.RequestPost[gin.H](t, baseUrl, "/send.php", tt.FormData{
		"user_key":  sendtok,
		"user_id":   fmt.Sprintf("%d", uid),
		"title":     "HelloWorld_003",
		"content":   content3,
		"priority":  "2",
		"msg_id":    "8a2c7e92-86f3-4d69-897a-571286954030",
		"timestamp": fmt.Sprintf("%d", ts3),
	})

	tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/messages/"+fmt.Sprintf("%v", msg3["scn_msg_id"]))

	tt.AssertEqual(t, "messageCount", 3, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.Title", "HelloWorld_003", pusher.Last().Message.Title)
	tt.AssertStrRepEqual(t, "msg.Content", content3, pusher.Last().Message.Content)
	tt.AssertStrRepEqual(t, "msg.SCNMessageID", msg3["scn_msg_id"], pusher.Last().Message.SCNMessageID)
	tt.AssertStrRepEqual(t, "msg.Priority", 2, pusher.Last().Message.Priority)
	tt.AssertStrRepEqual(t, "msg.UserMessageID", "8a2c7e92-86f3-4d69-897a-571286954030", pusher.Last().Message.UserMessageID)
	tt.AssertStrRepEqual(t, "msg.UserMessageID", ts3, pusher.Last().Message.Timestamp().Unix())

}

func TestSendToNewChannel(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := int(r0["user_id"].(float64))
	sendtok := r0["send_key"].(string)
	admintok := r0["admin_key"].(string)

	type chanlist struct {
		Channels []gin.H `json:"channels"`
	}

	{
		chan0 := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid))
		tt.AssertEqual(t, "chan-count", 0, len(chan0.Channels))
	}

	tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    "M0",
	})

	{
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid))
		tt.AssertMappedSet(t, "channels", []string{"main"}, clist.Channels, "display_name")
		tt.AssertMappedSet(t, "channels", []string{"main"}, clist.Channels, "internal_name")
	}

	tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    "M1",
		"content":  tt.ShortLipsum0(4),
		"channel":  "main",
	})

	{
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid))
		tt.AssertMappedSet(t, "channels", []string{"main"}, clist.Channels, "display_name")
		tt.AssertMappedSet(t, "channels", []string{"main"}, clist.Channels, "internal_name")
	}

	tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    "M2",
		"content":  tt.ShortLipsum0(4),
		"channel":  "test",
	})

	{
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid))
		tt.AssertMappedSet(t, "channels", []string{"main", "test"}, clist.Channels, "display_name")
		tt.AssertMappedSet(t, "channels", []string{"main", "test"}, clist.Channels, "internal_name")
	}

	tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    "M3",
		"channel":  "test",
	})

	{
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid))
		tt.AssertMappedSet(t, "channels", []string{"main", "test"}, clist.Channels, "display_name")
		tt.AssertMappedSet(t, "channels", []string{"main", "test"}, clist.Channels, "internal_name")
	}
}

func TestSendToManualChannel(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := int(r0["user_id"].(float64))
	sendtok := r0["send_key"].(string)
	admintok := r0["admin_key"].(string)

	type chanlist struct {
		Channels []gin.H `json:"channels"`
	}

	{
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid))
		tt.AssertMappedSet(t, "channels", []string{}, clist.Channels, "display_name")
		tt.AssertMappedSet(t, "channels", []string{}, clist.Channels, "internal_name")
	}

	tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    "M0",
	})

	{
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid))
		tt.AssertMappedSet(t, "channels", []string{"main"}, clist.Channels, "display_name")
		tt.AssertMappedSet(t, "channels", []string{"main"}, clist.Channels, "internal_name")
	}

	tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    "M1",
		"content":  tt.ShortLipsum0(4),
		"channel":  "main",
	})

	{
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid))
		tt.AssertEqual(t, "chan.len", 1, len(clist.Channels))
		tt.AssertEqual(t, "chan.internal_name", "main", clist.Channels[0]["internal_name"])
		tt.AssertEqual(t, "chan.display_name", "main", clist.Channels[0]["display_name"])
	}

	tt.RequestAuthPost[gin.H](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid), gin.H{
		"name": "test",
	})

	{
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid))
		tt.AssertMappedSet(t, "channels", []string{"main", "test"}, clist.Channels, "display_name")
		tt.AssertMappedSet(t, "channels", []string{"main", "test"}, clist.Channels, "internal_name")
	}

	tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    "M2",
		"content":  tt.ShortLipsum0(4),
		"channel":  "test",
	})

	{
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid))
		tt.AssertMappedSet(t, "channels", []string{"main", "test"}, clist.Channels, "display_name")
		tt.AssertMappedSet(t, "channels", []string{"main", "test"}, clist.Channels, "internal_name")
	}

	tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    "M3",
		"channel":  "test",
	})

	{
		clist := tt.RequestAuthGet[chanlist](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d/channels", uid))
		tt.AssertMappedSet(t, "channels", []string{"main", "test"}, clist.Channels, "display_name")
		tt.AssertMappedSet(t, "channels", []string{"main", "test"}, clist.Channels, "internal_name")
	}
}

func TestSendToTooLongChannel(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := int(r0["user_id"].(float64))
	sendtok := r0["send_key"].(string)

	tt.RequestPost[tt.Void](t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    "M3",
		"channel":  "12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
	})

	tt.RequestPost[tt.Void](t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    "M3",
		"channel":  "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
	})

	tt.RequestPostShouldFail(t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    "M3",
		"channel":  "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901",
	}, 400, apierr.CHANNEL_TOO_LONG)
}

func TestQuotaExceededNoPro(t *testing.T) {
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
	sendtok := r0["send_key"].(string)

	tt.AssertStrRepEqual(t, "quota.0", 0, r0["quota_used"])
	tt.AssertStrRepEqual(t, "quota.0", 50, r0["quota_max"])
	tt.AssertStrRepEqual(t, "quota.0", 50, r0["quota_remaining"])

	{
		msg1 := tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
			"user_key": sendtok,
			"user_id":  uid,
			"title":    tt.ShortLipsum0(2),
		})
		tt.AssertStrRepEqual(t, "quota.msg.1", 1, msg1["quota"])
		tt.AssertStrRepEqual(t, "quota.msg.1", 50, msg1["quota_max"])
	}

	{
		usr := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d", uid))

		tt.AssertStrRepEqual(t, "quota.1", 1, usr["quota_used"])
		tt.AssertStrRepEqual(t, "quota.1", 50, usr["quota_max"])
		tt.AssertStrRepEqual(t, "quota.1", 49, usr["quota_remaining"])
	}

	for i := 0; i < 48; i++ {

		tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
			"user_key": sendtok,
			"user_id":  uid,
			"title":    tt.ShortLipsum0(2),
		})
	}

	{
		usr := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d", uid))

		tt.AssertStrRepEqual(t, "quota.49", 49, usr["quota_used"])
		tt.AssertStrRepEqual(t, "quota.49", 50, usr["quota_max"])
		tt.AssertStrRepEqual(t, "quota.49", 1, usr["quota_remaining"])
	}

	msg50 := tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    tt.ShortLipsum0(2),
	})
	tt.AssertStrRepEqual(t, "quota.msg.50", 50, msg50["quota"])
	tt.AssertStrRepEqual(t, "quota.msg.50", 50, msg50["quota_max"])

	{
		usr := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d", uid))

		tt.AssertStrRepEqual(t, "quota.50", 50, usr["quota_used"])
		tt.AssertStrRepEqual(t, "quota.50", 50, usr["quota_max"])
		tt.AssertStrRepEqual(t, "quota.50", 0, usr["quota_remaining"])
	}

	tt.RequestPostShouldFail(t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    tt.ShortLipsum0(2),
	}, 403, apierr.QUOTA_REACHED)
}

func TestQuotaExceededPro(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
		"pro_token":     "ANDROID|v2|PURCHASED:DUMMY_TOK_XX",
	})

	uid := int(r0["user_id"].(float64))
	admintok := r0["admin_key"].(string)
	sendtok := r0["send_key"].(string)

	tt.AssertStrRepEqual(t, "quota.0", 0, r0["quota_used"])
	tt.AssertStrRepEqual(t, "quota.0", 1000, r0["quota_max"])
	tt.AssertStrRepEqual(t, "quota.0", 1000, r0["quota_remaining"])

	{
		msg1 := tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
			"user_key": sendtok,
			"user_id":  uid,
			"title":    tt.ShortLipsum0(2),
		})
		tt.AssertStrRepEqual(t, "quota.msg.1", 1, msg1["quota"])
		tt.AssertStrRepEqual(t, "quota.msg.1", 1000, msg1["quota_max"])
	}

	{
		usr := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d", uid))

		tt.AssertStrRepEqual(t, "quota.1", 1, usr["quota_used"])
		tt.AssertStrRepEqual(t, "quota.1", 1000, usr["quota_max"])
		tt.AssertStrRepEqual(t, "quota.1", 999, usr["quota_remaining"])
	}

	for i := 0; i < 998; i++ {

		tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
			"user_key": sendtok,
			"user_id":  uid,
			"title":    tt.ShortLipsum0(2),
		})
	}

	{
		usr := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d", uid))

		tt.AssertStrRepEqual(t, "quota.999", 999, usr["quota_used"])
		tt.AssertStrRepEqual(t, "quota.999", 1000, usr["quota_max"])
		tt.AssertStrRepEqual(t, "quota.999", 1, usr["quota_remaining"])
	}

	msg50 := tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    tt.ShortLipsum0(2),
	})
	tt.AssertStrRepEqual(t, "quota.msg.1000", 1000, msg50["quota"])
	tt.AssertStrRepEqual(t, "quota.msg.1000", 1000, msg50["quota_max"])

	{
		usr := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, fmt.Sprintf("/api/users/%d", uid))

		tt.AssertStrRepEqual(t, "quota.1000", 1000, usr["quota_used"])
		tt.AssertStrRepEqual(t, "quota.1000", 1000, usr["quota_max"])
		tt.AssertStrRepEqual(t, "quota.1000", 0, usr["quota_remaining"])
	}

	tt.RequestPostShouldFail(t, baseUrl, "/", gin.H{
		"user_key": sendtok,
		"user_id":  uid,
		"title":    tt.ShortLipsum0(2),
	}, 403, apierr.QUOTA_REACHED)
}

func TestSendParallel(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
		"pro_token":     "ANDROID|v2|PURCHASED:DUMMY_TOK_XX",
	})

	uid := int(r0["user_id"].(float64))
	sendtok := r0["send_key"].(string)

	count := 128

	sem := make(chan tt.Void, count) // semaphore pattern
	for i := 0; i < count; i++ {
		go func() {
			defer func() {
				sem <- tt.Void{}
			}()
			tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
				"user_key": sendtok,
				"user_id":  uid,
				"title":    tt.ShortLipsum0(2),
			})
		}()
	}
	// wait for goroutines to finish
	for i := 0; i < count; i++ {
		<-sem
	}
}

//TODO post to foreign channel via send-key

//TODO check message_counter + last_sent in channel

//TODO check message_counter + last_sent in user

//TODO test delivery-retry