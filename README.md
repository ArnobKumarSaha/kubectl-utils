# kubectl-utils

This is a helper tool for `kubectl` cli. It aims to implement the general-purpose utility features, those are currently missing.

## Install 
`git clone https://github.com/ArnobKumarSaha/kubectl-utils.git` <br>
`cd kubectl-utils` <br>
`go install`

## Example Commands
### rbac
To get the connections easily among serviceAccount(sa), role, roleBinding(rb), clusterRole(crole) & clusterRoleBinding(crb)

`kubectl utils rbac sa -n bb-r5z2w --name kube-binder --typ "crb,role,rb"` <br>
`kubectl utils rbac crole --name kube-binder --typ "crb,rb,sa" --oyaml` <br>
`kubectl utils rbac role -n kube-system --name kube-proxy --typ "rb,sa" --oyaml`
