package main

import (
	"fmt"
	"github.com/Arnobkumarsaha/kubectl-utils/client"
	"github.com/Arnobkumarsaha/kubectl-utils/rbac"
	"github.com/spf13/cobra"
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
	client.SetKubernetesClient(config)
}

func main() {
	rootCmd := NewRootCMD()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func NewRootCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use: "kubectl-utils",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	cmd.AddCommand(rbac.NewCMD())
	return cmd
}
