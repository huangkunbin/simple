package module

import (
	"simple/internal/mdb"
	"simple/internal/module/role"
)

type IRole interface {
	Login(userName, password string) string
}

var (
	Role IRole
)

func InitModule(db *mdb.Database) {
	Role = role.Init(db)
}
