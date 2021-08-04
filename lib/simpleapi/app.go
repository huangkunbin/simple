package simpleapi

import (
	"log"
	"net"
	"runtime/debug"
	"simple/lib/mynet"
	"time"
)

type Handler interface {
	InitSession(*mynet.Session) error
	Transaction(*mynet.Session, Message, func())
}

type App struct {
	serviceTypes []*ServiceType
	services     [256]Provider

	ReadBufSize int
	MaxRecvSize int
	MaxSendSize int
	RecvTimeout time.Duration
	SendTimeout time.Duration
	manager     *mynet.Manager
}

func New() *App {
	return &App{
		manager:     mynet.NewManger(),
		ReadBufSize: 1024,
		MaxRecvSize: 64 * 1024,
		MaxSendSize: 64 * 1024,
	}
}

func (app *App) handleSessoin(session *mynet.Session, handler Handler) {
	defer session.Close()

	if handler.InitSession(session) != nil {
		return
	}

	for {
		msg, err := session.Receive()
		if err != nil {
			return
		}

		req := msg.(Message)
		handler.Transaction(session, req, func() {
			app.services[req.ServiceID()].(Service).HandleRequest(session, req)
		})
	}
}

func (app *App) Dial(network, address string) (*mynet.Session, error) {
	return mynet.Dial(network, address, mynet.ProtocolFunc(app.newClientCodec))
}

func (app *App) Listen(network, address string, handler Handler) (*mynet.Server, error) {
	listener, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	return app.NewServer(listener, handler), nil
}

func (app *App) NewClient(conn net.Conn) *mynet.Session {
	codec, _ := app.newClientCodec(conn)
	return app.manager.NewSession(codec)
}

func (app *App) NewServer(listener net.Listener, handler Handler) *mynet.Server {
	if handler == nil {
		handler = &noHandler{}
	}
	return mynet.NewServer(
		listener, mynet.ProtocolFunc(app.newServerCodec),
		mynet.HandlerFunc(func(session *mynet.Session) {
			app.handleSessoin(session, handler)
		}),
	)
}

type noHandler struct {
}

func (t *noHandler) DropSession(session *mynet.Session) {
}

func (t *noHandler) InitSession(session *mynet.Session) error {
	return nil
}

func (t *noHandler) Transaction(session *mynet.Session, req Message, work func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("simpleapi: unhandled panic when processing '%s' - '%s'", req.Identity(), err)
			log.Println(string(debug.Stack()))
		}
	}()
	work()
}
