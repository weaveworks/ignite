package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Metadata of each individual image
type ImageData struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   ImageSpec   `json:"spec"`
	Status ImageStatus `json:"status"`
}

type ImageSpec struct {
	Source ImageSource `json:"source"`
}

type ImageSource struct {
	Type   string `json:"type"`
	Digest string `json:"digest"`
	Name   string `json:"name"`
	Size   uint64 `json:"size"`
}

type ImageStatus struct {
	Devices []PoolDevice `json:"devices"`
}

type PoolDevice struct {
	Name      string
	Type      string
	UID       string
	ParentUID string
	Blocks    uint64
}

// Metadata of the common pool
// The devices are separated in the per-image metadata
type PoolData struct {
	Blocks    uint64
	BlockSize uint64
}
