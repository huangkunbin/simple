package mdb

import (
	"database/sql"
	"sync/atomic"
)

type RoleData struct {
	_has     bool
	_version int
	Id       int64
	RoleId   int64
	Dianomd  int64
}

func (ld *loader) LoadRoleData() {
	rows := ld.LoadRoleTable("role_data", "role_id", &ld.db.rowIds.RoleData)
	defer rows.Close()
	for rows.Next() {
		var (
			vId      int64
			vRoleId  int64
			vDianomd int64
		)
		err := rows.Scan(
			&vId,
			&vRoleId,
			&vDianomd,
		)
		if err != nil {
			panic(err)
		}
		rdb := ld.db.getOrCreateTables(vRoleId)
		row := &RoleData{
			_has:     true,
			_version: 1,
			Id:       vId,
			RoleId:   vRoleId,
			Dianomd:  vDianomd,
		}
		appendRoleData(rdb.RoleData, row)
	}
}

func appendRoleData(s []RoleData, v *RoleData) []RoleData {
	return append(s, *v)
}

func (rdb *RoleDB) SelectRoleData(callback func(item *RoleData) (isBreak bool)) {
	for i := 0; i < len(rdb.tables.RoleData); i++ {
		if rdb.tables.RoleData[i]._has && callback(&rdb.tables.RoleData[i]) {
			break
		}
	}
}

func (rdb *RoleDB) LookupRoleData(id int64) *RoleData {
	for i := 0; i < len(rdb.tables.RoleData); i++ {
		if rdb.tables.RoleData[i]._has && rdb.tables.RoleData[i].RoleId == id {
			return &rdb.tables.RoleData[i]
		}
	}
	return nil
}

func (rdb *RoleDB) InsertRoleData(roleBase *RoleData) {
	if roleBase._version != 0 {
		panic("Dirty Insert RoleData")
	}
	roleBase._version++
	roleBase._has = true
	roleBase.RoleId = rdb.roleId
	newId := atomic.AddInt64(&(rdb.rowIds.RoleData), 1)
	roleBase.Id = newId
	var done bool
	for i := 0; i < len(rdb.tables.RoleData); i++ {
		if !rdb.tables.RoleData[i]._has {
			rdb.tables.RoleData[i] = *roleBase
			done = true
			break
		}
	}
	if !done {
		rdb.tables.RoleData = appendRoleData(rdb.tables.RoleData, roleBase)
	}
	rdb.addTransLog(&roleBaseTransLog{
		db:     rdb,
		Table:  "role_base",
		Action: TRANS_INSERT,
		New:    *roleBase,
	})
}

func (rdb *RoleDB) DeleteRoleData(roleBase *RoleData) {
	for i := 0; i < len(rdb.tables.RoleData); i++ {
		if !rdb.tables.RoleData[i]._has {
			continue
		}
		if rdb.tables.RoleData[i].RoleId == roleBase.RoleId {
			rdb.addTransLog(&roleBaseTransLog{
				db:     rdb,
				Action: TRANS_DELETE,
				Old:    rdb.tables.RoleData[i],
			})
			rdb.tables.RoleData[i]._has = false
			break
		}
	}
}

func (rdb *RoleDB) UpdateRoleData(roleBase *RoleData) {
	for i := 0; i < len(rdb.tables.RoleData); i++ {
		if !rdb.tables.RoleData[i]._has {
			continue
		}
		if rdb.tables.RoleData[i].RoleId == roleBase.RoleId {
			if roleBase._version != rdb.tables.RoleData[i]._version {
				panic("Dirty Update RoleData")
			}
			roleBase._version++
			rdb.addTransLog(&roleBaseTransLog{
				db:     rdb,
				Action: TRANS_UPDATE,
				New:    *roleBase,
				Old:    rdb.tables.RoleData[i],
			})
			rdb.tables.RoleData[i] = *roleBase
			return
		}
	}
	panic("Bad Update RoleData")
}

type roleBaseTransLog struct {
	db     *RoleDB
	Table  string
	Action string
	Old    RoleData
	New    RoleData
}

func (l *roleBaseTransLog) Commit(tx *sql.Tx, sql *syncSQL) error {
	switch l.Action {
	case TRANS_INSERT:
		stmt := sql.InsertRoleData
		_, err := tx.Stmt(stmt).Exec(
			l.New.Id,
			l.New.RoleId,
			l.New.Dianomd,
		)
		return err
	case TRANS_DELETE:
		stmt := sql.DeleteRoleData
		_, err := tx.Stmt(stmt).Exec(l.Old.Id)
		return err
	case TRANS_UPDATE:
		stmt := sql.UpdateRoleData
		_, err := tx.Stmt(stmt).Exec(
			l.New.Id,
			l.New.RoleId,
			l.New.Dianomd,
		)
		return err
	}
	return nil
}

func (l *roleBaseTransLog) Rollback() {
	switch l.Action {
	case TRANS_INSERT:
		for i := 0; i < len(l.db.tables.RoleData); i++ {
			if l.db.tables.RoleData[i].Id == l.New.Id {
				l.db.tables.RoleData[i]._has = false
				break
			}
		}
	case TRANS_DELETE:
		var done bool
		for i := 0; i < len(l.db.tables.RoleData); i++ {
			if !l.db.tables.RoleData[i]._has {
				l.db.tables.RoleData[i] = l.Old
				done = true
				break
			}
		}
		if !done {
			l.db.tables.RoleData = appendRoleData(l.db.tables.RoleData, &l.Old)
		}
	case TRANS_UPDATE:
		for i := 0; i < len(l.db.tables.RoleData); i++ {
			if l.db.tables.RoleData[i].Id == l.Old.Id {
				l.db.tables.RoleData[i] = l.Old
				break
			}
		}
	}
}
