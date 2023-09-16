package main

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"gomodules.xyz/oneliners"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

var (
	name, namespace     string
	typStr              = "crb,rb,crole,role"
	oyaml               bool
	role, crole         bool
	rb, crb             bool
	clusterRoleBindings []rbacv1.ClusterRoleBinding
	clusterRoles        []rbacv1.ClusterRole
	roleBindings        []rbacv1.RoleBinding
	roles               []rbacv1.Role
)

func ServiceAccountCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sa",
		Short: "list",
		Run: func(cmd *cobra.Command, args []string) {
			parse()
			_ = calcSA()
			printSA()
		},
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
	}

	cmd.Flags().StringVar(&name, "name", name, "name of sa")
	cmd.Flags().StringVarP(&namespace, "namespace", "n", namespace, "namespace of sa")
	cmd.Flags().StringVar(&typStr, "typ", typStr, "typ of relationships with sa you want")

	cmd.Flags().BoolVarP(&oyaml, "oyaml", "y", oyaml, "shows yaml too")
	cmd.Flags().Lookup("oyaml").NoOptDefVal = "true"
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func parse() {
	if typStr == "" {
		return
	}
	strs := strings.Split(typStr, ",")
	for _, str := range strs {
		switch str {
		case "crb":
			crb = true
			continue
		case "rb":
			rb = true
			continue
		case "role":
			role = true
			continue
		case "crole":
			crole = true
			continue
		default:
			_ = fmt.Errorf("Type %s not matched \n", str)
		}
	}
	//fmt.Printf("crb=%v rb=%v crole=%v role=%v  \n", crb, rb, crole, role)
}

func calcSA() error {
	if crb || crole {
		crbs, err := clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return err
		}
		for _, c := range crbs.Items {
			for _, sub := range c.Subjects {
				if !isOurSA(sub) {
					continue
				}
				clusterRoleBindings = append(clusterRoleBindings, c)
				err = collect(c.RoleRef, "")
				if err != nil {
					return err
				}
			}
		}
	}

	if crole || rb || role {
		rbs, err := clientset.RbacV1().RoleBindings("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return err
		}
		for _, c := range rbs.Items {
			for _, sub := range c.Subjects {
				if !isOurSA(sub) {
					continue
				}
				roleBindings = append(roleBindings, c)
				err = collect(c.RoleRef, c.Namespace)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func collect(ref rbacv1.RoleRef, ns string) error {
	if ref.Kind == "ClusterRole" {
		x, err := clientset.RbacV1().ClusterRoles().Get(context.TODO(), ref.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		clusterRoles = append(clusterRoles, *x)
	} else if ref.Kind == "Role" {
		x, err := clientset.RbacV1().Roles(ns).Get(context.TODO(), ref.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		roles = append(roles, *x)
	}
	return nil
}

func isOurSA(sub rbacv1.Subject) bool {
	return sub.Kind == "ServiceAccount" && sub.Name == name && sub.Namespace == namespace
}

func printSA() {
	fmt.Printf("::::::::::: Printing the resources connected with %s/%s ServiceAccount ::::::::::: \n", namespace, name)
	if crb {
		fmt.Printf("ClusterRoleBindings ==> ")
		for _, c := range clusterRoleBindings {
			pri(&c, fmt.Sprintf("ClusterRoleBinding %s", c.GetName()))
		}
	}
	if crole {
		fmt.Printf("\nClusterRoles ==> ")
		for _, c := range clusterRoles {
			pri(&c, fmt.Sprintf("ClusterRole %s", c.GetName()))
		}
	}
	if rb {
		fmt.Printf("\nRoleBindings ==> ")
		for _, c := range roleBindings {
			pri(&c, fmt.Sprintf("RoleBinding %s/%s", c.GetNamespace(), c.GetName()))
		}
	}
	if role {
		fmt.Printf("\nRoles ==> ")
		for _, c := range roles {
			pri(&c, fmt.Sprintf("Role %s/%s", c.GetNamespace(), c.GetName()))
		}
	}
}

func pri(c metav1.Object, header string) {
	if oyaml {
		oneliners.PrettyJson(c, header)
	} else {
		fmt.Printf("%s, ", c.GetName())
	}
}
