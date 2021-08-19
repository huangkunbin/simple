package mdb

import (
	"database/sql"
	"sync"
)

type fileCommitter struct {
	lock           sync.RWMutex
	mySqlCommitter *mySqlCommitter
}

func newfileCommitter(mySqlCommitter *mySqlCommitter, syncDir string) *fileCommitter {
	fc := &fileCommitter{
		mySqlCommitter: mySqlCommitter,
	}
	return fc
}

func (fc *fileCommitter) Commit(info interface{}, trans []transLog) {
	fc.lock.RLock()
	defer fc.lock.RUnlock()
	fc.mySqlCommitter.Commit(trans)
}

func (fc *fileCommitter) Stop() {
	fc.lock.Lock()
	defer fc.lock.Unlock()
	fc.mySqlCommitter.Stop()
}

func (ss *syncSQL) Prepare(db *sql.DB, sql string) *sql.Stmt {
	stmt, err := db.Prepare(sql)
	if err != nil {
		panic(err)
	}
	return stmt
}

type mySqlCommitter struct {
	commitConn   *sql.DB
	commitChan   chan []transLog
	waitStopChan chan int
	syncSQL      syncSQL
}

func newmySqlCommitter(conn *sql.DB) *mySqlCommitter {
	mc := &mySqlCommitter{
		commitConn:   conn,
		commitChan:   make(chan []transLog, 50000),
		waitStopChan: make(chan int),
	}
	mc.syncSQL.Init(conn)
	go mc.commitLoop()
	return mc
}

func (mc *mySqlCommitter) commitLoop() {
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
	}
	close(mc.waitStopChan)
}

func (mc *mySqlCommitter) Stop() {
	close(mc.commitChan)
	<-mc.waitStopChan
	mc.commitConn.Close()
}

func (mc *mySqlCommitter) Commit(trans []transLog) {
	mc.commitChan <- trans
}

func (mc *mySqlCommitter) QueueLength() int {
	return len(mc.commitChan)
}
