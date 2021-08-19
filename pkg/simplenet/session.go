package simplenet

import (
	"context"
	"sync"
	"sync/atomic"
)

type Session struct {
	ctx       context.Context
	id        uint64
	codec     Codec
	manager   *Manager
	recvMutex sync.Mutex
	sendMutex sync.RWMutex
	closeFlag int32

	State interface{}
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

func (session *Session) Get(key interface{}) interface{} {
	return session.ctx.Value(key)
}

func (session *Session) Set(key interface{}, val interface{}) {
	session.ctx = context.WithValue(session.ctx, key, val)
}
