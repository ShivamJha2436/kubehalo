package scalepolicy

import (
	"log"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
)

// ScalePolicyGVR defines the Group/Version/Resource for the CRD
var ScalePolicyGVR = schema.GroupVersionResource{
	Group:    "kubehalo.sh",
	Version:  "v1",
	Resource: "scalepolicies", // plural from CRD YAML
}

// Controller manages the informer and event handlers
type Controller struct {
	factory dynamicinformer.DynamicSharedInformerFactory
	handler *Handler
}

// NewController creates a Controller instance
func NewController(dynamicClient dynamic.Interface) *Controller {
	// Resync every 30s — watches all namespaces ("" means all)
	factory := dynamicinformer.NewDynamicSharedInformerFactory(dynamicClient, 30*time.Second)
	return &Controller{
		factory: factory,
		handler: NewHandler(),
	}
}

// Run starts the informer loop
func (c *Controller) Run(stopCh <-chan struct{}) {
	// Create informer for ScalePolicy CRD
	informer := c.factory.ForResource(ScalePolicyGVR).Informer()

	// Register event handlers
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			c.handler.OnAdd(obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			c.handler.OnUpdate(oldObj, newObj)
		},
		DeleteFunc: func(obj interface{}) {
			c.handler.OnDelete(obj)
		},
	})

	// Start informer factory
	log.Println("[controller] Starting informer...")
	c.factory.Start(stopCh)

	// Wait for caches to sync
	if !cache.WaitForCacheSync(stopCh, informer.HasSynced) {
		log.Println("[controller] Failed to sync caches.")
		return
	}

	log.Println("[controller] Cache synced. Watching ScalePolicy events...")
	<-stopCh
	log.Println("[controller] Stop signal received. Shutting down...")
}
