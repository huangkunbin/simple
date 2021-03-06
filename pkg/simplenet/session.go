package simplenet

import (
	"context"
	"sync"
	"sync/atomic"
)

type ISession interface {
	context.Context
	Codec
	IsClosed() bool
	SetState(interface{})
	GetState() interface{}
}

type Session struct {
	context.Context
	id        uint64
	codec     Codec
	manager   *Manager
	recvMutex sync.Mutex
	sendMutex sync.RWMutex
	closeFlag int32

	state interface{}
}

func (session *Session) Receive() (interface{}, error) {
	session.recvMutex.Lock()
	defer session.recvMutex.Unlock()
	msg, err := session.codec.Receive()
	if err != nil {
		session.Close()
	}
	return msg, err
}

func (session *Session) Send(msg interface{}) error {
	if session.IsClosed() {
		return nil
	}
	session.sendMutex.Lock()
	defer session.sendMutex.Unlock()
	err := session.codec.Send(msg)
	if err != nil {
		session.Close()
	}
	return err
}

func (session *Session) Close() error {
	if atomic.CompareAndSwapInt32(&session.closeFlag, 0, 1) {
		err := session.codec.Close()
		go func() {
			if session.manager != nil {
				delete(session.manager.sessions, session.id)
			}
		}()
		return err
	}
	return nil
}

func (session *Session) IsClosed() bool {
	return atomic.LoadInt32(&session.closeFlag) == 1
}

func (session *Session) SetState(state interface{}) {
	session.state = state
}

func (session *Session) GetState() interface{} {
	return session.state
}
