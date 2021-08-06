package simpleapi

import (
	"log"
	"runtime/debug"
	"simple/lib/mynet"
)

type Handler interface {
	InitSession(*mynet.Session) error
	Transaction(*mynet.Session, Message, func())
}

type defaultHandler struct {
}

func (t *defaultHandler) InitSession(session *mynet.Session) error {
	return nil
}

func (t *defaultHandler) Transaction(session *mynet.Session, req Message, work func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("simpleapi: unhandled panic when processing '%s' - '%s'", req.Identity(), err)
			log.Println(string(debug.Stack()))
		}
	}()
	work()
}
