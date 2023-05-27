package test

import (
	tt "blackforestbytes.com/simplecloudnotifier/test/util"
	"fmt"
	"testing"
)

func TestListUserKeys(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitSingleData(t, ws)

	type keylist struct {
		Tokens []struct {
			AllChannels  bool     `json:"all_channels"`
			Channels     []string `json:"channels"`
			KeytokenId   string   `json:"keytoken_id"`
			MessagesSent int      `json:"messages_sent"`
			Name         string   `json:"name"`
			OwnerUserId  string   `json:"owner_user_id"`
			Permissions  string   `json:"permissions"`
		} `json:"tokens"`
	}

	klist := tt.RequestAuthGet[keylist](t, data.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/keys", data.UserID))

	tt.AssertEqual(t, "len(keys)", 1, len(klist.Tokens))

	t.SkipNow() //TODO
}

func TestCreateUserKey(t *testing.T) {
	t.SkipNow() //TODO
}

func TestDeleteUserKey(t *testing.T) {
	t.SkipNow() //TODO
}

func TestGetUserKey(t *testing.T) {
	t.SkipNow() //TODO
}

func TestUpdateUserKey(t *testing.T) {
	t.SkipNow() //TODO
}

func TestUserKeyPermissions(t *testing.T) {
	t.SkipNow() //TODO
}

func TestUsedKeyInMessage(t *testing.T) {
	t.SkipNow() //TODO
}
