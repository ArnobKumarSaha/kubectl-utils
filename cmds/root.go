package cmds

import (
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
)

var (
	c         *kubernetes.Clientset
	name      string
	namespace string
	oyaml     bool
)

func NewRootCMD() *cobra.Command {
	return &cobra.Command{
		Use: "kubectl-utils",
		Run: func(cmd *cobra.Command, args []string) {},
	}
}

func NewRbacCMD() *cobra.Command {
	return &cobra.Command{
		Use: "rbac",
		Run: func(cmd *cobra.Command, args []string) {},
	}
}

func Execute(cmd *cobra.Command) {
	cmd.AddCommand(ServiceAccountCMD())
	cmd.AddCommand(RoleCMD())
	cmd.AddCommand(ClusterROleCMD())
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func SetKubernetesClient(config *rest.Config) {
	var err error
	c, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}
