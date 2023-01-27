package test

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	tt "blackforestbytes.com/simplecloudnotifier/test/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
	"testing"
)

func TestGetClient(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/v2/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := fmt.Sprintf("%v", r0["user_id"])

	tt.AssertEqual(t, "len(clients)", 1, len(r0["clients"].([]any)))

	admintok := r0["admin_key"].(string)

	fmt.Printf("uid       := %s\n", uid)
	fmt.Printf("admin_key := %s\n", admintok)

	r1 := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/v2/users/"+uid)

	tt.AssertEqual(t, "uid", uid, fmt.Sprintf("%v", r1["user_id"]))
	tt.AssertEqual(t, "admin_key", admintok, r1["admin_key"])
	tt.AssertEqual(t, "username", nil, r1["username"])

	type rt2 struct {
		Clients []gin.H `json:"clients"`
	}

	r2 := tt.RequestAuthGet[rt2](t, admintok, baseUrl, "/api/v2/users/"+uid+"/clients")

	tt.AssertEqual(t, "len(clients)", 1, len(r2.Clients))

	c0 := r2.Clients[0]

	tt.AssertEqual(t, "agent_model", "DUMMY_PHONE", c0["agent_model"])
	tt.AssertEqual(t, "agent_version", "4X", c0["agent_version"])
	tt.AssertEqual(t, "fcm_token", "DUMMY_FCM", c0["fcm_token"])
	tt.AssertEqual(t, "client_type", "ANDROID", c0["type"])
	tt.AssertEqual(t, "user_id", uid, fmt.Sprintf("%v", c0["user_id"]))

	cid := fmt.Sprintf("%v", c0["client_id"])

	r3 := tt.RequestAuthGet[gin.H](t, admintok, baseUrl, "/api/v2/users/"+uid+"/clients/"+cid)

	tt.AssertJsonMapEqual(t, "client", r3, c0)
}

func TestCreateAndDeleteClient(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/v2/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM",
	})

	uid := fmt.Sprintf("%v", r0["user_id"])

	tt.AssertEqual(t, "len(clients)", 1, len(r0["clients"].([]any)))

	admintok := r0["admin_key"].(string)

	fmt.Printf("uid       := %s\n", uid)
	fmt.Printf("admin_key := %s\n", admintok)

	r2 := tt.RequestAuthPost[gin.H](t, admintok, baseUrl, "/api/v2/users/"+uid+"/clients", gin.H{
		"agent_model":   "DUMMY_PHONE_2",
		"agent_version": "99X",
		"client_type":   "IOS",
		"fcm_token":     "DUMMY_FCM_2",
	})

	cid2 := fmt.Sprintf("%v", r2["client_id"])

	type rt3 struct {
		Clients []gin.H `json:"clients"`
	}

	r3 := tt.RequestAuthGet[rt3](t, admintok, baseUrl, "/api/v2/users/"+uid+"/clients")
	tt.AssertEqual(t, "len(clients)", 2, len(r3.Clients))

	r4 := tt.RequestAuthDelete[gin.H](t, admintok, baseUrl, "/api/v2/users/"+uid+"/clients/"+cid2, nil)
	tt.AssertEqual(t, "client_id", cid2, fmt.Sprintf("%v", r4["client_id"]))

	r5 := tt.RequestAuthGet[rt3](t, admintok, baseUrl, "/api/v2/users/"+uid+"/clients")
	tt.AssertEqual(t, "len(clients)", 1, len(r5.Clients))
}

func TestReuseFCM(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestPost[gin.H](t, baseUrl, "/api/v2/users", gin.H{
		"agent_model":   "DUMMY_PHONE",
		"agent_version": "4X",
		"client_type":   "ANDROID",
		"fcm_token":     "DUMMY_FCM_001",
	})

	uid := fmt.Sprintf("%v", r0["user_id"])

	tt.AssertEqual(t, "len(clients)", 1, len(r0["clients"].([]any)))

	admintok := r0["admin_key"].(string)

	fmt.Printf("uid       := %s\n", uid)
	fmt.Printf("admin_key := %s\n", admintok)

	type rt2 struct {
		Clients []gin.H `json:"clients"`
	}

	r1 := tt.RequestAuthGet[rt2](t, admintok, baseUrl, "/api/v2/users/"+uid+"/clients")

	tt.AssertEqual(t, "len(clients)", 1, len(r1.Clients))

	r2 := tt.RequestAuthPost[gin.H](t, admintok, baseUrl, "/api/v2/users/"+uid+"/clients", gin.H{
		"agent_model":   "DUMMY_PHONE_2",
		"agent_version": "99X",
		"client_type":   "IOS",
		"fcm_token":     "DUMMY_FCM_001",
	})

	cid2 := fmt.Sprintf("%v", r2["client_id"])

	type rt3 struct {
		Clients []gin.H `json:"clients"`
	}

	r3 := tt.RequestAuthGet[rt3](t, admintok, baseUrl, "/api/v2/users/"+uid+"/clients")
	tt.AssertEqual(t, "len(clients)", 1, len(r3.Clients))

	tt.AssertEqual(t, "clients->client_id", cid2, fmt.Sprintf("%v", r3.Clients[0]["client_id"]))
}

func TestListClients(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	type clientlist struct {
		Clients []gin.H `json:"clients"`
	}

	type T2 struct {
		CT []string
		AM []string
		AV []string
	}

	vals := map[int]T2{
		1:  {[]string{"ANDROID"}, []string{"Galaxy Quest"}, []string{"2022"}},
		2:  {[]string{"IOS", "IOS"}, []string{"GalaxySurfer", "Ocean Explorer"}, []string{"Triple-XXX", "737edc01"}},
		3:  {[]string{"ANDROID"}, []string{"Snow Leopard"}, []string{"1.0.1.99~3"}},
		8:  {[]string{"ANDROID"}, []string{"Galaxy Quest"}, []string{"2023.1"}},
		9:  {[]string{"ANDROID", "IOS", "IOS", "ANDROID"}, []string{"Galaxy Quest", "DreamWeaver", "GalaxySurfer", "Galaxy Quest"}, []string{"2023.2", "Triple-XXX", "Triple-XXX", "2023.1"}},
		5:  {[]string{"IOS"}, []string{"Ocean Explorer"}, []string{"737edc01"}},
		7:  {[]string{"ANDROID"}, []string{"Galaxy Quest"}, []string{"2023.1"}},
		10: {[]string{}, []string{}, []string{}},
		14: {[]string{"IOS"}, []string{"StarfireXX"}, []string{"1.x"}},
		11: {[]string{}, []string{}, []string{}},
		12: {[]string{"IOS"}, []string{"Ocean Explorer"}, []string{"737edc01"}},
		13: {[]string{}, []string{}, []string{}},
		0:  {[]string{"IOS"}, []string{"Starfire"}, []string{"2.0"}},
		4:  {[]string{"ANDROID"}, []string{"Thunder-Bolt-4$"}, []string{"#12"}},
		6:  {[]string{"IOS", "IOS"}, []string{"GalaxySurfer", "Cyber Nova"}, []string{"Triple-XXX", "Cyber 4"}},
		15: {[]string{"IOS"}, []string{"StarfireXX"}, []string{"1.x"}},
	}

	for k, v := range vals {
		clist1 := tt.RequestAuthGet[clientlist](t, data.User[k].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/clients", url.QueryEscape(data.User[k].UID)))
		tt.AssertMappedSet(t, fmt.Sprintf("clients[%d]->type", k), v.CT, clist1.Clients, "type")
		tt.AssertMappedSet(t, fmt.Sprintf("clients[%d]->agent_model", k), v.AM, clist1.Clients, "agent_model")
		tt.AssertMappedSet(t, fmt.Sprintf("clients[%d]->agent_version", k), v.AV, clist1.Clients, "agent_version")
	}

	tt.RequestAuthGetShouldFail(t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/clients", url.QueryEscape(data.User[1].UID)), 401, apierr.USER_AUTH_FAILED)
}
