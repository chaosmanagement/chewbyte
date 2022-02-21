package chewbyte

import (
	"fmt"
	"path/filepath"
	"strings"
)

type FormatHint int

const (
	String FormatHint = iota
	JSON
	YAML
	Jsonnet
)

func (format FormatHint) String() string {
	switch format {
	case String:
		return "string"
	case JSON:
		return "JSON"
	case YAML:
		return "YAML"
	case Jsonnet:
		return "Jsonnet"
	}
	return fmt.Sprintf("unknown format(%d)", format)
}

func DetectFormat(path string) FormatHint {
	ext := strings.ToLower(filepath.Ext(path))
	if len(ext) > 0 && ext[0] == '.' {
		ext = ext[1:]
	}

	switch ext {
	case "json":
		return JSON
	case "yml":
		return YAML
	case "yaml":
		return YAML
	case "jsonnet":
		return Jsonnet
	case "txt":
		return String
	default:
		return String
	}
}
