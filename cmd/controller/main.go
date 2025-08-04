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
	clientset, err := kube.NewClient()
	if err != nil {
		log.Fatalf("Error building k8s client: %v", err)
	}

	stopCh := make(chan struct{})
	defer close(stopCh)

	if err := scalepolicy.StartController(clientset, stopCh); err != nil {
		log.Fatalf("Controller failed: %v", err)
	}

	log.Println("[MAIN] Controller is running. Press Ctrl+C to stop...")
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	<-sigCh

	log.Println("[MAIN] Shutting down gracefully...")
}
