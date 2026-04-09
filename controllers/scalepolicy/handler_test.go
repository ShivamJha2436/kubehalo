package scalepolicy

import (
	"testing"
	"time"

	kubehalov1 "github.com/ShivamJha2436/kubehalo/api/kubehalo/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	kubefake "k8s.io/client-go/kubernetes/fake"
)

type fakePromClient struct {
	value float64
	err   error
}

func (f fakePromClient) QueryMetric(query string) (float64, error) {
	return f.value, f.err
}

func TestHandlerProcessScalesDeployment(t *testing.T) {
	kubeClient := kubefake.NewSimpleClientset(newDeployment("default", "demo", 2))
	handler := NewHandler(kubeClient, fakePromClient{value: 95})

	handler.process(mustUnstructuredPolicy(t, newScalePolicy(80, nil)))

	deployment, err := kubeClient.AppsV1().Deployments("default").Get(t.Context(), "demo", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("get deployment: %v", err)
	}

	if got := *deployment.Spec.Replicas; got != 5 {
		t.Fatalf("expected deployment replicas to be 5, got %d", got)
	}
}

func TestHandlerProcessHonorsBehaviorRules(t *testing.T) {
	stabilizationWindow := int32(60)
	maxScaleUpRate := int32(1)
	behavior := &kubehalov1.BehaviorSpec{
		StabilizationWindowSeconds: &stabilizationWindow,
		MaxScaleUpRate:             &maxScaleUpRate,
	}

	kubeClient := kubefake.NewSimpleClientset(newDeployment("default", "demo", 2))
	handler := NewHandler(kubeClient, fakePromClient{value: 95})

	now := time.Date(2026, time.April, 9, 10, 0, 0, 0, time.UTC)
	handler.now = func() time.Time { return now }

	policy := mustUnstructuredPolicy(t, newScalePolicy(80, behavior))

	handler.process(policy)
	assertDeploymentReplicas(t, kubeClient, 3)

	now = now.Add(30 * time.Second)
	handler.process(policy)
	assertDeploymentReplicas(t, kubeClient, 3)

	now = now.Add(31 * time.Second)
	handler.process(policy)
	assertDeploymentReplicas(t, kubeClient, 4)
}

func TestHandlerProcessRejectsInvalidPolicy(t *testing.T) {
	kubeClient := kubefake.NewSimpleClientset(newDeployment("default", "demo", 2))
	handler := NewHandler(kubeClient, fakePromClient{value: 95})

	policy := newScalePolicy(80, nil)
	policy.Spec.Metric.Query = ""

	handler.process(mustUnstructuredPolicy(t, policy))
	assertDeploymentReplicas(t, kubeClient, 2)
}

func newScalePolicy(threshold float64, behavior *kubehalov1.BehaviorSpec) *kubehalov1.ScalePolicy {
	return &kubehalov1.ScalePolicy{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "kubehalo.sh/v1",
			Kind:       "ScalePolicy",
		},
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
				Threshold: threshold,
			},
			ScaleUp: kubehalov1.ScaleAction{
				Step: 3,
			},
			ScaleDown: kubehalov1.ScaleAction{
				Step: 1,
			},
			MinReplicas: 1,
			MaxReplicas: 10,
			Behavior:    behavior,
		},
	}
}

func newDeployment(namespace, name string, replicas int32) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
		},
	}
}

func mustUnstructuredPolicy(t *testing.T, policy *kubehalov1.ScalePolicy) *unstructured.Unstructured {
	t.Helper()

	obj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(policy)
	if err != nil {
		t.Fatalf("convert policy: %v", err)
	}

	return &unstructured.Unstructured{Object: obj}
}

func assertDeploymentReplicas(t *testing.T, kubeClient *kubefake.Clientset, want int32) {
	t.Helper()

	deployment, err := kubeClient.AppsV1().Deployments("default").Get(t.Context(), "demo", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("get deployment: %v", err)
	}

	if got := *deployment.Spec.Replicas; got != want {
		t.Fatalf("expected deployment replicas to be %d, got %d", want, got)
	}
}
