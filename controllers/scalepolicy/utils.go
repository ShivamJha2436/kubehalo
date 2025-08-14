package scalepolicy

import (
	"fmt"
	"math"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// GetNestedString fetches a nested string from an Unstructured object, e.g. ("spec","targetRef","name")
func GetNestedString(u *unstructured.Unstructured, fields ...string) (string, error) {
	val, found, err := unstructured.NestedString(u.Object, fields...)
	if err != nil {
		return "", err
	}
	if !found {
		return "", fmt.Errorf("field %v not found", fields)
	}
	return val, nil
}

// CalculateReplicas decides new replica count based on metric value and threshold
func CalculateReplicas(currentReplicas int32, metricValue, threshold float64, scaleUpStep, scaleDownStep int32) int32 {
	if metricValue > threshold {
		return currentReplicas + scaleUpStep
	} else if metricValue < threshold {
		// minimum 1 replica
		return int32(math.Max(1, float64(currentReplicas-scaleDownStep)))
	}
	return currentReplicas
}
