package test

import (
	tt "blackforestbytes.com/simplecloudnotifier/test/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"testing"
)

func TestResponseChannel(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	response := tt.RequestAuthGetRaw(t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s", data.User[0].UID, data.User[0].Channels[0].ChannelID))

	tt.AssertJsonStructureMatch(t, "json[channel]", response, map[string]any{
		"channel_id":         "id",
		"owner_user_id":      "id",
		"internal_name":      "string",
		"display_name":       "string",
		"description_name":   "null",
		"subscribe_key":      "string",
		"timestamp_created":  "rfc3339",
		"timestamp_lastsent": "rfc3339",
		"messages_sent":      "int",
		"subscription": map[string]any{
			"subscription_id":       "id",
			"subscriber_user_id":    "id",
			"channel_owner_user_id": "id",
			"channel_id":            "id",
			"channel_internal_name": "string",
			"timestamp_created":     "rfc3339",
			"confirmed":             "bool",
		},
	})
}

func TestResponseClient(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	response := tt.RequestAuthGetRaw(t, data.User[2].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/clients/%s", data.User[2].UID, data.User[2].Clients[0]))

	tt.AssertJsonStructureMatch(t, "json[client]", response, map[string]any{
		"client_id":         "id",
		"user_id":           "id",
		"type":              "string",
		"fcm_token":         "string",
		"timestamp_created": "rfc3339",
		"agent_model":       "string",
		"agent_version":     "string",
		"name":              "string|null",
	})
}

func TestResponseKeyToken1(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	response := tt.RequestAuthGetRaw(t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys/%s", data.User[0].UID, data.User[0].Keys[0].KeyID))

	tt.AssertJsonStructureMatch(t, "json[key]", response, map[string]any{
		"keytoken_id":        "id",
		"name":               "string",
		"timestamp_created":  "rfc3339",
		"timestamp_lastused": "rfc3339|null",
		"owner_user_id":      "id",
		"all_channels":       "bool",
		"channels":           []any{"string"},
		"permissions":        "string",
		"messages_sent":      "int",
	})
}

func TestResponseKeyToken2(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	chan1 := tt.RequestAuthPost[gin.H](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", data.UID), gin.H{
		"name": "TestChan1asdf",
	})

	type keyobj struct {
		KeytokenId string `json:"keytoken_id"`
	}
	k0 := tt.RequestAuthPost[keyobj](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", data.UID), gin.H{
		"all_channels": false,
		"channels":     []string{chan1["channel_id"].(string)},
		"name":         "TKey1",
		"permissions":  "CS",
	})

	response := tt.RequestAuthGetRaw(t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys/%s", data.UID, k0.KeytokenId))

	tt.AssertJsonStructureMatch(t, "json[key]", response, map[string]any{
		"keytoken_id":        "id",
		"name":               "string",
		"timestamp_created":  "rfc3339",
		"timestamp_lastused": "rfc3339|null",
		"owner_user_id":      "id",
		"all_channels":       "bool",
		"channels":           []any{"string"},
		"permissions":        "string",
		"messages_sent":      "int",
	})
}

func TestResponseKeyToken3(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	response := tt.RequestAuthGetRaw(t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys/current", data.UID))

	tt.AssertJsonStructureMatch(t, "json[key]", response, map[string]any{
		"keytoken_id":        "id",
		"name":               "string",
		"timestamp_created":  "rfc3339",
		"timestamp_lastused": "rfc3339|null",
		"owner_user_id":      "id",
		"all_channels":       "bool",
		"channels":           []any{"string"},
		"permissions":        "string",
		"messages_sent":      "int",
		"token":              "string",
	})
}

func TestResponseKeyToken4(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	chan1 := tt.RequestAuthPost[gin.H](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", data.UID), gin.H{
		"name": "TestChan1asdf",
	})

	response := tt.RequestAuthPostRaw(t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", data.UID), gin.H{
		"all_channels": false,
		"channels":     []string{chan1["channel_id"].(string)},
		"name":         "TKey1",
		"permissions":  "CS",
	})

	tt.AssertJsonStructureMatch(t, "json[key]", response, map[string]any{
		"keytoken_id":        "id",
		"name":               "string",
		"timestamp_created":  "rfc3339",
		"timestamp_lastused": "rfc3339|null",
		"owner_user_id":      "id",
		"all_channels":       "bool",
		"channels":           []any{"string"},
		"permissions":        "string",
		"messages_sent":      "int",
		"token":              "string",
	})
}

func TestResponseMessage(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	response := tt.RequestAuthGetRaw(t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/messages/%s", data.User[0].Messages[0]))

	tt.AssertJsonStructureMatch(t, "json[message]", response, map[string]any{
		"message_id":            "id",
		"sender_user_id":        "id",
		"channel_internal_name": "string",
		"channel_id":            "id",
		"sender_name":           "string",
		"sender_ip":             "string",
		"timestamp":             "rfc3339",
		"title":                 "string",
		"content":               "null",
		"priority":              "int",
		"usr_message_id":        "null",
		"used_key_id":           "id",
		"trimmed":               "bool",
	})
}

func TestResponseSubscription(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	response := tt.RequestAuthGetRaw(t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data.User[0].UID, data.User[0].Subscriptions[0]))

	tt.AssertJsonStructureMatch(t, "json[subscription]", response, map[string]any{
		"subscription_id":       "id",
		"subscriber_user_id":    "id",
		"channel_owner_user_id": "id",
		"channel_id":            "id",
		"channel_internal_name": "string",
		"timestamp_created":     "rfc3339",
		"confirmed":             "bool",
	})
}

func TestResponseUser(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	response := tt.RequestAuthGetRaw(t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s", data.User[0].UID))

	tt.AssertJsonStructureMatch(t, "json[user]", response, map[string]any{
		"user_id":                        "id",
		"username":                       "null",
		"timestamp_created":              "rfc3339",
		"timestamp_lastread":             "null",
		"timestamp_lastsent":             "rfc3339",
		"messages_sent":                  "int",
		"quota_used":                     "int",
		"quota_remaining":                "int",
		"quota_max":                      "int",
		"is_pro":                         "bool",
		"default_channel":                "string",
		"max_body_size":                  "int",
		"max_title_length":               "int",
		"default_priority":               "int",
		"max_channel_name_length":        "int",
		"max_channel_description_length": "int",
		"max_sender_name_length":         "int",
		"max_user_message_id_length":     "int",
	})
}

func TestResponseChannelPreview(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	response := tt.RequestAuthGetRaw(t, data.User[1].AdminKey, baseUrl, fmt.Sprintf("/api/v2/preview/channels/%s", data.User[0].Channels[0].ChannelID))

	tt.AssertJsonStructureMatch(t, "json[channel]", response, map[string]any{
		"channel_id":       "id",
		"owner_user_id":    "id",
		"internal_name":    "string",
		"display_name":     "string",
		"description_name": "string|null",
	})
}

func TestResponseUserPreview(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	response := tt.RequestAuthGetRaw(t, data.User[1].AdminKey, baseUrl, fmt.Sprintf("/api/v2/preview/users/%s", data.User[0].UID))

	tt.AssertJsonStructureMatch(t, "json[user]", response, map[string]any{
		"user_id":  "id",
		"username": "string|null",
	})
}

func TestResponseKeyTokenPreview(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	response := tt.RequestAuthGetRaw(t, data.User[1].AdminKey, baseUrl, fmt.Sprintf("/api/v2/preview/keys/%s", data.User[0].Keys[0].KeyID))

	tt.AssertJsonStructureMatch(t, "json[key]", response, map[string]any{
		"keytoken_id":   "id",
		"name":          "string",
		"owner_user_id": "id",
		"all_channels":  "bool",
		"channels":      []any{"id"},
		"permissions":   "string",
	})
}
