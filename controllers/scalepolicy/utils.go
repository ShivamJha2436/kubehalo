package scalepolicy

import (
	"fmt"
	"math"

	kubehalov1 "github.com/ShivamJha2436/kubehalo/api/kubehalo/v1"
	"github.com/ShivamJha2436/kubehalo/internal/validation"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// ParseScalePolicy converts an unstructured CR into the typed API model.
func ParseScalePolicy(u *unstructured.Unstructured) (*kubehalov1.ScalePolicy, error) {
	var policy kubehalov1.ScalePolicy
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &policy); err != nil {
		return nil, fmt.Errorf("convert ScalePolicy: %w", err)
	}

	if err := ValidateScalePolicy(&policy); err != nil {
		return nil, err
	}

	if policy.Namespace == "" {
		policy.Namespace = policy.Spec.TargetRef.Namespace
	}

	return &policy, nil
}

// ValidateScalePolicy performs lightweight runtime validation for controller use.
func ValidateScalePolicy(policy *kubehalov1.ScalePolicy) error {
	return validation.ValidateScalePolicy(policy)
}

// CalculateReplicas decides the new replica count based on a metric value and threshold.
func CalculateReplicas(currentReplicas int32, metricValue, threshold float64, scaleUpStep, scaleDownStep int32) int32 {
	if metricValue > threshold {
		return currentReplicas + scaleUpStep
	}
	if metricValue < threshold {
		return int32(math.Max(1, float64(currentReplicas-scaleDownStep)))
	}
	return currentReplicas
}

// ClampReplicas enforces min/max limits on the desired replica count.
func ClampReplicas(replicas, minReplicas, maxReplicas int32) int32 {
	if replicas < minReplicas {
		return minReplicas
	}
	if replicas > maxReplicas {
		return maxReplicas
	}
	return replicas
}
