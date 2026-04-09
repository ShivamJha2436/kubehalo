package scalepolicy

import (
	"fmt"
	"math"

	kubehalov1 "github.com/ShivamJha2436/kubehalo/api/kubehalo/v1"
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
	switch {
	case policy.Spec.TargetRef.Kind == "":
		return fmt.Errorf("spec.targetRef.kind must not be empty")
	case policy.Spec.TargetRef.Name == "":
		return fmt.Errorf("spec.targetRef.name must not be empty")
	case policy.Spec.TargetRef.Namespace == "":
		return fmt.Errorf("spec.targetRef.namespace must not be empty")
	case policy.Spec.Metric.Query == "":
		return fmt.Errorf("spec.metric.query must not be empty")
	case policy.Spec.Metric.Threshold <= 0:
		return fmt.Errorf("spec.metric.threshold must be greater than zero")
	case policy.Spec.MinReplicas <= 0:
		return fmt.Errorf("spec.minReplicas must be greater than zero")
	case policy.Spec.MaxReplicas < policy.Spec.MinReplicas:
		return fmt.Errorf("spec.maxReplicas must be greater than or equal to spec.minReplicas")
	case policy.Spec.ScaleUp.Step <= 0:
		return fmt.Errorf("spec.scaleUp.step must be greater than zero")
	case policy.Spec.ScaleDown.Step <= 0:
		return fmt.Errorf("spec.scaleDown.step must be greater than zero")
	}

	return nil
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
