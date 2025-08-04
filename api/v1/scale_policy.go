package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ScalePolicySpec defines the desired state of ScalePolicy
type ScalePolicySpec struct {
	DeploymentName string `json:"deploymentName"`
	Namespace      string `json:"namespace"`
	MinReplicas    int32  `json:"minReplicas"`
	MaxReplicas    int32  `json:"maxReplicas"`
	Query          string `json:"query"`
	Threshold      float64 `json:"threshold"`
	ScaleUp        bool    `json:"scaleUp"`
}

// ScalePolicyStatus defines the observed state of ScalePolicy
type ScalePolicyStatus struct {
	LastScaleTime metav1.Time `json:"lastScaleTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ScalePolicy is the Schema for the scalepolicies API
type ScalePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ScalePolicySpec   `json:"spec,omitempty"`
	Status ScalePolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ScalePolicyList contains a list of ScalePolicy
type ScalePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ScalePolicy `json:"items"`
}
