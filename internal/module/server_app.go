package module

import (
	"simple/internal/log"
	"simple/internal/mdb"
	"simple/pkg/simpleapi"
	"simple/pkg/simplenet"
	"time"
)

type serverApp struct {
	app      *simpleapi.App
	server   *simplenet.Server
	serverId int
	db       *mdb.Database
}

func NewServer() *serverApp {
	return &serverApp{}
}

func (app *serverApp) Start(network, address string, serverId int, db *mdb.Database, apiAPP *simpleapi.App) {
	app.serverId = serverId
	app.db = db
	app.app = apiAPP

	var err error
	app.server, err = app.app.Listen(network, address)
	if err != nil {
		panic(err)
	}

	go app.server.Serve()
}

func (app *serverApp) Stop() {
	if app.server != nil {
		app.server.Stop()
		log.Info("api stop.")
	}
}

func (app *serverApp) InitSession(session simplenet.ISession) error {
	newSessionState(session, app.serverId)
	return nil
}

func (app *serverApp) Transaction(session simplenet.ISession, msg simpleapi.Message, work func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("Session failed:", err)
		}
	}()
	var (
		beginTime time.Time
		endTime   time.Time
	)
	beginTime = time.Now()
	app.db.Transaction(msg, func() {
		work()
	})
	endTime = time.Now()
	log.Info("cost time:", endTime.Sub(beginTime))
}
