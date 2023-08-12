package test

import (
	"blackforestbytes.com/simplecloudnotifier/push"
	tt "blackforestbytes.com/simplecloudnotifier/test/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"testing"
	"time"
)

func TestUptimeKumaDown(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	pusher := ws.Pusher.(*push.TestSink)

	suffix := fmt.Sprintf("/external/v1/uptime-kuma?user_id=%v&key=%v", data.UID, data.SendKey)
	_ = tt.RequestPost[gin.H](t, baseUrl, suffix, gin.H{
		"msg": "[test] [🔴 Down] getaddrinfo ENOTFOUND exampleasdsda.com",
		"heartbeat": gin.H{
			"status": 0,
			"msg":    "getaddrinfo ENOTFOUND exampleasdsda.com",
		},
		"monitor": gin.H{
			"name": "test",
		},
	})

	tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.title", "Monitor test went down!", pusher.Last().Message.Title)

	type mglist struct {
		Messages []gin.H `json:"messages"`
	}

	msgList1 := tt.RequestAuthGet[mglist](t, data.AdminKey, baseUrl, "/api/v2/messages")
	tt.AssertEqual(t, "len(messages)", 1, len(msgList1.Messages))
	tt.AssertStrRepEqual(t, "msg.title", "Monitor test went down!", msgList1.Messages[0]["title"])
	tt.AssertStrRepEqual(t, "msg.content", "getaddrinfo ENOTFOUND exampleasdsda.com", msgList1.Messages[0]["content"])
}

func TestUptimeKumaUp(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	pusher := ws.Pusher.(*push.TestSink)

	suffix := fmt.Sprintf("/external/v1/uptime-kuma?user_id=%v&key=%v", data.UID, data.SendKey)
	_ = tt.RequestPost[gin.H](t, baseUrl, suffix, gin.H{
		"msg": "[test] [✅ Up] 200 - OK",
		"heartbeat": gin.H{
			"status": 1,
			"msg":    "200 - OK",
		},
		"monitor": gin.H{
			"name": "test",
		},
	})

	tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.title", "Monitor test is back online", pusher.Last().Message.Title)

	type mglist struct {
		Messages []gin.H `json:"messages"`
	}

	msgList1 := tt.RequestAuthGet[mglist](t, data.AdminKey, baseUrl, "/api/v2/messages")
	tt.AssertEqual(t, "len(messages)", 1, len(msgList1.Messages))
	tt.AssertStrRepEqual(t, "msg.title", "Monitor test is back online", msgList1.Messages[0]["title"])
	tt.AssertStrRepEqual(t, "msg.content", "200 - OK", msgList1.Messages[0]["content"])
}

func TestUptimeKumaFullDown(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	pusher := ws.Pusher.(*push.TestSink)

	ts := time.Now().Add(-time.Hour).Format("2006-01-02 15:04:05")

	suffix := fmt.Sprintf("/external/v1/uptime-kuma?user_id=%v&key=%v", data.UID, data.SendKey)
	_ = tt.RequestPost[gin.H](t, baseUrl, suffix, tt.RawJSON{ContentType: "application/json", Body: `{"heartbeat":{"monitorID":89,"status":0,"time":"` + ts + `","msg":"timeout of 16000ms exceeded","important":true,"duration":36,"timezone":"Europe/Berlin","timezoneOffset":"+02:00","localDateTime":"` + ts + `"},"monitor":{"id":89,"name":"test","description":null,"pathName":"test","parent":null,"childrenIDs":[],"url":"https://exampleXYZ.com","method":"GET","hostname":null,"port":null,"maxretries":1,"weight":2000,"active":true,"forceInactive":false,"type":"http","interval":20,"retryInterval":20,"resendInterval":0,"keyword":null,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"2":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"httpBodyEncoding":"json","includeSensitiveData":false},"msg":"[test] [🔴 Down] timeout of 16000ms exceeded"}`})

	tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.title", "Monitor test went down!", pusher.Last().Message.Title)
	tt.AssertStrRepEqual(t, "msg.title", "timeout of 16000ms exceeded", pusher.Last().Message.Content)

	type mglist struct {
		Messages []gin.H `json:"messages"`
	}

	msgList1 := tt.RequestAuthGet[mglist](t, data.AdminKey, baseUrl, "/api/v2/messages")
	tt.AssertEqual(t, "len(messages)", 1, len(msgList1.Messages))
	tt.AssertStrRepEqual(t, "msg.title", "Monitor test went down!", msgList1.Messages[0]["title"])
	tt.AssertStrRepEqual(t, "msg.content", "timeout of 16000ms exceeded", msgList1.Messages[0]["content"])
}

func TestUptimeKumaFullUp(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	pusher := ws.Pusher.(*push.TestSink)

	ts := time.Now().Add(-time.Hour).Format("2006-01-02 15:04:05")

	suffix := fmt.Sprintf("/external/v1/uptime-kuma?user_id=%v&key=%v", data.UID, data.SendKey)
	_ = tt.RequestPost[gin.H](t, baseUrl, suffix, tt.RawJSON{ContentType: "application/json", Body: `{"heartbeat":{"monitorID":89,"status":1,"time":"` + ts + `","msg":"200 - OK","ping":55,"important":true,"duration":41,"timezone":"Europe/Berlin","timezoneOffset":"+02:00","localDateTime":"` + ts + `"},"monitor":{"id":89,"name":"test","description":null,"pathName":"test","parent":null,"childrenIDs":[],"url":"https://example.com","method":"GET","hostname":null,"port":null,"maxretries":1,"weight":2000,"active":true,"forceInactive":false,"type":"http","interval":20,"retryInterval":20,"resendInterval":0,"keyword":null,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"2":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"httpBodyEncoding":"json","includeSensitiveData":false},"msg":"[test] [✅ Up] 200 - OK"}`})

	tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.title", "Monitor test is back online", pusher.Last().Message.Title)

	type mglist struct {
		Messages []gin.H `json:"messages"`
	}

	msgList1 := tt.RequestAuthGet[mglist](t, data.AdminKey, baseUrl, "/api/v2/messages")
	tt.AssertEqual(t, "len(messages)", 1, len(msgList1.Messages))
	tt.AssertStrRepEqual(t, "msg.title", "Monitor test is back online", msgList1.Messages[0]["title"])
	tt.AssertStrRepEqual(t, "msg.content", "200 - OK", msgList1.Messages[0]["content"])
}

func TestUptimeKumaChannelNone(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	pusher := ws.Pusher.(*push.TestSink)

	suffix := fmt.Sprintf("/external/v1/uptime-kuma?user_id=%v&key=%v", data.UID, data.SendKey)
	_ = tt.RequestPost[gin.H](t, baseUrl, suffix, gin.H{
		"msg": "[test] [🔴 Down] getaddrinfo ENOTFOUND exampleasdsda.com",
		"heartbeat": gin.H{
			"status": 0,
			"msg":    "getaddrinfo ENOTFOUND exampleasdsda.com",
		},
		"monitor": gin.H{
			"name": "test",
		},
	})

	tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.channel", "main", pusher.Last().Message.ChannelInternalName)
}

func TestUptimeKumaChannelSingle(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	pusher := ws.Pusher.(*push.TestSink)

	suffix := fmt.Sprintf("/external/v1/uptime-kuma?user_id=%v&key=%v&channel=CTEST", data.UID, data.SendKey)
	_ = tt.RequestPost[gin.H](t, baseUrl, suffix, gin.H{
		"msg": "[test] [🔴 Down] getaddrinfo ENOTFOUND exampleasdsda.com",
		"heartbeat": gin.H{
			"status": 0,
			"msg":    "getaddrinfo ENOTFOUND exampleasdsda.com",
		},
		"monitor": gin.H{
			"name": "test",
		},
	})

	tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.channel", "CTEST", pusher.Last().Message.ChannelInternalName)
}

func TestUptimeKumaChannelAllDown(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	pusher := ws.Pusher.(*push.TestSink)

	suffix := fmt.Sprintf("/external/v1/uptime-kuma?user_id=%v&key=%v&channel=CTEST&channel_up=CTEST_UP&channel_down=CTEST_DOWN", data.UID, data.SendKey)
	_ = tt.RequestPost[gin.H](t, baseUrl, suffix, gin.H{
		"msg": "[test] [🔴 Down] getaddrinfo ENOTFOUND exampleasdsda.com",
		"heartbeat": gin.H{
			"status": 0,
			"msg":    "getaddrinfo ENOTFOUND exampleasdsda.com",
		},
		"monitor": gin.H{
			"name": "test",
		},
	})

	tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.channel", "CTEST_DOWN", pusher.Last().Message.ChannelInternalName)
}

func TestUptimeKumaChannelSpecDown(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	pusher := ws.Pusher.(*push.TestSink)

	suffix := fmt.Sprintf("/external/v1/uptime-kuma?user_id=%v&key=%v&channel_up=CTEST_UP&channel_down=CTEST_DOWN", data.UID, data.SendKey)
	_ = tt.RequestPost[gin.H](t, baseUrl, suffix, gin.H{
		"msg": "[test] [🔴 Down] getaddrinfo ENOTFOUND exampleasdsda.com",
		"heartbeat": gin.H{
			"status": 0,
			"msg":    "getaddrinfo ENOTFOUND exampleasdsda.com",
		},
		"monitor": gin.H{
			"name": "test",
		},
	})

	tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.channel", "CTEST_DOWN", pusher.Last().Message.ChannelInternalName)
}

func TestUptimeKumaChannelAllUp(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	pusher := ws.Pusher.(*push.TestSink)

	suffix := fmt.Sprintf("/external/v1/uptime-kuma?user_id=%v&key=%v&channel=CTEST&channel_up=CTEST_UP&channel_down=CTEST_DOWN", data.UID, data.SendKey)
	_ = tt.RequestPost[gin.H](t, baseUrl, suffix, gin.H{
		"msg": "[test] [✅ Up] 200 - OK",
		"heartbeat": gin.H{
			"status": 1,
			"msg":    "200 - OK",
		},
		"monitor": gin.H{
			"name": "test",
		},
	})

	tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.channel", "CTEST_UP", pusher.Last().Message.ChannelInternalName)
}

func TestUptimeKumaChannelSpecUp(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	pusher := ws.Pusher.(*push.TestSink)

	suffix := fmt.Sprintf("/external/v1/uptime-kuma?user_id=%v&key=%v&channel_up=CTEST_UP&channel_down=CTEST_DOWN", data.UID, data.SendKey)
	_ = tt.RequestPost[gin.H](t, baseUrl, suffix, gin.H{
		"msg": "[test] [✅ Up] 200 - OK",
		"heartbeat": gin.H{
			"status": 1,
			"msg":    "200 - OK",
		},
		"monitor": gin.H{
			"name": "test",
		},
	})

	tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.channel", "CTEST_UP", pusher.Last().Message.ChannelInternalName)
}

func TestUptimeKumaPriorityNone(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	pusher := ws.Pusher.(*push.TestSink)

	suffix := fmt.Sprintf("/external/v1/uptime-kuma?user_id=%v&key=%v", data.UID, data.SendKey)
	_ = tt.RequestPost[gin.H](t, baseUrl, suffix, gin.H{
		"msg": "[test] [🔴 Down] getaddrinfo ENOTFOUND exampleasdsda.com",
		"heartbeat": gin.H{
			"status": 0,
			"msg":    "getaddrinfo ENOTFOUND exampleasdsda.com",
		},
		"monitor": gin.H{
			"name": "test",
		},
	})

	tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.channel", 1, pusher.Last().Message.Priority)
}

func TestUptimeKumaPrioritySingle(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	pusher := ws.Pusher.(*push.TestSink)

	suffix0 := fmt.Sprintf("/external/v1/uptime-kuma?user_id=%v&key=%v&priority=0", data.UID, data.SendKey)
	_ = tt.RequestPost[gin.H](t, baseUrl, suffix0, gin.H{
		"msg": "[test] [🔴 Down] getaddrinfo ENOTFOUND exampleasdsda.com",
		"heartbeat": gin.H{
			"status": 0,
			"msg":    "getaddrinfo ENOTFOUND exampleasdsda.com",
		},
		"monitor": gin.H{
			"name": "test",
		},
	})

	tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.prio", 0, pusher.Last().Message.Priority)

	suffix1 := fmt.Sprintf("/external/v1/uptime-kuma?user_id=%v&key=%v&priority=1", data.UID, data.SendKey)
	_ = tt.RequestPost[gin.H](t, baseUrl, suffix1, gin.H{
		"msg": "[test] [🔴 Down] getaddrinfo ENOTFOUND exampleasdsda.com",
		"heartbeat": gin.H{
			"status": 0,
			"msg":    "getaddrinfo ENOTFOUND exampleasdsda.com",
		},
		"monitor": gin.H{
			"name": "test",
		},
	})

	tt.AssertEqual(t, "messageCount", 2, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.prio", 1, pusher.Last().Message.Priority)

	suffix2 := fmt.Sprintf("/external/v1/uptime-kuma?user_id=%v&key=%v&priority=2", data.UID, data.SendKey)
	_ = tt.RequestPost[gin.H](t, baseUrl, suffix2, gin.H{
		"msg": "[test] [🔴 Down] getaddrinfo ENOTFOUND exampleasdsda.com",
		"heartbeat": gin.H{
			"status": 0,
			"msg":    "getaddrinfo ENOTFOUND exampleasdsda.com",
		},
		"monitor": gin.H{
			"name": "test",
		},
	})

	tt.AssertEqual(t, "messageCount", 3, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.prio", 2, pusher.Last().Message.Priority)
}

func TestUptimeKumaPriorityAllDown(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	pusher := ws.Pusher.(*push.TestSink)

	suffix := fmt.Sprintf("/external/v1/uptime-kuma?user_id=%v&key=%v&priority=1&priority_up=2&priority_down=0", data.UID, data.SendKey)
	_ = tt.RequestPost[gin.H](t, baseUrl, suffix, gin.H{
		"msg": "[test] [🔴 Down] getaddrinfo ENOTFOUND exampleasdsda.com",
		"heartbeat": gin.H{
			"status": 0,
			"msg":    "getaddrinfo ENOTFOUND exampleasdsda.com",
		},
		"monitor": gin.H{
			"name": "test",
		},
	})

	tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.prio", 0, pusher.Last().Message.Priority)
}

func TestUptimeKumaPrioritySpecDown(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	pusher := ws.Pusher.(*push.TestSink)

	suffix := fmt.Sprintf("/external/v1/uptime-kuma?user_id=%v&key=%v&priority_up=2&priority_down=0", data.UID, data.SendKey)
	_ = tt.RequestPost[gin.H](t, baseUrl, suffix, gin.H{
		"msg": "[test] [🔴 Down] getaddrinfo ENOTFOUND exampleasdsda.com",
		"heartbeat": gin.H{
			"status": 0,
			"msg":    "getaddrinfo ENOTFOUND exampleasdsda.com",
		},
		"monitor": gin.H{
			"name": "test",
		},
	})

	tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.prio", 0, pusher.Last().Message.Priority)
}

func TestUptimeKumaPriorityAllUp(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	pusher := ws.Pusher.(*push.TestSink)

	suffix := fmt.Sprintf("/external/v1/uptime-kuma?user_id=%v&key=%v&priority=1&priority_up=2&priority_down=0", data.UID, data.SendKey)
	_ = tt.RequestPost[gin.H](t, baseUrl, suffix, gin.H{
		"msg": "[test] [✅ Up] 200 - OK",
		"heartbeat": gin.H{
			"status": 1,
			"msg":    "200 - OK",
		},
		"monitor": gin.H{
			"name": "test",
		},
	})

	tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.prio", 2, pusher.Last().Message.Priority)
}

func TestUptimeKumaPrioritySpecUp(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	pusher := ws.Pusher.(*push.TestSink)

	suffix := fmt.Sprintf("/external/v1/uptime-kuma?user_id=%v&key=%v&priority_up=2&priority_down=0", data.UID, data.SendKey)
	_ = tt.RequestPost[gin.H](t, baseUrl, suffix, gin.H{
		"msg": "[test] [✅ Up] 200 - OK",
		"heartbeat": gin.H{
			"status": 1,
			"msg":    "200 - OK",
		},
		"monitor": gin.H{
			"name": "test",
		},
	})

	tt.AssertEqual(t, "messageCount", 1, len(pusher.Data))
	tt.AssertStrRepEqual(t, "msg.prio", 2, pusher.Last().Message.Priority)
}
