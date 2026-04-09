package scaling

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ScalingEngine handles scaling operations for Deployments
type ScalingEngine struct {
	client kubernetes.Interface
}

// NewScalingEngine accepts kubernetes.Interface instead of *Clientset
func NewScalingEngine(client kubernetes.Interface) *ScalingEngine {
	return &ScalingEngine{client: client}
}

// CurrentReplicas returns the current replica count for a Deployment.
func (s *ScalingEngine) CurrentReplicas(ctx context.Context, namespace, name string) (int32, error) {
	deploy, err := s.client.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return 0, err
	}
	if deploy.Spec.Replicas == nil {
		return 0, nil
	}
	return *deploy.Spec.Replicas, nil
}

// ScaleDeployment scales a Deployment to the desired number of replicas.
func (s *ScalingEngine) ScaleDeployment(ctx context.Context, namespace, name string, replicas int32) error {
	deploy, err := s.client.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	deploy.Spec.Replicas = &replicas
	_, err = s.client.AppsV1().Deployments(namespace).Update(ctx, deploy, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}
