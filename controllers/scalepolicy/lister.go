package scalepolicy

import (
	"log"
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	v1 "kubehalo/api/v1"
)

var (
	scalePolicyCache = make(map[string]v1.ScalePolicy)
	mu               sync.RWMutex
)

func StartScalePolicyInformer(clientset *kubernetes.Clientset, stopCh <-chan struct{}) {
	// TODO: Use dynamic client / generated CRD client
	log.Println("[INFORMER] Starting ScalePolicy informer...")

	// Dummy: In real implementation, this will watch CRD changes
}

func GetCurrentScalePolicies() []v1.ScalePolicy {
	mu.RLock()
	defer mu.RUnlock()

	var list []v1.ScalePolicy
	for _, p := range scalePolicyCache {
		list = append(list, p)
	}
	return list
}
