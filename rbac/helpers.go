package rbac

import (
	"context"
	"github.com/Arnobkumarsaha/kubectl-utils/client"
	"github.com/Arnobkumarsaha/kubectl-utils/rbac/store"
	"github.com/spf13/cobra"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func addCommonFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&name, "name", name, "")
	cmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", namespace, "")

	cmd.PersistentFlags().BoolVarP(&oyaml, "oyaml", "y", oyaml, "shows yaml too")
	cmd.PersistentFlags().Lookup("oyaml").NoOptDefVal = "true"
	_ = cmd.MarkFlagRequired("name")
}

func collectSubjects(subjects []rbacv1.Subject) error {
	for _, s := range subjects {
		sub, err := client.Client.CoreV1().ServiceAccounts(s.Namespace).Get(context.TODO(), s.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		store.ServiceAccounts = append(store.ServiceAccounts, *sub)
	}
	return nil
}

func isOurSA(sub rbacv1.Subject) bool {
	return sub.Kind == "ServiceAccount" && sub.Name == name && sub.Namespace == namespace
}

func isOurRole(r rbacv1.RoleRef) bool {
	return r.Kind == "Role" && r.Name == name
}

func isOurClusterRole(r rbacv1.RoleRef) bool {
	return r.Kind == "ClusterRole" && r.Name == name
}
