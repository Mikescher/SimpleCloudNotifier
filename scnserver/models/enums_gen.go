// Code generated by enum-generate.go DO NOT EDIT.

package models

import "gogs.mikescher.com/BlackForestBytes/goext/langext"
import "gogs.mikescher.com/BlackForestBytes/goext/enums"

const ChecksumEnumGenerator = "5b115c5f107801af608630d2c5adce57cd4b050d176c8cd3db5c132020bf153c" // GoExtVersion: 0.0.463

// ================================ ClientType ================================
//
// File:       client.go
// StringEnum: true
// DescrEnum:  false
// DataEnum:   false
//

var __ClientTypeValues = []ClientType{
	ClientTypeAndroid,
	ClientTypeIOS,
	ClientTypeLinux,
	ClientTypeMacOS,
	ClientTypeWindows,
}

var __ClientTypeVarnames = map[ClientType]string{
	ClientTypeAndroid: "ClientTypeAndroid",
	ClientTypeIOS:     "ClientTypeIOS",
	ClientTypeLinux:   "ClientTypeLinux",
	ClientTypeMacOS:   "ClientTypeMacOS",
	ClientTypeWindows: "ClientTypeWindows",
}

func (e ClientType) Valid() bool {
	return langext.InArray(e, __ClientTypeValues)
}

func (e ClientType) Values() []ClientType {
	return __ClientTypeValues
}

func (e ClientType) ValuesAny() []any {
	return langext.ArrCastToAny(__ClientTypeValues)
}

func (e ClientType) ValuesMeta() []enums.EnumMetaValue {
	return ClientTypeValuesMeta()
}

func (e ClientType) String() string {
	return string(e)
}

func (e ClientType) VarName() string {
	if d, ok := __ClientTypeVarnames[e]; ok {
		return d
	}
	return ""
}

func (e ClientType) TypeName() string {
	return "ClientType"
}

func (e ClientType) PackageName() string {
	return "models"
}

func (e ClientType) Meta() enums.EnumMetaValue {
	return enums.EnumMetaValue{VarName: e.VarName(), Value: e, Description: nil}
}

func ParseClientType(vv string) (ClientType, bool) {
	for _, ev := range __ClientTypeValues {
		if string(ev) == vv {
			return ev, true
		}
	}
	return "", false
}

func ClientTypeValues() []ClientType {
	return __ClientTypeValues
}

func ClientTypeValuesMeta() []enums.EnumMetaValue {
	return []enums.EnumMetaValue{
		ClientTypeAndroid.Meta(),
		ClientTypeIOS.Meta(),
		ClientTypeLinux.Meta(),
		ClientTypeMacOS.Meta(),
		ClientTypeWindows.Meta(),
	}
}

// ================================ DeliveryStatus ================================
//
// File:       delivery.go
// StringEnum: true
// DescrEnum:  false
// DataEnum:   false
//

var __DeliveryStatusValues = []DeliveryStatus{
	DeliveryStatusRetry,
	DeliveryStatusSuccess,
	DeliveryStatusFailed,
}

var __DeliveryStatusVarnames = map[DeliveryStatus]string{
	DeliveryStatusRetry:   "DeliveryStatusRetry",
	DeliveryStatusSuccess: "DeliveryStatusSuccess",
	DeliveryStatusFailed:  "DeliveryStatusFailed",
}

func (e DeliveryStatus) Valid() bool {
	return langext.InArray(e, __DeliveryStatusValues)
}

func (e DeliveryStatus) Values() []DeliveryStatus {
	return __DeliveryStatusValues
}

func (e DeliveryStatus) ValuesAny() []any {
	return langext.ArrCastToAny(__DeliveryStatusValues)
}

func (e DeliveryStatus) ValuesMeta() []enums.EnumMetaValue {
	return DeliveryStatusValuesMeta()
}

func (e DeliveryStatus) String() string {
	return string(e)
}

func (e DeliveryStatus) VarName() string {
	if d, ok := __DeliveryStatusVarnames[e]; ok {
		return d
	}
	return ""
}

func (e DeliveryStatus) TypeName() string {
	return "DeliveryStatus"
}

func (e DeliveryStatus) PackageName() string {
	return "models"
}

func (e DeliveryStatus) Meta() enums.EnumMetaValue {
	return enums.EnumMetaValue{VarName: e.VarName(), Value: e, Description: nil}
}

func ParseDeliveryStatus(vv string) (DeliveryStatus, bool) {
	for _, ev := range __DeliveryStatusValues {
		if string(ev) == vv {
			return ev, true
		}
	}
	return "", false
}

func DeliveryStatusValues() []DeliveryStatus {
	return __DeliveryStatusValues
}

func DeliveryStatusValuesMeta() []enums.EnumMetaValue {
	return []enums.EnumMetaValue{
		DeliveryStatusRetry.Meta(),
		DeliveryStatusSuccess.Meta(),
		DeliveryStatusFailed.Meta(),
	}
}

// ================================ TokenPerm ================================
//
// File:       keytoken.go
// StringEnum: true
// DescrEnum:  true
// DataEnum:   false
//

var __TokenPermValues = []TokenPerm{
	PermAdmin,
	PermChannelRead,
	PermChannelSend,
	PermUserRead,
}

var __TokenPermDescriptions = map[TokenPerm]string{
	PermAdmin:       "Edit userdata (+ includes all other permissions)",
	PermChannelRead: "Read messages",
	PermChannelSend: "Send messages",
	PermUserRead:    "Read userdata",
}

var __TokenPermVarnames = map[TokenPerm]string{
	PermAdmin:       "PermAdmin",
	PermChannelRead: "PermChannelRead",
	PermChannelSend: "PermChannelSend",
	PermUserRead:    "PermUserRead",
}

func (e TokenPerm) Valid() bool {
	return langext.InArray(e, __TokenPermValues)
}

func (e TokenPerm) Values() []TokenPerm {
	return __TokenPermValues
}

func (e TokenPerm) ValuesAny() []any {
	return langext.ArrCastToAny(__TokenPermValues)
}

func (e TokenPerm) ValuesMeta() []enums.EnumMetaValue {
	return TokenPermValuesMeta()
}

func (e TokenPerm) String() string {
	return string(e)
}

func (e TokenPerm) Description() string {
	if d, ok := __TokenPermDescriptions[e]; ok {
		return d
	}
	return ""
}

func (e TokenPerm) VarName() string {
	if d, ok := __TokenPermVarnames[e]; ok {
		return d
	}
	return ""
}

func (e TokenPerm) TypeName() string {
	return "TokenPerm"
}

func (e TokenPerm) PackageName() string {
	return "models"
}

func (e TokenPerm) Meta() enums.EnumMetaValue {
	return enums.EnumMetaValue{VarName: e.VarName(), Value: e, Description: langext.Ptr(e.Description())}
}

func (e TokenPerm) DescriptionMeta() enums.EnumDescriptionMetaValue {
	return enums.EnumDescriptionMetaValue{VarName: e.VarName(), Value: e, Description: e.Description()}
}

func ParseTokenPerm(vv string) (TokenPerm, bool) {
	for _, ev := range __TokenPermValues {
		if string(ev) == vv {
			return ev, true
		}
	}
	return "", false
}

func TokenPermValues() []TokenPerm {
	return __TokenPermValues
}

func TokenPermValuesMeta() []enums.EnumMetaValue {
	return []enums.EnumMetaValue{
		PermAdmin.Meta(),
		PermChannelRead.Meta(),
		PermChannelSend.Meta(),
		PermUserRead.Meta(),
	}
}

func TokenPermValuesDescriptionMeta() []enums.EnumDescriptionMetaValue {
	return []enums.EnumDescriptionMetaValue{
		PermAdmin.DescriptionMeta(),
		PermChannelRead.DescriptionMeta(),
		PermChannelSend.DescriptionMeta(),
		PermUserRead.DescriptionMeta(),
	}
}

// ================================ ================= ================================

func AllPackageEnums() []enums.Enum {
	return []enums.Enum{
		ClientTypeAndroid,   // ClientType
		DeliveryStatusRetry, // DeliveryStatus
		PermAdmin,           // TokenPerm
	}
}
