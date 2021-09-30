package role

import "simple/internal/mdb"

type RoleMod struct {
	db *mdb.Database
}

func Init(db *mdb.Database) *RoleMod {
	return &RoleMod{db: db}
}

func (mod *RoleMod) Login(userName, password string) string {
	return userName
}

func (mod *RoleMod) Create(userName, password string) string {
	// mod.db.GetRoleDB().InsertRoleBase(&mdb.RoleBase{})

	return userName
}
