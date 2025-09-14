package scaling

import (
	"context"
	"fmt"

	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ScalingEngine handles scaling operations for Deployments
type ScalingEngine struct {
	client kubernetes.Interface
}

// NewScalingEngine accepts kubernetes.Interface instead of *Clientset
func NewScalingEngine(client kubernetes.Interface) *ScalingEngine {
	return &ScalingEngine{client: client}
}

// ScaleDeployment scales a Deployment to the desired number of replicas
func (s *ScalingEngine) ScaleDeployment(namespace, name string, replicas int32) error {
	deploy, err := s.client.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	deploy.Spec.Replicas = &replicas
	_, err = s.client.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("Deployment %s scaled to %d replicas\n", name, replicas)
	return nil
}
