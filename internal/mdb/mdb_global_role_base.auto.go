package mdb

type GlobalRoleBase struct {
	_version int
	Id       int64
	Nickname string
	Password string
}

func (ld *loader) LoadGlobalRoleBase() {
	rows := ld.LoadGlobalTable("global_role_base", &ld.db.rowIds.GlobalRoleBase)
	defer rows.Close()
	for rows.Next() {
		var (
			vId       int64
			vNickname string
			vPassword string
		)
		err := rows.Scan(
			&vId,
			&vNickname,
			&vPassword,
		)
		if err != nil {
			panic(err)
		}
		row := &GlobalRoleBase{
			_version: 1,
			Id:       vId,
			Nickname: vNickname,
			Password: vPassword,
		}
		ld.db.globalTables.GlobalRoleBase[row.Id] = row
	}
}

func (db *Database) FetchGlobalRoleBase(callback func(item GlobalRoleBase) bool) {
	if len(db.globalTables.GlobalRoleBase) == 0 {
		return
	}
	for _, v := range db.globalTables.GlobalRoleBase {
		if callback(*v) {
			break
		}
	}
}

func (db *Database) LookupGlobalRoleBase(id int64) (*GlobalRoleBase, bool) {
	globalRoleBase, exists := db.globalTables.GlobalRoleBase[id]
	if exists {
		globalRoleBase := *globalRoleBase
		return &globalRoleBase, true
	}
	return nil, false
}

func (db *Database) InsertGlobalGlobalRoleBase(globalRoleBase *GlobalRoleBase) {
	if globalRoleBase._version != 0 {
		panic("Dirty Insert GlobalRoleBase")
	}
	globalRoleBase._version++
	newGlobalRoleBase := *globalRoleBase
	db.globalTables.GlobalRoleBase[newGlobalRoleBase.Id] = &newGlobalRoleBase
	// db.addTransLog()

}
