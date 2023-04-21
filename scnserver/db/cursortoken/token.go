package cursortoken

import (
	"encoding/base32"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type Mode string //@enum:type

const (
	CTMStart  = "START"
	CTMNormal = "NORMAL"
	CTMEnd    = "END"
)

type CursorToken struct {
	Mode       Mode
	Timestamp  int64
	Id         string
	Direction  string
	FilterHash string
}

type cursorTokenSerialize struct {
	Timestamp  *int64  `json:"ts,omitempty"`
	Id         *string `json:"id,omitempty"`
	Direction  *string `json:"dir,omitempty"`
	FilterHash *string `json:"f,omitempty"`
}

func Start() CursorToken {
	return CursorToken{
		Mode:       CTMStart,
		Timestamp:  0,
		Id:         "",
		Direction:  "",
		FilterHash: "",
	}
}

func End() CursorToken {
	return CursorToken{
		Mode:       CTMEnd,
		Timestamp:  0,
		Id:         "",
		Direction:  "",
		FilterHash: "",
	}
}

func Normal(ts time.Time, id string, dir string, filter string) CursorToken {
	return CursorToken{
		Mode:       CTMNormal,
		Timestamp:  ts.UnixMilli(),
		Id:         id,
		Direction:  dir,
		FilterHash: filter,
	}
}

func (c *CursorToken) Token() string {
	if c.Mode == CTMStart {
		return "@start"
	}
	if c.Mode == CTMEnd {
		return "@end"
	}

	// We kinda manually implement omitempty for the CursorToken here
	// because omitempty does not work for time.Time and otherwise we would always
	// get weird time values when decoding a token that initially didn't have an Timestamp set
	// For this usecase we treat Unix=0 as an empty timestamp

	sertok := cursorTokenSerialize{}

	if c.Id != "" {
		sertok.Id = &c.Id
	}

	if c.Timestamp != 0 {
		sertok.Timestamp = &c.Timestamp
	}

	if c.Direction != "" {
		sertok.Direction = &c.Direction
	}

	if c.FilterHash != "" {
		sertok.FilterHash = &c.FilterHash
	}

	body, err := json.Marshal(sertok)
	if err != nil {
		panic(err)
	}

	return "tok_" + base32.StdEncoding.EncodeToString(body)
}

func Decode(tok string) (CursorToken, error) {
	if tok == "" {
		return Start(), nil
	}
	if strings.ToLower(tok) == "@start" {
		return Start(), nil
	}
	if strings.ToLower(tok) == "@end" {
		return End(), nil
	}

	if !strings.HasPrefix(tok, "tok_") {
		return CursorToken{}, errors.New("could not decode token, missing prefix")
	}

	body, err := base32.StdEncoding.DecodeString(tok[len("tok_"):])
	if err != nil {
		return CursorToken{}, err
	}

	var tokenDeserialize cursorTokenSerialize
	err = json.Unmarshal(body, &tokenDeserialize)
	if err != nil {
		return CursorToken{}, err
	}

	token := CursorToken{Mode: CTMNormal}

	if tokenDeserialize.Timestamp != nil {
		token.Timestamp = *tokenDeserialize.Timestamp
	}
	if tokenDeserialize.Id != nil {
		token.Id = *tokenDeserialize.Id
	}
	if tokenDeserialize.Direction != nil {
		token.Direction = *tokenDeserialize.Direction
	}
	if tokenDeserialize.FilterHash != nil {
		token.FilterHash = *tokenDeserialize.FilterHash
	}

	return token, nil
}
