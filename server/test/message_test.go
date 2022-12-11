package test

import (
	tt "blackforestbytes.com/simplecloudnotifier/test/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"net/url"
	"testing"
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

//TODO test missing message-xx methods
