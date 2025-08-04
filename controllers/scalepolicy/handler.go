package controller

import (
	"context"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	kubehalov1 "github.com/ShivamJha2436/kubehalo/api/v1"
)

func HandleAdd(obj interface{}) {
	policy, ok := obj.(*kubehalov1.ScalePolicy)
	if !ok {
		runtime.HandleError(nil)
		return
	}
	// your logic here
	_ = policy
}

func HandleUpdate(oldObj, newObj interface{}) {
	newPolicy, ok := newObj.(*kubehalov1.ScalePolicy)
	if !ok {
		runtime.HandleError(nil)
		return
	}
	// your logic here
	_ = newPolicy
}

func HandleDelete(obj interface{}) {
	// your logic here
}
