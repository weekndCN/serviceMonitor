package k8s

import (
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// KubeConfig define config file path
type KubeConfig struct {
	Path string
}

// NewClient resturn kubernetes client
func (c *KubeConfig) NewClient() *kubernetes.Clientset {

	config, err := clientcmd.BuildConfigFromFlags("", c.Path)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return clientset
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func int32Ptr(i int32) *int32 {
	return &i
}
