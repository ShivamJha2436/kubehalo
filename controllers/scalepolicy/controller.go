package scalepolicy

import (
	"log"
	"time"

	"kubehalo/internal/metrics"
	"kubehalo/internal/scaling"
	"kubehalo/controllers/scalepolicy/handler"
	"kubehalo/controllers/scalepolicy/lister"

	"k8s.io/client-go/kubernetes"
)

func StartController(clientset *kubernetes.Clientset, stopCh <-chan struct{}) error {
	log.Println("[CONTROLLER] Starting ScalePolicy controller...")

	// Start informers to watch ScalePolicy CRD
	go lister.StartScalePolicyInformer(clientset, stopCh)

	// Main controller loop
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				scalePolicies := lister.GetCurrentScalePolicies()
				for _, policy := range scalePolicies {
					metricsValue, err := metrics.FetchMetric(policy)
					if err != nil {
						log.Printf("[ERROR] Fetching metric: %v", err)
						continue
					}

					if err := scaling.ApplyScaling(policy, metricsValue, clientset); err != nil {
						log.Printf("[ERROR] Scaling failed: %v", err)
					}
				}
			case <-stopCh:
				log.Println("[CONTROLLER] Stopping ScalePolicy controller...")
				return
			}
		}
	}()

	return nil
}
