package mdb

func (ld *loader) LoadGlobalTables() {
	ld.LoadGlobalRoleBase()
}

func (ld *loader) LoadRoleTables() {
	ld.LoadRoleData()
}
