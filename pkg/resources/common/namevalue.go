package common

import (
	apiextensionsV1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"strconv"
)

type NameValuePair struct {
	Name  string       `json:"name"`
	Value DynamicValue `json:"value,omitempty"`
}

type DynamicValue struct {
	apiextensionsV1.JSON `json:",inline"`
}

// String returns the string representation of the raw value.
// If the value is a string, it will be unquoted as the string is guaranteed to be a JSON serialized string.
func (d *DynamicValue) String() string {
	s := string(d.Raw)
	c, err := strconv.Unquote(s)
	if err == nil {
		s = c
	}
	return s
}
