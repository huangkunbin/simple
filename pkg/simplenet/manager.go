package simplenet

import (
	"context"
	"sync"
	"sync/atomic"
)

type Manager struct {
	sync.RWMutex
	sessions map[uint64]*Session
}

func NewManger() *Manager {
	return &Manager{
		sessions: map[uint64]*Session{},
	}
}

func (smap *Manager) NewSession(codec Codec) ISession {
	smap.Lock()
	defer smap.Unlock()
	session := NewSession(codec, smap)
	smap.sessions[session.id] = session
	return session
}

func NewSession(codec Codec, smap *Manager) *Session {
	session := &Session{
		codec:   codec,
		manager: smap,
		id:      atomic.AddUint64(&globalSessionId, 1),
	}
	session.Context = context.TODO()
	return session
}

func (smap *Manager) GetSession(sessionID uint64) *Session {
	smap.RLock()
	defer smap.RUnlock()
	return smap.sessions[sessionID]
}
