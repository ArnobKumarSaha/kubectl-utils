package main

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"strings"
)

var (
	name, namespace             string
	typStr                      string
	oyaml                       bool
	role, crole                 bool
	rb, crb                     bool
	filteredClusterRoleBindings []rbacv1.ClusterRoleBinding
	filteredClusterRoles        []rbacv1.ClusterRole
	filteredRoleBindings        []rbacv1.RoleBinding
	filteredRoles               []rbacv1.Role
)

func ServiceAccountCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sa",
		Short: "list",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("sa called !")
			parse()
			_ = calcSA()
		},
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
	}

	cmd.Flags().StringVar(&name, "name", name, "name of sa")
	cmd.Flags().StringVarP(&namespace, "namespace", "n", namespace, "namespace of sa")
	cmd.Flags().StringVar(&typStr, "typ", typStr, "typ of relationships with sa you want")

	cmd.Flags().BoolVarP(&oyaml, "oyaml", "y", oyaml, "shows yaml too")
	cmd.Flags().Lookup("oyaml").NoOptDefVal = "true"
	cmd.MarkFlagRequired("name")
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
			fmt.Errorf("Type %s not matched \n", str)
		}
	}
	klog.Infof("crb=%v rb=%v crole=%v role=%v  \n", crb, rb, crole, role)
}

func calcSA() error {
	crbs, err := clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, c := range crbs.Items {
		for _, sub := range c.Subjects {
			if sub.Kind == "ServiceAccount" && sub.Name == name && sub.Namespace == namespace {
				filteredClusterRoleBindings = append(filteredClusterRoleBindings, c)
				collect(c.RoleRef, "")
			}
		}
	}

	rbs, err := clientset.RbacV1().RoleBindings("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, c := range rbs.Items {
		for _, sub := range c.Subjects {
			if sub.Kind == "ServiceAccount" && sub.Name == name && sub.Namespace == namespace {
				filteredRoleBindings = append(filteredRoleBindings, c)
				collect(c.RoleRef, c.Namespace)
			}
		}
	}

	klog.Infof("crbs=%+v, crs=%+v, rbs=%+v, roles=%+v \n",
		filteredClusterRoleBindings, filteredClusterRoles, filteredRoleBindings, filteredRoles)
	return nil
}

func collect(ref rbacv1.RoleRef, ns string) error {
	if ref.Kind == "ClusterRole" {
		x, err := clientset.RbacV1().ClusterRoles().Get(context.TODO(), ref.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		filteredClusterRoles = append(filteredClusterRoles, *x)
	} else if ref.Kind == "Role" {
		x, err := clientset.RbacV1().Roles(ns).Get(context.TODO(), ref.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		filteredRoles = append(filteredRoles, *x)
	}
	return nil
}
