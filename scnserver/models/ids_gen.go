// Code generated by csid-generate.go DO NOT EDIT.

package models

import "crypto/rand"
import "crypto/sha256"
import "fmt"
import "github.com/go-playground/validator/v10"
import "github.com/rs/zerolog/log"
import "gogs.mikescher.com/BlackForestBytes/goext/exerr"
import "gogs.mikescher.com/BlackForestBytes/goext/langext"
import "gogs.mikescher.com/BlackForestBytes/goext/rext"
import "math/big"
import "reflect"
import "regexp"
import "strings"

const ChecksumCharsetIDGenerator = "ba14f2f5d0b0357f248dcbd12933de102c80f1e61be697a37ebb723609fc0c59" // GoExtVersion: 0.0.485

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
	prefixKeyTokenID     = "TOK"
)

var (
	regexUserID         = generateRegex(prefixUserID)
	regexChannelID      = generateRegex(prefixChannelID)
	regexDeliveryID     = generateRegex(prefixDeliveryID)
	regexMessageID      = generateRegex(prefixMessageID)
	regexSubscriptionID = generateRegex(prefixSubscriptionID)
	regexClientID       = generateRegex(prefixClientID)
	regexRequestID      = generateRegex(prefixRequestID)
	regexKeyTokenID     = generateRegex(prefixKeyTokenID)
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
	csMax := big.NewInt(int64(idCharsetLen))
	checksum := 0
	for i := 0; i < idlen-len(prefix)-checklen; i++ {
		v, err := rand.Int(rand.Reader, csMax)
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

func generateIDFromSeed(prefix string, seed string) string {
	h := sha256.New()

	iddata := ""
	for len(iddata) < idlen-len(prefix)-checklen {
		h.Write([]byte(seed))
		bs := h.Sum(nil)
		iddata += langext.NewAnyBaseConverter(idCharset).Encode(bs)
	}

	checksum := 0
	for i := 0; i < idlen-len(prefix)-checklen; i++ {
		ichr := int(iddata[i])
		checksum = (checksum + charSetReverseMap[ichr]) % (idCharsetLen)
	}

	checkstr := string(idCharset[checksum%idCharsetLen])

	return prefix + iddata[:(idlen-len(prefix)-checklen)] + checkstr
}

func validateID(prefix string, value string) error {
	if len(value) != idlen {
		return exerr.New(exerr.TypeInvalidCSID, "id has the wrong length").Str("value", value).Build()
	}

	if !strings.HasPrefix(value, prefix) {
		return exerr.New(exerr.TypeInvalidCSID, "id is missing the correct prefix").Str("value", value).Str("prefix", prefix).Build()
	}

	checksum := 0
	for i := len(prefix); i < len(value)-checklen; i++ {
		ichr := int(value[i])
		if ichr < 0 || ichr >= len(charSetReverseMap) || charSetReverseMap[ichr] == -1 {
			return exerr.New(exerr.TypeInvalidCSID, "id contains invalid characters").Str("value", value).Build()
		}
		checksum = (checksum + charSetReverseMap[ichr]) % (idCharsetLen)
	}

	checkstr := string(idCharset[checksum%idCharsetLen])

	if !strings.HasSuffix(value, checkstr) {
		return exerr.New(exerr.TypeInvalidCSID, "id checkstring is invalid").Str("value", value).Str("checkstr", checkstr).Build()
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

// ================================ UserID (ids.go) ================================

func NewUserID() UserID {
	return UserID(generateID(prefixUserID))
}

func (id UserID) Valid() error {
	return validateID(prefixUserID, string(id))
}

func (i UserID) String() string {
	return string(i)
}

func (i UserID) Prefix() string {
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

// ================================ ChannelID (ids.go) ================================

func NewChannelID() ChannelID {
	return ChannelID(generateID(prefixChannelID))
}

func (id ChannelID) Valid() error {
	return validateID(prefixChannelID, string(id))
}

func (i ChannelID) String() string {
	return string(i)
}

func (i ChannelID) Prefix() string {
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

// ================================ DeliveryID (ids.go) ================================

func NewDeliveryID() DeliveryID {
	return DeliveryID(generateID(prefixDeliveryID))
}

func (id DeliveryID) Valid() error {
	return validateID(prefixDeliveryID, string(id))
}

func (i DeliveryID) String() string {
	return string(i)
}

func (i DeliveryID) Prefix() string {
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

// ================================ MessageID (ids.go) ================================

func NewMessageID() MessageID {
	return MessageID(generateID(prefixMessageID))
}

func (id MessageID) Valid() error {
	return validateID(prefixMessageID, string(id))
}

func (i MessageID) String() string {
	return string(i)
}

func (i MessageID) Prefix() string {
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

// ================================ SubscriptionID (ids.go) ================================

func NewSubscriptionID() SubscriptionID {
	return SubscriptionID(generateID(prefixSubscriptionID))
}

func (id SubscriptionID) Valid() error {
	return validateID(prefixSubscriptionID, string(id))
}

func (i SubscriptionID) String() string {
	return string(i)
}

func (i SubscriptionID) Prefix() string {
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

// ================================ ClientID (ids.go) ================================

func NewClientID() ClientID {
	return ClientID(generateID(prefixClientID))
}

func (id ClientID) Valid() error {
	return validateID(prefixClientID, string(id))
}

func (i ClientID) String() string {
	return string(i)
}

func (i ClientID) Prefix() string {
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

// ================================ RequestID (ids.go) ================================

func NewRequestID() RequestID {
	return RequestID(generateID(prefixRequestID))
}

func (id RequestID) Valid() error {
	return validateID(prefixRequestID, string(id))
}

func (i RequestID) String() string {
	return string(i)
}

func (i RequestID) Prefix() string {
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

// ================================ KeyTokenID (ids.go) ================================

func NewKeyTokenID() KeyTokenID {
	return KeyTokenID(generateID(prefixKeyTokenID))
}

func (id KeyTokenID) Valid() error {
	return validateID(prefixKeyTokenID, string(id))
}

func (i KeyTokenID) String() string {
	return string(i)
}

func (i KeyTokenID) Prefix() string {
	return prefixKeyTokenID
}

func (id KeyTokenID) Raw() string {
	return getRawData(prefixKeyTokenID, string(id))
}

func (id KeyTokenID) CheckString() string {
	return getCheckString(prefixKeyTokenID, string(id))
}

func (id KeyTokenID) Regex() rext.Regex {
	return regexKeyTokenID
}
