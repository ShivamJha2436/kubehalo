package scalepolicy

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ShivamJha2436/kubehalo/internal/scaling"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes"
)

// PromClientInterface defines the methods used by Handler
type PromClientInterface interface {
	QueryMetric(query string) (float64, error)
}

// Handler processes ScalePolicy events
type Handler struct {
	promClient     PromClientInterface
	engine         *scaling.ScalingEngine
	KubeClient     kubernetes.Interface
	lastScaleTimes map[string]time.Time // key = namespace/name
}

// NewHandler returns a new Handler instance
func NewHandler(kubeClient kubernetes.Interface, promClient PromClientInterface) *Handler {
	return &Handler{
		promClient:     promClient,
		engine:         scaling.NewScalingEngine(kubeClient),
		KubeClient:     kubeClient,
		lastScaleTimes: make(map[string]time.Time),
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

// process calculates scaling decision and applies BehaviorSpec rules
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

	// assume current replicas = 1 (TODO: fetch actual replicas from Deployment)
	currentReplicas := int32(1)
	desiredReplicas := CalculateReplicas(currentReplicas, metricValue, threshold, 1, 1)

	// --- Apply BehaviorSpec ---
	behavior := u.Object["spec"].(map[string]interface{})["behavior"]
	if behavior != nil {
		if bMap, ok := behavior.(map[string]interface{}); ok {
			key := namespace + "/" + deploymentName

			// Stabilization window
			if secs, ok := bMap["stabilizationWindowSeconds"].(int64); ok && secs > 0 {
				last, exists := h.lastScaleTimes[key]
				if exists && time.Since(last) < time.Duration(secs)*time.Second {
					log.Printf("[behavior] Skipping scale for %s due to stabilization window (%ds)",
						key, secs)
					return
				}
				h.lastScaleTimes[key] = time.Now()
			}

			// Max scale-up rate
			if maxUp, ok := bMap["maxScaleUpRate"].(int64); ok && desiredReplicas-currentReplicas > int32(maxUp) {
				log.Printf("[behavior] Capping scale-up: requested %d, capped to %d",
					desiredReplicas, currentReplicas+int32(maxUp))
				desiredReplicas = currentReplicas + int32(maxUp)
			}

			// Max scale-down rate
			if maxDown, ok := bMap["maxScaleDownRate"].(int64); ok && currentReplicas-desiredReplicas > int32(maxDown) {
				log.Printf("[behavior] Capping scale-down: requested %d, capped to %d",
					desiredReplicas, currentReplicas-int32(maxDown))
				desiredReplicas = currentReplicas - int32(maxDown)
			}

			// Policy (absolute vs percent) — log only for now
			if policy, ok := bMap["policy"].(string); ok {
				log.Printf("[behavior] Policy set to: %s", policy)
			}
		}
	}

	log.Printf("[Phase3] Would scale %s/%s from %d -> %d replicas (metric=%.2f, threshold=%.2f)\n",
		namespace, deploymentName, currentReplicas, desiredReplicas, metricValue, threshold)
}
