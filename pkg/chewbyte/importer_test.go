package chewbyte

import (
	"fmt"
	"reflect"
	"testing"
)

func compareStringSlices(d1, d2 []string) bool {
	if len(d1) != len(d2) {
		return false
	}

	for i, v := range d1 {
		if v != d2[i] {
			return false
		}
	}

	return true
}

func compareStringKeyedMaps(m1, m2 map[string]interface{}) bool {
	if len(m1) != len(m2) {
		return false // unequal number of items in maps
	}

	for k, v1 := range m1 {
		v2, ok := m2[k]
		if !ok {
			return false // key exists only in first map
		}
		if reflect.TypeOf(v1) != reflect.TypeOf(v2) {
			return false // values are of diffrent types
		}
		if v1 != v2 {
			return false // values differ
		}
	}

	return true
}

func TestCompareStringSlices(t *testing.T) {
	testTable := []struct {
		d1, d2   []string
		expected bool
	}{
		{[]string{}, []string{}, true},
		{[]string{"foo"}, []string{"foo"}, true},
		{[]string{"foo"}, []string{"bar"}, false},
	}

	for _, testData := range testTable {
		testName := fmt.Sprintf("%s == %s = %t", testData.d1, testData.d2, testData.expected)
		t.Run(testName, func(t *testing.T) {
			output := compareStringSlices(testData.d1, testData.d2)
			if output != testData.expected {
				t.Errorf("%t != %t", output, testData.expected)
			}
		})
	}
}

func TestCompareStringKeyedMaps(t *testing.T) {
	testTable := []struct {
		m1, m2   map[string]interface{}
		expected bool
	}{
		{map[string]interface{}{}, map[string]interface{}{}, true},
		{map[string]interface{}{"a": 1}, map[string]interface{}{"a": 1}, true},
		{map[string]interface{}{"a": 1}, map[string]interface{}{"b": 2}, false},
		{map[string]interface{}{"a": 1}, map[string]interface{}{"a": 2}, false},
		{map[string]interface{}{"a": 1}, map[string]interface{}{}, false},
		{map[string]interface{}{}, map[string]interface{}{"a": 1}, false},
	}

	for _, testData := range testTable {
		testName := fmt.Sprintf("%s == %s = %t", testData.m1, testData.m2, testData.expected)
		t.Run(testName, func(t *testing.T) {
			output := compareStringKeyedMaps(testData.m1, testData.m2)
			if output != testData.expected {
				t.Errorf("%t != %t", output, testData.expected)
			}
		})
	}
}

func TestMerge(t *testing.T) {
	testTable := []struct {
		d1, d2        interface{}
		strategy      MergeStrategy
		expected      interface{}
		expectedError bool
	}{
		{[]string{}, []string{}, NoStrategy, []string{}, false},
		{[]string{"Ala", "ma"}, []string{"kota"}, NoStrategy, []string{"Ala", "ma", "kota"}, false},
		{map[string]interface{}{}, map[string]interface{}{}, NoStrategy, map[string]interface{}{}, false},
		{map[string]interface{}{"foo": 1}, map[string]interface{}{"bar": 2}, NoStrategy, map[string]interface{}{"foo": 1, "bar": 2}, false},
	}

	for _, testData := range testTable {
		testName := fmt.Sprintf("%s merge %s == %s", testData.d1, testData.d2, testData.expected)
		t.Run(testName, func(t *testing.T) {
			output, err := mergeData(testData.d1, testData.d2, testData.strategy)

			if testData.expectedError {
				if err == nil {
					t.Fatalf("expected an error but the error was nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("mergeData returned error %s", err)
			}

			outputTypeStr := reflect.TypeOf(output)
			expectedTypeStr := reflect.TypeOf(testData.expected)

			switch output.(type) {
			case []string:
				expectedSlice, expectedValid := testData.expected.([]string)
				if expectedValid {
					if !compareStringSlices(output.([]string), expectedSlice) {
						t.Errorf("%s != %s", output, testData.expected)
					}
				} else {
					t.Errorf("output type %s != expected type %s", outputTypeStr, expectedTypeStr)
				}
			case map[string]interface{}:
				expectedSlice, expectedValid := testData.expected.(map[string]interface{})
				if expectedValid {
					if !compareStringKeyedMaps(output.(map[string]interface{}), expectedSlice) {
						t.Errorf("%s != %s", output, testData.expected)
					}
				} else {
					t.Errorf("output type %s != expected type %s", outputTypeStr, expectedTypeStr)
				}
			default:
				if output != testData.expected {
					t.Errorf("%s != %s", output, testData.expected)
				}
			}
		})
	}
}
