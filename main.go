package main

import (
	"fmt"
	"github.com/Arnobkumarsaha/kubectl-utils/cmds"
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
	rbacCmd := cmds.NewRbacCMD()
	rootCmd.AddCommand(rbacCmd)
	cmds.Execute(rbacCmd)
}
