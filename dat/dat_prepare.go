package dat

import "database/sql"

type syncSQL struct {
	InsertGlobalRole *sql.Stmt
	DeleteGlobalRole *sql.Stmt
	UpdateGlobalRole *sql.Stmt
}

func (ss *syncSQL) Init(db *sql.DB) {
	ss.InsertGlobalRole = ss.Prepare(db, "INSERT INTO `global_role` SET `id`=?, `user`=?, `group`=?, `nick`=?, `create_time`=?, `max_distance`=?, `max_distance_time`=?, `best_alive_time`=?, `best_golds`=?, `best_group`=?, `last_distance`=?")
	ss.DeleteGlobalRole = ss.Prepare(db, "DELETE FROM `global_role` WHERE `id`=?")
	ss.UpdateGlobalRole = ss.Prepare(db, "UPDATE `global_role` SET `id`=?, `user`=?, `group`=?, `nick`=?, `create_time`=?, `max_distance`=?, `max_distance_time`=?, `best_alive_time`=?, `best_golds`=?, `best_group`=?, `last_distance`=? WHERE `id`=?")
}
