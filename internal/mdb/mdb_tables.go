package mdb

type rowIds struct {
	Role           int64
	GlobalRoleBase int64
	RoleData       int64
}

func (ids *rowIds) Init(shardID int64) {
	ids.Role = shardID
	ids.GlobalRoleBase = shardID
	ids.RoleData = shardID
}

type indexes struct {
}

func (indexes *indexes) Init() {
}

type globalTables struct {
	GlobalRoleBase           map[int64]*GlobalRoleBase
	GlobalRoleBaseByUserName map[string]*GlobalRoleBase
}

func NewGlobalTables() *globalTables {
	return &globalTables{
		GlobalRoleBase:           make(map[int64]*GlobalRoleBase, 10000),
		GlobalRoleBaseByUserName: make(map[string]*GlobalRoleBase, 10000),
	}
}

type roleTables struct {
	RoleId   int64
	RoleData []RoleData
}

func NewRoleTables() *roleTables {
	return &roleTables{
		RoleData: []RoleData{},
	}
}
