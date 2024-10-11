package cmds

import (
	"context"
	"fmt"
	"github.com/Arnobkumarsaha/kubectl-utils/formatter"
	"github.com/Arnobkumarsaha/kubectl-utils/parser"
	"github.com/Arnobkumarsaha/kubectl-utils/store"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ClusterROleCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "crole",
		Aliases: []string{"cr", "clusterrole"},
		Example: `rbac crole --name kube-binder --typ "crb,rb,sa" --oyaml`,
		Run: func(cmd *cobra.Command, args []string) {
			parser.Parse()
			_ = calcClusterRole()
			formatter.Print(fmt.Sprintf("%s ClusterRole", name), oyaml)
		},
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
	}

	addCommonFlags(cmd)
	cmd.Flags().StringVar(&parser.TypeStr, "typ", "crb,rb,sa", "typ of relationships with crole, you want")
	return cmd
}

func calcClusterRole() error {
	if parser.Crb || parser.Sa {
		crbs, err := c.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return err
		}
		for _, c := range crbs.Items {
			if !isOurClusterRole(c.RoleRef) {
				continue
			}
			store.ClusterRoleBindings = append(store.ClusterRoleBindings, c)
			err = collectSubjects(c.Subjects)
			if err != nil {
				return err
			}
		}
	}

	if parser.Rb || parser.Sa {
		rbs, err := c.RbacV1().RoleBindings("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return err
		}
		for _, c := range rbs.Items {
			if !isOurClusterRole(c.RoleRef) {
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
