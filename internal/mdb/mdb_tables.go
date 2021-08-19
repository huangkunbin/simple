package mdb

type rowIds struct {
	Role     int64
	RoleBase int64
}

func (ids *rowIds) Init(serverId int64) {
	ids.Role = serverId
	ids.RoleBase = serverId
}

type indexes struct {
}

func (indexes *indexes) Init() {
}

type globalTables struct {
}

func NewGlobalTables() *globalTables {
	return &globalTables{}
}

type roleTables struct {
	RoleId   int64
	RoleBase []RoleBase
}

func NewRoleTables() *roleTables {
	return &roleTables{
		RoleBase: []RoleBase{},
	}
}
