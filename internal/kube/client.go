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

// NewClients builds both typed and dynamic clients.
func NewClients() (*kubernetes.Clientset, dynamic.Interface, *rest.Config, error) {
	cfg, err := GetRestConfig()
	if err != nil {
		return nil, nil, nil, err
	}

	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, nil, nil, err
	}

	dc, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, nil, nil, err
	}

	return cs, dc, cfg, nil
}
