package test

import (
	tt "blackforestbytes.com/simplecloudnotifier/test/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"testing"
)

func TestCreateCompatUser(t *testing.T) {
	_, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	r0 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/register.php?fcm_token=%s&pro=%s&pro_token=%s", "DUMMY_FCM", "0", ""))

	tt.AssertEqual(t, "success", true, r0["success"])

	userid := int64(r0["user_id"].(float64))
	userkey := r0["user_key"].(string)

	r1 := tt.RequestGet[gin.H](t, baseUrl, fmt.Sprintf("/api/info.php?user_id=%d&user_key=%s", userid, userkey))

	tt.AssertEqual(t, "success", true, r1["success"])
}

//TODO test compat methods

//TODO also test compat_id mapping
