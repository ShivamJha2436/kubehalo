package config

import (
	"os"
)

const (
	defaultAPIAddress        = ":8080"
	defaultWebhookAddress    = ":8443"
	defaultWebhookCertFile   = "/tls/tls.crt"
	defaultWebhookKeyFile    = "/tls/tls.key"
	defaultPrometheusAddress = "http://localhost:9090"
)

// APIConfig defines the HTTP API runtime configuration.
type APIConfig struct {
	Address string
}

// ControllerConfig defines the controller runtime configuration.
type ControllerConfig struct {
	PrometheusAddress string
}

// WebhookConfig defines the admission webhook runtime configuration.
type WebhookConfig struct {
	Address  string
	CertFile string
	KeyFile  string
}

// LoadAPIConfig reads the API configuration from the environment.
func LoadAPIConfig() APIConfig {
	return APIConfig{
		Address: getEnv("KUBEHALO_API_ADDR", defaultAPIAddress),
	}
}

// LoadControllerConfig reads the controller configuration from the environment.
func LoadControllerConfig() ControllerConfig {
	return ControllerConfig{
		PrometheusAddress: getEnv("KUBEHALO_PROMETHEUS_ADDR", defaultPrometheusAddress),
	}
}

// LoadWebhookConfig reads the webhook configuration from the environment.
func LoadWebhookConfig() WebhookConfig {
	return WebhookConfig{
		Address:  getEnv("KUBEHALO_WEBHOOK_ADDR", defaultWebhookAddress),
		CertFile: getEnv("KUBEHALO_WEBHOOK_CERT_FILE", defaultWebhookCertFile),
		KeyFile:  getEnv("KUBEHALO_WEBHOOK_KEY_FILE", defaultWebhookKeyFile),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
