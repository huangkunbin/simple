package mdb

import (
	"database/sql"
	"sync/atomic"
)

type RoleBase struct {
	_has     bool
	_version int
	RoleId   int64
	Nickname string
	Password string
	Dianomd  int64
}

func (ld *loader) LoadRoleBase() {
	rows := ld.LoadRoleTable("role_base", "role_id", &ld.db.rowIds.RoleBase)
	defer rows.Close()
	for rows.Next() {
		var (
			vRoleId   int64
			vNickname sql.NullString
			vPassword sql.NullString
			vDianomd  int64
		)
		err := rows.Scan(
			&vRoleId,
			&vNickname,
			&vPassword,
			&vDianomd,
		)
		if err != nil {
			panic(err)
		}
		rdb := ld.db.getOrCreateTables(vRoleId)
		row := &RoleBase{
			_has:     true,
			_version: 1,
			RoleId:   vRoleId,
			Nickname: vNickname.String,
			Password: vPassword.String,
			Dianomd:  vDianomd,
		}
		appendRoleBase(rdb.RoleBase, row)
	}
}

func appendRoleBase(s []RoleBase, v *RoleBase) []RoleBase {
	return append(s, *v)
}

func (rdb *RoleDB) SelectRoleBase(callback func(item *RoleBase) (isBreak bool)) {
	for i := 0; i < len(rdb.tables.RoleBase); i++ {
		if rdb.tables.RoleBase[i]._has && callback(&rdb.tables.RoleBase[i]) {
			break
		}
	}
}

func (rdb *RoleDB) LookupRoleBase(id int64) *RoleBase {
	for i := 0; i < len(rdb.tables.RoleBase); i++ {
		if rdb.tables.RoleBase[i]._has && rdb.tables.RoleBase[i].RoleId == id {
			return &rdb.tables.RoleBase[i]
		}
	}
	return nil
}

func (rdb *RoleDB) InsertRoleBase(roleBase *RoleBase) {
	if roleBase._version != 0 {
		panic("Dirty Insert RoleBase")
	}
	roleBase._version++
	roleBase._has = true
	roleBase.RoleId = rdb.roleId
	newId := atomic.AddInt64(&(rdb.rowIds.RoleBase), 1)
	roleBase.RoleId = newId
	var done bool
	for i := 0; i < len(rdb.tables.RoleBase); i++ {
		if !rdb.tables.RoleBase[i]._has {
			rdb.tables.RoleBase[i] = *roleBase
			done = true
			break
		}
	}
	if !done {
		rdb.tables.RoleBase = appendRoleBase(rdb.tables.RoleBase, roleBase)
	}
	rdb.addTransLog(&roleBaseTransLog{
		db:     rdb,
		Table:  "role_base",
		Action: TRANS_INSERT,
		New:    *roleBase,
	})
}

func (rdb *RoleDB) DeleteRoleBase(roleBase *RoleBase) {
	for i := 0; i < len(rdb.tables.RoleBase); i++ {
		if !rdb.tables.RoleBase[i]._has {
			continue
		}
		if rdb.tables.RoleBase[i].RoleId == roleBase.RoleId {
			rdb.addTransLog(&roleBaseTransLog{
				db:     rdb,
				Action: TRANS_DELETE,
				Old:    rdb.tables.RoleBase[i],
			})
			rdb.tables.RoleBase[i]._has = false
			break
		}
	}
}

func (rdb *RoleDB) UpdateRoleBase(roleBase *RoleBase) {
	for i := 0; i < len(rdb.tables.RoleBase); i++ {
		if !rdb.tables.RoleBase[i]._has {
			continue
		}
		if rdb.tables.RoleBase[i].RoleId == roleBase.RoleId {
			if roleBase._version != rdb.tables.RoleBase[i]._version {
				panic("Dirty Update RoleBase")
			}
			roleBase._version++
			rdb.addTransLog(&roleBaseTransLog{
				db:     rdb,
				Action: TRANS_UPDATE,
				New:    *roleBase,
				Old:    rdb.tables.RoleBase[i],
			})
			rdb.tables.RoleBase[i] = *roleBase
			return
		}
	}
	panic("Bad Update RoleBase")
}

type roleBaseTransLog struct {
	db     *RoleDB
	Table  string
	Action string
	Old    RoleBase
	New    RoleBase
}

func (l *roleBaseTransLog) Commit(tx *sql.Tx, sql *syncSQL) error {
	switch l.Action {
	case TRANS_INSERT:
		stmt := sql.InsertRoleBase
		_, err := tx.Stmt(stmt).Exec(
			l.New.RoleId,
			l.New.Nickname,
			l.New.Password,
			l.New.Dianomd,
		)
		return err
	case TRANS_DELETE:
		stmt := sql.DeleteRoleBase
		_, err := tx.Stmt(stmt).Exec(l.Old.RoleId)
		return err
	case TRANS_UPDATE:
		stmt := sql.UpdateRoleBase
		_, err := tx.Stmt(stmt).Exec(
			l.New.RoleId,
			l.New.Nickname,
			l.New.Password,
			l.New.Dianomd,
		)
		return err
	}
	return nil
}

func (l *roleBaseTransLog) Rollback() {
	switch l.Action {
	case TRANS_INSERT:
		for i := 0; i < len(l.db.tables.RoleBase); i++ {
			if l.db.tables.RoleBase[i].RoleId == l.New.RoleId {
				l.db.tables.RoleBase[i]._has = false
				break
			}
		}
	case TRANS_DELETE:
		var done bool
		for i := 0; i < len(l.db.tables.RoleBase); i++ {
			if !l.db.tables.RoleBase[i]._has {
				l.db.tables.RoleBase[i] = l.Old
				done = true
				break
			}
		}
		if !done {
			l.db.tables.RoleBase = appendRoleBase(l.db.tables.RoleBase, &l.Old)
		}
	case TRANS_UPDATE:
		for i := 0; i < len(l.db.tables.RoleBase); i++ {
			if l.db.tables.RoleBase[i].RoleId == l.Old.RoleId {
				l.db.tables.RoleBase[i] = l.Old
				break
			}
		}
	}
}
