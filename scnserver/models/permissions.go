package models

type PermissionSet struct {
	Token *KeyToken // KeyToken.Permissions
}

func NewEmptyPermissions() PermissionSet {
	return PermissionSet{
		Token: nil,
	}
}
