package parser

import (
	"fmt"
	"strings"
)

var (
	TypeStr     string
	Role, CRole bool
	Rb, Crb     bool
)

func Parse() {
	if TypeStr == "" {
		return
	}
	strs := strings.Split(TypeStr, ",")
	for _, str := range strs {
		switch str {
		case "crb":
			Crb = true
			continue
		case "rb":
			Rb = true
			continue
		case "role":
			Role = true
			continue
		case "crole":
			CRole = true
			continue
		default:
			_ = fmt.Errorf("Type %s not matched \n", str)
		}
	}
}
