package webhook

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ShivamJha2436/kubehalo/api/kubehalo/v1"
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	codecs = serializer.NewCodecFactory(runtime.NewScheme())
	deserializer = codecs.UniversalDeserializer()
)

// validateScalePolicy checks ScalePolicy rules
func validateScalePolicy(sp *v1.ScalePolicy) error {
	if sp.Spec.MetricQuery == "" {
		return fmt.Errorf("metricQuery must not be empty")
	}
	if sp.Spec.Threshold <= 0 {
		return fmt.Errorf("threshold must be > 0")
	}
	if sp.Spec.ScaleTargetRef.Name == "" || sp.Spec.ScaleTargetRef.Kind == "" {
		return fmt.Errorf("scaleTargetRef must include kind and name")
	}
	return nil
}

// Serve handles admission requests
func Serve(w http.ResponseWriter, r *http.Request) {
	var review admissionv1.AdmissionReview
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := review.Request
	var sp v1.ScalePolicy
	if _, _, err := deserializer.Decode(req.Object.Raw, nil, &sp); err != nil {
		writeResponse(w, review, false, fmt.Sprintf("cannot decode object: %v", err))
		return
	}

	// Run validation
	if err := validateScalePolicy(&sp); err != nil {
		writeResponse(w, review, false, err.Error())
		return
	}

	writeResponse(w, review, true, "allowed")
}

func writeResponse(w http.ResponseWriter, review admissionv1.AdmissionReview, allowed bool, msg string) {
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
	_ = json.NewEncoder(w).Encode(response)
}
