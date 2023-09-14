package components

import (
	"strings"
)

func IsInitialVersion(version string) bool {
	v := strings.ToLower(version)
	return v == "" || v == UnstableVersion || v == FirstStableVersion
}

const (
	// Unstable version (v0).
	UnstableVersion = "v0"

	// First stable version (v1).
	FirstStableVersion = "v1"
)
