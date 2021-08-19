package module

import "simple/internal/module/role"

type IRole interface {
	Login(userName, password string) string
}

var (
	Role IRole = &role.RoleMod{}
)
