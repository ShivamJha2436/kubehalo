package scaling

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubefake "k8s.io/client-go/kubernetes/fake"
)

func TestCurrentReplicas(t *testing.T) {
	engine := NewScalingEngine(kubefake.NewSimpleClientset(newDeployment(3)))

	replicas, err := engine.CurrentReplicas(t.Context(), "default", "demo")
	if err != nil {
		t.Fatalf("current replicas: %v", err)
	}
	if replicas != 3 {
		t.Fatalf("expected 3 replicas, got %d", replicas)
	}
}

func TestScaleDeployment(t *testing.T) {
	client := kubefake.NewSimpleClientset(newDeployment(2))
	engine := NewScalingEngine(client)

	if err := engine.ScaleDeployment(t.Context(), "default", "demo", 5); err != nil {
		t.Fatalf("scale deployment: %v", err)
	}

	deployment, err := client.AppsV1().Deployments("default").Get(t.Context(), "demo", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("get deployment: %v", err)
	}
	if got := *deployment.Spec.Replicas; got != 5 {
		t.Fatalf("expected 5 replicas, got %d", got)
	}
}

func newDeployment(replicas int32) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "demo",
			Namespace: "default",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
		},
	}
}
