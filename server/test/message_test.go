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
	tt.AssertEqual(t, "msg.title", "HelloWorld_001", pusher.Last().Message.Title)
	tt.AssertStrRepEqual(t, "msg.scn_msg_id", msg1["scn_msg_id"], pusher.Last().Message.SCNMessageID)

	type mglist struct {
		Messages []gin.H `json:"messages"`
	}

	msgList1 := tt.RequestAuthGet[mglist](t, admintok, baseUrl, "/api/messages")
	tt.AssertEqual(t, "len(messages)", 1, len(msgList1.Messages))

	msg1Get := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/messages/"+fmt.Sprintf("%v", msg1["scn_msg_id"]))
	tt.AssertEqual(t, "msg.title", "HelloWorld_001", msg1Get["title"])
	tt.AssertEqual(t, "msg.channel_name", "main", msg1Get["channel_name"])
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
	tt.AssertEqual(t, "msg.title", "Hello World 2134", pusher.Last().Message.Title)
	tt.AssertStrRepEqual(t, "msg.scn_msg_id", msg1["scn_msg_id"], pusher.Last().Message.SCNMessageID)

	type mglist struct {
		Messages []gin.H `json:"messages"`
	}

	msgList1 := tt.RequestAuthGet[mglist](t, admintok, baseUrl, "/api/messages")
	tt.AssertEqual(t, "len(messages)", 1, len(msgList1.Messages))

	msg1Get := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/messages/"+fmt.Sprintf("%v", msg1["scn_msg_id"]))
	tt.AssertEqual(t, "msg.title", "Hello World 2134", msg1Get["title"])
	tt.AssertEqual(t, "msg.channel_name", "main", msg1Get["channel_name"])
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
	tt.AssertEqual(t, "msg.title", "Hello World 9999 [$$$]", pusher.Last().Message.Title)
	tt.AssertStrRepEqual(t, "msg.scn_msg_id", msg1["scn_msg_id"], pusher.Last().Message.SCNMessageID)

	type mglist struct {
		Messages []gin.H `json:"messages"`
	}

	msgList1 := tt.RequestAuthGet[mglist](t, admintok, baseUrl, "/api/messages")
	tt.AssertEqual(t, "len(messages)", 1, len(msgList1.Messages))

	msg1Get := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/messages/"+fmt.Sprintf("%v", msg1["scn_msg_id"]))
	tt.AssertEqual(t, "msg.title", "Hello World 9999 [$$$]", msg1Get["title"])
	tt.AssertEqual(t, "msg.channel_name", "main", msg1Get["channel_name"])
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
	tt.AssertEqual(t, "msg.title", "1111111", pusher.Last().Message.Title)
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
	tt.AssertEqual(t, "msg.title", "1111111", pusher.Last().Message.Title)
	tt.AssertStrRepEqual(t, "msg.scn_msg_id", msg1["scn_msg_id"], pusher.Last().Message.SCNMessageID)
}

//TODO send content
//TODO sendername

//TODO trim too-long content

//TODO compat route

//TODO post to channel
//TODO post to newly-created-channel
//TODO post to foreign channel via send-key

//TODO usr_msg_id

//TODO quota exceed (+ quota counter)

//TODO invalid priority
//TODO chan_naem too long
//TODO chan_name normalization
//TODO custom_timestamp
//TODO invalid time
//TODO title too long
//TODO content too long

//TODO check message_counter + last_sent in channel
//TODO check message_counter + last_sent in user
