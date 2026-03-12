package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
	"os"
	"path/filepath"
)

type Clients struct {
	Kubernetes *kubernetes.Clientset
	Metrics    *metricsv.Clientset
}

func NewClients() (*Clients, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		// Fallback to kubeconfig for local development
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	metricsClient, err := metricsv.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Clients{
		Kubernetes: k8sClient,
		Metrics:    metricsClient,
	}, nil
}
