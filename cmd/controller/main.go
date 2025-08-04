package controller

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"kubehalo/internal/kube"
	"kubehalo/controllers/scalepolicy"
)

func main() {
	// Step 1: Create Kubernetes client (works both inside and outside the cluster)
	clientset, err := kube.NewClient()
	if err != nil {
		log.Fatalf("[Startup Error] Failed to create Kubernetes client: %v", err)
	}

	// Step 2: Create a channel to handle graceful shutdowns
	stopCh := make(chan struct{})

	// Step 3: Start the controller loop for ScalePolicy
	err = scalepolicy.StartController(clientset, stopCh)
	if err != nil {
		log.Fatalf("[Startup Error] Failed to start ScalePolicy controller: %v", err)
	}

	log.Println("[INFO] ScalePolicy controller started successfully")

	// Step 4: Wait for OS signals to shut down gracefully
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs // blocking until signal is received

	log.Println("[INFO] Received termination signal, shutting down...")
	close(stopCh)
}