package scalepolicy

import (
	"testing"

	kubehalov1 "github.com/ShivamJha2436/kubehalo/api/kubehalo/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestParseScalePolicy(t *testing.T) {
	policy := &kubehalov1.ScalePolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "demo-policy",
			Namespace: "default",
		},
		Spec: kubehalov1.ScalePolicySpec{
			TargetRef: kubehalov1.TargetRefSpec{
				Kind:      "Deployment",
				Name:      "demo",
				Namespace: "default",
			},
			Metric: kubehalov1.MetricSpec{
				Name:      "cpu",
				Query:     "demo_metric",
				Threshold: 0.8,
			},
			ScaleUp: kubehalov1.ScaleAction{Step: 2},
			ScaleDown: kubehalov1.ScaleAction{
				Step: 1,
			},
			MinReplicas: 1,
			MaxReplicas: 5,
		},
	}

	object, err := runtime.DefaultUnstructuredConverter.ToUnstructured(policy)
	if err != nil {
		t.Fatalf("convert policy: %v", err)
	}

	parsed, err := ParseScalePolicy(&unstructured.Unstructured{Object: object})
	if err != nil {
		t.Fatalf("parse policy: %v", err)
	}

	if parsed.Spec.Metric.Threshold != 0.8 {
		t.Fatalf("expected threshold 0.8, got %v", parsed.Spec.Metric.Threshold)
	}
}

func TestValidateScalePolicy(t *testing.T) {
	policy := &kubehalov1.ScalePolicy{
		Spec: kubehalov1.ScalePolicySpec{
			TargetRef: kubehalov1.TargetRefSpec{
				Kind:      "Deployment",
				Name:      "demo",
				Namespace: "default",
			},
			Metric: kubehalov1.MetricSpec{
				Query:     "demo_metric",
				Threshold: 0.8,
			},
			ScaleUp: kubehalov1.ScaleAction{Step: 2},
			ScaleDown: kubehalov1.ScaleAction{
				Step: 1,
			},
			MinReplicas: 5,
			MaxReplicas: 2,
		},
	}

	if err := ValidateScalePolicy(policy); err == nil {
		t.Fatal("expected validation error for maxReplicas < minReplicas")
	}
}

func TestCalculateReplicasAndClamp(t *testing.T) {
	if got := CalculateReplicas(2, 100, 80, 2, 1); got != 4 {
		t.Fatalf("expected scale up to 4, got %d", got)
	}

	if got := CalculateReplicas(2, 40, 80, 2, 1); got != 1 {
		t.Fatalf("expected scale down to 1, got %d", got)
	}

	if got := ClampReplicas(9, 1, 5); got != 5 {
		t.Fatalf("expected clamp result 5, got %d", got)
	}
}
