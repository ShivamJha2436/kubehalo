package scalepolicy

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	kubehalov1 "github.com/ShivamJha2436/kubehalo/api/kubehalo/v1"
	"github.com/ShivamJha2436/kubehalo/internal/scaling"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

// PromClientInterface defines the metric query behavior used by the handler.
type PromClientInterface interface {
	QueryMetric(query string) (float64, error)
}

// Handler processes ScalePolicy events.
type Handler struct {
	promClient     PromClientInterface
	engine         *scaling.ScalingEngine
	lastScaleTimes map[string]time.Time
	mu             sync.Mutex
	now            func() time.Time
}

// NewHandler returns a new handler instance.
func NewHandler(kubeClient kubernetes.Interface, promClient PromClientInterface) *Handler {
	return &Handler{
		promClient:     promClient,
		engine:         scaling.NewScalingEngine(kubeClient),
		lastScaleTimes: make(map[string]time.Time),
		now:            time.Now,
	}
}

// OnAdd is called when a ScalePolicy is created.
func (h *Handler) OnAdd(obj interface{}) {
	u, ok := toUnstructured(obj)
	if !ok {
		log.Println("[handler] add event ignored: expected Unstructured object")
		return
	}
	h.prettyLog("ADD", u)
	h.process(u)
}

// OnUpdate is called when a ScalePolicy is updated.
func (h *Handler) OnUpdate(oldObj, newObj interface{}) {
	u, ok := toUnstructured(newObj)
	if !ok {
		log.Println("[handler] update event ignored: expected Unstructured object")
		return
	}
	h.prettyLog("UPDATE", u)
	h.process(u)
}

// OnDelete is called when a ScalePolicy is deleted.
func (h *Handler) OnDelete(obj interface{}) {
	u, ok := toUnstructured(obj)
	if !ok {
		log.Println("[handler] delete event ignored: expected Unstructured object")
		return
	}
	h.prettyLog("DELETE", u)
}

// prettyLog prints the ScalePolicy name and spec payload for debugging.
func (h *Handler) prettyLog(eventType string, u *unstructured.Unstructured) {
	name := u.GetName()
	namespace := u.GetNamespace()
	spec := u.Object["spec"]
	specJSON, _ := json.MarshalIndent(spec, "", "  ")

	log.Printf("[handler] %s ScalePolicy: %s/%s\nSpec:\n%s\n",
		eventType, namespace, name, string(specJSON))
}

// process calculates a scaling decision and applies it to the target workload.
func (h *Handler) process(u *unstructured.Unstructured) {
	policy, err := ParseScalePolicy(u)
	if err != nil {
		log.Printf("[handler] invalid ScalePolicy %s/%s: %v", u.GetNamespace(), u.GetName(), err)
		return
	}

	ctx := context.Background()
	namespace := policy.Spec.TargetRef.Namespace
	deploymentName := policy.Spec.TargetRef.Name

	currentReplicas, err := h.engine.CurrentReplicas(ctx, namespace, deploymentName)
	if err != nil {
		log.Printf("[handler] unable to read current replicas for %s/%s: %v", namespace, deploymentName, err)
		return
	}

	metricValue, err := h.promClient.QueryMetric(policy.Spec.Metric.Query)
	if err != nil {
		log.Printf("[handler] unable to query metric %q: %v", policy.Spec.Metric.Query, err)
		return
	}

	desiredReplicas := CalculateReplicas(
		currentReplicas,
		metricValue,
		policy.Spec.Metric.Threshold,
		policy.Spec.ScaleUp.Step,
		policy.Spec.ScaleDown.Step,
	)
	desiredReplicas = ClampReplicas(desiredReplicas, policy.Spec.MinReplicas, policy.Spec.MaxReplicas)
	if desiredReplicas == currentReplicas {
		log.Printf("[handler] no scaling required for %s/%s (metric=%.2f threshold=%.2f replicas=%d)",
			namespace, deploymentName, metricValue, policy.Spec.Metric.Threshold, currentReplicas)
		return
	}

	desiredReplicas = h.applyBehavior(policy, currentReplicas, desiredReplicas)
	if desiredReplicas == currentReplicas {
		return
	}

	if err := h.engine.ScaleDeployment(ctx, namespace, deploymentName, desiredReplicas); err != nil {
		log.Printf("[handler] failed to scale %s/%s from %d to %d: %v",
			namespace, deploymentName, currentReplicas, desiredReplicas, err)
		return
	}

	log.Printf("[handler] scaled %s/%s from %d to %d replicas (metric=%.2f threshold=%.2f)",
		namespace, deploymentName, currentReplicas, desiredReplicas, metricValue, policy.Spec.Metric.Threshold)
}

func (h *Handler) applyBehavior(policy *kubehalov1.ScalePolicy, currentReplicas, desiredReplicas int32) int32 {
	behavior := policy.Spec.Behavior
	if behavior == nil {
		return desiredReplicas
	}

	key := policy.Spec.TargetRef.Namespace + "/" + policy.Spec.TargetRef.Name
	if behavior.StabilizationWindowSeconds != nil && *behavior.StabilizationWindowSeconds > 0 {
		if h.withinStabilizationWindow(key, *behavior.StabilizationWindowSeconds) {
			log.Printf("[behavior] skipping scale for %s due to stabilization window (%ds)",
				key, *behavior.StabilizationWindowSeconds)
			return currentReplicas
		}
	}

	if behavior.MaxScaleUpRate != nil && desiredReplicas > currentReplicas {
		desiredReplicas = min(desiredReplicas, currentReplicas+*behavior.MaxScaleUpRate)
	}

	if behavior.MaxScaleDownRate != nil && desiredReplicas < currentReplicas {
		desiredReplicas = max(desiredReplicas, currentReplicas-*behavior.MaxScaleDownRate)
	}

	return desiredReplicas
}

func (h *Handler) withinStabilizationWindow(key string, windowSeconds int32) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	lastScaledAt, ok := h.lastScaleTimes[key]
	now := h.now()
	if ok && now.Sub(lastScaledAt) < time.Duration(windowSeconds)*time.Second {
		return true
	}

	h.lastScaleTimes[key] = now
	return false
}

func toUnstructured(obj interface{}) (*unstructured.Unstructured, bool) {
	switch typed := obj.(type) {
	case *unstructured.Unstructured:
		return typed, true
	case cache.DeletedFinalStateUnknown:
		u, ok := typed.Obj.(*unstructured.Unstructured)
		return u, ok
	case *cache.DeletedFinalStateUnknown:
		u, ok := typed.Obj.(*unstructured.Unstructured)
		return u, ok
	default:
		return nil, false
	}
}

func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
