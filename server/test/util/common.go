package util

import (
	"fmt"
	"reflect"
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

		strExp := fmt.Sprintf("%v", expected)
		strAct := fmt.Sprintf("%v", actual)

		if strings.Contains(strAct, "\n") {
			t.Errorf("Actual:\n~~~~~~~~~~~~~~~~\n%v\n~~~~~~~~~~~~~~~~\n\n", actual)
		} else {
			t.Errorf("Actual    := \"%v\"\n", actual)
		}

		if strings.Contains(strExp, "\n") {
			t.Errorf("Expected:\n~~~~~~~~~~~~~~~~\n%v\n~~~~~~~~~~~~~~~~\n\n", expected)
		} else {
			t.Errorf("Expected  := \"%v\"\n", expected)
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

func AssertStrRepEqual(t *testing.T, key string, expected any, actual any) {
	strExp := fmt.Sprintf("%v", unpointer(expected))
	strAct := fmt.Sprintf("%v", unpointer(actual))

	if strAct != strExp {
		t.Errorf("Value [%s] differs (%T <-> %T):\n", key, expected, actual)

		if strings.Contains(strAct, "\n") {
			t.Errorf("Actual:\n~~~~~~~~~~~~~~~~\n%v\n~~~~~~~~~~~~~~~~\n\n", strAct)
		} else {
			t.Errorf("Actual    := \"%v\"\n", strAct)
		}

		if strings.Contains(strExp, "\n") {
			t.Errorf("Expected:\n~~~~~~~~~~~~~~~~\n%v\n~~~~~~~~~~~~~~~~\n\n", strExp)
		} else {
			t.Errorf("Expected  := \"%v\"\n", strExp)
		}

		t.Error(string(debug.Stack()))

		t.FailNow()
	}
}

func AssertNotStrRepEqual(t *testing.T, key string, expected any, actual any) {
	strExp := fmt.Sprintf("%v", unpointer(expected))
	strAct := fmt.Sprintf("%v", unpointer(actual))

	if strAct == strExp {
		t.Errorf("Value [%s] does not differ (%T <-> %T):\n", key, expected, actual)

		if strings.Contains(strAct, "\n") {
			t.Errorf("Actual:\n~~~~~~~~~~~~~~~~\n%v\n~~~~~~~~~~~~~~~~\n\n", strAct)
		} else {
			t.Errorf("Actual    := \"%v\"\n", strAct)
		}

		if strings.Contains(strExp, "\n") {
			t.Errorf("Expected:\n~~~~~~~~~~~~~~~~\n%v\n~~~~~~~~~~~~~~~~\n\n", strExp)
		} else {
			t.Errorf("Expected  := \"%v\"\n", strExp)
		}

		t.Error(string(debug.Stack()))

		t.FailNow()
	}
}

func TestFail(t *testing.T, msg string) {
	t.Error(msg)
	t.Error(string(debug.Stack()))
	t.FailNow()
}

func TestFailFmt(t *testing.T, format string, args ...any) {
	t.Errorf(format, args...)
	t.Error(string(debug.Stack()))
	t.FailNow()
}

func TestFailErr(t *testing.T, e error) {
	t.Error(fmt.Sprintf("Failed with error:\n%s\n\nError:\n%+v\n\nTrace:\n%s", e.Error(), e, string(debug.Stack())))
	t.Error(string(debug.Stack()))
	t.FailNow()
}

func TestFailIfErr(t *testing.T, e error) {
	if e != nil {
		TestFailErr(t, e)
	}
}

func unpointer(v any) any {
	if v == nil {
		return v
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return v
		}
		val = val.Elem()
		return unpointer(val.Interface())
	}
	return v
}

func AssertMultiNonEmpty(t *testing.T, key string, args ...any) {
	for i := 0; i < len(args); i++ {

		reflval := reflect.ValueOf(args[i])

		if args[i] == nil || reflval.IsZero() {
			t.Errorf("Value %s[%d] is empty (AssertMultiNonEmpty)", key, i)
			t.FailNow()
		}
	}
}
