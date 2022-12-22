package util

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
)

func RequestGet[TResult any](t *testing.T, baseURL string, urlSuffix string) TResult {
	return RequestAny[TResult](t, "", "GET", baseURL, urlSuffix, nil)
}

func RequestAuthGet[TResult any](t *testing.T, akey string, baseURL string, urlSuffix string) TResult {
	return RequestAny[TResult](t, akey, "GET", baseURL, urlSuffix, nil)
}

func RequestPost[TResult any](t *testing.T, baseURL string, urlSuffix string, body any) TResult {
	return RequestAny[TResult](t, "", "POST", baseURL, urlSuffix, body)
}

func RequestAuthPost[TResult any](t *testing.T, akey string, baseURL string, urlSuffix string, body any) TResult {
	return RequestAny[TResult](t, akey, "POST", baseURL, urlSuffix, body)
}

func RequestPut[TResult any](t *testing.T, baseURL string, urlSuffix string, body any) TResult {
	return RequestAny[TResult](t, "", "PUT", baseURL, urlSuffix, body)
}

func RequestAuthPUT[TResult any](t *testing.T, akey string, baseURL string, urlSuffix string, body any) TResult {
	return RequestAny[TResult](t, akey, "PUT", baseURL, urlSuffix, body)
}

func RequestPatch[TResult any](t *testing.T, baseURL string, urlSuffix string, body any) TResult {
	return RequestAny[TResult](t, "", "PATCH", baseURL, urlSuffix, body)
}

func RequestAuthPatch[TResult any](t *testing.T, akey string, baseURL string, urlSuffix string, body any) TResult {
	return RequestAny[TResult](t, akey, "PATCH", baseURL, urlSuffix, body)
}

func RequestDelete[TResult any](t *testing.T, baseURL string, urlSuffix string, body any) TResult {
	return RequestAny[TResult](t, "", "DELETE", baseURL, urlSuffix, body)
}

func RequestAuthDelete[TResult any](t *testing.T, akey string, baseURL string, urlSuffix string, body any) TResult {
	return RequestAny[TResult](t, akey, "DELETE", baseURL, urlSuffix, body)
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

func RequestAny[TResult any](t *testing.T, akey string, method string, baseURL string, urlSuffix string, body any) TResult {
	client := http.Client{}

	TPrintf("[-> REQUEST] (%s) %s%s [%s] [%s]\n", method, baseURL, urlSuffix, langext.Conditional(akey == "", "NO AUTH", "AUTH"), langext.Conditional(body == nil, "NO BODY", "BODY"))

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

	TPrintln("")
	TPrintf("----------------  RESPONSE (%d) ----------------\n", resp.StatusCode)
	TPrintln(langext.TryPrettyPrintJson(string(respBodyBin)))
	TryPrintTraceObj("----------------  --------  ----------------", respBodyBin, "")
	TPrintln("----------------  --------  ----------------")
	TPrintln("")

	if resp.StatusCode != 200 {
		TestFailFmt(t, "Statuscode != 200 (actual = %d)", resp.StatusCode)
	}

	var data TResult
	if err := json.Unmarshal(respBodyBin, &data); err != nil {
		TestFailErr(t, err)
	}

	return data
}

func RequestAuthAnyShouldFail(t *testing.T, akey string, method string, baseURL string, urlSuffix string, body any, expectedStatusCode int, errcode apierr.APIError) {
	client := http.Client{}

	TPrintf("[-> REQUEST] (%s) %s%s [%s] (should-fail with %d/%d)\n", method, baseURL, urlSuffix, langext.Conditional(akey == "", "NO AUTH", "AUTH"), expectedStatusCode, errcode)

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

	TPrintln("")
	TPrintf("----------------  RESPONSE (%d) ----------------\n", resp.StatusCode)
	TPrintln(langext.TryPrettyPrintJson(string(respBodyBin)))
	if (expectedStatusCode != 0 && resp.StatusCode != expectedStatusCode) || (expectedStatusCode == 0 && resp.StatusCode == 200) {
		TryPrintTraceObj("----------------  --------  ----------------", respBodyBin, "")
	}
	TPrintln("----------------  --------  ----------------")
	TPrintln("")

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

func TryPrintTraceObj(prefix string, body []byte, suffix string) {
	v1 := gin.H{}
	if err := json.Unmarshal(body, &v1); err == nil {
		if v2, ok := v1["traceObj"]; ok {
			if v3, ok := v2.(string); ok {
				if prefix != "" {
					TPrintln(prefix)
				}

				TPrintln(strings.TrimSpace(v3))

				if suffix != "" {
					TPrintln(suffix)
				}
			}
		}
	}
}
