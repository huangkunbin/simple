package mdb

import (
	"database/sql"
	"sync/atomic"
)

type GlobalRoleBase struct {
	_version int
	Id       int64
	UserName string
	Password string
}

func (ld *loader) LoadGlobalRoleBase() {
	rows := ld.LoadGlobalTable("global_role_base", &ld.db.rowIds.GlobalRoleBase)
	defer rows.Close()
	for rows.Next() {
		var (
			vId       int64
			vUserName string
			vPassword string
		)
		err := rows.Scan(
			&vId,
			&vUserName,
			&vPassword,
		)
		if err != nil {
			panic(err)
		}
		row := &GlobalRoleBase{
			_version: 1,
			Id:       vId,
			UserName: vUserName,
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

func (db *Database) LookupGlobalRoleBaseByUserName(userName string) (*GlobalRoleBase, bool) {
	globalRoleBase, exists := db.globalTables.GlobalRoleBaseByUserName[userName]
	if exists {
		globalRoleBase := *globalRoleBase
		return &globalRoleBase, true
	}
	return nil, false
}

func (db *Database) InsertGlobalRoleBase(globalRoleBase *GlobalRoleBase) {
	if globalRoleBase._version != 0 {
		panic("Dirty Insert GlobalRoleBase")
	}
	globalRoleBase.Id = atomic.AddInt64(&db.rowIds.GlobalRoleBase, 1<<13)
	globalRoleBase._version++
	newGlobalRoleBase := *globalRoleBase
	db.globalTables.GlobalRoleBase[newGlobalRoleBase.Id] = &newGlobalRoleBase
	db.globalTables.GlobalRoleBaseByUserName[newGlobalRoleBase.UserName] = &newGlobalRoleBase
	db.addTransLog(&gloablRoleBaseTransLog{
		db:     db,
		Table:  "global_role_base",
		Action: TRANS_INSERT,
		New:    &newGlobalRoleBase,
	})
}

func (db *Database) UpdateGlobalRoleBase(globalRoleBase *GlobalRoleBase) {
	id := globalRoleBase.Id
	oldGlobalRoleBase := db.globalTables.GlobalRoleBase[id]
	if globalRoleBase._version != oldGlobalRoleBase._version {
		panic("Dirty Update GlobalRoleBase")
	} else {
		globalRoleBase._version++
	}
	newGlobalRoleBase := *globalRoleBase
	db.globalTables.GlobalRoleBase[id] = &newGlobalRoleBase
	db.globalTables.GlobalRoleBaseByUserName[globalRoleBase.UserName] = &newGlobalRoleBase
	db.addTransLog(&gloablRoleBaseTransLog{
		db:     db,
		Table:  "global_role_base",
		Action: TRANS_UPDATE,
		Old:    oldGlobalRoleBase,
		New:    &newGlobalRoleBase,
	})
}

func (db *Database) DeleteGlobalRoleBase(globalRoleBase *GlobalRoleBase) {
	id := globalRoleBase.Id
	oldGlobalRoleBase := db.globalTables.GlobalRoleBase[id]
	delete(db.globalTables.GlobalRoleBase, id)
	delete(db.globalTables.GlobalRoleBaseByUserName, globalRoleBase.UserName)
	db.addTransLog(&gloablRoleBaseTransLog{
		db:     db,
		Table:  "global_role_base",
		Action: TRANS_DELETE,
		Old:    oldGlobalRoleBase,
	})
}

type gloablRoleBaseTransLog struct {
	db     *Database
	Table  string
	Action string
	Old    *GlobalRoleBase
	New    *GlobalRoleBase
}

func (l *gloablRoleBaseTransLog) Commit(tx *sql.Tx, sql *syncSQL) error {
	switch l.Action {
	case TRANS_INSERT:
		stmt := sql.InsertGloabalRoleBase
		_, err := tx.Stmt(stmt).Exec(
			l.New.Id,
			l.New.UserName,
			l.New.Password,
		)
		return err
	case TRANS_DELETE:
		stmt := sql.DeleteGloabalRoleBase
		_, err := tx.Stmt(stmt).Exec(l.Old.Id)
		return err
	case TRANS_UPDATE:
		stmt := sql.UpdateGloabalRoleBase
		_, err := tx.Stmt(stmt).Exec(
			l.New.Id,
			l.New.UserName,
			l.New.Password,
		)
		return err
	}
	return nil
}

func (l *gloablRoleBaseTransLog) Rollback() {
	switch l.Action {
	case TRANS_INSERT:
		delete(l.db.globalTables.GlobalRoleBase, l.New.Id)
		delete(l.db.globalTables.GlobalRoleBaseByUserName, l.New.UserName)
	case TRANS_DELETE:
		l.db.globalTables.GlobalRoleBase[l.Old.Id] = l.Old
		l.db.globalTables.GlobalRoleBaseByUserName[l.Old.UserName] = l.Old
	case TRANS_UPDATE:
		l.db.globalTables.GlobalRoleBase[l.Old.Id] = l.Old
		l.db.globalTables.GlobalRoleBaseByUserName[l.Old.UserName] = l.Old
	}
}
