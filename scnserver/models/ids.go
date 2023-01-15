package models

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/rext"
	"math/big"
	"reflect"
	"regexp"
	"strings"
)

type EntityID interface {
	String() string
	Valid() error
	Prefix() string
	Raw() string
	CheckString() string
	Regex() rext.Regex
}

const idlen = 24

const checklen = 1

const idCharset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const idCharsetLen = len(idCharset)

var charSetReverseMap = generateCharsetMap()

const (
	prefixUserID         = "USR"
	prefixChannelID      = "CHA"
	prefixDeliveryID     = "DEL"
	prefixMessageID      = "MSG"
	prefixSubscriptionID = "SUB"
	prefixClientID       = "CLN"
	prefixRequestID      = "REQ"
)

var (
	regexUserID         = generateRegex(prefixUserID)
	regexChannelID      = generateRegex(prefixChannelID)
	regexDeliveryID     = generateRegex(prefixDeliveryID)
	regexMessageID      = generateRegex(prefixMessageID)
	regexSubscriptionID = generateRegex(prefixSubscriptionID)
	regexClientID       = generateRegex(prefixClientID)
	regexRequestID      = generateRegex(prefixRequestID)
)

func generateRegex(prefix string) rext.Regex {
	return rext.W(regexp.MustCompile(fmt.Sprintf("^%s[%s]{%d}[%s]{%d}$", prefix, idCharset, idlen-len(prefix)-checklen, idCharset, checklen)))
}

func generateCharsetMap() []int {
	result := make([]int, 128)
	for i := 0; i < len(result); i++ {
		result[i] = -1
	}
	for idx, chr := range idCharset {
		result[int(chr)] = idx
	}
	return result
}

func generateID(prefix string) string {
	k := ""
	max := big.NewInt(int64(idCharsetLen))
	checksum := 0
	for i := 0; i < idlen-len(prefix)-checklen; i++ {
		v, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(err)
		}
		v64 := v.Int64()
		k += string(idCharset[v64])
		checksum = (checksum + int(v64)) % (idCharsetLen)
	}
	checkstr := string(idCharset[checksum%idCharsetLen])
	return prefix + k + checkstr
}

func validateID(prefix string, value string) error {
	if len(value) != idlen {
		return errors.New("id has the wrong length")
	}

	if !strings.HasPrefix(value, prefix) {
		return errors.New("id is missing the correct prefix")
	}

	checksum := 0
	for i := len(prefix); i < len(value)-checklen; i++ {
		ichr := int(value[i])
		if ichr < 0 || ichr >= len(charSetReverseMap) || charSetReverseMap[ichr] == -1 {
			return errors.New("id contains invalid characters")
		}
		checksum = (checksum + charSetReverseMap[ichr]) % (idCharsetLen)
	}

	checkstr := string(idCharset[checksum%idCharsetLen])

	if !strings.HasSuffix(value, checkstr) {
		return errors.New("id checkstring is invalid")
	}

	return nil
}

func getRawData(prefix string, value string) string {
	if len(value) != idlen {
		return ""
	}
	return value[len(prefix) : idlen-checklen]
}

func getCheckString(prefix string, value string) string {
	if len(value) != idlen {
		return ""
	}
	return value[idlen-checklen:]
}

func ValidateEntityID(vfl validator.FieldLevel) bool {
	if !vfl.Field().CanInterface() {
		log.Error().Msgf("Failed to validate EntityID (cannot interface ?!?)")
		return false
	}

	ifvalue := vfl.Field().Interface()

	if value1, ok := ifvalue.(EntityID); ok {

		if vfl.Field().Type().Kind() == reflect.Pointer && langext.IsNil(value1) {
			return true
		}

		if err := value1.Valid(); err != nil {
			log.Debug().Msgf("Failed to validate EntityID '%s' (%s)", value1.String(), err.Error())
			return false
		} else {
			return true
		}

	} else {
		log.Error().Msgf("Failed to validate EntityID (wrong type: %T)", ifvalue)
		return false
	}
}

// ------------------------------------------------------------

type UserID string

func NewUserID() UserID {
	return UserID(generateID(prefixUserID))
}

func (id UserID) Valid() error {
	return validateID(prefixUserID, string(id))
}

func (id UserID) String() string {
	return string(id)
}

func (id UserID) Prefix() string {
	return prefixUserID
}

func (id UserID) Raw() string {
	return getRawData(prefixUserID, string(id))
}

func (id UserID) CheckString() string {
	return getCheckString(prefixUserID, string(id))
}

func (id UserID) Regex() rext.Regex {
	return regexUserID
}

// ------------------------------------------------------------

type ChannelID string

func NewChannelID() ChannelID {
	return ChannelID(generateID(prefixChannelID))
}

func (id ChannelID) Valid() error {
	return validateID(prefixChannelID, string(id))
}

func (id ChannelID) String() string {
	return string(id)
}

func (id ChannelID) Prefix() string {
	return prefixChannelID
}

func (id ChannelID) Raw() string {
	return getRawData(prefixChannelID, string(id))
}

func (id ChannelID) CheckString() string {
	return getCheckString(prefixChannelID, string(id))
}

func (id ChannelID) Regex() rext.Regex {
	return regexChannelID
}

// ------------------------------------------------------------

type DeliveryID string

func NewDeliveryID() DeliveryID {
	return DeliveryID(generateID(prefixDeliveryID))
}

func (id DeliveryID) Valid() error {
	return validateID(prefixDeliveryID, string(id))
}

func (id DeliveryID) String() string {
	return string(id)
}

func (id DeliveryID) Prefix() string {
	return prefixDeliveryID
}

func (id DeliveryID) Raw() string {
	return getRawData(prefixDeliveryID, string(id))
}

func (id DeliveryID) CheckString() string {
	return getCheckString(prefixDeliveryID, string(id))
}

func (id DeliveryID) Regex() rext.Regex {
	return regexDeliveryID
}

// ------------------------------------------------------------

type MessageID string

func NewMessageID() MessageID {
	return MessageID(generateID(prefixMessageID))
}

func (id MessageID) Valid() error {
	return validateID(prefixMessageID, string(id))
}

func (id MessageID) String() string {
	return string(id)
}

func (id MessageID) Prefix() string {
	return prefixMessageID
}

func (id MessageID) Raw() string {
	return getRawData(prefixMessageID, string(id))
}

func (id MessageID) CheckString() string {
	return getCheckString(prefixMessageID, string(id))
}

func (id MessageID) Regex() rext.Regex {
	return regexMessageID
}

// ------------------------------------------------------------

type SubscriptionID string

func NewSubscriptionID() SubscriptionID {
	return SubscriptionID(generateID(prefixSubscriptionID))
}

func (id SubscriptionID) Valid() error {
	return validateID(prefixSubscriptionID, string(id))
}

func (id SubscriptionID) String() string {
	return string(id)
}

func (id SubscriptionID) Prefix() string {
	return prefixSubscriptionID
}

func (id SubscriptionID) Raw() string {
	return getRawData(prefixSubscriptionID, string(id))
}

func (id SubscriptionID) CheckString() string {
	return getCheckString(prefixSubscriptionID, string(id))
}

func (id SubscriptionID) Regex() rext.Regex {
	return regexSubscriptionID
}

// ------------------------------------------------------------

type ClientID string

func NewClientID() ClientID {
	return ClientID(generateID(prefixClientID))
}

func (id ClientID) Valid() error {
	return validateID(prefixClientID, string(id))
}

func (id ClientID) String() string {
	return string(id)
}

func (id ClientID) Prefix() string {
	return prefixClientID
}

func (id ClientID) Raw() string {
	return getRawData(prefixClientID, string(id))
}

func (id ClientID) CheckString() string {
	return getCheckString(prefixClientID, string(id))
}

func (id ClientID) Regex() rext.Regex {
	return regexClientID
}

// ------------------------------------------------------------

type RequestID string

func NewRequestID() RequestID {
	return RequestID(generateID(prefixRequestID))
}

func (id RequestID) Valid() error {
	return validateID(prefixRequestID, string(id))
}

func (id RequestID) String() string {
	return string(id)
}

func (id RequestID) Prefix() string {
	return prefixRequestID
}

func (id RequestID) Raw() string {
	return getRawData(prefixRequestID, string(id))
}

func (id RequestID) CheckString() string {
	return getCheckString(prefixRequestID, string(id))
}

func (id RequestID) Regex() rext.Regex {
	return regexRequestID
}
