package chewbyte

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"gopkg.in/yaml.v2"
)

type MergeStrategy int

const (
	NoStrategy MergeStrategy = iota
	UniqueJoin
)

func (strategy MergeStrategy) String() string {
	switch strategy {
	case NoStrategy:
		return "no stategy"
	case UniqueJoin:
		return "unique join"
	}
	return fmt.Sprintf("unknown strategy(%d)", strategy)
}

func mergeMaps(m1, m2 map[string]interface{}, mergeStrategy MergeStrategy) (map[string]interface{}, error) {
	switch mergeStrategy {
	case NoStrategy:
		result := make(map[string]interface{})
		for k, v := range m1 {
			result[k] = v
		}
		for k, v := range m2 {
			result[k] = v
		}
		return result, nil
	case UniqueJoin:
		result := make(map[string]interface{})
		for k, v := range m1 {
			result[k] = v
		}
		for k, v := range m2 {
			_, exists := m1[k]
			if exists {
				return nil, fmt.Errorf("key %s is not unique", k)
			}
			result[k] = v
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unknown merge strategy: %s", mergeStrategy)
	}
}

func mergeData(d1, d2 interface{}, mergeStrategy MergeStrategy) (interface{}, error) {
	stringSlice1, valid1 := d1.([]string)
	stringSlice2, valid2 := d2.([]string)
	if valid1 && valid2 {
		return append(stringSlice1, stringSlice2...), nil
	}

	map1, valid1 := d1.(map[string]interface{})
	map2, valid2 := d2.(map[string]interface{})
	if valid1 && valid2 {
		merged, err := mergeMaps(map1, map2, mergeStrategy)
		if err != nil {
			return nil, err
		}

		return merged, nil
	}

	type1 := reflect.TypeOf(d1)
	type2 := reflect.TypeOf(d2)

	if type1 == type2 {
		return nil, fmt.Errorf("unsupported datatype for merging: %s", type1)
	}

	return nil, fmt.Errorf("unable to merge %s with %s", type1, type2)
}

func ImportFiles(paths []string, mergeStrategy MergeStrategy) (interface{}, error) {
	var mergedData interface{}

	for _, path := range paths {
		fileData, err := ImportFile(path)
		if err != nil {
			return mergedData, err
		}

		mergedData, err = mergeData(mergedData, fileData, mergeStrategy)
		if err != nil {
			return mergedData, err
		}
	}

	return mergedData, nil
}

func ImportFile(path string) (interface{}, error) {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	data, err := ImportStr(string(fileContent), path)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func importJson(content string) (interface{}, error) {
	var data interface{}
	err := json.Unmarshal([]byte(content), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func importYaml(content string) (interface{}, error) {
	var data interface{}
	err := yaml.Unmarshal([]byte(content), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func ImportStr(content string, filenameHint string) (interface{}, error) {
	format := DetectFormat(filenameHint)
	switch format {
	case String:
		return content, nil
	case JSON:
		data, err := importJson(content)
		if err != nil {
			return nil, err
		}
		return data, nil
	case YAML:
		data, err := importYaml(content)
		if err != nil {
			return nil, err
		}
		return data, nil
	case Jsonnet:
		data, err := importJsonnet(content, filenameHint)
		if err != nil {
			return nil, err
		}
		return data, nil
	default:
		return nil, fmt.Errorf("unsupported input format: %s", format)
	}
}
