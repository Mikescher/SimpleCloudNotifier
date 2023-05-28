package util

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"math"
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

	// try to fix types, kinda hacky, but its only unit tests...
	switch vex := expected.(type) {
	case int:
		switch vac := actual.(type) {
		case int:
			// same
		case int32:
			expected = int64(vex)
			actual = int64(vac)
		case int64:
			expected = int64(vex)
		case float32:
			if IsWholeFloat(vac) {
				expected = int64(vex)
				actual = int64(vac)
			}
		case float64:
			if IsWholeFloat(vac) {
				expected = int64(vex)
				actual = int64(vac)
			}
		}
	case int32:
		switch vac := actual.(type) {
		case int:
			expected = int64(vex)
			actual = int64(vac)
		case int32:
			// same
		case int64:
			expected = int64(vex)
		case float32:
			if IsWholeFloat(vac) {
				expected = int64(vex)
				actual = int64(vac)
			}
		case float64:
			if IsWholeFloat(vac) {
				expected = int64(vex)
				actual = int64(vac)
			}
		}
	case int64:
		switch vac := actual.(type) {
		case int:
			actual = int64(vac)
		case int32:
			actual = int64(vac)
		case int64:
			// same
		case float32:
			if IsWholeFloat(vac) {
				actual = int64(vac)
			}
		case float64:
			if IsWholeFloat(vac) {
				actual = int64(vac)
			}
		}
	case float32:
		switch vac := actual.(type) {
		case int:
			if IsWholeFloat(vex) {
				expected = int64(vex)
				actual = int64(vac)
			}
		case int32:
			if IsWholeFloat(vex) {
				expected = int64(vex)
				actual = int64(vac)
			}
		case int64:
			if IsWholeFloat(vex) {
				expected = int64(vex)
			}
		case float32:
			// same
		case float64:
			expected = float64(vex)
		}
	case float64:
		switch vac := actual.(type) {
		case int:
			if IsWholeFloat(vex) {
				expected = int64(vex)
				actual = int64(vac)
			}
		case int32:
			if IsWholeFloat(vex) {
				expected = int64(vex)
				actual = int64(vac)
			}
		case int64:
			if IsWholeFloat(vex) {
				expected = int64(vex)
			}
		case float32:
			actual = float64(vac)
		case float64:
			// same
		}

	}

	if langext.IsNil(expected) && langext.IsNil(actual) {
		return
	}

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

func AssertTrue(t *testing.T, key string, v bool) {
	if !v {
		t.Errorf("AssertTrue(%s) failed", key)
		t.Error(string(debug.Stack()))
		t.FailNow()
	}
}

func AssertNotDefault[T comparable](t *testing.T, key string, v T) {
	if v == *new(T) {
		t.Errorf("AssertNotDefault(%s) failed", key)
		t.Error(ljson(v))
		t.Error(string(debug.Stack()))
		t.FailNow()
	}
}

func AssertNotDefaultAny[T any](t *testing.T, key string, v T) {
	if ljson(v) == ljson(*new(T)) {
		t.Errorf("AssertNotDefault(%s) failed", key)
		t.Error(ljson(v))
		t.Error(string(debug.Stack()))
		t.FailNow()
	}
}

func AssertNotEqual(t *testing.T, key string, expected any, actual any) {
	if expected == actual || (langext.IsNil(expected) && langext.IsNil(actual)) {
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

func AssertArrAny[T any](t *testing.T, key string, arr []T, fn func(T) bool) {
	if !langext.ArrAny(arr, fn) {
		t.Errorf("AssertArrAny(%s) failed", key)
		t.Error(string(debug.Stack()))
		t.FailNow()
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

func AssertMappedSet[T langext.OrderedConstraint](t *testing.T, key string, expected []T, values []gin.H, objkey string) {

	actual := make([]T, 0)
	for idx, vv := range values {
		if tv, ok := vv[objkey].(T); ok {
			actual = append(actual, tv)
		} else {
			TestFailFmt(t, "[%s]->[%d] is wrong type (expected: %T, actual: %T)", key, idx, *new(T), vv)
		}
	}

	langext.Sort(actual)
	langext.Sort(expected)

	if !langext.ArrEqualsExact(actual, expected) {
		t.Errorf("Value [%s] differs (%T <-> %T):\n", key, expected, actual)

		t.Errorf("Actual    := [%v]\n", ljson(actual))
		t.Errorf("Expected  := [%v]\n", ljson(expected))

		t.Error(string(debug.Stack()))

		t.FailNow()
	}
}

func AssertMappedArr[T langext.OrderedConstraint](t *testing.T, key string, expected []T, values []gin.H, objkey string) {

	actual := make([]T, 0)
	for idx, vv := range values {
		if tv, ok := vv[objkey].(T); ok {
			actual = append(actual, tv)
		} else {
			TestFailFmt(t, "[%s]->[%d] is wrong type (expected: %T, actual: %T)", key, idx, *new(T), vv)
		}
	}

	if !langext.ArrEqualsExact(actual, expected) {
		t.Errorf("Value [%s] differs (%T <-> %T):\n", key, expected, actual)

		t.Errorf("Actual    := [%v]\n", actual)
		t.Errorf("Expected  := [%v]\n", expected)

		t.Error(string(debug.Stack()))

		t.FailNow()
	}
}

func IsWholeFloat[T langext.FloatConstraint](v T) bool {
	_, frac := math.Modf(math.Abs(float64(v)))
	return frac == 0.0
}

func ljson(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
