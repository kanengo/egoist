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

	CloudEventMetadataType   = "cloudevent.type"
	CloudEventMetadataSource = "cloudevent.source"
	CloudEventMetadataId     = "cloudevent.id"
)

var (
	CloudEventMetadataKeys = map[string]struct{}{
		CloudEventMetadataType:   {},
		CloudEventMetadataSource: {},
		CloudEventMetadataId:     {},
	}
)
