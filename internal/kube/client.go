package kube

import (
	"flag"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// NewClient New client returns a Kubernetes clientset, using in-cluster config if available,
// otherwise falling back to the local kubeconfig (for development).
func NewClient() (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error

	// Try in-cluster config
	config, err = rest.InClusterConfig()
	if err != nil {
		// Fall back to kubeconfig (useful for local development)
		kubeconfig := filepath.Join(homeDir(), ".kube", "config")

		// Allow overriding path via --kubeconfig flag (optional)
		kubeconfigFlag := flag.String("kubeconfig", kubeconfig, "(optional) absolute path to the kubeconfig file")
		flag.Parse()

		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfigFlag)
		if err != nil {
			return nil, err
		}
	}
	return kubernetes.NewForConfig(config)
}

// homeDir returns the home directory of the current user.
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // for Windows
}