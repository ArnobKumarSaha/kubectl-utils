package cmds

import (
	"context"
	"fmt"
	"github.com/Arnobkumarsaha/rbac/formatter"
	"github.com/Arnobkumarsaha/rbac/parser"
	"github.com/Arnobkumarsaha/rbac/store"
	"github.com/spf13/cobra"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ServiceAccountCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sa",
		Example: `rbac sa -n bb-r5z2w --name kube-binder --typ "crb,role,rb" --oyaml`,
		Run: func(cmd *cobra.Command, args []string) {
			parser.Parse()
			_ = calcSA()
			formatter.Print(fmt.Sprintf("%s/%s ServiceAccount", namespace, name), oyaml)
		},
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
	}

	addCommonFlags(cmd)
	cmd.Flags().StringVar(&parser.TypeStr, "typ", "crb,rb,crole,role", "typ of relationships with sa, you want")
	return cmd
}

func calcSA() error {
	if parser.Crb || parser.CRole {
		crbs, err := c.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return err
		}
		for _, c := range crbs.Items {
			for _, sub := range c.Subjects {
				if !isOurSA(sub) {
					continue
				}
				store.ClusterRoleBindings = append(store.ClusterRoleBindings, c)
				err = collectForSA(c.RoleRef, "")
				if err != nil {
					return err
				}
			}
		}
	}

	if parser.CRole || parser.Rb || parser.Role {
		rbs, err := c.RbacV1().RoleBindings("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return err
		}
		for _, c := range rbs.Items {
			for _, sub := range c.Subjects {
				if !isOurSA(sub) {
					continue
				}
				store.RoleBindings = append(store.RoleBindings, c)
				err = collectForSA(c.RoleRef, c.Namespace)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func collectForSA(ref rbacv1.RoleRef, ns string) error {
	if ref.Kind == "ClusterRole" {
		x, err := c.RbacV1().ClusterRoles().Get(context.TODO(), ref.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		store.ClusterRoles = append(store.ClusterRoles, *x)
	} else if ref.Kind == "Role" {
		x, err := c.RbacV1().Roles(ns).Get(context.TODO(), ref.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		store.Roles = append(store.Roles, *x)
	}
	return nil
}


