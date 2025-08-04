package kube

import (
	"flag"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// NewClient returns a Kubernetes clientset (in-cluster or out-of-cluster)
func NewClient() (*kubernetes.Clientset, error) {
	// Try in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		// Fallback to kubeconfig (for local dev)
		kubeconfig := filepath.Join(homeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	return kubernetes.NewForConfig(config)
}

func homeDir() string {
	if h := filepath.Dir(flag.Lookup("test.v").Value.String()); h != "" {
		return h
	}
	return filepath.Dir(".")
}