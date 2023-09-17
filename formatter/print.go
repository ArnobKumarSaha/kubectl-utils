package formatter

import (
	"fmt"
	"github.com/Arnobkumarsaha/rbac/parser"
	"github.com/Arnobkumarsaha/rbac/store"
	"gomodules.xyz/oneliners"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Print(title string, yaml bool) {
	fmt.Printf("::::::::::: Printing the resources connected with %s ::::::::::: \n", title)
	if parser.Crb {
		fmt.Printf("ClusterRoleBindings ==> ")
		for _, c := range store.ClusterRoleBindings {
			pri(&c, fmt.Sprintf("ClusterRoleBinding %s", c.GetName()), yaml)
		}
	}
	if parser.CRole {
		fmt.Printf("\nClusterRoles ==> ")
		for _, c := range store.ClusterRoles {
			pri(&c, fmt.Sprintf("ClusterRole %s", c.GetName()), yaml)
		}
	}
	if parser.Rb {
		fmt.Printf("\nRoleBindings ==> ")
		for _, c := range store.RoleBindings {
			pri(&c, fmt.Sprintf("RoleBinding %s/%s", c.GetNamespace(), c.GetName()), yaml)
		}
	}
	if parser.Role {
		fmt.Printf("\nRoles ==> ")
		for _, c := range store.Roles {
			pri(&c, fmt.Sprintf("Role %s/%s", c.GetNamespace(), c.GetName()), yaml)
		}
	}
	if parser.Sa {
		fmt.Printf("\nServiceAccounts ==> ")
		for _, c := range store.ServiceAccounts {
			pri(&c, fmt.Sprintf("ServiceAccount %s/%s", c.GetNamespace(), c.GetName()), yaml)
		}
	}
}

func pri(c metav1.Object, header string, yaml bool) {
	if yaml {
		oneliners.PrettyJson(c, header)
	} else {
		fmt.Printf("%s, ", c.GetName())
	}
}
