package test

import (
	tt "blackforestbytes.com/simplecloudnotifier/test/util"
	"fmt"
	"testing"
)

func TestResponseChannel(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	response := tt.RequestAuthGetRaw(t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s", data.User[0].UID, data.User[0].Channels[0]))

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

	response := tt.RequestAuthGetRaw(t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/clients/%s", data.User[2].UID, data.User[2].Clients[2]))

	tt.AssertJsonStructureMatch(t, "json[client]", response, map[string]any{})
}

func TestResponseKeyToken(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	response := tt.RequestAuthGetRaw(t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys/%s", data.User[0].UID, data.User[0].Keys[0]))

	tt.AssertJsonStructureMatch(t, "json[key]", response, map[string]any{})
}

func TestResponseMessage(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	response := tt.RequestAuthGetRaw(t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/messages/%s", data.User[0].UID, data.User[0].Messages[0]))

	tt.AssertJsonStructureMatch(t, "json[message]", response, map[string]any{})
}

func TestResponseSubscription(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	response := tt.RequestAuthGetRaw(t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data.User[0].UID, data.User[0].Subscriptions[0]))

	tt.AssertJsonStructureMatch(t, "json[subscription]", response, map[string]any{})
}

func TestResponseUser(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	response := tt.RequestAuthGetRaw(t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s", data.User[0].UID))

	tt.AssertJsonStructureMatch(t, "json[user]", response, map[string]any{})
}

func TestResponseChannelPreview(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	response := tt.RequestAuthGetRaw(t, data.User[1].AdminKey, baseUrl, fmt.Sprintf("/api/v2/preview/channels/%s", data.User[0].Channels[0]))

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

func TestResponseUserPreview(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	response := tt.RequestAuthGetRaw(t, data.User[1].AdminKey, baseUrl, fmt.Sprintf("/api/v2/preview/users/%s", data.User[0].UID))

	tt.AssertJsonStructureMatch(t, "json[user]", response, map[string]any{})
}

func TestResponseKeyTokenPreview(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	response := tt.RequestAuthGetRaw(t, data.User[1].AdminKey, baseUrl, fmt.Sprintf("/api/v2/preview/keys/%s", data.User[0].Keys[0]))

	tt.AssertJsonStructureMatch(t, "json[key]", response, map[string]any{})
}
