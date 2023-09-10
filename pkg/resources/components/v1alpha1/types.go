package v1alpha1

import (
	"github.com/kanengo/egoist/pkg/resources/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Component struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ComponentSpec `json:"spec,omitempty"`
}

func (Component) Kind() string {
	return "Component"
}

type ComponentSpec struct {
	Type    string `json:"type"`
	Version string `json:"version"`
	//+optional
	IgnoreErrors bool                   `json:"ignoreErrors"`
	InitTimeout  string                 `json:"initTimeout"`
	Metadata     []common.NameValuePair `json:"metadata"`
}
