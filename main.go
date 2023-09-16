package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

var clientset *kubernetes.Clientset

func init() {
	home := homedir.HomeDir()
	path := filepath.Join(home, ".kube", "config")
	if val, exists := os.LookupEnv("KUBECONFIG"); exists {
		path = val
	}

	fmt.Println("kubeconfig path = ", path)
	config, err := clientcmd.BuildConfigFromFlags("", path)
	if err != nil {
		panic(err.Error())
	}

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}

func NewRootCMD() *cobra.Command {
	return &cobra.Command{
		Use: "rbac",
		Run: func(cmd *cobra.Command, args []string) {},
	}
}

func main() {
	rootCmd := NewRootCMD()
	Execute(rootCmd)
}

func Execute(rootCmd *cobra.Command) {
	rootCmd.AddCommand(ServiceAccountCMD())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
