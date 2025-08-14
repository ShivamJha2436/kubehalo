package scalepolicy

import (
	"encoding/json"
	"fmt"
	"log"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes"
	"github.com/ShivamJha2436/kubehalo/internal/metrics"
	"github.com/ShivamJha2436/kubehalo/internal/scaling"
)

// Handler processes ScalePolicy events
type Handler struct{
	promClient *metrics.PrometheusClient
	engine	   *scaling.ScalingEngine
	KubeClient *kubernetes.Clientset
}

// NewHandler returns a new Handler instance
func NewHandler(kubeClient *kubernetes.Clientset, promClient *metrics.PrometheusClient) *Handler {
	return &Handler{
		promClient: promClient,
		engine: scaling.NewScalingEngine(kubeClient),
		KubeClient: kubeClient,
	}
}

// OnAdd is called when a ScalePolicy is created
func (h *Handler) OnAdd(obj interface{}) {
	u, ok := obj.(*unstructured.Unstructured)
	if !ok {
		log.Println("[handler] OnAdd: unable to cast object to Unstructured")
		return
	}
	h.prettyLog("ADD", u)
	h.process(u)
}

// OnUpdate is called when a ScalePolicy is updated
func (h *Handler) OnUpdate(oldObj, newObj interface{}) {
	u, ok := newObj.(*unstructured.Unstructured)
	if !ok {
		log.Println("[handler] OnUpdate: unable to cast object to Unstructured")
		return
	}
	h.prettyLog("UPDATE", u)
	h.process(u)
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
	name := u.GetName()
	namespace := u.GetNamespace()
	spec := u.Object["spec"]
	specJSON, _ := json.MarshalIndent(spec, "", "  ")

	log.Printf("[handler] %s ScalePolicy: %s/%s\nSpec:\n%s\n",
		eventType, namespace, name, string(specJSON))
}

// process calculates scaling decision (Phase 3: log only)
func (h *Handler) process(u *unstructured.Unstructured) {
	deploymentName, err := GetNestedString(u, "spec", "targetRef", "name")
	if err != nil {
		log.Println("Error reading deployment name:", err)
		return
	}
	namespace, err := GetNestedString(u, "spec", "targetRef", "namespace")
	if err != nil {
		log.Println("Error reading namespace:", err)
		return
	}

	metricQuery, err := GetNestedString(u, "spec", "metric", "query")
	if err != nil {
		log.Println("Error reading metric query:", err)
		return
	}

	thresholdStr, err := GetNestedString(u, "spec", "metric", "threshold")
	if err != nil {
		log.Println("Error reading threshold:", err)
		return
	}

	// Convert threshold to float64
	var threshold float64
	fmt.Sscanf(thresholdStr, "%f", &threshold)

	// Fetch metric from Prometheus
	metricValue, err := h.promClient.QueryMetric(metricQuery)
	if err != nil {
		log.Println("Error querying Prometheus:", err)
		return
	}

	// assume current replicas = 1 (will replace later with actual fetch)
	currentReplicas := int32(1)
	// Calculate desired replicas
	newReplicas := CalculateReplicas(currentReplicas, metricValue, threshold, 1, 1)
	// Log scaling decision (Phase 3)
	log.Printf("[Phase3] Would scale %s/%s from %d -> %d replicas (metric=%.2f, threshold=%.2f)\n",
		namespace, deploymentName, currentReplicas, newReplicas, metricValue, threshold)
}