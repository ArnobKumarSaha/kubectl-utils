package client

import (
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	Client          *kubernetes.Clientset
	DynamicClient   dynamic.Interface
	DiscoveryClient discovery.DiscoveryInterface
)

func SetClients(config *rest.Config) {
	SetKubernetesClient(config)
	SetDynamicClient(config)
	SetDiscoveryClient(config)

}

func SetKubernetesClient(config *rest.Config) {
	var err error
	Client, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}

func SetDynamicClient(config *rest.Config) {
	var err error
	DynamicClient, err = dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}

func SetDiscoveryClient(config *rest.Config) {
	var err error
	DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}
