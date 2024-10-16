package parser

import (
	"fmt"
	"strings"
)

var (
	TypeStr     string
	Role, CRole bool
	Rb, Crb, Sa bool
)

func Parse() {
	if TypeStr == "" {
		return
	}
	strs := strings.Split(TypeStr, ",")
	for _, str := range strs {
		switch strings.ToLower(str) {
		case "crb", "croleb", "crolebinding", "clusterrolebinding":
			Crb = true
			continue
		case "rb", "rbinding", "rolebinding":
			Rb = true
			continue
		case "role":
			Role = true
			continue
		case "crole", "clusterrole":
			CRole = true
			continue
		case "sa", "serviceaccount":
			Sa = true
			continue
		default:
			_ = fmt.Errorf("Type %s not matched \n", str)
		}
	}
}
