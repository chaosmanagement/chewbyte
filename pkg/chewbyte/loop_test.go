package chewbyte

import (
	"testing"
)

func TestLoops(t *testing.T) {
	testTable := []struct {
		filenameHint string
		content      string
	}{
		{"simple_string.txt", "foo"},
		{"empty_json.json", "{}"},
		{"yml_scalar.yml", "123\n"},
	}

	for _, testData := range testTable {
		t.Run(testData.filenameHint, func(t *testing.T) {
			data, err := ImportStr(testData.content, testData.filenameHint)
			if err != nil {
				t.Fatalf("import returned error: %v", err)
			}
			format := DetectFormat(testData.filenameHint)
			result, err := ExportStr(data, format)
			if err != nil {
				t.Fatalf("export returned error: %v", err)
			}
			if testData.content != result {
				t.Errorf("test data != result\n%v\n%v", testData.content, result)
			}
		})
	}
}
