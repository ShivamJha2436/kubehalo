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
	TargetRef   TargetRefSpec   `json:"targetRef"`
	Metrics     []MetricSpec    `json:"metrics"` // support multiple metrics
	ScaleUp     ScaleAction     `json:"scaleUp"`
	ScaleDown   ScaleAction     `json:"scaleDown"`
	MinReplicas int32           `json:"minReplicas"`
	MaxReplicas int32           `json:"maxReplicas"`
	EvaluationIntervalSeconds int32 `json:"evaluationIntervalSeconds"`
	Enabled     bool            `json:"enabled"`
	Behavior    *BehaviorSpec   `json:"behavior,omitempty"` // optional advanced behavior
}

// TargetRefSpec tells which deployment to scale
type TargetRefSpec struct {
	// +kubebuilder:validation:Enum=Deployment;StatefulSet
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// MetricSpec defines what metric to watch
type MetricSpec struct {
	Name       string  `json:"name"`
	Type       string  `json:"type"`
	Query      string  `json:"query,omitempty"`
	Threshold  int32   `json:"threshold"`
	TargetType string  `json:"targetType,omitempty"`
}

// ScaleAction defines scale-up or scale-down behavior
type ScaleAction struct {
	Step            int32 `json:"step"`
	CooldownSeconds int32 `json:"cooldownSeconds"`
}

// BehaviourSpec defines advanced scaling behaviour (optional)
type BehaviorSpec struct {
	// StabilizationWindowSeconds prevents flapping by requiring
	// metrics to stay above/below threshold for this duration before scaling.
	// +optional
	StabilizationWindowSeconds *int32 `json:"stabilizationWindowSeconds,omitempty"`

	// MaxScaleUpRate limits how quickly replicas can increase (absolute or percent).
	// +optional
	MaxScaleUpRate *int32 `json:"maxScaleUpRate,omitempty"`

	// MaxScaleDownRate limits how quickly replicas can decrease (absolute or percent).
	// +optional
	MaxScaleDownRate *int32 `json:"maxScaleDownRate,omitempty"`

	// Policy defines whether step values are absolute numbers or percentages.
	// +kubebuilder:validation:Enum=absolute;percent
	// +optional
	Policy string `json:"policy,omitempty"`
}

// ScalePolicyStatus reflects observed state
type ScalePolicyStatus struct {
	LastReconcileTime *metav1.Time      `json:"lastReconcileTime,omitempty"`
	LastScaleTime     *metav1.Time      `json:"lastScaleTime,omitempty"`
	CurrentReplicas   int32             `json:"currentReplicas"`
	DesiredReplicas   int32             `json:"desiredReplicas"`
	Conditions        []metav1.Condition `json:"conditions,omitempty"`
	LastMetricValues  map[string]int32   `json:"lastMetricValues,omitempty"`
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