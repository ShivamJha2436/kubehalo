package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ShivamJha2436/kubehalo/controllers/scalepolicy"
	"github.com/ShivamJha2436/kubehalo/internal/kube"
)

func main() {
	// Build clients (works both in-cluster and locally with kubeconfig)
	_, dyn, _, err := kube.NewClients()
	if err != nil {
		log.Fatalf("failed to build kubernetes clients: %v", err)
	}

	ctrl := scalepolicy.NewController(dyn)

	// Handle shutdown signals gracefully
	stopCh := make(chan struct{})
	defer close(stopCh)

	go ctrl.Run(stopCh)

	// Block until SIGINT/SIGTERM
	sigCh := make(chan os.Signal, 2)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	s := <-sigCh
	log.Printf("[main] received signal %s, exiting...", s)
}
