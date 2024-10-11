package store

import (
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
)

var (
	ClusterRoleBindings []rbacv1.ClusterRoleBinding
	ClusterRoles        []rbacv1.ClusterRole
	RoleBindings        []rbacv1.RoleBinding
	Roles               []rbacv1.Role
	ServiceAccounts     []corev1.ServiceAccount
)
