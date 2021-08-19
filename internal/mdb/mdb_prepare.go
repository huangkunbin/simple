package mdb

import "database/sql"

type syncSQL struct {
	InsertRoleBase *sql.Stmt
	DeleteRoleBase *sql.Stmt
	UpdateRoleBase *sql.Stmt
}

func (ss *syncSQL) Init(db *sql.DB) {
	ss.InsertRoleBase = ss.Prepare(db, "INSERT INTO `role_base` SET `role_id`=?, `nickname`=?, `password`=?, `diamond`=?")
	ss.DeleteRoleBase = ss.Prepare(db, "DELETE FROM `role_base` WHERE `role_id`=?")
	ss.UpdateRoleBase = ss.Prepare(db, "UPDATE `role_base` SET `role_id`=?, `nickname`=?, `password`=?, `diamond`=?")
}
