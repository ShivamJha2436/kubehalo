package scalepolicy

import (
	"encoding/json"
	"log"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Handler processes ScalePolicy events
type Handler struct{}

// NewHandler returns a new Handler instance
func NewHandler() *Handler {
	return &Handler{}
}

// OnAdd is called when a ScalePolicy is created
func (h *Handler) OnAdd(obj interface{}) {
	u, ok := obj.(*unstructured.Unstructured)
	if !ok {
		log.Println("[handler] OnAdd: unable to cast object to Unstructured")
		return
	}
	h.prettyLog("ADD", u)
}

// OnUpdate is called when a ScalePolicy is updated
func (h *Handler) OnUpdate(oldObj, newObj interface{}) {
	u, ok := newObj.(*unstructured.Unstructured)
	if !ok {
		log.Println("[handler] OnUpdate: unable to cast object to Unstructured")
		return
	}
	h.prettyLog("UPDATE", u)
}

// OnDelete is called when a ScalePolicy is deleted
func (h *Handler) OnDelete(obj interface{}) {
	u, ok := obj.(*unstructured.Unstructured)
	if !ok {
		log.Println("[handler] OnDelete: unable to cast object to Unstructured")
		return
	}
	h.prettyLog("DELETE", u)
}

// prettyLog prints ScalePolicy name + JSON spec
func (h *Handler) prettyLog(eventType string, u *unstructured.Unstructured) {
	// Extract name & namespace
	name := u.GetName()
	namespace := u.GetNamespace()

	// Get the spec field
	spec := u.Object["spec"]

	// Marshal for pretty printing
	specJSON, _ := json.MarshalIndent(spec, "", "  ")

	log.Printf("[handler] %s ScalePolicy: %s/%s\nSpec:\n%s\n",
		eventType, namespace, name, string(specJSON))
}
