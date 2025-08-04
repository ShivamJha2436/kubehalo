package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ScalePolicySpec defines the desired state of ScalePolicy
type ScalePolicySpec struct {
	TargetDeployment string `json:"targetDeployment"`
	TargetNamespace  string `json:"targetNamespace"`
	MetricQuery      string `json:"metricQuery"`
	Thresholds       struct {
		ScaleUp   float64 `json:"scaleUp"`
		ScaleDown float64 `json:"scaleDown"`
	} `json:"thresholds"`
}

// ScalePolicy is the Schema for the scalepolicies API
type ScalePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ScalePolicySpec `json:"spec,omitempty"`
}

// ScalePolicyList contains a list of ScalePolicy
type ScalePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ScalePolicy `json:"items"`
}
