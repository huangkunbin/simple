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

	globalRole, exists := mod.db.LookupGlobalRoleBaseByUserName(userName)
	if !exists {
		globalRole = &mdb.GlobalRoleBase{
			UserName: userName,
			Password: password,
		}
		mod.db.InsertGlobalRoleBase(globalRole)
	}
	roleDB := mod.db.GetRoleDB(globalRole.Id)
	state.Database = roleDB
	state.RoleId = globalRole.Id

	roleData := roleDB.LookupRoleData(globalRole.Id)
	if roleData == nil {
		roleData = &mdb.RoleData{}
		roleDB.InsertRoleData(roleData)
	}
	roleData.Dianomd++
	roleDB.UpdateRoleData(roleData)

	return userName
}
