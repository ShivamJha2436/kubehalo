package scalepolicy

import (
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/tools/cache"
)

type stubEventHandler struct{}

func (stubEventHandler) OnAdd(obj interface{})               {}
func (stubEventHandler) OnUpdate(oldObj, newObj interface{}) {}
func (stubEventHandler) OnDelete(obj interface{})            {}

func TestNewControllerWithHandler(t *testing.T) {
	dynamicClient := fake.NewSimpleDynamicClient(runtime.NewScheme())
	handler := stubEventHandler{}

	controller := NewControllerWithHandler(dynamicClient, handler, time.Second)

	if controller.handler != handler {
		t.Fatal("expected controller to keep the injected handler")
	}
	if controller.informer == nil {
		t.Fatal("expected controller informer to be initialized")
	}
}

func TestToUnstructuredSupportsDeletedFinalStateUnknown(t *testing.T) {
	obj := &unstructured.Unstructured{Object: map[string]interface{}{"kind": "ScalePolicy"}}

	got, ok := toUnstructured(cache.DeletedFinalStateUnknown{Obj: obj})
	if !ok {
		t.Fatal("expected tombstone object to be handled")
	}
	if got != obj {
		t.Fatal("expected original unstructured object back")
	}
}
