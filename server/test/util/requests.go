package util

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"io"
	"net/http"
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

func RequestAny[TResult any](t *testing.T, akey string, method string, baseURL string, urlSuffix string, body any) TResult {
	client := http.Client{}

	fmt.Printf("[-> REQUEST] (%s) %s%s [%s] [%s]\n", method, baseURL, urlSuffix, langext.Conditional(akey == "", "NO AUTH", "AUTH"), langext.Conditional(body == nil, "NO BODY", "BODY"))

	bytesbody := make([]byte, 0)
	if body != nil {
		bjson, err := json.Marshal(body)
		if err != nil {
			TestFailErr(t, err)
		}
		bytesbody = bjson
	}

	req, err := http.NewRequest(method, baseURL+urlSuffix, bytes.NewReader(bytesbody))
	if err != nil {
		TestFailErr(t, err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
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

	fmt.Println("")
	fmt.Printf("----------------  RESPONSE (%d) ----------------\n", resp.StatusCode)
	fmt.Println(langext.TryPrettyPrintJson(string(respBodyBin)))
	fmt.Println("----------------  --------  ----------------")
	fmt.Println("")

	if resp.StatusCode != 200 {
		TestFail(t, "Statuscode != 200")
	}

	var data TResult
	if err := json.Unmarshal(respBodyBin, &data); err != nil {
		TestFailErr(t, err)
	}

	return data
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

func RequestAuthAnyShouldFail(t *testing.T, akey string, method string, baseURL string, urlSuffix string, body any, statusCode int, errcode apierr.APIError) {
	client := http.Client{}

	fmt.Printf("[-> REQUEST] (%s) %s%s [%s] (should-fail with %d/%d)\n", method, baseURL, urlSuffix, langext.Conditional(akey == "", "NO AUTH", "AUTH"), statusCode, errcode)

	bytesbody := make([]byte, 0)
	if body != nil {
		bjson, err := json.Marshal(body)
		if err != nil {
			TestFailErr(t, err)
		}
		bytesbody = bjson
	}

	req, err := http.NewRequest(method, baseURL+urlSuffix, bytes.NewReader(bytesbody))
	if err != nil {
		TestFailErr(t, err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
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

	fmt.Println("")
	fmt.Printf("----------------  RESPONSE (%d) ----------------\n", resp.StatusCode)
	fmt.Println(langext.TryPrettyPrintJson(string(respBodyBin)))
	fmt.Println("----------------  --------  ----------------")
	fmt.Println("")

	if resp.StatusCode != statusCode {
		fmt.Println("Request: " + method + " :: " + baseURL + urlSuffix)
		fmt.Println(string(respBodyBin))
		TestFailFmt(t, "Statuscode != %d (expected failure)", statusCode)
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

	if v, ok := data["error"]; ok {
		if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", errcode) {
			TestFailFmt(t, "wrong errorcode (expected: %d), (actual: %v)", errcode, v)
		}
	} else {
		TestFail(t, "missing response['error']")
	}
}
