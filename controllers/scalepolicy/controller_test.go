package scalepolicy

import (
	"sync"
	"testing"
	"time"

	"github.com/ShivamJha2436/kubehalo/internal/metrics"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes/fake"
)

// mockHandler embeds Handler but adds counters for calls
type mockHandler struct {
	addCount    int
	updateCount int
	deleteCount int
	mu          sync.Mutex
}

func (m *mockHandler) OnAdd(obj interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.addCount++
}

func (m *mockHandler) OnUpdate(oldObj, newObj interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.updateCount++
}

func (m *mockHandler) OnDelete(obj interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.deleteCount++
}

func TestController_Run(t *testing.T) {
	// Fake k8s + dynamic clients
	kubeClient := fake.NewSimpleClientset()
	dynamicClient := fake.NewSimpleDynamicClient(kubeClient.Discovery().RESTMapper())

	// Fake Prometheus client
	promClient := &metrics.PrometheusClient{}

	// Create controller
	ctrl := NewController(dynamicClient, kubeClient, promClient)

	// Replace real handler with mock handler
	mock := &mockHandler{}
	ctrl.handler = &Handler{
		OnAdd:    mock.OnAdd,
		OnUpdate: mock.OnUpdate,
		OnDelete: mock.OnDelete,
	}

	stopCh := make(chan struct{})
	defer close(stopCh)

	// Run controller in goroutine
	go ctrl.Run(stopCh)

	// Wait briefly to allow informer to start
	time.Sleep(200 * time.Millisecond)

	// Create a fake ScalePolicy with BehaviorSpec
	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "kubehalo.sh/v1",
			"kind":       "ScalePolicy",
			"metadata": map[string]interface{}{
				"name": "test-policy",
			},
			"spec": map[string]interface{}{
				"behavior": map[string]interface{}{
					"stabilizationWindowSeconds": 60,
					"maxScaleUpRate":             2,
					"maxScaleDownRate":           1,
					"policy":                     "absolute",
				},
			},
		},
	}

	// Simulate Add event
	ctrl.informer.GetStore().Add(obj)

	// Simulate Update event
	updated := obj.DeepCopy()
	updated.Object["spec"].(map[string]interface{})["behavior"].(map[string]interface{})["maxScaleUpRate"] = 5
	ctrl.informer.GetStore().Update(updated)

	// Simulate Delete event
	ctrl.informer.GetStore().Delete(updated)

	// Allow handlers to be called
	time.Sleep(200 * time.Millisecond)

	mock.mu.Lock()
	defer mock.mu.Unlock()

	if mock.addCount != 1 {
		t.Errorf("expected addCount=1, got %d", mock.addCount)
	}
	if mock.updateCount != 1 {
		t.Errorf("expected updateCount=1, got %d", mock.updateCount)
	}
	if mock.deleteCount != 1 {
		t.Errorf("expected deleteCount=1, got %d", mock.deleteCount)
	}
}

