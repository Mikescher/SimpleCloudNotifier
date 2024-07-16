package util

import (
	"encoding/json"
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

	AssertJsonStructureMatchOfMap(t, key, realData, expected)
}

func AssertJsonStructureMatchOfMap(t *testing.T, key string, realData map[string]any, expected map[string]any) {

	for k := range expected {
		if _, ok := realData[k]; !ok {
			t.Errorf("Missing Key in data '%s': [[%s]]", key, k)
		}
	}

	for k := range realData {
		if _, ok := expected[k]; !ok {
			t.Errorf("Additional key in data '%s': [[%s]]", key, k)
		}
	}

	for k, v := range realData {

		schema, ok := expected[k]

		if !ok {
			continue
		}

		if strschema, ok := schema.(string); ok {
			switch strschema {
			case "id":
				if _, ok := v.(string); !ok {
					t.Errorf("Key [[%s]] in data '%s' is not a string<id> (its actually %T: '%v')", k, key, v, v)
					continue
				}
				if len(v.(string)) != 24 { //TODO validate checksum?
					t.Errorf("Key [[%s]] in data '%s' is not a valid entity-id date (its '%v')", k, key, v)
					continue
				}
			case "string":
				if _, ok := v.(string); !ok {
					t.Errorf("Key [[%s]] in data '%s' is not a string (its actually %T: '%v')", k, key, v, v)
					continue
				}
			case "null":
				if !langext.IsNil(v) {
					t.Errorf("Key [[%s]] in data '%s' is not a NULL (its actually %T: '%v')", k, key, v, v)
					continue
				}
			case "rfc3339":
				if _, ok := v.(string); !ok {
					t.Errorf("Key [[%s]] in data '%s' is not a string<rfc3339> (its actually %T: '%v')", k, key, v, v)
					continue
				}
				if _, err := time.Parse(time.RFC3339, v.(string)); err != nil {
					t.Errorf("Key [[%s]] in data '%s' is not a valid rfc3339 date (its '%v')", k, key, v)
					continue
				}
			case "int":
				if _, ok := v.(float64); !ok {
					t.Errorf("Key [[%s]] in data '%s' is not a int (its actually %T: '%v')", k, key, v, v)
					continue
				}
				if v.(float64) != float64(int(v.(float64))) {
					t.Errorf("Key [[%s]] in data '%s' is not a int (its actually %T: '%v')", k, key, v, v)
					continue
				}
			case "float":
				if _, ok := v.(float64); !ok {
					t.Errorf("Key [[%s]] in data '%s' is not a int (its actually %T: '%v')", k, key, v, v)
					continue
				}
			case "bool":
				if _, ok := v.(bool); !ok {
					t.Errorf("Key [[%s]] in data '%s' is not a int (its actually %T: '%v')", k, key, v, v)
					continue
				}
			case "object":
				if reflect.ValueOf(v).Kind() != reflect.Map {
					t.Errorf("Key [[%s]] in data '%s' is not a object (its actually %T: '%v')", k, key, v, v)
					continue
				}
			case "array":
				if reflect.ValueOf(v).Kind() != reflect.Array {
					t.Errorf("Key [[%s]] in data '%s' is not a array (its actually %T: '%v')", k, key, v, v)
					continue
				}
			}
		} else if mapschema, ok := schema.(map[string]any); ok {
			if reflect.ValueOf(v).Kind() != reflect.Map {
				t.Errorf("Key [[%s]] in data '%s' is not a object (its actually %T: '%v')", k, key, v, v)
				continue
			}
			if _, ok := v.(map[string]any); !ok {
				t.Errorf("Key [[%s]] in data '%s' is not a object[recursive] (its actually %T: '%v')", k, key, v, v)
				continue
			}
			AssertJsonStructureMatchOfMap(t, key+".["+k+"]", v.(map[string]any), mapschema)
		}
	}
}
