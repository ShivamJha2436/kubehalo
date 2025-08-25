package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"github.com/ShivamJha2436/kubehalo/internal/webhook"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/validate", webhook.Serve)

	// TLS cert and key (mounted from Secret)
	certFile := "/tls/tls.crt"
	keyFile := "/tls/tls.key"

	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12},
		Handler:   mux,
	}

	log.Println("[webhook] starting on :8443/validate")
	if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
		log.Fatal(err)
	}
}
