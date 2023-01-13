package models

type PermKeyType string

const (
	PermKeyTypeNone      PermKeyType = "NONE"       // (nothing)
	PermKeyTypeUserSend  PermKeyType = "USER_SEND"  // send-messages
	PermKeyTypeUserRead  PermKeyType = "USER_READ"  // send-messages, list-messages, read-user
	PermKeyTypeUserAdmin PermKeyType = "USER_ADMIN" // send-messages, list-messages, read-user, delete-messages, update-user
)

type PermissionSet struct {
	UserID  *UserID
	KeyType PermKeyType
}

func NewEmptyPermissions() PermissionSet {
	return PermissionSet{
		UserID:  nil,
		KeyType: PermKeyTypeNone,
	}
}
