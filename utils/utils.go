package utils

import (
	"path/filepath"
	"strings"
)

func IsYaml(fileName string) bool {
	ext := strings.ToLower(filepath.Ext(fileName))
	if ext == ".yaml" || ext == ".yml" {
		return true
	}

	return false
}

func IsTruthy(val string) bool {
	switch strings.ToLower(strings.TrimSpace(val)) {
	case "y", "yes", "true", "t", "on", "1":
		return true
	default:
		return false
	}
}
