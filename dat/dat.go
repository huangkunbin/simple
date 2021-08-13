package dat

import (
	"fmt"
	"runtime"
	"sync"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var ConfigRWMutex *sync.RWMutex = &sync.RWMutex{}

type Database struct {
	mutex           sync.Mutex
	gRowIds         rowIds
	gIncrement      int64
	gMySqlCommitter *mysqlCommitter
	gFileCommitter  *fileCommitter
	gGlobalTables   *GlobalTables
	gRoleTables     map[int64]*RoleTables
	gIndexes        indexes
	*transLogs
}

func New(shardID int) *Database {
	db := &Database{
		transLogs: &transLogs{},
	}

	// 新数据库初始化ID
	db.gRowIds.Init(int64(shardID))
	db.gIndexes.Init()
	db.gGlobalTables = NewGlobalTables()
	return db
}

func (db *Database) Start(connStr, syncDir string) {
	conn, err1 := sql.Open("mysql", connStr+"?autocommit=0&charset=utf8mb4&collation=utf8mb4_unicode_ci")
	if err1 != nil {
		panic(err1)
	}

	db.gMySqlCommitter = newMySqlCommitter(conn)
	db.gFileCommitter = newFileCommitter(db.gMySqlCommitter, syncDir)
	db.gRoleTables = make(map[int64]*RoleTables)

	// 加载全局表数据
	(&loader{db, conn, 0, 0, 0}).LoadGlobalTables()

	// 根据玩家表数据初始化数据库切片
	// for row := db.gGlobalTables.GlobalRole; row != nil; row = row.next[0] {
	// 	tables := NewRoleTables()
	// 	tables.Pid = row.Id
	// 	db.setRoleTables(tables)
	// }

	// 并行加载玩家数据
	wg := new(sync.WaitGroup)
	workerNum := runtime.GOMAXPROCS(-1)
	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go func(workerId int) {
			// startTime := time.Now()
			(&loader{db, conn, 0, workerNum, workerId}).LoadRoleTables()
			// usedTime := time.Since(startTime)

			wg.Done()
			// log.Info("mdb load data cost time", log.M{
			// 	"worker": workerId,
			// 	"time":   usedTime,
			// })
		}(i)
	}
	wg.Wait()
}

func (db *Database) Stop() {
	db.gFileCommitter.Stop()
}

type loader struct {
	db        *Database
	conn      *sql.DB
	RoleId    int64
	workerNum int
	workerId  int
}

func (ld *loader) LoadGlobalTables() {
}

func (ld *loader) LoadRoleTables() {
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

func (ld *loader) LoadGlobalTable(table, pk string, id *int64) *sql.Rows {
	// log.Info("load " + table)
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
		// log.Info("load " + table)
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
	InvokeCallback(db *Database)
	Commit(*sql.Tx, *syncSQL) error
	Rollback()
	Free()
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

	for _, log := range logs3 {
		log.InvokeCallback(db)
	}
	db.gFileCommitter.Commit(info, logs3)
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

func (db *Database) newRoleTables(pid int64) *RoleTables {
	tables := NewRoleTables()
	tables.Pid = pid
	db.gRoleTables[pid] = tables
	return tables
}

func (db *Database) getRoleTables(pid int64) *RoleTables {
	if t, ok := db.gRoleTables[pid]; ok {
		return t
	}
	return nil
}

func (db *Database) setRoleTables(tables *RoleTables) {
	if _, exists := db.gRoleTables[tables.Pid]; exists {
		panic("duplicate set Role tables")
	}
	db.gRoleTables[int64(tables.Pid)] = tables
}

func (db *Database) delRoleTables(pid int64) {
	if _, exists := db.gRoleTables[pid]; !exists {
		panic("delete not exists Role tables")
	}
	delete(db.gRoleTables, pid)
}

type RoleDB struct {
	pid        int64
	gIncrement int64
	gRowIds    *rowIds
	gIndexes   *indexes
	tables     *RoleTables
	*transLogs
}

func (RoleDB *RoleDB) RoleId() int64 {
	return RoleDB.pid
}

func (db *Database) GetRoleDB(pid int64) *RoleDB {
	tables := db.getRoleTables(pid)
	if tables == nil {
		return nil
	}
	return &RoleDB{
		pid:        pid,
		gIncrement: db.gIncrement,
		tables:     tables,
		gRowIds:    &db.gRowIds,
		gIndexes:   &db.gIndexes,
		transLogs:  db.transLogs,
	}
}
