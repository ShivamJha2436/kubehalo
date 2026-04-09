package webhook

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	kubehalov1 "github.com/ShivamJha2436/kubehalo/api/kubehalo/v1"
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

type stubQueryValidator struct {
	err error
}

func (s stubQueryValidator) ValidateQuery(query string) error {
	return s.err
}

func TestValidateScalePolicy(t *testing.T) {
	validator := NewValidator(nil)
	policy := validScalePolicy()
	if err := validator.ValidateScalePolicy(policy); err != nil {
		t.Fatalf("expected valid policy, got error: %v", err)
	}

	policy.Spec.Metric.Query = ""
	if err := validator.ValidateScalePolicy(policy); err == nil {
		t.Fatal("expected validation error for empty metric query")
	}
}

func TestServeAllowsValidPolicy(t *testing.T) {
	response := serveAdmissionReview(t, NewValidator(stubQueryValidator{}), validScalePolicy())
	if response.Response == nil || !response.Response.Allowed {
		t.Fatalf("expected admission request to be allowed, got %+v", response.Response)
	}
}

func TestServeRejectsInvalidPolicy(t *testing.T) {
	policy := validScalePolicy()
	policy.Spec.MaxReplicas = 0

	response := serveAdmissionReview(t, NewValidator(stubQueryValidator{}), policy)
	if response.Response == nil || response.Response.Allowed {
		t.Fatalf("expected admission request to be denied, got %+v", response.Response)
	}
}

func TestServeRejectsInvalidPrometheusQuery(t *testing.T) {
	response := serveAdmissionReview(t, NewValidator(stubQueryValidator{err: errors.New("parse error")}), validScalePolicy())
	if response.Response == nil || response.Response.Allowed {
		t.Fatalf("expected invalid query to be rejected, got %+v", response.Response)
	}
}

func TestServeRejectsOverlappingSchedules(t *testing.T) {
	policy := validScalePolicy()
	policy.Spec.Schedules = []kubehalov1.ScheduleSpec{
		{
			Name:      "morning",
			Days:      []string{"Mon", "Tue"},
			StartTime: "09:00",
			EndTime:   "12:00",
		},
		{
			Name:      "overlap",
			Days:      []string{"Tue"},
			StartTime: "11:00",
			EndTime:   "13:00",
		},
	}

	response := serveAdmissionReview(t, NewValidator(stubQueryValidator{}), policy)
	if response.Response == nil || response.Response.Allowed {
		t.Fatalf("expected overlapping schedules to be rejected, got %+v", response.Response)
	}
}

func serveAdmissionReview(t *testing.T, validator *Validator, policy *kubehalov1.ScalePolicy) admissionv1.AdmissionReview {
	t.Helper()

	rawPolicy, err := json.Marshal(policy)
	if err != nil {
		t.Fatalf("marshal policy: %v", err)
	}

	review := admissionv1.AdmissionReview{
		Request: &admissionv1.AdmissionRequest{
			UID:    types.UID("test-review"),
			Object: runtime.RawExtension{Raw: rawPolicy},
		},
	}

	body, err := json.Marshal(review)
	if err != nil {
		t.Fatalf("marshal review: %v", err)
	}

	request := httptest.NewRequest(http.MethodPost, "/validate", bytes.NewReader(body))
	recorder := httptest.NewRecorder()

	NewHandler(validator)(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	var response admissionv1.AdmissionReview
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	return response
}

func validScalePolicy() *kubehalov1.ScalePolicy {
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
				Threshold: 0,
			},
			ScaleUp: kubehalov1.ScaleAction{Step: 2},
			ScaleDown: kubehalov1.ScaleAction{
				Step: 1,
			},
			MinReplicas: 1,
			MaxReplicas: 5,
		},
	}
}
