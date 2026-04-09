package scalepolicy

import (
	"log"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

const defaultInformerResync = 30 * time.Second

// ScalePolicyGVR defines the Group/Version/Resource for the CRD.
var ScalePolicyGVR = schema.GroupVersionResource{
	Group:    "kubehalo.sh",
	Version:  "v1",
	Resource: "scalepolicies",
}

// EventHandler describes the callbacks used by the controller.
type EventHandler interface {
	OnAdd(obj interface{})
	OnUpdate(oldObj, newObj interface{})
	OnDelete(obj interface{})
}

// Controller manages the informer and event handlers.
type Controller struct {
	factory  dynamicinformer.DynamicSharedInformerFactory
	handler  EventHandler
	informer cache.SharedIndexInformer
}

// NewController creates a controller instance with the default handler.
func NewController(dynamicClient dynamic.Interface, kubeClient kubernetes.Interface, promClient PromClientInterface) *Controller {
	return NewControllerWithHandler(dynamicClient, NewHandler(kubeClient, promClient), defaultInformerResync)
}

// NewControllerWithHandler creates a controller with a custom event handler.
func NewControllerWithHandler(dynamicClient dynamic.Interface, handler EventHandler, resyncPeriod time.Duration) *Controller {
	factory := dynamicinformer.NewDynamicSharedInformerFactory(dynamicClient, resyncPeriod)

	return &Controller{
		factory:  factory,
		handler:  handler,
		informer: factory.ForResource(ScalePolicyGVR).Informer(),
	}
}

// Run starts the informer loop.
func (c *Controller) Run(stopCh <-chan struct{}) {
	c.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			log.Println("[controller] Add event received for ScalePolicy")
			c.handler.OnAdd(obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			log.Println("[controller] Update event received for ScalePolicy")
			c.handler.OnUpdate(oldObj, newObj)
		},
		DeleteFunc: func(obj interface{}) {
			log.Println("[controller] Delete event received for ScalePolicy")
			c.handler.OnDelete(obj)
		},
	})

	log.Println("[controller] Starting informer...")
	c.factory.Start(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		log.Println("[controller] Failed to sync caches.")
		return
	}

	log.Println("[controller] Cache synced. Watching ScalePolicy events...")
	<-stopCh
	log.Println("[controller] Stop signal received. Shutting down...")
}
