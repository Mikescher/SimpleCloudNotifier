package test

import (
	"blackforestbytes.com/simplecloudnotifier/push"
	tt "blackforestbytes.com/simplecloudnotifier/test/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
	"testing"
	"time"
)

func TestSendCompatWithOldUser(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	pusher := ws.Pusher.(*push.TestSink)

	r0 := tt.RequestGet[gin.H](t, baseUrl, "/api/register.php?fcm_token=DUMMY_FCM&pro=0&pro_token=")
	tt.AssertEqual(t, "success", true, r0["success"])

	uidold := int64(r0["user_id"].(float64))
	admintok := r0["user_key"].(string)

	msg1 := tt.RequestPost[gin.H](t, baseUrl, "/send.php", tt.FormData{
		"user_key": admintok,
		"user_id":  fmt.Sprintf("%d", uidold),
		"title":    "HelloWorld_001",
	})

	// does not allow json - only form & query
	tt.RequestPostShouldFail(t, baseUrl, "/send.php", gin.H{
		"user_key": admintok,
		"user_id":  uidold,
		"title":    "HelloWorld_001",
	}, 400, 0)

	tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.title", "HelloWorld_001", pusher.Last().Message.Title)
	tt.AssertStrRepEqual(t, "msg.content", nil, pusher.Last().Message.Content)

	exp1 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/expand.php?user_id=%d&user_key=%s&scn_msg_id=%d", uidold, admintok, int64(msg1["scn_msg_id"].(float64))))
	tt.AssertEqual(t, "success", true, exp1["success"])

	type mglist struct {
		Messages []gin.H `json:"messages"`
	}

	msgList1 := tt.RequestAuthGet[mglist](t, admintok, baseUrl, "/api/v2/messages")
	tt.AssertEqual(t, "len(messages)", 1, len(msgList1.Messages))

	msg1Get := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/v2/messages/"+fmt.Sprintf("%v", pusher.Last().Message.MessageID))
	tt.AssertStrRepEqual(t, "msg.title", "HelloWorld_001", msg1Get["title"])
	tt.AssertStrRepEqual(t, "msg.channel_internal_name", "main", msg1Get["channel_internal_name"])

	msg2 := tt.RequestPost[gin.H](t, baseUrl, fmt.Sprintf("/send.php?user_key=%s&user_id=%d&title=%s", admintok, uidold, "HelloWorld_002"), nil)

	tt.AssertEqual(t, "messageCount", 2, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.title", "HelloWorld_002", pusher.Last().Message.Title)
	tt.AssertStrRepEqual(t, "msg.content", nil, pusher.Last().Message.Content)

	exp2 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/expand.php?user_id=%d&user_key=%s&scn_msg_id=%d", uidold, admintok, int64(msg2["scn_msg_id"].(float64))))
	tt.AssertEqual(t, "success", true, exp2["success"])

	tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/v2/messages/"+fmt.Sprintf("%v", pusher.Last().Message.MessageID))

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
		"user_key":  admintok,
		"user_id":   fmt.Sprintf("%d", uidold),
		"title":     "HelloWorld_003",
		"content":   content3,
		"priority":  "2",
		"msg_id":    "8a2c7e92-86f3-4d69-897a-571286954030",
		"timestamp": fmt.Sprintf("%d", ts3),
	})

	tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/v2/messages/"+fmt.Sprintf("%v", pusher.Last().Message.MessageID))

	tt.AssertEqual(t, "messageCount", 3, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.Title", "HelloWorld_003", pusher.Last().Message.Title)
	tt.AssertStrRepEqual(t, "msg.Content", content3, pusher.Last().Message.Content)
	tt.AssertStrRepEqual(t, "msg.Priority", 2, pusher.Last().Message.Priority)
	tt.AssertStrRepEqual(t, "msg.UserMessageID", "8a2c7e92-86f3-4d69-897a-571286954030", pusher.Last().Message.UserMessageID)
	tt.AssertStrRepEqual(t, "msg.UserMessageID", ts3, pusher.Last().Message.Timestamp().Unix())

	exp3 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/expand.php?user_id=%d&user_key=%s&scn_msg_id=%d", uidold, admintok, int64(msg3["scn_msg_id"].(float64))))
	tt.AssertEqual(t, "success", true, exp3["success"])
}

func TestSendCompatWithNewUser(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	pusher := ws.Pusher.(*push.TestSink)

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/v2/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := r0["user_id"].(string)
	admintok := r0["admin_key"].(string)
	readtok := r0["read_key"].(string)
	sendtok := r0["send_key"].(string)

	uidold := tt.CreateCompatID(t, ws, "userid", uid)

	msg1 := tt.RequestPost[gin.H](t, baseUrl, "/send.php", tt.FormData{
		"user_key": sendtok,
		"user_id":  fmt.Sprintf("%d", uidold),
		"title":    "HelloWorld_001",
	})

	// does not allow json - only form & query
	tt.RequestPostShouldFail(t, baseUrl, "/send.php", gin.H{
		"user_key": readtok,
		"user_id":  uidold,
		"title":    "HelloWorld_001",
	}, 400, 0)

	tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.title", "HelloWorld_001", pusher.Last().Message.Title)
	tt.AssertStrRepEqual(t, "msg.content", nil, pusher.Last().Message.Content)

	exp1 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/expand.php?user_id=%d&user_key=%s&scn_msg_id=%d", uidold, admintok, int64(msg1["scn_msg_id"].(float64))))
	tt.AssertEqual(t, "success", true, exp1["success"])

	type mglist struct {
		Messages []gin.H `json:"messages"`
	}

	msgList1 := tt.RequestAuthGet[mglist](t, admintok, baseUrl, "/api/v2/messages")
	tt.AssertEqual(t, "len(messages)", 1, len(msgList1.Messages))

	msg1Get := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/v2/messages/"+fmt.Sprintf("%v", pusher.Last().Message.MessageID))
	tt.AssertStrRepEqual(t, "msg.title", "HelloWorld_001", msg1Get["title"])
	tt.AssertStrRepEqual(t, "msg.channel_internal_name", "main", msg1Get["channel_internal_name"])

	msg2 := tt.RequestPost[gin.H](t, baseUrl, fmt.Sprintf("/send.php?user_key=%s&user_id=%d&title=%s", sendtok, uidold, "HelloWorld_002"), nil)

	tt.AssertEqual(t, "messageCount", 2, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.title", "HelloWorld_002", pusher.Last().Message.Title)
	tt.AssertStrRepEqual(t, "msg.content", nil, pusher.Last().Message.Content)

	exp2 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/expand.php?user_id=%d&user_key=%s&scn_msg_id=%d", uidold, admintok, int64(msg2["scn_msg_id"].(float64))))
	tt.AssertEqual(t, "success", true, exp2["success"])

	tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/v2/messages/"+fmt.Sprintf("%v", pusher.Last().Message.MessageID))

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
		"user_id":   fmt.Sprintf("%d", uidold),
		"title":     "HelloWorld_003",
		"content":   content3,
		"priority":  "2",
		"msg_id":    "8a2c7e92-86f3-4d69-897a-571286954030",
		"timestamp": fmt.Sprintf("%d", ts3),
	})

	exp3 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/expand.php?user_id=%d&user_key=%s&scn_msg_id=%d", uidold, admintok, int64(msg3["scn_msg_id"].(float64))))
	tt.AssertEqual(t, "success", true, exp3["success"])
	tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/v2/messages/"+fmt.Sprintf("%v", pusher.Last().Message.MessageID))

	tt.AssertEqual(t, "messageCount", 3, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.Title", "HelloWorld_003", pusher.Last().Message.Title)
	tt.AssertStrRepEqual(t, "msg.Content", content3, pusher.Last().Message.Content)
	tt.AssertStrRepEqual(t, "msg.Priority", 2, pusher.Last().Message.Priority)
	tt.AssertStrRepEqual(t, "msg.UserMessageID", "8a2c7e92-86f3-4d69-897a-571286954030", pusher.Last().Message.UserMessageID)
	tt.AssertStrRepEqual(t, "msg.UserMessageID", ts3, pusher.Last().Message.Timestamp().Unix())

}

func TestSendCompatMessageByQuery(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/register.php?fcm_token=%s&pro=%s&pro_token=%s", "DUMMY_FCM", "0", ""))
	tt.AssertEqual(t, "success", true, r0["success"])

	userid := int64(r0["user_id"].(float64))
	userkey := r0["user_key"].(string)

	r1 := tt.RequestPost[gin.H](t, baseUrl, fmt.Sprintf("/send.php?user_id=%d&user_key=%s&title=%s",
		userid,
		userkey,
		url.QueryEscape("my title 11 & x")), nil)
	tt.AssertEqual(t, "success", true, r1["success"])
	tt.AssertEqual(t, "suppress_send", false, r1["suppress_send"])

	r1scnid := int64(r1["scn_msg_id"].(float64))

	r1x := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/expand.php?user_id=%d&user_key=%s&scn_msg_id=%d", userid, userkey, r1scnid))
	tt.AssertEqual(t, "success", true, r1x["success"])
	tt.AssertEqual(t, "success", "my title 11 & x", (r1x["data"].(map[string]any))["title"])

	r2 := tt.RequestPost[gin.H](t, baseUrl, fmt.Sprintf("/send.php?user_id=%d&user_key=%s&title=%s&content=%s&priority=%s&msg_id=%s&timestamp=%s",
		userid,
		userkey,
		url.QueryEscape("my title"),
		url.QueryEscape("message content"),
		url.QueryEscape("2"),
		url.QueryEscape("624dbe5e-6d03-47cd-9a0e-a306faa2e977"),
		url.QueryEscape(fmt.Sprintf("%d", time.Now().Unix()+666))), nil)
	tt.AssertEqual(t, "success", true, r2["success"])
	tt.AssertEqual(t, "suppress_send", false, r2["suppress_send"])

	r2scnid := int64(r2["scn_msg_id"].(float64))

	r2x := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/expand.php?user_id=%d&user_key=%s&scn_msg_id=%d", userid, userkey, r2scnid))
	tt.AssertEqual(t, "success", true, r2x["success"])
	tt.AssertEqual(t, "success", "my title", (r2x["data"].(map[string]any))["title"])

	r3 := tt.RequestPost[gin.H](t, baseUrl, fmt.Sprintf("/send.php?user_id=%d&user_key=%s&title=%s&content=%s&priority=%s&msg_id=%s&timestamp=%s",
		userid,
		userkey,
		url.QueryEscape("my title"),
		url.QueryEscape("message content"),
		url.QueryEscape("2"),
		url.QueryEscape("624dbe5e-6d03-47cd-9a0e-a306faa2e977"),
		url.QueryEscape(fmt.Sprintf("%d", time.Now().Unix()+666))), nil)
	tt.AssertEqual(t, "success", true, r3["success"])
	tt.AssertEqual(t, "suppress_send", true, r3["suppress_send"])
}

func TestSendCompatMessageByFormData(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/register.php?fcm_token=%s&pro=%s&pro_token=%s", "DUMMY_FCM", "0", ""))
	tt.AssertEqual(t, "success", true, r0["success"])

	userid := int64(r0["user_id"].(float64))
	userkey := r0["user_key"].(string)

	r1 := tt.RequestPost[gin.H](t, baseUrl, "/send.php", tt.FormData{
		"user_id":  fmt.Sprintf("%d", userid),
		"user_key": userkey,
		"title":    "my title 11 & x",
	})
	tt.AssertEqual(t, "success", true, r1["success"])
	tt.AssertEqual(t, "suppress_send", false, r1["suppress_send"])

	r1scnid := int64(r1["scn_msg_id"].(float64))

	r1x := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/expand.php?user_id=%d&user_key=%s&scn_msg_id=%d", userid, userkey, r1scnid))
	tt.AssertEqual(t, "success", true, r1x["success"])
	tt.AssertEqual(t, "title", "my title 11 & x", (r1x["data"].(map[string]any))["title"])

	r2 := tt.RequestPost[gin.H](t, baseUrl, "/send.php", tt.FormData{
		"user_id":   fmt.Sprintf("%d", userid),
		"user_key":  userkey,
		"title":     "my title",
		"content":   "message content",
		"priority":  "2",
		"msg_id":    "624dbe5e-6d03-47cd-9a0e-a306faa2e977",
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()-666),
	})
	tt.AssertEqual(t, "success", true, r2["success"])
	tt.AssertEqual(t, "suppress_send", false, r2["suppress_send"])

	r2scnid := int64(r2["scn_msg_id"].(float64))

	r2x := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/expand.php?user_id=%d&user_key=%s&scn_msg_id=%d", userid, userkey, r2scnid))
	tt.AssertEqual(t, "success", true, r2x["success"])
	tt.AssertEqual(t, "title", "my title", (r2x["data"].(map[string]any))["title"])

	r3 := tt.RequestPost[gin.H](t, baseUrl, "/send.php", tt.FormData{
		"user_id":   fmt.Sprintf("%d", userid),
		"user_key":  userkey,
		"title":     "my title",
		"content":   "message content",
		"priority":  "2",
		"msg_id":    "624dbe5e-6d03-47cd-9a0e-a306faa2e977",
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()-666),
	})
	tt.AssertEqual(t, "success", true, r3["success"])
	tt.AssertEqual(t, "suppress_send", true, r3["suppress_send"])
}

func TestCompatRegister(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/register.php?fcm_token=%s&pro=%s&pro_token=%s", "DUMMY_FCM", "0", ""))
	tt.AssertEqual(t, "success", true, r0["success"])
	tt.AssertEqual(t, "message", "New user registered", r0["message"])
	tt.AssertEqual(t, "quota", 0, r0["quota"])
	tt.AssertEqual(t, "quota_max", 50, r0["quota_max"])
	tt.AssertEqual(t, "is_pro", false, r0["is_pro"])
}

func TestCompatRegisterPro(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/register.php?fcm_token=%s&pro=%s&pro_token=%s", "DUMMY_FCM", "true", url.QueryEscape("PURCHASED:000")))
	tt.AssertEqual(t, "success", true, r0["success"])
	tt.AssertEqual(t, "message", "New user registered", r0["message"])
	tt.AssertEqual(t, "quota", 0, r0["quota"])
	tt.AssertEqual(t, "quota_max", 5000, r0["quota_max"])
	tt.AssertEqual(t, "is_pro", true, r0["is_pro"])

	r1 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/register.php?fcm_token=%s&pro=%s&pro_token=%s", "DUMMY_FCM", "true", url.QueryEscape("INVALID")))
	tt.AssertEqual(t, "success", false, r1["success"])
}

func TestCompatInfo(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/register.php?fcm_token=%s&pro=%s&pro_token=%s", "DUMMY_FCM", "0", ""))
	tt.AssertEqual(t, "success", true, r0["success"])

	userid := int64(r0["user_id"].(float64))
	userkey := r0["user_key"].(string)

	r1 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/info.php?user_id=%d&user_key=%s", userid, userkey))

	tt.AssertEqual(t, "success", true, r1["success"])
	tt.AssertEqual(t, "fcm_token_set", true, r1["fcm_token_set"])
	tt.AssertEqual(t, "is_pro", 0, r1["is_pro"])
	tt.AssertEqual(t, "message", "ok", r1["message"])
	tt.AssertEqual(t, "quota", 0, r1["quota"])
	tt.AssertEqual(t, "quota_max", 50, r1["quota_max"])
	tt.AssertEqual(t, "unack_count", 0, r1["unack_count"])
	tt.AssertEqual(t, "user_id", userid, r1["user_id"])
	tt.AssertEqual(t, "user_key", userkey, r1["user_key"])

	tt.RequestPost[gin.H](t, baseUrl, "/send.php", tt.FormData{
		"user_key": userkey,
		"user_id":  fmt.Sprintf("%d", userid),
		"title":    tt.ShortLipsum0(1),
	})

	r2 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/info.php?user_id=%d&user_key=%s", userid, userkey))

	tt.AssertEqual(t, "success", true, r2["success"])
	tt.AssertEqual(t, "fcm_token_set", true, r2["fcm_token_set"])
	tt.AssertEqual(t, "is_pro", 0, r2["is_pro"])
	tt.AssertEqual(t, "message", "ok", r2["message"])
	tt.AssertEqual(t, "quota", 1, r2["quota"])
	tt.AssertEqual(t, "quota_max", 50, r2["quota_max"])
	tt.AssertEqual(t, "unack_count", 1, r2["unack_count"])
	tt.AssertEqual(t, "user_id", userid, r2["user_id"])
	tt.AssertEqual(t, "user_key", userkey, r2["user_key"])

}

func TestCompatAck(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/register.php?fcm_token=%s&pro=%s&pro_token=%s", "DUMMY_FCM", "0", ""))
	tt.AssertEqual(t, "success", true, r0["success"])

	userid := int64(r0["user_id"].(float64))
	userkey := r0["user_key"].(string)

	r1 := tt.RequestPost[gin.H](t, baseUrl, "/send.php", tt.FormData{
		"user_id":  fmt.Sprintf("%d", userid),
		"user_key": userkey,
		"title":    "my title 11 & x",
	})
	tt.AssertEqual(t, "success", true, r1["success"])
	r1scnid := int64(r1["scn_msg_id"].(float64))

	ack := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/ack.php?user_id=%d&user_key=%s&scn_msg_id=%d", userid, userkey, r1scnid))
	tt.AssertEqual(t, "success", true, ack["success"])
	tt.AssertEqual(t, "prev_ack", 0, ack["prev_ack"])
	tt.AssertEqual(t, "new_ack", 1, ack["new_ack"])
	tt.AssertEqual(t, "message", "ok", ack["message"])

}

func TestCompatExpand(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/register.php?fcm_token=%s&pro=%s&pro_token=%s", "DUMMY_FCM", "0", ""))
	tt.AssertEqual(t, "success", true, r0["success"])

	userid := int64(r0["user_id"].(float64))
	userkey := r0["user_key"].(string)

	{
		ts := time.Now().Unix()

		r1 := tt.RequestPost[gin.H](t, baseUrl, "/send.php", tt.FormData{
			"user_id":   fmt.Sprintf("%d", userid),
			"user_key":  userkey,
			"title":     "_title_",
			"content":   "_content_",
			"timestamp": fmt.Sprintf("%d", ts),
		})
		tt.AssertEqual(t, "success", true, r1["success"])
		r1scnid := int64(r1["scn_msg_id"].(float64))

		exp1 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/expand.php?user_id=%d&user_key=%s&scn_msg_id=%d", userid, userkey, r1scnid))
		tt.AssertEqual(t, "success", true, exp1["success"])

		exp1data := exp1["data"].(map[string]any)

		tt.AssertEqual(t, "title", "_title_", exp1data["title"])
		tt.AssertEqual(t, "body", "_content_", exp1data["body"])
		tt.AssertEqual(t, "priority", 1, exp1data["priority"])
		tt.AssertEqual(t, "timestamp", ts, exp1data["timestamp"])
		tt.AssertEqual(t, "usr_msg_id", nil, exp1data["usr_msg_id"])
		tt.AssertEqual(t, "scn_msg_id", r1scnid, exp1data["scn_msg_id"])
		tt.AssertEqual(t, "trimmed", false, exp1data["trimmed"])
	}

	{
		ts := time.Now().Unix()

		r1 := tt.RequestPost[gin.H](t, baseUrl, "/send.php", tt.FormData{
			"user_id":   fmt.Sprintf("%d", userid),
			"user_key":  userkey,
			"title":     "_title_",
			"timestamp": fmt.Sprintf("%d", ts),
			"priority":  "0",
			"msg_id":    "36aa8281-4bcd-4973-9368-e1d1ca5e21cb",
		})
		tt.AssertEqual(t, "success", true, r1["success"])
		r1scnid := int64(r1["scn_msg_id"].(float64))

		exp1 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/expand.php?user_id=%d&user_key=%s&scn_msg_id=%d", userid, userkey, r1scnid))
		tt.AssertEqual(t, "success", true, exp1["success"])

		exp1data := exp1["data"].(map[string]any)

		tt.AssertEqual(t, "title", "_title_", exp1data["title"])
		tt.AssertEqual(t, "body", nil, exp1data["body"])
		tt.AssertEqual(t, "priority", 0, exp1data["priority"])
		tt.AssertEqual(t, "timestamp", ts, exp1data["timestamp"])
		tt.AssertEqual(t, "usr_msg_id", "36aa8281-4bcd-4973-9368-e1d1ca5e21cb", exp1data["usr_msg_id"])
		tt.AssertEqual(t, "scn_msg_id", r1scnid, exp1data["scn_msg_id"])
		tt.AssertEqual(t, "trimmed", false, exp1data["trimmed"])
	}

}

func TestCompatUpdateUserKey(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	pusher := ws.Pusher.(*push.TestSink)

	r0 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/register.php?fcm_token=%s&pro=%s&pro_token=%s", "DUMMY_FCM", "0", ""))
	tt.AssertEqual(t, "success", true, r0["success"])

	userid := int64(r0["user_id"].(float64))
	userkey := r0["user_key"].(string)

	s0 := tt.RequestPost[gin.H](t, baseUrl, fmt.Sprintf("/send.php?user_id=%d&user_key=%s&title=%s", userid, userkey, url.QueryEscape("msg_1")), nil)
	tt.AssertEqual(t, "success", true, s0["success"])
	tt.AssertEqual(t, "fcm", "DUMMY_FCM", pusher.Last().Client.FCMToken)

	upd := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/update.php?user_id=%d&user_key=%s", userid, userkey))
	tt.AssertEqual(t, "success", true, upd["success"])

	newkey := upd["user_key"].(string)

	r1 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/info.php?user_id=%d&user_key=%s", userid, newkey))
	tt.AssertEqual(t, "success", true, r1["success"])
	tt.AssertEqual(t, "fcm_token_set", true, r1["fcm_token_set"])
	tt.AssertEqual(t, "is_pro", 0, r1["is_pro"])
	tt.AssertEqual(t, "message", "ok", r1["message"])
	tt.AssertEqual(t, "quota", 1, r1["quota"])
	tt.AssertEqual(t, "quota_max", 50, r1["quota_max"])
	tt.AssertEqual(t, "unack_count", 1, r1["unack_count"])
	tt.AssertEqual(t, "user_id", userid, r1["user_id"])
	tt.AssertEqual(t, "user_key", newkey, r1["user_key"])

	s1 := tt.RequestPost[gin.H](t, baseUrl, fmt.Sprintf("/send.php?user_id=%d&user_key=%s&title=%s", userid, newkey, url.QueryEscape("msg_2")), nil)
	tt.AssertEqual(t, "success", true, s1["success"])
	tt.AssertEqual(t, "fcm", "DUMMY_FCM", pusher.Last().Client.FCMToken)
}

func TestCompatUpdateFCM(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	pusher := ws.Pusher.(*push.TestSink)

	r0 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/register.php?fcm_token=%s&pro=%s&pro_token=%s", "DUMMY_FCM", "0", ""))
	tt.AssertEqual(t, "success", true, r0["success"])

	userid := int64(r0["user_id"].(float64))
	userkey := r0["user_key"].(string)

	s0 := tt.RequestPost[gin.H](t, baseUrl, fmt.Sprintf("/send.php?user_id=%d&user_key=%s&title=%s", userid, userkey, url.QueryEscape("msg_1")), nil)
	tt.AssertEqual(t, "success", true, s0["success"])
	tt.AssertEqual(t, "fcm", "DUMMY_FCM", pusher.Last().Client.FCMToken)

	upd := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/update.php?user_id=%d&user_key=%s&fcm_token=%s", userid, userkey, "NEW_FCM"))
	tt.AssertEqual(t, "success", true, upd["success"])

	newkey := upd["user_key"].(string)

	r1 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/info.php?user_id=%d&user_key=%s", userid, newkey))
	tt.AssertEqual(t, "success", true, r1["success"])
	tt.AssertEqual(t, "fcm_token_set", true, r1["fcm_token_set"])
	tt.AssertEqual(t, "is_pro", 0, r1["is_pro"])
	tt.AssertEqual(t, "message", "ok", r1["message"])
	tt.AssertEqual(t, "quota", 1, r1["quota"])
	tt.AssertEqual(t, "quota_max", 50, r1["quota_max"])
	tt.AssertEqual(t, "unack_count", 1, r1["unack_count"])
	tt.AssertEqual(t, "user_id", userid, r1["user_id"])
	tt.AssertEqual(t, "user_key", newkey, r1["user_key"])

	s1 := tt.RequestPost[gin.H](t, baseUrl, fmt.Sprintf("/send.php?user_id=%d&user_key=%s&title=%s", userid, newkey, url.QueryEscape("msg_2")), nil)
	tt.AssertEqual(t, "success", true, s1["success"])
	tt.AssertEqual(t, "fcm", "NEW_FCM", pusher.Last().Client.FCMToken)

	upd2 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/update.php?user_id=%d&user_key=%s&fcm_token=%s", userid, newkey, "NEW_FCM_2"))
	tt.AssertEqual(t, "success", true, upd2["success"])

	newkey = upd2["user_key"].(string)

	r1 = tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/info.php?user_id=%d&user_key=%s", userid, newkey))
	tt.AssertEqual(t, "success", true, r1["success"])
	tt.AssertEqual(t, "fcm_token_set", true, r1["fcm_token_set"])
}

func TestCompatUpgrade(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/register.php?fcm_token=%s&pro=%s&pro_token=%s", "FCM_DUMMY_FCM", "0", ""))
	tt.AssertEqual(t, "success", true, r0["success"])
	tt.AssertEqual(t, "message", "New user registered", r0["message"])
	tt.AssertEqual(t, "quota", 0, r0["quota"])
	tt.AssertEqual(t, "quota_max", 50, r0["quota_max"])
	tt.AssertEqual(t, "is_pro", false, r0["is_pro"])

	userid := int64(r0["user_id"].(float64))
	userkey := r0["user_key"].(string)

	r1 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/upgrade.php?user_id=%d&user_key=%s&pro=%s&pro_token=%s", userid, userkey, "true", url.QueryEscape("PURCHASED:000")))
	tt.AssertEqual(t, "success", true, r1["success"])
	tt.AssertEqual(t, "message", "user updated", r1["message"])
	tt.AssertEqual(t, "quota", 0, r1["quota"])
	tt.AssertEqual(t, "quota_max", 5000, r1["quota_max"])
	tt.AssertEqual(t, "is_pro", true, r1["is_pro"])
}

func TestCompatRequery(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/register.php?fcm_token=%s&pro=%s&pro_token=%s", "DUMMY_FCM", "0", ""))
	tt.AssertEqual(t, "success", true, r0["success"])

	userid := int64(r0["user_id"].(float64))
	userkey := r0["user_key"].(string)
	useridnew := tt.ConvertCompatID(t, ws, userid, "userid")

	rq1 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/requery.php?user_id=%d&user_key=%s", userid, userkey))
	tt.AssertEqual(t, "success", true, rq1["success"])
	tt.AssertEqual(t, "count", 0, rq1["count"])
	tt.AssertStrRepEqual(t, "data", make([]any, 0), rq1["data"])

	r1 := tt.RequestPost[gin.H](t, baseUrl, "/send.php", tt.FormData{
		"user_id":  fmt.Sprintf("%d", userid),
		"user_key": userkey,
		"title":    "_title_",
		"msg_id":   "r1",
	})
	tt.AssertEqual(t, "success", true, r1["success"])

	type respRequery struct {
		Success bool    `json:"success"`
		Message string  `json:"message"`
		Count   int     `json:"count"`
		Data    []gin.H `json:"data"`
	}

	rq2 := tt.RequestGet[respRequery](t, baseUrl, fmt.Sprintf("/api/requery.php?user_id=%d&user_key=%s", userid, userkey))
	tt.AssertEqual(t, "success", true, rq2.Success)
	tt.AssertEqual(t, "count", 1, rq2.Count)
	tt.AssertMappedSet(t, "data", []string{"r1"}, rq2.Data, "usr_msg_id")

	rq3 := tt.RequestGet[respRequery](t, baseUrl, fmt.Sprintf("/api/requery.php?user_id=%d&user_key=%s", userid, userkey))
	tt.AssertEqual(t, "success", true, rq3.Success)
	tt.AssertEqual(t, "count", 1, rq3.Count)
	tt.AssertMappedSet(t, "data", []string{"r1"}, rq3.Data, "usr_msg_id")

	a2 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/ack.php?user_id=%d&user_key=%s&scn_msg_id=%d", userid, userkey, int(r1["scn_msg_id"].(float64))))
	tt.AssertEqual(t, "success", true, a2["success"])

	rq31 := tt.RequestGet[respRequery](t, baseUrl, fmt.Sprintf("/api/requery.php?user_id=%d&user_key=%s", userid, userkey))
	tt.AssertEqual(t, "success", true, rq31.Success)
	tt.AssertEqual(t, "count", 0, rq31.Count)

	r2 := tt.RequestPost[gin.H](t, baseUrl, "/send.php", tt.FormData{
		"user_id":  fmt.Sprintf("%d", userid),
		"user_key": userkey,
		"title":    "_title_",
		"msg_id":   "r2",
	})
	tt.AssertEqual(t, "success", true, r2["success"])

	rq4 := tt.RequestGet[respRequery](t, baseUrl, fmt.Sprintf("/api/requery.php?user_id=%d&user_key=%s", userid, userkey))
	tt.AssertEqual(t, "success", true, rq4.Success)
	tt.AssertEqual(t, "count", 1, rq4.Count)
	tt.AssertMappedSet(t, "data", []string{"r2"}, rq4.Data, "usr_msg_id")

	r3 := tt.RequestPost[gin.H](t, baseUrl, "/send.php", tt.FormData{
		"user_id":  fmt.Sprintf("%d", userid),
		"user_key": userkey,
		"title":    "_title_",
		"msg_id":   "r3",
	})
	tt.AssertEqual(t, "success", true, r3["success"])

	r4 := tt.RequestPost[gin.H](t, baseUrl, "/send.php", tt.FormData{
		"user_id":  fmt.Sprintf("%d", userid),
		"user_key": userkey,
		"title":    "_title_",
		"msg_id":   "r4",
	})
	tt.AssertEqual(t, "success", true, r4["success"])

	r5 := tt.RequestPost[gin.H](t, baseUrl, "/send.php", tt.FormData{
		"user_id":  fmt.Sprintf("%d", userid),
		"user_key": userkey,
		"title":    "_title_",
		"msg_id":   "r5",
	})
	tt.AssertEqual(t, "success", true, r5["success"])

	a1 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/ack.php?user_id=%d&user_key=%s&scn_msg_id=%d", userid, userkey, int(r4["scn_msg_id"].(float64))))
	tt.AssertEqual(t, "success", true, a1["success"])

	rq5 := tt.RequestGet[respRequery](t, baseUrl, fmt.Sprintf("/api/requery.php?user_id=%d&user_key=%s", userid, userkey))
	tt.AssertEqual(t, "success", true, rq5.Success)
	tt.AssertEqual(t, "count", 3, rq5.Count)
	tt.AssertMappedSet(t, "data", []string{"r2", "r3", "r5"}, rq5.Data, "usr_msg_id")

	a7 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/ack.php?user_id=%d&user_key=%s&scn_msg_id=%d", userid, userkey, int(r2["scn_msg_id"].(float64))))
	tt.AssertEqual(t, "success", true, a7["success"])

	a3 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/ack.php?user_id=%d&user_key=%s&scn_msg_id=%d", userid, userkey, int(r3["scn_msg_id"].(float64))))
	tt.AssertEqual(t, "success", true, a3["success"])
	tt.AssertEqual(t, "prev_ack", 0, a3["prev_ack"])
	tt.AssertEqual(t, "new_ack", 1, a3["new_ack"])

	a4 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/ack.php?user_id=%d&user_key=%s&scn_msg_id=%d", userid, userkey, int(r3["scn_msg_id"].(float64))))
	tt.AssertEqual(t, "success", true, a4["success"])
	tt.AssertEqual(t, "prev_ack", 1, a4["prev_ack"])
	tt.AssertEqual(t, "new_ack", 1, a4["new_ack"])

	a5 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/ack.php?user_id=%d&user_key=%s&scn_msg_id=%d", userid, userkey, int(r5["scn_msg_id"].(float64))))
	tt.AssertEqual(t, "success", true, a5["success"])
	tt.AssertEqual(t, "prev_ack", 0, a5["prev_ack"])
	tt.AssertEqual(t, "new_ack", 1, a5["new_ack"])

	r6 := tt.RequestPost[gin.H](t, baseUrl, "/", gin.H{
		"user_id": useridnew,
		"key":     userkey,
		"title":   "HelloWorld_001",
		"msg_id":  "r6",
	})
	tt.AssertEqual(t, "success", true, r6["success"])

	rq6 := tt.RequestGet[respRequery](t, baseUrl, fmt.Sprintf("/api/requery.php?user_id=%d&user_key=%s", userid, userkey))
	tt.AssertEqual(t, "success", true, rq6.Success)
	tt.AssertEqual(t, "count", 1, rq6.Count)
	tt.AssertMappedSet(t, "data", []string{"r6"}, rq6.Data, "usr_msg_id")

	a6 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/ack.php?user_id=%d&user_key=%s&scn_msg_id=%d", userid, userkey, tt.ConvertToCompatID(t, ws, r6["scn_msg_id"].(string))))
	tt.AssertEqual(t, "success", true, a6["success"])
	tt.AssertEqual(t, "prev_ack", 0, a6["prev_ack"])
	tt.AssertEqual(t, "new_ack", 1, a6["new_ack"])
	tt.AssertEqual(t, "message", "ok", a6["message"])

	rq7 := tt.RequestGet[respRequery](t, baseUrl, fmt.Sprintf("/api/requery.php?user_id=%d&user_key=%s", userid, userkey))
	tt.AssertEqual(t, "success", true, rq7.Success)
	tt.AssertEqual(t, "count", 0, rq7.Count)

}

func TestCompatAckCount(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/register.php?fcm_token=%s&pro=%s&pro_token=%s", "DUMMY_FCM", "0", ""))
	tt.AssertEqual(t, "success", true, r0["success"])

	userid := int64(r0["user_id"].(float64))
	userkey := r0["user_key"].(string)

	{
		ri1 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/info.php?user_id=%d&user_key=%s", userid, userkey))
		tt.AssertEqual(t, "unack_count", 0, ri1["unack_count"])
	}

	r1 := tt.RequestPost[gin.H](t, baseUrl, "/send.php", tt.FormData{
		"user_id":  fmt.Sprintf("%d", userid),
		"user_key": userkey,
		"title":    "my title 11 & x",
	})
	tt.AssertEqual(t, "success", true, r1["success"])
	r1scnid := int64(r1["scn_msg_id"].(float64))

	{
		ri1 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/info.php?user_id=%d&user_key=%s", userid, userkey))
		tt.AssertEqual(t, "unack_count", 1, ri1["unack_count"])
	}

	r2 := tt.RequestPost[gin.H](t, baseUrl, "/send.php", tt.FormData{
		"user_id":  fmt.Sprintf("%d", userid),
		"user_key": userkey,
		"title":    "my title 11 & x",
	})
	tt.AssertEqual(t, "success", true, r2["success"])
	r2scnid := int64(r2["scn_msg_id"].(float64))

	{
		ri1 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/info.php?user_id=%d&user_key=%s", userid, userkey))
		tt.AssertEqual(t, "unack_count", 2, ri1["unack_count"])
	}

	{
		ack := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/ack.php?user_id=%d&user_key=%s&scn_msg_id=%d", userid, userkey, r1scnid))
		tt.AssertEqual(t, "success", true, ack["success"])
		tt.AssertEqual(t, "prev_ack", 0, ack["prev_ack"])
		tt.AssertEqual(t, "new_ack", 1, ack["new_ack"])
		tt.AssertEqual(t, "message", "ok", ack["message"])
	}

	{
		ri1 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/info.php?user_id=%d&user_key=%s", userid, userkey))
		tt.AssertEqual(t, "unack_count", 1, ri1["unack_count"])
	}

	{
		ack := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/ack.php?user_id=%d&user_key=%s&scn_msg_id=%d", userid, userkey, r1scnid))
		tt.AssertEqual(t, "success", true, ack["success"])
		tt.AssertEqual(t, "prev_ack", 1, ack["prev_ack"])
		tt.AssertEqual(t, "new_ack", 1, ack["new_ack"])
		tt.AssertEqual(t, "message", "ok", ack["message"])
	}

	{
		ri1 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/info.php?user_id=%d&user_key=%s", userid, userkey))
		tt.AssertEqual(t, "unack_count", 1, ri1["unack_count"])
	}

	{
		ri1 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/info.php?user_id=%d&user_key=%s", userid, userkey))
		tt.AssertEqual(t, "unack_count", 1, ri1["unack_count"])
	}

	{
		ack := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/ack.php?user_id=%d&user_key=%s&scn_msg_id=%d", userid, userkey, r2scnid))
		tt.AssertEqual(t, "success", true, ack["success"])
		tt.AssertEqual(t, "prev_ack", 0, ack["prev_ack"])
		tt.AssertEqual(t, "new_ack", 1, ack["new_ack"])
		tt.AssertEqual(t, "message", "ok", ack["message"])
	}

	{
		ri1 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/info.php?user_id=%d&user_key=%s", userid, userkey))
		tt.AssertEqual(t, "unack_count", 0, ri1["unack_count"])
	}

}
