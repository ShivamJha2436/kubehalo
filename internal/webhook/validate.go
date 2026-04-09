package webhook

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	kubehalov1 "github.com/ShivamJha2436/kubehalo/api/kubehalo/v1"
	"github.com/ShivamJha2436/kubehalo/internal/validation"
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// QueryValidator performs Prometheus query validation for admission checks.
type QueryValidator interface {
	ValidateQuery(query string) error
}

// Validator validates ScalePolicy objects before they are admitted.
type Validator struct {
	queryValidator QueryValidator
}

// NewValidator builds a Validator with an optional Prometheus query validator.
func NewValidator(queryValidator QueryValidator) *Validator {
	return &Validator{queryValidator: queryValidator}
}

// ValidateScalePolicy checks ScalePolicy rules.
func (v *Validator) ValidateScalePolicy(sp *kubehalov1.ScalePolicy) error {
	if err := validation.ValidateScalePolicy(sp); err != nil {
		return err
	}

	if v.queryValidator == nil {
		return nil
	}

	if err := v.queryValidator.ValidateQuery(sp.Spec.Metric.Query); err != nil {
		return fmt.Errorf("spec.metric.query failed Prometheus validation: %w", err)
	}

	return nil
}

// NewHandler returns an admission webhook handler.
func NewHandler(validator *Validator) http.HandlerFunc {
	if validator == nil {
		validator = NewValidator(nil)
	}

	return func(w http.ResponseWriter, r *http.Request) {
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

		if err := validator.ValidateScalePolicy(&sp); err != nil {
			writeResponse(w, review, false, err.Error())
			return
		}

		writeResponse(w, review, true, "allowed")
	}
}

// Serve handles admission requests with structural validation only.
func Serve(w http.ResponseWriter, r *http.Request) {
	NewHandler(nil)(w, r)
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
