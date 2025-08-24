package scalepolicy

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
)


// Lister wraps the dynamic client for ScalePolicy listing
type Lister struct {
	dyn dynamic.Interface
}

// NewLister builds a new Lister
func NewLister(dyn dynamic.Interface) *Lister {
	return &Lister{dyn: dyn}
}

// ListScalePolicies returns all ScalePolicy objects across namespaces
func (l *Lister) ListScalePolicies(ctx context.Context) ([]*unstructured.Unstructured, error) {
	list, err := l.dyn.Resource(ScalePolicyGVR).Namespace(metav1.NamespaceAll).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list ScalePolicies: %w", err)
	}

	var items []*unstructured.Unstructured
	for i := range list.Items {
		items = append(items, &list.Items[i])
	}
	return items, nil
}
