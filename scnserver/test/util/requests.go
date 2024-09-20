package util

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
)

func RequestRaw(t *testing.T, baseURL string, urlSuffix string) {
	RequestAny[Void](t, "", "GET", baseURL, urlSuffix, nil, false)
}

func RequestGet[TResult any](t *testing.T, baseURL string, urlSuffix string) TResult {
	return RequestAny[TResult](t, "", "GET", baseURL, urlSuffix, nil, true)
}

func RequestAuthGet[TResult any](t *testing.T, akey string, baseURL string, urlSuffix string) TResult {
	return RequestAny[TResult](t, akey, "GET", baseURL, urlSuffix, nil, true)
}

func RequestAuthGetRaw(t *testing.T, akey string, baseURL string, urlSuffix string) string {
	return RequestAny[string](t, akey, "GET", baseURL, urlSuffix, nil, false)
}

func RequestPost[TResult any](t *testing.T, baseURL string, urlSuffix string, body any) TResult {
	return RequestAny[TResult](t, "", "POST", baseURL, urlSuffix, body, true)
}

func RequestAuthPostRaw(t *testing.T, akey string, baseURL string, urlSuffix string, body any) string {
	return RequestAny[string](t, akey, "POST", baseURL, urlSuffix, body, false)
}

func RequestAuthPost[TResult any](t *testing.T, akey string, baseURL string, urlSuffix string, body any) TResult {
	return RequestAny[TResult](t, akey, "POST", baseURL, urlSuffix, body, true)
}

func RequestPut[TResult any](t *testing.T, baseURL string, urlSuffix string, body any) TResult {
	return RequestAny[TResult](t, "", "PUT", baseURL, urlSuffix, body, true)
}

func RequestAuthPUT[TResult any](t *testing.T, akey string, baseURL string, urlSuffix string, body any) TResult {
	return RequestAny[TResult](t, akey, "PUT", baseURL, urlSuffix, body, true)
}

func RequestPatch[TResult any](t *testing.T, baseURL string, urlSuffix string, body any) TResult {
	return RequestAny[TResult](t, "", "PATCH", baseURL, urlSuffix, body, true)
}

func RequestAuthPatch[TResult any](t *testing.T, akey string, baseURL string, urlSuffix string, body any) TResult {
	return RequestAny[TResult](t, akey, "PATCH", baseURL, urlSuffix, body, true)
}

func RequestDelete[TResult any](t *testing.T, baseURL string, urlSuffix string, body any) TResult {
	return RequestAny[TResult](t, "", "DELETE", baseURL, urlSuffix, body, true)
}

func RequestAuthDelete[TResult any](t *testing.T, akey string, baseURL string, urlSuffix string, body any) TResult {
	return RequestAny[TResult](t, akey, "DELETE", baseURL, urlSuffix, body, true)
}

func RequestGetShouldFail(t *testing.T, baseURL string, urlSuffix string, statusCode int, errcode apierr.APIError) {
	RequestAuthAnyShouldFail(t, "", "GET", baseURL, urlSuffix, nil, statusCode, errcode)
}

func RequestPostShouldFail(t *testing.T, baseURL string, urlSuffix string, body any, statusCode int, errcode apierr.APIError) {
	RequestAuthAnyShouldFail(t, "", "POST", baseURL, urlSuffix, body, statusCode, errcode)
}

func RequestPatchShouldFail(t *testing.T, baseURL string, urlSuffix string, body any, statusCode int, errcode apierr.APIError) {
	RequestAuthAnyShouldFail(t, "", "PATCH", baseURL, urlSuffix, body, statusCode, errcode)
}

func RequestDeleteShouldFail(t *testing.T, baseURL string, urlSuffix string, body any, statusCode int, errcode apierr.APIError) {
	RequestAuthAnyShouldFail(t, "", "DELETE", baseURL, urlSuffix, body, statusCode, errcode)
}

func RequestAuthGetShouldFail(t *testing.T, akey string, baseURL string, urlSuffix string, statusCode int, errcode apierr.APIError) {
	RequestAuthAnyShouldFail(t, akey, "GET", baseURL, urlSuffix, nil, statusCode, errcode)
}

func RequestAuthPostShouldFail(t *testing.T, akey string, baseURL string, urlSuffix string, body any, statusCode int, errcode apierr.APIError) {
	RequestAuthAnyShouldFail(t, akey, "POST", baseURL, urlSuffix, body, statusCode, errcode)
}

func RequestAuthPatchShouldFail(t *testing.T, akey string, baseURL string, urlSuffix string, body any, statusCode int, errcode apierr.APIError) {
	RequestAuthAnyShouldFail(t, akey, "PATCH", baseURL, urlSuffix, body, statusCode, errcode)
}

func RequestAuthDeleteShouldFail(t *testing.T, akey string, baseURL string, urlSuffix string, body any, statusCode int, errcode apierr.APIError) {
	RequestAuthAnyShouldFail(t, akey, "DELETE", baseURL, urlSuffix, body, statusCode, errcode)
}

func RequestAny[TResult any](t *testing.T, akey string, method string, baseURL string, urlSuffix string, body any, deserialize bool) TResult {
	client := http.Client{}

	TPrintf(zerolog.InfoLevel, "[-> REQUEST] (%s) %s%s [%s] [%s]\n", method, baseURL, urlSuffix, langext.Conditional(akey == "", "NO AUTH", "AUTH"), langext.Conditional(body == nil, "NO BODY", "BODY"))

	bytesbody := make([]byte, 0)
	contentType := ""
	if body != nil {
		switch bd := body.(type) {
		case FormData:
			bodybuffer := &bytes.Buffer{}
			writer := multipart.NewWriter(bodybuffer)
			for bdk, bdv := range bd {
				err := writer.WriteField(bdk, bdv)
				if err != nil {
					TestFailErr(t, err)
				}
			}
			err := writer.Close()
			if err != nil {
				TestFailErr(t, err)
			}
			bytesbody = bodybuffer.Bytes()
			contentType = writer.FormDataContentType()
		case RawJSON:
			bytesbody = []byte(body.(RawJSON).Body)
			contentType = "application/json"
		default:
			bjson, err := json.Marshal(body)
			if err != nil {
				TestFailErr(t, err)
			}
			bytesbody = bjson
			contentType = "application/json"
		}
	}

	req, err := http.NewRequest(method, baseURL+urlSuffix, bytes.NewReader(bytesbody))
	if err != nil {
		TestFailErr(t, err)
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	if akey != "" {
		req.Header.Set("Authorization", "SCN "+akey)
	}

	resp, err := client.Do(req)
	if err != nil {
		TestFailErr(t, err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBodyBin, err := io.ReadAll(resp.Body)
	if err != nil {
		TestFailErr(t, err)
	}

	TPrintln(zerolog.DebugLevel, "")
	TPrintf(zerolog.DebugLevel, "----------------  RESPONSE (%d) ----------------\n", resp.StatusCode)
	if len(respBodyBin) > 100_000 {
		TPrintln(zerolog.DebugLevel, "[[RESPONSE TOO LONG]]")
	} else {
		TPrintln(zerolog.DebugLevel, langext.TryPrettyPrintJson(string(respBodyBin)))
	}
	TryPrintTraceObj(zerolog.DebugLevel, "----------------  --------  ----------------", respBodyBin, "")
	TPrintln(zerolog.DebugLevel, "----------------  --------  ----------------")
	TPrintln(zerolog.DebugLevel, "")

	if resp.StatusCode != 200 {
		TestFailFmt(t, "Statuscode != 200 (actual = %d)", resp.StatusCode)
	}

	if deserialize {
		var data TResult
		if err := json.Unmarshal(respBodyBin, &data); err != nil {
			TestFailErr(t, err)
			return data
		}
		return data
	} else {
		if _, ok := (any(*new(TResult))).([]byte); ok {
			return any(respBodyBin).(TResult)
		} else if _, ok := (any(*new(TResult))).(string); ok {
			return any(string(respBodyBin)).(TResult)
		} else {
			return *new(TResult)
		}
	}
}

func RequestAuthAnyShouldFail(t *testing.T, akey string, method string, baseURL string, urlSuffix string, body any, expectedStatusCode int, errcode apierr.APIError) {
	client := http.Client{}

	TPrintf(zerolog.InfoLevel, "[-> REQUEST] (%s) %s%s [%s] (should-fail with %d/%d)\n", method, baseURL, urlSuffix, langext.Conditional(akey == "", "NO AUTH", "AUTH"), expectedStatusCode, errcode)

	bytesbody := make([]byte, 0)
	contentType := ""
	if body != nil {
		switch bd := body.(type) {
		case FormData:
			bodybuffer := &bytes.Buffer{}
			writer := multipart.NewWriter(bodybuffer)
			for bdk, bdv := range bd {
				err := writer.WriteField(bdk, bdv)
				if err != nil {
					TestFailErr(t, err)
				}
			}
			err := writer.Close()
			if err != nil {
				TestFailErr(t, err)
			}
			bytesbody = bodybuffer.Bytes()
			contentType = writer.FormDataContentType()
		default:
			bjson, err := json.Marshal(body)
			if err != nil {
				TestFailErr(t, err)
			}
			bytesbody = bjson
			contentType = "application/json"
		}
	}

	req, err := http.NewRequest(method, baseURL+urlSuffix, bytes.NewReader(bytesbody))
	if err != nil {
		TestFailErr(t, err)
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	if akey != "" {
		req.Header.Set("Authorization", "SCN "+akey)
	}

	resp, err := client.Do(req)
	if err != nil {
		TestFailErr(t, err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBodyBin, err := io.ReadAll(resp.Body)
	if err != nil {
		TestFailErr(t, err)
	}

	TPrintln(zerolog.DebugLevel, "")
	TPrintf(zerolog.DebugLevel, "----------------  RESPONSE (%d) ----------------\n", resp.StatusCode)
	TPrintln(zerolog.DebugLevel, langext.TryPrettyPrintJson(string(respBodyBin)))
	if (expectedStatusCode != 0 && resp.StatusCode != expectedStatusCode) || (expectedStatusCode == 0 && resp.StatusCode == 200) {
		TryPrintTraceObj(zerolog.DebugLevel, "----------------  --------  ----------------", respBodyBin, "")
	}
	TPrintln(zerolog.DebugLevel, "----------------  --------  ----------------")
	TPrintln(zerolog.DebugLevel, "")

	if expectedStatusCode != 0 && resp.StatusCode != expectedStatusCode {
		TestFailFmt(t, "Statuscode != %d (expected failure, but got %d)", expectedStatusCode, resp.StatusCode)
	}
	if expectedStatusCode == 0 && resp.StatusCode == 200 {
		TestFailFmt(t, "Statuscode == %d (expected any failure, but got %d)", resp.StatusCode, resp.StatusCode)
	}

	var data gin.H
	if err := json.Unmarshal(respBodyBin, &data); err != nil {
		TestFailErr(t, err)
	}

	if v, ok := data["success"]; ok {
		if v.(bool) {
			TestFail(t, "Success == true (expected failure)")
		}
	} else {
		TestFail(t, "missing response['success']")
	}

	if errcode != 0 {
		if v, ok := data["error"]; ok {
			if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", errcode) {
				TestFailFmt(t, "wrong errorcode (expected: %d), (actual: %v)", errcode, v)
			}
		} else {
			TestFail(t, "missing response['error']")
		}
	}
}

func TryPrintTraceObj(lvl zerolog.Level, prefix string, body []byte, suffix string) {
	v1 := gin.H{}
	if err := json.Unmarshal(body, &v1); err == nil {
		if v2, ok := v1["traceObj"]; ok {
			if v3, ok := v2.(string); ok {
				if prefix != "" {
					TPrintln(lvl, prefix)
				}

				TPrintln(lvl, strings.TrimSpace(v3))

				if suffix != "" {
					TPrintln(lvl, suffix)
				}
			}
		}
	}
}
