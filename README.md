# rbac-cli

This is a helper cli tool for kubernetes to get the connections easily among serviceAccount(sa), role, roleBinding(rb), clusterRole(crole) & clusterRoleBinding(crb).

## Install 
`git clone https://github.com/ArnobKumarSaha/rbac.git` <br>
`cd rbac` <br>
`go install`

## Example Commands
`rbac sa -n bb-r5z2w --name kube-binder --typ "crb,role,rb"` <br>
`rbac crole --name kube-binder --typ "crb,rb,sa" --oyaml` <br>
`rbac role -n kube-system --name kube-proxy --typ "rb,sa" --oyaml`
