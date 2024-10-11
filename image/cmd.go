package image

import (
	"context"
	"github.com/Arnobkumarsaha/kubectl-utils/client"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

var (
	name      string
	namespace string
	oyaml     bool
)

func NewCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use: "image",
		Run: func(cmd *cobra.Command, args []string) {
			klog.Infof("image name: %s, namespace: %s", name, namespace)
			pods, err := client.Client.CoreV1().Pods(corev1.NamespaceDefault).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				panic(err)
			}
			for _, item := range pods.Items {
				klog.Infof(item.Name)
			}
		},
	}
	return cmd
}
