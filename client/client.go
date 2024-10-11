package client

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	Client *kubernetes.Clientset
)

func SetKubernetesClient(config *rest.Config) {
	var err error
	Client, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}
