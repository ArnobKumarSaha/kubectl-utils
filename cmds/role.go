package cmds

import (
	"context"
	"fmt"
	"github.com/Arnobkumarsaha/rbac/formatter"
	"github.com/Arnobkumarsaha/rbac/parser"
	"github.com/Arnobkumarsaha/rbac/store"
	"github.com/spf13/cobra"
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

	addCommonFlags(cmd)
	cmd.Flags().StringVar(&parser.TypeStr, "typ", "rb,sa", "typ of relationships with role, you want")
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
			err = collectSubjects(c.Subjects)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
