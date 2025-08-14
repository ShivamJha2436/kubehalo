package scaling

import (
	"context"
	"fmt"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ScalingEngine struct {
	clientset *kubernetes.Clientset
}

func NewScalingEngine(clientset *kubernetes.Clientset) *ScalingEngine {
	return &ScalingEngine{clientset: clientset}
}

func (s *ScalingEngine) ScaleDeployment(namespace, name string, replicas int32) error {
	deploy, err := s.clientset.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	deploy.Spec.Replicas = &replicas
	_, err = s.clientset.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("Deployment %s scaled to %d replicas\n", name, replicas)
	return nil
}
