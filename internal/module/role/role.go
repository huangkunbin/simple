package role

import "simple/internal/mdb"

type RoleMod struct {
	db *mdb.Database
}

func Init(db *mdb.Database) *RoleMod {
	return &RoleMod{db: db}
}

func (b *RoleMod) Login(userName, password string) string {
	return userName
}
