package logic

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/common/ginresp"
	"blackforestbytes.com/simplecloudnotifier/models"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
)

type PermKeyType string

const (
	PermKeyTypeNone      PermKeyType = "NONE"       // (nothing)
	PermKeyTypeUserSend  PermKeyType = "USER_SEND"  // send-messages
	PermKeyTypeUserRead  PermKeyType = "USER_READ"  // send-messages, list-messages, read-user
	PermKeyTypeUserAdmin PermKeyType = "USER_ADMIN" // send-messages, list-messages, read-user, delete-messages, update-user
)

type PermissionSet struct {
	UserID  *int64
	KeyType PermKeyType
}

func NewEmptyPermissions() PermissionSet {
	return PermissionSet{
		UserID:  nil,
		KeyType: PermKeyTypeNone,
	}
}

var respoNotAuthorized = ginresp.InternAPIError(401, apierr.USER_AUTH_FAILED, "You are not authorized for this action", nil)

func (ac *AppContext) CheckPermissionUserRead(userid int64) *ginresp.HTTPResponse {
	p := ac.permissions
	if p.UserID != nil && *p.UserID == userid && p.KeyType == PermKeyTypeUserRead {
		return nil
	}
	if p.UserID != nil && *p.UserID == userid && p.KeyType == PermKeyTypeUserAdmin {
		return nil
	}

	return langext.Ptr(respoNotAuthorized)
}

func (ac *AppContext) CheckPermissionUserAdmin(userid int64) *ginresp.HTTPResponse {
	p := ac.permissions
	if p.UserID != nil && *p.UserID == userid && p.KeyType == PermKeyTypeUserAdmin {
		return nil
	}

	return langext.Ptr(respoNotAuthorized)
}

func (ac *AppContext) CheckPermissionAny() *ginresp.HTTPResponse {
	p := ac.permissions
	if p.KeyType == PermKeyTypeNone {
		return langext.Ptr(respoNotAuthorized)
	}

	return nil
}

func (ac *AppContext) CheckPermissionMessageReadDirect(msg models.Message) bool {
	p := ac.permissions
	if p.UserID != nil && msg.OwnerUserID == *p.UserID && p.KeyType == PermKeyTypeUserRead {
		return true
	}
	if p.UserID != nil && msg.OwnerUserID == *p.UserID && p.KeyType == PermKeyTypeUserAdmin {
		return true
	}

	return false
}

func (ac *AppContext) GetPermissionUserID() *int64 {
	if ac.permissions.UserID == nil {
		return nil
	} else {
		return langext.Ptr(*ac.permissions.UserID)
	}
}

func (ac *AppContext) IsPermissionUserRead() bool {
	p := ac.permissions
	return p.KeyType == PermKeyTypeUserRead || p.KeyType == PermKeyTypeUserAdmin
}

func (ac *AppContext) IsPermissionUserSend() bool {
	p := ac.permissions
	return p.KeyType == PermKeyTypeUserSend || p.KeyType == PermKeyTypeUserAdmin
}

func (ac *AppContext) IsPermissionUserAdmin() bool {
	p := ac.permissions
	return p.KeyType == PermKeyTypeUserAdmin
}
