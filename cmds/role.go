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

func RoleCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "role",
		Example: `rbac role -n kube-system --name kube-proxy --typ "rb,sa" --oyaml`,
		Run: func(cmd *cobra.Command, args []string) {
			parser.Parse()
			_ = calcRole()
			formatter.Print(fmt.Sprintf("%s/%s Role", namespace, name), oyaml)
		},
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
	}

	cmd.Flags().StringVar(&name, "name", name, "name of role")
	cmd.Flags().StringVarP(&namespace, "namespace", "n", namespace, "namespace of role")
	cmd.Flags().StringVar(&parser.TypeStr, "typ", "crb,rb,crole,role", "typ of relationships with role, you want")

	cmd.Flags().BoolVarP(&oyaml, "oyaml", "y", oyaml, "shows yaml too")
	cmd.Flags().Lookup("oyaml").NoOptDefVal = "true"
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func calcRole() error {
	if parser.Rb || parser.Sa {
		rbs, err := c.RbacV1().RoleBindings(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return err
		}
		for _, c := range rbs.Items {
			if !isOurRole(c.RoleRef) {
				continue
			}
			store.RoleBindings = append(store.RoleBindings, c)
			err = collectForRole(c.Subjects)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func collectForRole(subjects []rbacv1.Subject) error {
	for _, s := range subjects {
		sub, err := c.CoreV1().ServiceAccounts(s.Namespace).Get(context.TODO(), s.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		store.ServiceAccounts = append(store.ServiceAccounts, *sub)
	}
	return nil
}

func isOurRole(r rbacv1.RoleRef) bool {
	return r.Kind == "Role" && r.Name == name
}