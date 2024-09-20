package test

import (
	tt "blackforestbytes.com/simplecloudnotifier/test/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"testing"
)

func TestListSenderNames(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	type sn struct {
		SenderName     string `json:"name"`
		LastTimestamp  string `json:"last_timestamp"`
		FirstTimestamp string `json:"first_timestamp"`
		Count          int    `json:"count"`
	}
	type snlistS struct {
		SNList []sn `json:"sender_names"`
	}
	type snlistH struct {
		SNList []gin.H `json:"sender_names"`
	}

	responses := []struct {
		Idx  int
		Resp []string
	}{
		{0, []string{"Pocket Pal", "Cellular Confidant", "Mobile Mate"}},
		{1, []string{}},
		{2, []string{}},
		{3, []string{}},
		{4, []string{"Server0"}},
		{5, []string{"example.org", "example.com", "localhost"}},
		{6, []string{"server1", "server2"}},
		{7, []string{"localhost"}},
		{8, []string{}},
		{9, []string{"Vincent", "Tim", "Max"}},
		{10, []string{}},
		{11, []string{"192.168.0.1", "#S0", "localhost"}},
		{12, []string{}},
		{13, []string{}},
		{14, []string{"dummy-man"}},
		{15, []string{"dummy-man"}},
		{16, []string{}},
	}

	for _, resp := range responses {
		msgList := tt.RequestAuthGet[snlistH](t, data.User[resp.Idx].AdminKey, baseUrl, "/api/v2/sender-names")
		tt.AssertMappedArr(t, "sender_names_"+strconv.Itoa(resp.Idx), resp.Resp, msgList.SNList, "name")
	}
}

func TestListUserSenderNames(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	type sn struct {
		SenderName     string `json:"name"`
		LastTimestamp  string `json:"last_timestamp"`
		FirstTimestamp string `json:"first_timestamp"`
		Count          int    `json:"count"`
	}
	type snlistS struct {
		SNList []sn `json:"sender_names"`
	}
	type snlistH struct {
		SNList []gin.H `json:"sender_names"`
	}

	responses := []struct {
		Idx  int
		Resp []string
	}{
		{0, []string{"Pocket Pal", "Cellular Confidant", "Mobile Mate"}},
		{1, []string{}},
		{2, []string{}},
		{3, []string{}},
		{4, []string{"Server0"}},
		{5, []string{"example.org", "example.com", "localhost"}},
		{6, []string{"server1", "server2"}},
		{7, []string{"localhost"}},
		{8, []string{}},
		{9, []string{"Vincent", "Tim", "Max"}},
		{10, []string{}},
		{11, []string{"192.168.0.1", "#S0", "localhost"}},
		{12, []string{}},
		{13, []string{}},
		{14, []string{}},
		{15, []string{"dummy-man"}},
		{16, []string{}},
	}

	for _, resp := range responses {
		msgList := tt.RequestAuthGet[snlistH](t, data.User[resp.Idx].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/sender-names", data.User[resp.Idx].UID))
		tt.AssertMappedArr(t, "sender_names_"+strconv.Itoa(resp.Idx), resp.Resp, msgList.SNList, "name")
	}
}
