package scalepolicy

import (
	"log"
	"time"

	"github.com/ShivamJha2436/kubehalo/internal/metrics"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

// ScalePolicyGVR defines the Group/Version/Resource for the CRD
var ScalePolicyGVR = schema.GroupVersionResource{
	Group:    "kubehalo.sh",
	Version:  "v1",
	Resource: "scalepolicies",
}

// Controller manages the informer and event handlers
type Controller struct {
	factory dynamicinformer.DynamicSharedInformerFactory
	handler *Handler
	informer cache.SharedIndexInformer
}

// NewController creates a Controller instance
func NewController(dynamicClient dynamic.Interface, kubeClient *kubernetes.Clientset, promClient *metrics.PrometheusClient) *Controller {
	factory := dynamicinformer.NewDynamicSharedInformerFactory(dynamicClient, 30*time.Second)
	//Create the controller instance
	ctrl := &Controller{
		factory: factory,
		handler: NewHandler(kubeClient,promClient),
		informer: factory.ForResource(ScalePolicyGVR).Informer(),
	}
	return ctrl
}

// Run starts the informer loop
func (c *Controller) Run(stopCh <-chan struct{}) {
	// Register event handlers
	c.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
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
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		log.Println("[controller] Failed to sync caches.")
		return
	}

	log.Println("[controller] Cache synced. Watching ScalePolicy events...")
	<-stopCh
	log.Println("[controller] Stop signal received. Shutting down...")
}
