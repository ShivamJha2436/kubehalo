package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ShivamJha2436/kubehalo/controllers/scalepolicy"
	"github.com/ShivamJha2436/kubehalo/internal/config"
	"github.com/ShivamJha2436/kubehalo/internal/kube"
	"github.com/ShivamJha2436/kubehalo/internal/metrics"
)

func main() {
	cfg := config.LoadControllerConfig()

	clients, err := kube.NewClients()
	if err != nil {
		log.Fatalf("failed to build kubernetes clients: %v", err)
	}

	promClient, err := metrics.NewPrometheusClient(cfg.PrometheusAddress)
	if err != nil {
		log.Fatalf("failed to create Prometheus client: %v", err)
	}

	ctrl := scalepolicy.NewController(clients.Dynamic, clients.Kube, promClient)

	stopCh := make(chan struct{})
	defer close(stopCh)

	go ctrl.Run(stopCh)

	sigCh := make(chan os.Signal, 2)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	s := <-sigCh
	log.Printf("[main] received signal %s, exiting...", s)
}
