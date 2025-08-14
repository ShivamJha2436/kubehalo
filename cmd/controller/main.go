package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ShivamJha2436/kubehalo/controllers/scalepolicy"
	"github.com/ShivamJha2436/kubehalo/internal/kube"
	"github.com/ShivamJha2436/kubehalo/internal/metrics"
)

func main() {
	kubeClient, dynClient, _, err := kube.NewClients()
	if err != nil {
		log.Fatalf("failed to build kubernetes clients: %v", err)
	}
	promClient, err := metrics.NewPrometheusClient("http://localhost:9000")
	if err != nil {
		log.Fatalf("failed to create Prometheus client: %v", err)
	}

	ctrl := scalepolicy.NewController(dynClient, kubeClient, promClient)

	// Handle shutdown signals
	stopCh := make(chan struct{})
	defer close(stopCh)

	go ctrl.Run(stopCh)

	// Block until SIGINT/SIGTERM
	sigCh := make(chan os.Signal, 2)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	s := <-sigCh
	log.Printf("[main] received signal %s, exiting...", s)
}
