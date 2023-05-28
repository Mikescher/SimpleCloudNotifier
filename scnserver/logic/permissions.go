package logic

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"blackforestbytes.com/simplecloudnotifier/api/ginresp"
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
)

func (ac *AppContext) CheckPermissionUserRead(userid models.UserID) *ginresp.HTTPResponse {
	p := ac.permissions
	if p.Token != nil && p.Token.IsUserRead(userid) {
		return nil
	}

	return langext.Ptr(ginresp.APIError(ac.ginContext, 401, apierr.USER_AUTH_FAILED, "You are not authorized for this action", nil))
}

func (ac *AppContext) CheckPermissionSelfAllMessagesRead() *ginresp.HTTPResponse {
	p := ac.permissions
	if p.Token != nil && p.Token.IsAllMessagesRead(p.Token.OwnerUserID) {
		return nil
	}

	return langext.Ptr(ginresp.APIError(ac.ginContext, 401, apierr.USER_AUTH_FAILED, "You are not authorized for this action", nil))
}

func (ac *AppContext) CheckPermissionAllMessagesRead(userid models.UserID) *ginresp.HTTPResponse {
	p := ac.permissions
	if p.Token != nil && p.Token.IsAllMessagesRead(userid) {
		return nil
	}

	return langext.Ptr(ginresp.APIError(ac.ginContext, 401, apierr.USER_AUTH_FAILED, "You are not authorized for this action", nil))
}

func (ac *AppContext) CheckPermissionChanMessagesRead(channel models.Channel) *ginresp.HTTPResponse {
	p := ac.permissions
	if p.Token != nil && p.Token.IsChannelMessagesRead(channel.ChannelID) {

		if channel.OwnerUserID == p.Token.OwnerUserID {
			return nil // owned channel
		} else {
			sub, err := ac.app.Database.Primary.GetSubscriptionBySubscriber(ac, p.Token.OwnerUserID, channel.ChannelID)
			if err == sql.ErrNoRows {
				return langext.Ptr(ginresp.APIError(ac.ginContext, 401, apierr.USER_AUTH_FAILED, "You are not authorized for this action", nil))
			}
			if err != nil {
				return langext.Ptr(ginresp.APIError(ac.ginContext, 500, apierr.DATABASE_ERROR, "Failed to query subscription", err))
			}
			if sub == nil {
				return langext.Ptr(ginresp.APIError(ac.ginContext, 401, apierr.USER_AUTH_FAILED, "You are not authorized for this action (no subscription)", nil))
			}
			if !sub.Confirmed {
				return langext.Ptr(ginresp.APIError(ac.ginContext, 401, apierr.USER_AUTH_FAILED, "You are not authorized for this action (subscription not confirmed)", nil))
			}
			return nil // subscribed channel
		}
	}

	return langext.Ptr(ginresp.APIError(ac.ginContext, 401, apierr.USER_AUTH_FAILED, "You are not authorized for this action", nil))
}

func (ac *AppContext) CheckPermissionUserAdmin(userid models.UserID) *ginresp.HTTPResponse {
	p := ac.permissions
	if p.Token != nil && p.Token.IsAdmin(userid) {
		return nil
	}

	return langext.Ptr(ginresp.APIError(ac.ginContext, 401, apierr.USER_AUTH_FAILED, "You are not authorized for this action", nil))
}

func (ac *AppContext) CheckPermissionSend(channel models.Channel, key string) (*models.KeyToken, *ginresp.HTTPResponse) {

	keytok, err := ac.app.Database.Primary.GetKeyTokenByToken(ac, key)
	if err != nil {
		return nil, langext.Ptr(ginresp.APIError(ac.ginContext, 500, apierr.DATABASE_ERROR, "Failed to query token", err))
	}
	if keytok == nil {
		return nil, langext.Ptr(ginresp.APIError(ac.ginContext, 401, apierr.USER_AUTH_FAILED, "You are not authorized for this action", nil))
	}

	if keytok.IsChannelMessagesSend(channel) {
		return keytok, nil
	}

	return nil, langext.Ptr(ginresp.APIError(ac.ginContext, 401, apierr.USER_AUTH_FAILED, "You are not authorized for this action", nil))
}

func (ac *AppContext) CheckPermissionMessageRead(msg models.Message) bool {
	p := ac.permissions
	if p.Token != nil && p.Token.IsChannelMessagesRead(msg.ChannelID) {
		return true
	}

	return false
}

func (ac *AppContext) CheckPermissionAny() *ginresp.HTTPResponse {
	p := ac.permissions
	if p.Token == nil {
		return langext.Ptr(ginresp.APIError(ac.ginContext, 401, apierr.USER_AUTH_FAILED, "You are not authorized for this action", nil))
	}

	return nil
}

func (ac *AppContext) GetPermissionUserID() *models.UserID {
	if ac.permissions.Token == nil {
		return nil
	} else {
		return langext.Ptr(ac.permissions.Token.OwnerUserID)
	}
}

func (ac *AppContext) GetPermissionKeyTokenID() *models.KeyTokenID {
	if ac.permissions.Token == nil {
		return nil
	} else {
		return langext.Ptr(ac.permissions.Token.KeyTokenID)
	}
}
