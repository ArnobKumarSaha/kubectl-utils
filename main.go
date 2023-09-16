package main

import (
	"fmt"
	"github.com/Arnobkumarsaha/rbac/cmds"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

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
	cmds.SetKubernetesClient(config)
}

func main() {
	rootCmd := cmds.NewRootCMD()
	cmds.Execute(rootCmd)
}
