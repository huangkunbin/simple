package module

import (
	"simple/internal/mdb"
	"simple/pkg/simplenet"
	"simple/pkg/util"
)

type RoleMod struct {
	db *mdb.Database
}

func InitRole(db *mdb.Database) *RoleMod {
	return &RoleMod{db: db}
}

func (mod *RoleMod) Login(session simplenet.ISession, userName, password string) string {
	userlen := len(userName)
	util.Assert(userlen > 0, "user is empty.")
	util.Assert(userlen <= 100, "user length too long.")

	state := State(session)

	if globalRole, exists := mod.db.LookupGlobalRoleBaseByUserName(userName); exists {
		db := mod.db.GetRoleDB(globalRole.Id)
		state.Database = db
		state.RoleId = globalRole.Id
	} else {
		globalRole = &mdb.GlobalRoleBase{
			UserName: userName,
			Password: password,
		}
		mod.db.InsertGlobalRoleBase(globalRole)
		db := mod.db.GetRoleDB(globalRole.Id)
		state.Database = db
		state.RoleId = globalRole.Id
	}
	return userName
}
