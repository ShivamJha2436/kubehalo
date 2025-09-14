package scalepolicy

import (
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes/fake"
)

// MockPromClient implements PromClientInterface
type MockPromClient struct{}

func (m *MockPromClient) QueryMetric(query string) (float64, error) {
	return 100.0, nil
}

func TestHandlerBehaviorSpec(t *testing.T) {
	kubeClient := fake.NewSimpleClientset() // implements kubernetes.Interface
	promClient := &MockPromClient{}

	handler := NewHandler(kubeClient, promClient)

	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": map[string]interface{}{
				"name":      "test-policy",
				"namespace": "default",
			},
			"spec": map[string]interface{}{
				"targetRef": map[string]interface{}{
					"name":      "test-deployment",
					"namespace": "default",
				},
				"metric": map[string]interface{}{
					"query":     "cpu_usage",
					"threshold": 50,
				},
				"behavior": map[string]interface{}{
					"stabilizationWindowSeconds": int64(2),
					"maxScaleUpRate":             int64(2),
					"maxScaleDownRate":           int64(1),
					"policy":                      "absolute",
				},
			},
		},
	}

	// First scale should execute
	handler.process(obj)

	// Second scale within stabilization window should be skipped
	handler.process(obj)

	time.Sleep(3 * time.Second)

	// After stabilization window, scaling should occur again
	handler.process(obj)
}
