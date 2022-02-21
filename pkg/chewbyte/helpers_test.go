package chewbyte

import (
	"fmt"
	"testing"
)

func TestDetectFormat(t *testing.T) {
	testTable := []struct {
		path           string
		expectedFormat FormatHint
	}{
		{"/home/foo.json", JSON},
		{"/home/foo.yaml", YAML},
		{"/home/foo.yml", YAML},
		{"/home/foo.jsonnet", Jsonnet},
	}

	for _, testData := range testTable {
		testName := fmt.Sprintf("%s == %s", testData.path, testData.expectedFormat)
		t.Run(testName, func(t *testing.T) {
			output := DetectFormat(testData.path)
			if output != testData.expectedFormat {
				t.Errorf("%s != %s", output, testData.expectedFormat)
			}
		})
	}
}
