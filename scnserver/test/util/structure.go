package util

import (
	"encoding/json"
	"fmt"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"reflect"
	"testing"
	"time"
)

func AssertJsonStructureMatch(t *testing.T, key string, jsonData string, expected map[string]any) {

	realData := make(map[string]any)

	err := json.Unmarshal([]byte(jsonData), &realData)
	if err != nil {
		t.Errorf("Failed to decode json of [%s]: %s", key, err.Error())
		return
	}

	assertjsonStructureMatchMapObject(t, expected, realData, key)
}

func assertJsonStructureMatch(t *testing.T, schema any, realValue any, keyPath string) {

	if strschema, ok := schema.(string); ok {

		assertjsonStructureMatchSingleValue(t, strschema, realValue, keyPath)

	} else if mapschema, ok := schema.(map[string]any); ok {

		if reflect.ValueOf(realValue).Kind() != reflect.Map {
			t.Errorf("Key < %s > is not a object (its actually %T: '%v')", keyPath, realValue, realValue)
			return
		}
		if _, ok := realValue.(map[string]any); !ok {
			t.Errorf("Key < %s > is not a object[recursive] (its actually %T: '%v')", keyPath, realValue, realValue)
			return
		}

		assertjsonStructureMatchMapObject(t, mapschema, realValue.(map[string]any), keyPath)

	} else if arrschema, ok := schema.([]any); ok && len(arrschema) == 1 {

		if _, ok := realValue.([]any); !ok {
			t.Errorf("Key < %s > is not a array[recursive] (its actually %T: '%v')", keyPath, realValue, realValue)
			return
		}

		assertjsonStructureMatchArray(t, arrschema, realValue.([]any), keyPath)

	} else {
		t.Errorf("Unknown schema type '%s' for key < %s >", schema, keyPath)
	}
}

func assertjsonStructureMatchSingleValue(t *testing.T, strschema string, realValue any, keyPath string) {
	switch strschema {
	case "id":
		if _, ok := realValue.(string); !ok {
			t.Errorf("Key < %s > is not a string<id> (its actually %T: '%v')", keyPath, realValue, realValue)
			return
		}
		if len(realValue.(string)) != 24 { //TODO validate checksum?
			t.Errorf("Key < %s > is not a valid entity-id date (its '%v')", keyPath, realValue)
			return
		}
	case "string":
		if _, ok := realValue.(string); !ok {
			t.Errorf("Key < %s > is not a string (its actually %T: '%v')", keyPath, realValue, realValue)
			return
		}
	case "null":
		if !langext.IsNil(realValue) {
			t.Errorf("Key < %s > is not a NULL (its actually %T: '%v')", keyPath, realValue, realValue)
			return
		}
	case "string|null":
		if langext.IsNil(realValue) {
			return // OK
		} else if _, ok := realValue.(string); !ok {
			return // OK
		} else {
			t.Errorf("Key < %s > is not a string|null (its actually %T: '%v')", keyPath, realValue, realValue)
			return
		}
	case "rfc3339":
		if _, ok := realValue.(string); !ok {
			t.Errorf("Key < %s > is not a string<rfc3339> (its actually %T: '%v')", keyPath, realValue, realValue)
			return
		}
		if _, err := time.Parse(time.RFC3339, realValue.(string)); err != nil {
			t.Errorf("Key < %s > is not a valid rfc3339 date (its '%v')", keyPath, realValue)
			return
		}
	case "rfc3339|null":
		if langext.IsNil(realValue) {
			return // OK
		}
		if _, ok := realValue.(string); !ok {
			t.Errorf("Key < %s > is not a string<rfc3339> (its actually %T: '%v')", keyPath, realValue, realValue)
			return
		}
		if _, err := time.Parse(time.RFC3339, realValue.(string)); err != nil {
			t.Errorf("Key < %s > is not a valid rfc3339 date (its '%v')", keyPath, realValue)
			return
		}
	case "int":
		if _, ok := realValue.(float64); !ok {
			t.Errorf("Key < %s > is not a int (its actually %T: '%v')", keyPath, realValue, realValue)
			return
		}
		if realValue.(float64) != float64(int(realValue.(float64))) {
			t.Errorf("Key < %s > is not a int (its actually %T: '%v')", keyPath, realValue, realValue)
			return
		}
	case "float":
		if _, ok := realValue.(float64); !ok {
			t.Errorf("Key < %s > is not a int (its actually %T: '%v')", keyPath, realValue, realValue)
			return
		}
	case "bool":
		if _, ok := realValue.(bool); !ok {
			t.Errorf("Key < %s > is not a int (its actually %T: '%v')", keyPath, realValue, realValue)
			return
		}
	default:
		t.Errorf("Unknown schema type '%s' for key < %s >", strschema, keyPath)
		return
	}
}

func assertjsonStructureMatchMapObject(t *testing.T, mapschema map[string]any, realValue map[string]any, keyPath string) {

	for k := range mapschema {
		if _, ok := realValue[k]; !ok {
			t.Errorf("Missing Key: < %s >", keyPath)
		}
	}

	for k := range realValue {
		if _, ok := mapschema[k]; !ok {
			t.Errorf("Additional key: < %s >", keyPath)
		}
	}

	for k, v := range realValue {

		kpath := keyPath + "." + k

		schema, ok := mapschema[k]

		if !ok {
			t.Errorf("Key < %s > is missing in response", kpath)
			continue
		}

		assertJsonStructureMatch(t, schema, v, kpath)

	}

}

func assertjsonStructureMatchArray(t *testing.T, arrschema []any, realValue []any, keyPath string) {

	if len(arrschema) != 1 {
		t.Errorf("Array schema must have exactly one element, but got %d", len(arrschema))
		return
	}

	for i, realArrVal := range realValue {
		assertJsonStructureMatch(t, arrschema[0], realArrVal, fmt.Sprintf("%s[%d]", keyPath, i))
	}

}
