package logic

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/common/ginresp"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
)

type PermKeyType string

const (
	PermKeyTypeNone        PermKeyType = "NONE"           // (nothing)
	PermKeyTypeUserSend    PermKeyType = "USER_SEND"      // send-messages
	PermKeyTypeUserRead    PermKeyType = "USER_READ"      // send-messages, list-messages, read-user
	PermKeyTypeUserAdmin   PermKeyType = "USER_ADMIN"     // send-messages, list-messages, read-user, delete-messages, update-user
	PermKeyTypeChannelSub  PermKeyType = "CHAN_SUBSCRIBE" // subscribe-channel
	PermKeyTypeChannelSend PermKeyType = "CHAN_SEND"      // send-messages
)

type PermissionSet struct {
	ReferenceID *int64
	KeyType     PermKeyType
}

func NewEmptyPermissions() PermissionSet {
	return PermissionSet{
		ReferenceID: nil,
		KeyType:     PermKeyTypeNone,
	}
}

var respoNotAuthorized = ginresp.InternAPIError(apierr.USER_AUTH_FAILED, "You are not authorized for this action", nil)

func (ac *AppContext) CheckPermissionUserRead(userid int64) *ginresp.HTTPResponse {
	p := ac.permissions
	if p.ReferenceID != nil && *p.ReferenceID == userid && p.KeyType == PermKeyTypeUserRead {
		return nil
	}
	if p.ReferenceID != nil && *p.ReferenceID == userid && p.KeyType == PermKeyTypeUserAdmin {
		return nil
	}

	return langext.Ptr(respoNotAuthorized)
}
