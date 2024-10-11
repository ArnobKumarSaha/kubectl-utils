# kubectl-utils

This is a helper cli tool for kubernetes to get the connections easily among serviceAccount(sa), role, roleBinding(rb), clusterRole(crole) & clusterRoleBinding(crb).

## Install 
`git clone https://github.com/ArnobKumarSaha/kubectl-utils.git` <br>
`cd kubectl-utils` <br>
`go install`

## Example Commands
### rbac
`kubectl utils rbac sa -n bb-r5z2w --name kube-binder --typ "crb,role,rb"` <br>
`kubectl utils rbac crole --name kube-binder --typ "crb,rb,sa" --oyaml` <br>
`kubectl utils rbac role -n kube-system --name kube-proxy --typ "rb,sa" --oyaml`
