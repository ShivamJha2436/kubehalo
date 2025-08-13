package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	GroupVersion = schema.GroupVersion{Group: "kubehalo.sh", Version: "v1"}
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}
	AddToScheme = SchemeBuilder.AddToScheme
)

// ScalePolicySpec defines the desired state
type ScalePolicySpec struct {
	TargetRef TargetRefSpec `json:"targetRef"`
	Metric    MetricSpec    `json:"metric"`
	ScaleUp   ScaleAction   `json:"scaleUp"`
	ScaleDown ScaleAction   `json:"scaleDown"`
	MinReplicas int32       `json:"minReplicas"`
	MaxReplicas int32       `json:"maxReplicas"`
}

// TargetRefSpec tells which deployment to scale
type TargetRefSpec struct {
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// MetricSpec defines what metric to watch
type MetricSpec struct {
	Name      string  `json:"name"`
	Query     string  `json:"query"`
	Threshold float64 `json:"threshold"`
}

// ScaleAction defines scale-up or scale-down behavior
type ScaleAction struct {
	Step            int32 `json:"step"`
	CooldownSeconds int32 `json:"cooldownSeconds"`
}

// ScalePolicyStatus reflects observed state
type ScalePolicyStatus struct {
	LastScaleTime *metav1.Time `json:"lastScaleTime,omitempty"`
	CurrentReplicas int32      `json:"currentReplicas"`
	DesiredReplicas int32      `json:"desiredReplicas"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type ScalePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ScalePolicySpec   `json:"spec,omitempty"`
	Status ScalePolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type ScalePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ScalePolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ScalePolicy{}, &ScalePolicyList{})
}