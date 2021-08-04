package mynet

import (
	"errors"
	"io"
	"net"
)

var globalSessionId uint64

type Server struct {
	listener net.Listener
	protocol Protocol
	handler  Handler
	manager  *Manager
}

type Handler interface {
	HandleSession(*Session)
}

var _ Handler = HandlerFunc(nil)

type HandlerFunc func(*Session)

func (f HandlerFunc) HandleSession(session *Session) {
	f(session)
}

func NewServer(listener net.Listener, protocol Protocol, handler Handler) *Server {
	return &Server{
		manager: &Manager{
			sessions: map[uint64]*Session{},
		},
		listener: listener,
		protocol: protocol,
		handler:  handler,
	}
}

func Listen(network, address string, protocol Protocol, handler Handler) (*Server, error) {
	listener, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	return NewServer(listener, protocol, handler), nil
}

func Dial(network, address string, protocol Protocol) (*Session, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	codec, err := protocol.NewCodec(conn)
	if err != nil {
		return nil, err
	}
	return NewSession(codec, nil), nil
}

func (server *Server) Listener() net.Listener {
	return server.listener
}

func (server *Server) Serve() error {
	for {
		conn, err := Accept(server.listener)
		if err != nil {
			return err
		}
		go func() {
			codec, err := server.protocol.NewCodec(conn)
			if err != nil {
				conn.Close()
				return
			}
			session := server.manager.NewSession(codec)
			server.handler.HandleSession(session)
		}()
	}
}

func Accept(listener net.Listener) (net.Conn, error) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil, io.EOF
			}
			return nil, err
		} else {
			return conn, nil
		}
	}
}

func (server *Server) GetSession(sessionID uint64) *Session {
	return server.manager.GetSession(sessionID)
}

func (server *Server) Stop() {
	server.listener.Close()
}
