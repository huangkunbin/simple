package mdb

import (
	"fmt"
	"runtime"
	"simple/internal/log"
	"sync"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var ConfigRWMutex *sync.RWMutex = &sync.RWMutex{}

type Database struct {
	mutex          sync.Mutex
	rowIds         rowIds
	increment      int64
	mySqlCommitter *mySqlCommitter
	fileCommitter  *fileCommitter
	globalTables   *globalTables
	roleTables     map[int64]*roleTables
	indexes        indexes
	*transLogs
}

func New(shardID int) *Database {
	db := &Database{
		transLogs: &transLogs{},
	}

	db.rowIds.Init(int64(shardID))
	db.indexes.Init()
	db.globalTables = NewGlobalTables()
	return db
}

func (db *Database) Start(connStr, syncDir string) {
	conn, err := sql.Open("mysql", connStr+"?autocommit=0&charset=utf8mb4&collation=utf8mb4_unicode_ci")
	if err != nil {
		panic(err)
	}

	db.mySqlCommitter = newmySqlCommitter(conn)
	db.fileCommitter = newfileCommitter(db.mySqlCommitter, syncDir)
	db.roleTables = make(map[int64]*roleTables)

	(&loader{db, conn, 0, 0, 0}).LoadGlobalTables()

	for roleId := range db.globalTables.GlobalRoleBase {
		db.NewRoleTables(roleId)
	}

	wg := &sync.WaitGroup{}
	workerNum := runtime.GOMAXPROCS(-1)
	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go func(workerId int) {
			startTime := time.Now()
			(&loader{db, conn, 0, workerNum, workerId}).LoadRoleTables()
			usedTime := time.Since(startTime)
			wg.Done()
			log.Infof(" db load data cost time:worker[%d]time[%v]", workerId, usedTime)
		}(i)
	}
	wg.Wait()
}

func (db *Database) Stop() {
	db.fileCommitter.Stop()
}

type loader struct {
	db        *Database
	conn      *sql.DB
	RoleId    int64
	workerNum int
	workerId  int
}

func (ld *loader) InitRowId(id *int64, q string) {
	var v sql.NullInt64
	err := ld.conn.QueryRow(q).Scan(&v)
	if err == sql.ErrNoRows || !v.Valid {
		return
	}
	*id = v.Int64
	if err != nil {
		panic(err)
	}
}

func (ld *loader) LoadGlobalTable(table string, id *int64) *sql.Rows {
	log.Info("load " + table)
	if id != nil && ld.workerId == 0 {
		ld.InitRowId(id, "SELECT MAX(`id`) FROM `"+table+"`")
	}
	t, err := ld.conn.Query("SELECT * FROM `" + table + "`")
	if err != nil {
		panic(err)
	}
	return t
}

func (ld *loader) LoadRoleTable(table, pk string, id *int64) *sql.Rows {
	sql := "SELECT * FROM `" + table + "`"
	if ld.RoleId == 0 {
		log.Info("load " + table)
		if id != nil && ld.workerId == 0 {
			ld.InitRowId(id, "SELECT MAX(`id`) FROM `"+table+"`")
		}
		sql += fmt.Sprintf(" WHERE %s %% %d = %d", pk, ld.workerNum, ld.workerId)
	} else {
		sql += fmt.Sprintf(" WHERE `role_id` = %d", ld.RoleId)
	}
	t, err := ld.conn.Query(sql)
	if err != nil {
		panic(err)
	}
	return t
}

type transLog interface {
	Commit(*sql.Tx, *syncSQL) error
	Rollback()
}

type transLogs []transLog

func (logs *transLogs) commit(db *Database, info interface{}) {
	logs2 := *logs
	if len(logs2) == 0 {
		return
	}

	logs3 := make([]transLog, len(logs2))
	copy(logs3, logs2)
	*logs = logs2[0:0]

	db.fileCommitter.Commit(info, logs3)
}

func (logs *transLogs) rollback() {
	logs2 := *logs
	if len(logs2) == 0 {
		return
	}

	for i := len(logs2) - 1; i >= 0; i-- {
		logs2[i].Rollback()
	}
	*logs = logs2[0:0]
}

func (logs *transLogs) addTransLog(log transLog) {
	*logs = append(*logs, log)
}

func (db *Database) Transaction(info interface{}, work func()) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	defer func() {
		if err := recover(); err == nil {
			db.transLogs.commit(db, info)
		} else {
			db.transLogs.rollback()
			panic(err)
		}
	}()
	work()
}

func (db *Database) NewRoleTables(roleId int64) *roleTables {
	tables := NewRoleTables()
	tables.RoleId = roleId
	db.roleTables[roleId] = tables
	return tables
}

func (db *Database) getRoleTables(roleId int64) *roleTables {
	if t, ok := db.roleTables[roleId]; ok {
		return t
	}
	return nil
}

func (db *Database) getOrCreateTables(roleId int64) *roleTables {
	if t, ok := db.roleTables[roleId]; ok {
		return t
	}
	return db.NewRoleTables(roleId)
}

func (db *Database) setRoleTables(tables *roleTables) {
	if _, exists := db.roleTables[tables.RoleId]; exists {
		panic("duplicate set Role tables")
	}
	db.roleTables[tables.RoleId] = tables
}

func (db *Database) delRoleTables(roleId int64) {
	if _, exists := db.roleTables[roleId]; !exists {
		panic("delete not exists Role tables")
	}
	delete(db.roleTables, roleId)
}

type RoleDB struct {
	db        *Database
	roleId    int64
	increment int64
	rowIds    *rowIds
	indexes   *indexes
	tables    *roleTables
	*transLogs
}

func (RoleDB *RoleDB) RoleId() int64 {
	return RoleDB.roleId
}

func (db *Database) GetRoleDB(roleId int64) *RoleDB {
	tables := db.getOrCreateTables(roleId)
	return &RoleDB{
		db:        db,
		roleId:    roleId,
		increment: db.increment,
		tables:    tables,
		rowIds:    &db.rowIds,
		indexes:   &db.indexes,
		transLogs: db.transLogs,
	}
}
