package simpleapi

import (
	"log"
	"runtime/debug"
	"simple/pkg/simplenet"
)

type Handler interface {
	InitSession(simplenet.ISession) error
	Transaction(simplenet.ISession, Message, func())
}

type defaultHandler struct {
}

func (t *defaultHandler) InitSession(session simplenet.ISession) error {
	return nil
}

func (t *defaultHandler) Transaction(session simplenet.ISession, req Message, work func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("simpleapi: unhandled panic when processing '%s' - '%s'", req.Identity(), err)
			log.Println(string(debug.Stack()))
		}
	}()
	work()
}
