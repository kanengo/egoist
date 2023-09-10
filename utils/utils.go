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
