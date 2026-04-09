package webhook

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	kubehalov1 "github.com/ShivamJha2436/kubehalo/api/kubehalo/v1"
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// validateScalePolicy checks ScalePolicy rules.
func validateScalePolicy(sp *kubehalov1.ScalePolicy) error {
	if sp.Spec.TargetRef.Kind == "" || sp.Spec.TargetRef.Name == "" || sp.Spec.TargetRef.Namespace == "" {
		return fmt.Errorf("targetRef must include kind, name, and namespace")
	}
	if sp.Spec.Metric.Query == "" {
		return fmt.Errorf("metric.query must not be empty")
	}
	if sp.Spec.Metric.Threshold <= 0 {
		return fmt.Errorf("metric.threshold must be greater than zero")
	}
	if sp.Spec.MinReplicas <= 0 {
		return fmt.Errorf("minReplicas must be greater than zero")
	}
	if sp.Spec.MaxReplicas < sp.Spec.MinReplicas {
		return fmt.Errorf("maxReplicas must be greater than or equal to minReplicas")
	}
	if sp.Spec.ScaleUp.Step <= 0 || sp.Spec.ScaleDown.Step <= 0 {
		return fmt.Errorf("scaleUp.step and scaleDown.step must be greater than zero")
	}
	return nil
}

// Serve handles admission requests.
func Serve(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var review admissionv1.AdmissionReview
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if review.Request == nil {
		http.Error(w, "admission review request is required", http.StatusBadRequest)
		return
	}

	var sp kubehalov1.ScalePolicy
	if err := json.Unmarshal(review.Request.Object.Raw, &sp); err != nil {
		writeResponse(w, review, false, fmt.Sprintf("cannot decode object: %v", err))
		return
	}

	if err := validateScalePolicy(&sp); err != nil {
		writeResponse(w, review, false, err.Error())
		return
	}

	writeResponse(w, review, true, "allowed")
}

func writeResponse(w http.ResponseWriter, review admissionv1.AdmissionReview, allowed bool, msg string) {
	w.Header().Set("Content-Type", "application/json")

	response := admissionv1.AdmissionReview{
		TypeMeta: review.TypeMeta,
		Response: &admissionv1.AdmissionResponse{
			UID:     review.Request.UID,
			Allowed: allowed,
			Result: &metav1.Status{
				Message: msg,
			},
		},
	}
	if err := json.NewEncoder(w).Encode(response); err != nil && err != io.EOF {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
