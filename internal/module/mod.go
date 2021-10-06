package module

import (
	"simple/internal/mdb"

	"simple/pkg/simplenet"
)

type IRole interface {
	Login(session simplenet.ISession, userName, password string) string
}

var (
	Role IRole
)

func InitModule(db *mdb.Database) {
	Role = InitRole(db)
}
