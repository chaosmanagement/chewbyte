package chewbyte

import (
	"fmt"
	"reflect"
	"testing"
)

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

			if !reflect.DeepEqual(output, testData.expected) {
				t.Errorf("%v != %v", output, testData.expected)
			}
		})
	}
}
