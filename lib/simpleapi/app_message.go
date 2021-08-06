package simpleapi

import (
	"simple/lib/mynet"
)

type Service interface {
	ServiceID() byte
	NewRequest(byte) Message
	NewResponse(byte) Message
	HandleRequest(*mynet.Session, Message)
}

type Message interface {
	ServiceID() byte
	MessageID() byte
	Identity() string
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}
