package role

import "simple/internal/mdb"

type RoleMod struct {
	db *mdb.Database
}

func (b *RoleMod) Login(userName, password string) string {
	return userName
}
