package scaling

import (
	"context"
	"fmt"
	"log"

	"kubehalo/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ApplyScaling(policy v1.ScalePolicy, metric float64, clientset *kubernetes.Clientset) error {
	namespace := policy.Spec.TargetNamespace
	name := policy.Spec.TargetDeployment

	deploy, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	currentReplicas := *deploy.Spec.Replicas
	var desired int32

	if metric > policy.Spec.Thresholds.ScaleUp {
		desired = currentReplicas + 1
	} else if metric < policy.Spec.Thresholds.ScaleDown && currentReplicas > 1 {
		desired = currentReplicas - 1
	} else {
		log.Printf("[SCALING] No change for %s/%s", namespace, name)
		return nil
	}

	deploy.Spec.Replicas = &desired
	_, err = clientset.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	log.Printf("[SCALING] Scaled %s/%s to %d replicas", namespace, name, desired)
	return nil
}
