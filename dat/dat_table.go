package dat

import "sync"

type rowIds struct {
	Role       int64
	GlobalRole int64
}

func (ids *rowIds) Init(serverId int64) {
	ids.Role = serverId
	ids.GlobalRole = serverId
}

type indexes struct {
	idxGlobalRoleMutex sync.Mutex
}

func (indexes *indexes) Init() {

}

type GlobalTables struct {
	*GlobalRole
}

func NewGlobalTables() *GlobalTables {
	return &GlobalTables{
		&GlobalRole{},
	}
}

type GlobalRole struct {
}

type RoleTables struct {
	Pid int64
}

func NewRoleTables() *RoleTables {
	return &RoleTables{}
}
