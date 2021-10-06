package module

import (
	"simple/internal/mdb"
	"simple/pkg/simplenet"
)

type SessionState struct {
	ServerId int
	RoleId   int64
	Database *mdb.RoleDB
}

func newSessionState(session simplenet.ISession, serverId int) *SessionState {
	state := &SessionState{
		ServerId: serverId,
	}
	session.SetState(state)
	return state
}

func State(session simplenet.ISession) *SessionState {
	return session.GetState().(*SessionState)
}

func deleteSessionState(session simplenet.ISession) {

}
