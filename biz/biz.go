package biz

type IRole interface {
	Login(userName, password string) string
}

var (
	Role IRole = &RoleBiz{}
)
