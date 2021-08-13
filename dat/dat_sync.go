package dat

import (
	"database/sql"
	"sync"
)

type fileCommitter struct {
	lock           sync.RWMutex
	mysqlCommitter *mysqlCommitter
}

func newFileCommitter(mysqlCommitter *mysqlCommitter, syncDir string) *fileCommitter {
	fc := &fileCommitter{
		mysqlCommitter: mysqlCommitter,
	}
	return fc
}

func (fc *fileCommitter) Commit(info interface{}, trans []transLog) {
	fc.lock.RLock()
	defer fc.lock.RUnlock()
	fc.mysqlCommitter.Commit(trans)
}

func (fc *fileCommitter) Stop() {
	fc.lock.Lock()
	defer fc.lock.Unlock()
	fc.mysqlCommitter.Stop()
}

func (ss *syncSQL) Prepare(db *sql.DB, sql string) *sql.Stmt {
	stmt, err := db.Prepare(sql)
	if err != nil {
		panic(err)
	}
	return stmt
}

type mysqlCommitter struct {
	commitConn   *sql.DB
	commitChan   chan []transLog
	waitStopChan chan int
	syncSQL      syncSQL
}

func newMySqlCommitter(conn *sql.DB) *mysqlCommitter {
	mc := &mysqlCommitter{
		commitConn:   conn,
		commitChan:   make(chan []transLog, 50000),
		waitStopChan: make(chan int),
	}
	mc.syncSQL.Init(conn)
	go mc.commitLoop()
	return mc
}

func (mc *mysqlCommitter) commitLoop() {
	for transLog := range mc.commitChan {
		tx, err := mc.commitConn.Begin()
		if err != nil {
			panic(err)
		}
		for _, item := range transLog {
			err = item.Commit(tx, &mc.syncSQL)
			if err != nil {
				break
			}
		}
		if err != nil {
			panic(err)
		}
		err = tx.Commit()
		if err != nil {
			panic(err)
		}
		for _, item := range transLog {
			item.Free()
		}
	}
	close(mc.waitStopChan)
}

func (mc *mysqlCommitter) Stop() {
	close(mc.commitChan)
	<-mc.waitStopChan
	mc.commitConn.Close()
}

func (mc *mysqlCommitter) Commit(trans []transLog) {
	mc.commitChan <- trans
}

func (mc *mysqlCommitter) QueueLength() int {
	return len(mc.commitChan)
}
