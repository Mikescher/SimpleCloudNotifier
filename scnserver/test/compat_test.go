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

	msgList1 := tt.RequestAuthGet[mglist](t, admintok, baseUrl, "/api/messages")
	tt.AssertEqual(t, "len(messages)", 1, len(msgList1.Messages))

	msg1Get := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/messages/"+fmt.Sprintf("%v", pusher.Last().Message.MessageID))
	tt.AssertStrRepEqual(t, "msg.title", "HelloWorld_001", msg1Get["title"])
	tt.AssertStrRepEqual(t, "msg.channel_internal_name", "main", msg1Get["channel_internal_name"])

	msg2 := tt.RequestPost[gin.H](t, baseUrl, fmt.Sprintf("/send.php?user_key=%s&user_id=%d&title=%s", admintok, uidold, "HelloWorld_002"), nil)

	tt.AssertEqual(t, "messageCount", 2, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.title", "HelloWorld_002", pusher.Last().Message.Title)
	tt.AssertStrRepEqual(t, "msg.content", nil, pusher.Last().Message.Content)

	exp2 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/expand.php?user_id=%d&user_key=%s&scn_msg_id=%d", uidold, admintok, int64(msg2["scn_msg_id"].(float64))))
	tt.AssertEqual(t, "success", true, exp2["success"])

	tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/messages/"+fmt.Sprintf("%v", pusher.Last().Message.MessageID))

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

	tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/messages/"+fmt.Sprintf("%v", pusher.Last().Message.MessageID))

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

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/users", gin.H{
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

	msgList1 := tt.RequestAuthGet[mglist](t, admintok, baseUrl, "/api/messages")
	tt.AssertEqual(t, "len(messages)", 1, len(msgList1.Messages))

	msg1Get := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/messages/"+fmt.Sprintf("%v", pusher.Last().Message.MessageID))
	tt.AssertStrRepEqual(t, "msg.title", "HelloWorld_001", msg1Get["title"])
	tt.AssertStrRepEqual(t, "msg.channel_internal_name", "main", msg1Get["channel_internal_name"])

	msg2 := tt.RequestPost[gin.H](t, baseUrl, fmt.Sprintf("/send.php?user_key=%s&user_id=%d&title=%s", sendtok, uidold, "HelloWorld_002"), nil)

	tt.AssertEqual(t, "messageCount", 2, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.title", "HelloWorld_002", pusher.Last().Message.Title)
	tt.AssertStrRepEqual(t, "msg.content", nil, pusher.Last().Message.Content)

	exp2 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/expand.php?user_id=%d&user_key=%s&scn_msg_id=%d", uidold, admintok, int64(msg2["scn_msg_id"].(float64))))
	tt.AssertEqual(t, "success", true, exp2["success"])

	tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/messages/"+fmt.Sprintf("%v", pusher.Last().Message.MessageID))

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
	tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/messages/"+fmt.Sprintf("%v", pusher.Last().Message.MessageID))

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
	tt.AssertEqual(t, "success", "my title 11 & x", (r1x["data"].(map[string]any))["title"])

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
	tt.AssertEqual(t, "success", "my title", (r2x["data"].(map[string]any))["title"])

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
}

func TestCompatAck(t *testing.T) {
	t.SkipNow() //TODO
}

func TestCompatRequery(t *testing.T) {
	t.SkipNow() //TODO
}

func TestCompatUpdate(t *testing.T) {
	t.SkipNow() //TODO
}

func TestCompatExpand(t *testing.T) {
	t.SkipNow() //TODO
}

func TestCompatUpgrade(t *testing.T) {
	t.SkipNow() //TODO
}
