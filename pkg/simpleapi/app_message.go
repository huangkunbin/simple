package simpleapi

import (
	"simple/pkg/simplenet"
)

type Service interface {
	ServiceID() byte
	NewRequest(byte) Message
	NewResponse(byte) Message
	HandleRequest(*simplenet.Session, Message)
}

type Message interface {
	ServiceID() byte
	MessageID() byte
	Identity() string
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}
