package validation

import (
	"testing"

	kubehalov1 "github.com/ShivamJha2436/kubehalo/api/kubehalo/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestValidateScalePolicyAcceptsZeroThreshold(t *testing.T) {
	policy := validPolicy()
	policy.Spec.Metric.Threshold = 0

	if err := ValidateScalePolicy(policy); err != nil {
		t.Fatalf("expected zero threshold to be allowed, got %v", err)
	}
}

func TestValidateScalePolicyRejectsNegativeThreshold(t *testing.T) {
	policy := validPolicy()
	policy.Spec.Metric.Threshold = -1

	if err := ValidateScalePolicy(policy); err == nil {
		t.Fatal("expected negative threshold to be rejected")
	}
}

func TestValidateScalePolicyRejectsOverlappingSchedules(t *testing.T) {
	policy := validPolicy()
	policy.Spec.Schedules = []kubehalov1.ScheduleSpec{
		{
			Name:      "morning",
			Days:      []string{"Mon", "Tue"},
			StartTime: "09:00",
			EndTime:   "12:00",
		},
		{
			Name:      "midday",
			Days:      []string{"Tue"},
			StartTime: "11:00",
			EndTime:   "13:00",
		},
	}

	if err := ValidateScalePolicy(policy); err == nil {
		t.Fatal("expected overlapping schedules to be rejected")
	}
}

func TestValidateScalePolicyRejectsInvalidScheduleReplicaBounds(t *testing.T) {
	policy := validPolicy()
	policy.Spec.Schedules = []kubehalov1.ScheduleSpec{
		{
			Days:        []string{"Mon"},
			StartTime:   "09:00",
			EndTime:     "10:00",
			MinReplicas: 5,
			MaxReplicas: 3,
		},
	}

	if err := ValidateScalePolicy(policy); err == nil {
		t.Fatal("expected invalid schedule replica bounds to be rejected")
	}
}

func validPolicy() *kubehalov1.ScalePolicy {
	return &kubehalov1.ScalePolicy{
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
			ScaleUp: kubehalov1.ScaleAction{Step: 1},
			ScaleDown: kubehalov1.ScaleAction{
				Step: 1,
			},
			MinReplicas: 1,
			MaxReplicas: 5,
		},
	}
}
