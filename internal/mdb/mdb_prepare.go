package mdb

import "database/sql"

type syncSQL struct {
	InsertGloabalRoleBase *sql.Stmt
	DeleteGloabalRoleBase *sql.Stmt
	UpdateGloabalRoleBase *sql.Stmt
	InsertRoleData        *sql.Stmt
	DeleteRoleData        *sql.Stmt
	UpdateRoleData        *sql.Stmt
}

func (ss *syncSQL) Init(db *sql.DB) {
	ss.InsertGloabalRoleBase = ss.Prepare(db, "INSERT INTO `global_role_base` SET `id`=?, `user_name`=?, `password`=?")
	ss.DeleteGloabalRoleBase = ss.Prepare(db, "DELETE FROM `global_role_base` WHERE `id`=?")
	ss.UpdateGloabalRoleBase = ss.Prepare(db, "UPDATE `global_role_base` SET `id`=?, `user_name`=?, `password`=? WHERE `id`=?")
	ss.InsertRoleData = ss.Prepare(db, "INSERT INTO `role_data` SET `id`=?, `role_id`=?, `diamond`=?")
	ss.DeleteRoleData = ss.Prepare(db, "DELETE FROM `role_data` WHERE `id`=?")
	ss.UpdateRoleData = ss.Prepare(db, "UPDATE `role_data` SET `id`=?, `role_id`=?, `diamond`=? WHERE `id`=?")
}
