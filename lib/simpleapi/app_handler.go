package simpleapi

import (
	"log"
	"runtime/debug"
	"simple/lib/simplenet"
)

type Handler interface {
	InitSession(*simplenet.Session) error
	Transaction(*simplenet.Session, Message, func())
}

type defaultHandler struct {
}

func (t *defaultHandler) InitSession(session *simplenet.Session) error {
	return nil
}

func (t *defaultHandler) Transaction(session *simplenet.Session, req Message, work func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("simpleapi: unhandled panic when processing '%s' - '%s'", req.Identity(), err)
			log.Println(string(debug.Stack()))
		}
	}()
	work()
}
