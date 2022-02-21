package chewbyte

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"gopkg.in/yaml.v2"
)

func exportJson(data interface{}) (string, error) {
	str, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(str), nil
}

func exportYaml(data interface{}) (string, error) {
	str, err := yaml.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(str), nil
}

func ExportStr(data interface{}, format FormatHint) (string, error) {
	switch format {
	case String:
		str, ok := data.(string)
		if ok {
			return str, nil
		}
		return "", fmt.Errorf("unable to convert %s to string", reflect.TypeOf(data))
	case JSON:
		str, err := exportJson(data)
		if err != nil {
			return "", err
		}
		return str, nil
	case YAML:
		str, err := exportYaml(data)
		if err != nil {
			return "", err
		}
		return str, nil
	default:
		return "", fmt.Errorf("unsupported output format: %s", format)
	}
}

func ExportFile(data interface{}, path string, format FormatHint) error {
	str, err := ExportStr(data, format)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, []byte(str), 0644)
	if err != nil {
		return err
	}

	return nil
}
