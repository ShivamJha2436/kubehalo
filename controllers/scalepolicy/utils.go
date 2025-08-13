package scalepolicy

import (
	"fmt"

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
