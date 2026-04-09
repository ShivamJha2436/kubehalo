package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/ShivamJha2436/kubehalo/internal/config"
	"github.com/ShivamJha2436/kubehalo/internal/metrics"
	"github.com/ShivamJha2436/kubehalo/internal/webhook"
)

func main() {
	cfg := config.LoadWebhookConfig()
	promClient, err := metrics.NewPrometheusClient(cfg.PrometheusAddress)
	if err != nil {
		log.Fatalf("failed to create Prometheus client: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/validate", webhook.NewHandler(webhook.NewValidator(promClient)))

	server := &http.Server{
		Addr:      cfg.Address,
		TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12},
		Handler:   mux,
	}

	log.Printf("[webhook] starting on %s/validate", cfg.Address)
	if err := server.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile); err != nil {
		log.Fatal(err)
	}
}
