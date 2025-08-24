package kube

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetRestConfig returns in-cluster config if present, otherwise local kubeconfig.
// - If running inside a Pod, it uses the in-cluster config.
// - Otherwise, it falls back to the local kubeconfig file (~/.kube/config or $KUBECONFIG).
func GetRestConfig() (*rest.Config, error) {
	// Try in-cluster config (when running inside a Pod)
	if cfg, err := rest.InClusterConfig(); err == nil {
		return cfg, nil
	}

	// Fallback to local kubeconfig
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		home, _ := os.UserHomeDir()
		kubeconfig = filepath.Join(home, ".kube", "config")
	}
	if _, err := os.Stat(kubeconfig); err != nil {
		return nil, fmt.Errorf("kubeconfig not found at %s", kubeconfig)
	}

	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig}
	configOverrides := &clientcmd.ConfigOverrides{}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides).ClientConfig()
}

// Clients holds all Kubernetes clients that we commonly use.
type Clients struct {
	Kube    *kubernetes.Clientset // typed client
	Dynamic dynamic.Interface     // dynamic client (for CRDs like ScalePolicy)
	Config  *rest.Config          // raw config
}

// NewClients builds a Clients struct with typed and dynamic clients.
func NewClients() (*Clients, error) {
	cfg, err := GetRestConfig()
	if err != nil {
		return nil, err
	}

	// Typed client
	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	// Dynamic client
	dc, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &Clients{
		Kube:    cs,
		Dynamic: dc,
		Config:  cfg,
	}, nil
}
