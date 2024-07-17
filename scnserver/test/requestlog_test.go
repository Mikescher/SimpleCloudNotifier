package test

import (
	ct "blackforestbytes.com/simplecloudnotifier/db/cursortoken"
	"blackforestbytes.com/simplecloudnotifier/models"
	tt "blackforestbytes.com/simplecloudnotifier/test/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
	"testing"
	"time"
)

func TestRequestLogSimple(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	ctx := ws.NewSimpleTransactionContext(5 * time.Second)
	defer ctx.Cancel()

	// Ping
	{
		tt.RequestGet[tt.Void](t, baseUrl, fmt.Sprintf("/api/ping"))
		time.Sleep(100 * time.Millisecond)

		rl, _, err := ws.Database.Requests.ListRequestLogs(ctx, models.RequestLogFilter{}, nil, ct.Start())
		tt.TestFailIfErr(t, err)

		tt.AssertEqual(t, "requestlog.count", 1, len(rl))

		tt.AssertEqual(t, "requestlog[0].Method", "GET", rl[0].Method)
		tt.AssertEqual(t, "requestlog[0].KeyID", nil, rl[0].KeyID)
		tt.AssertEqual(t, "requestlog[0].UserID", nil, rl[0].UserID)
		tt.AssertEqual(t, "requestlog[0].Panicked", false, rl[0].Panicked)
		tt.AssertEqual(t, "requestlog[0].URI", "/api/ping", rl[0].URI)
		tt.AssertEqual(t, "requestlog[0].ResponseContentType", "application/json", rl[0].ResponseContentType)
	}

	// HTMl request
	{
		tt.RequestRaw(t, baseUrl, fmt.Sprintf("/"))
		time.Sleep(100 * time.Millisecond)

		rl, _, err := ws.Database.Requests.ListRequestLogs(ctx, models.RequestLogFilter{}, nil, ct.Start())
		tt.TestFailIfErr(t, err)

		tt.AssertEqual(t, "requestlog.count", 2, len(rl))

		tt.AssertEqual(t, "requestlog[0].Method", "GET", rl[0].Method)
		tt.AssertEqual(t, "requestlog[0].KeyID", nil, rl[0].KeyID)
		tt.AssertEqual(t, "requestlog[0].UserID", nil, rl[0].UserID)
		tt.AssertEqual(t, "requestlog[0].Panicked", false, rl[0].Panicked)
		tt.AssertEqual(t, "requestlog[0].URI", "/", rl[0].URI)
		tt.AssertEqual(t, "requestlog[0].ResponseContentType", "text/html", rl[0].ResponseContentType)
	}

	type R struct {
		Clients []struct {
			ClientId string `json:"client_id"`
			UserId   string `json:"user_id"`
		} `json:"clients"`
		ReadKey  string `json:"read_key"`
		SendKey  string `json:"send_key"`
		AdminKey string `json:"admin_key"`
		UserId   string `json:"user_id"`
	}
	usr := tt.RequestPost[R](t, baseUrl, "/api/v2/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})
	time.Sleep(100 * time.Millisecond)

	// API request
	{

		tt.RequestAuthGet[R](t, usr.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s", usr.UserId))
		time.Sleep(100 * time.Millisecond)

		rl, _, err := ws.Database.Requests.ListRequestLogs(ctx, models.RequestLogFilter{}, nil, ct.Start())
		tt.TestFailIfErr(t, err)

		tt.AssertEqual(t, "requestlog.count", 4, len(rl))

		tt.AssertEqual(t, "requestlog[0].Method", "GET", rl[0].Method)
		tt.AssertNotEqual(t, "requestlog[0].KeyID", nil, rl[0].KeyID)
		tt.AssertStrRepEqual(t, "requestlog[0].UserID", usr.UserId, rl[0].UserID)
		tt.AssertEqual(t, "requestlog[0].Panicked", false, rl[0].Panicked)
		tt.AssertStrRepEqual(t, "requestlog[0].Permissions", "A", rl[0].Permissions)
		tt.AssertEqual(t, "requestlog[0].URI", fmt.Sprintf("/api/v2/users/%s", usr.UserId), rl[0].URI)
		tt.AssertEqual(t, "requestlog[0].ResponseContentType", "application/json", rl[0].ResponseContentType)
	}

	// Send request
	{
		tt.RequestPost[gin.H](t, baseUrl, fmt.Sprintf("/?user_id=%s&key=%s&title=%s", usr.UserId, usr.SendKey, url.QueryEscape("Hello World 2134")), nil)
		time.Sleep(100 * time.Millisecond)

		rl, _, err := ws.Database.Requests.ListRequestLogs(ctx, models.RequestLogFilter{}, nil, ct.Start())
		tt.TestFailIfErr(t, err)

		tt.AssertEqual(t, "requestlog.count", 5, len(rl))

		tt.AssertEqual(t, "requestlog[0].Method", "POST", rl[0].Method)
		tt.AssertEqual(t, "requestlog[0].UserID", nil, rl[0].UserID)
		tt.AssertEqual(t, "requestlog[0].ResponseContentType", "application/json", rl[0].ResponseContentType)
	}

	// Failed request
	{
		tt.RequestAuthGetShouldFail(t, usr.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s", models.NewUserID()), 0, 0)
		time.Sleep(100 * time.Millisecond)

		rl, _, err := ws.Database.Requests.ListRequestLogs(ctx, models.RequestLogFilter{}, nil, ct.Start())
		tt.TestFailIfErr(t, err)

		tt.AssertEqual(t, "requestlog.count", 6, len(rl))

		tt.AssertEqual(t, "requestlog[0].Method", "GET", rl[0].Method)
		tt.AssertStrRepEqual(t, "requestlog[0].UserID", usr.UserId, rl[0].UserID)
		tt.AssertEqual(t, "requestlog[0].ResponseContentType", "application/json", rl[0].ResponseContentType)
		tt.AssertStrRepEqual(t, "requestlog[0].ResponseStatuscode", 401, rl[0].ResponseStatuscode)
	}

}

func TestRequestLogAPI(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)
	time.Sleep(900 * time.Millisecond)

	ctx := ws.NewSimpleTransactionContext(5 * time.Second)
	defer ctx.Cancel()

	rl1, _, err := ws.Database.Requests.ListRequestLogs(ctx, models.RequestLogFilter{}, nil, ct.Start())
	tt.TestFailIfErr(t, err)

	tt.RequestAuthGet[gin.H](t, data.User[0].ReadKey, baseUrl, "/api/v2/users/"+data.User[0].UID)
	time.Sleep(900 * time.Millisecond)

	rl2, _, err := ws.Database.Requests.ListRequestLogs(ctx, models.RequestLogFilter{}, nil, ct.Start())
	tt.TestFailIfErr(t, err)

	tt.AssertEqual(t, "requestlog.count", len(rl1)+1, len(rl2))
}
