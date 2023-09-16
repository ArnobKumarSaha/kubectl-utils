package store

import rbacv1 "k8s.io/api/rbac/v1"

var (
	ClusterRoleBindings []rbacv1.ClusterRoleBinding
	ClusterRoles        []rbacv1.ClusterRole
	RoleBindings        []rbacv1.RoleBinding
	Roles               []rbacv1.Role
)
