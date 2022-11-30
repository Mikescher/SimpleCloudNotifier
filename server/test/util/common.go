package util

import (
	"fmt"
	"runtime/debug"
	"strings"
	"testing"
)

func AssertJsonMapEqual(t *testing.T, key string, expected map[string]any, actual map[string]any) {
	mkeys := make(map[string]string)
	for k := range expected {
		mkeys[k] = k
	}
	for k := range actual {
		mkeys[k] = k
	}

	for mapkey := range mkeys {

		if _, ok := expected[mapkey]; !ok {
			TestFailFmt(t, "Missing Key expected['%s'] ( assertJsonMapEqual[%s] )", mapkey, key)
		}
		if _, ok := actual[mapkey]; !ok {
			TestFailFmt(t, "Missing Key actual['%s'] ( assertJsonMapEqual[%s] )", mapkey, key)
		}

		AssertEqual(t, key+"."+mapkey, expected[mapkey], actual[mapkey])
	}

}

func AssertEqual(t *testing.T, key string, expected any, actual any) {
	if expected != actual {
		t.Errorf("Value [%s] differs (%T <-> %T):\n", key, expected, actual)

		str1 := fmt.Sprintf("%v", expected)
		str2 := fmt.Sprintf("%v", actual)

		if strings.Contains(str1, "\n") {
			t.Errorf("Actual:\n~~~~~~~~~~~~~~~~\n%v\n~~~~~~~~~~~~~~~~\n\n", expected)
		} else {
			t.Errorf("Actual    := \"%v\"\n", expected)
		}

		if strings.Contains(str2, "\n") {
			t.Errorf("Expected:\n~~~~~~~~~~~~~~~~\n%v\n~~~~~~~~~~~~~~~~\n\n", actual)
		} else {
			t.Errorf("Expected  := \"%v\"\n", actual)
		}

		t.Error(string(debug.Stack()))

		t.FailNow()
	}
}

func AssertStrRepEqual(t *testing.T, key string, expected any, actual any) {
	str1 := fmt.Sprintf("%v", expected)
	str2 := fmt.Sprintf("%v", actual)

	if str1 != str2 {
		t.Errorf("Value [%s] differs (%T <-> %T):\n", key, expected, actual)

		if strings.Contains(str1, "\n") {
			t.Errorf("Actual:\n~~~~~~~~~~~~~~~~\n%v\n~~~~~~~~~~~~~~~~\n\n", expected)
		} else {
			t.Errorf("Actual    := \"%v\"\n", expected)
		}

		if strings.Contains(str2, "\n") {
			t.Errorf("Expected:\n~~~~~~~~~~~~~~~~\n%v\n~~~~~~~~~~~~~~~~\n\n", actual)
		} else {
			t.Errorf("Expected  := \"%v\"\n", actual)
		}

		t.Error(string(debug.Stack()))

		t.FailNow()
	}
}

func AssertNotEqual(t *testing.T, key string, expected any, actual any) {
	if expected == actual {
		t.Errorf("Value [%s] does not differ (%T <-> %T):\n", key, expected, actual)

		str1 := fmt.Sprintf("%v", expected)
		str2 := fmt.Sprintf("%v", actual)

		if strings.Contains(str1, "\n") {
			t.Errorf("Actual:\n~~~~~~~~~~~~~~~~\n%v\n~~~~~~~~~~~~~~~~\n\n", expected)
		} else {
			t.Errorf("Actual       := \"%v\"\n", expected)
		}

		if strings.Contains(str2, "\n") {
			t.Errorf("Not Expected:\n~~~~~~~~~~~~~~~~\n%v\n~~~~~~~~~~~~~~~~\n\n", actual)
		} else {
			t.Errorf("Not Expected := \"%v\"\n", actual)
		}

		t.Error(string(debug.Stack()))

		t.FailNow()
	}
}

func TestFail(t *testing.T, msg string) {
	t.Error(msg)
	t.FailNow()
}

func TestFailFmt(t *testing.T, format string, args ...any) {
	t.Errorf(format, args...)
	t.FailNow()
}

func TestFailErr(t *testing.T, e error) {
	t.Error(fmt.Sprintf("Failed with error:\n%s\n\nError:\n%+v\n\nTrace:\n%s", e.Error(), e, string(debug.Stack())))
	t.FailNow()
}
