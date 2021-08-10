package biz

type RoleBiz struct {
}

func (b *RoleBiz) Login(userName, password string) string {
	return userName
}
