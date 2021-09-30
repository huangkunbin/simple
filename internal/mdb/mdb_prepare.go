package mdb

import "database/sql"

type syncSQL struct {
	InsertRoleData *sql.Stmt
	DeleteRoleData *sql.Stmt
	UpdateRoleData *sql.Stmt
}

func (ss *syncSQL) Init(db *sql.DB) {
	ss.InsertRoleData = ss.Prepare(db, "INSERT INTO `role_data` SET `id`=?, `role_id`=?, `diamond`=?")
	ss.DeleteRoleData = ss.Prepare(db, "DELETE FROM `role_data` WHERE `id`=?")
	ss.UpdateRoleData = ss.Prepare(db, "UPDATE `role_data` SET `id`=?, `role_id`=?, `diamond`=?")
}
