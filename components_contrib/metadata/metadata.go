package metadata

type Base struct {
	Name       string
	Properties map[string]string `json:"properties,omitempty"`
}
