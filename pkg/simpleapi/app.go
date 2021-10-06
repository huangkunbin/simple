package simpleapi

import (
	"net"
	"simple/pkg/simplenet"
	"time"
)

type App struct {
	serviceTypes []*ServiceType
	services     [256]Provider

	ReadBufSize int
	MaxRecvSize int
	MaxSendSize int
	RecvTimeout time.Duration
	SendTimeout time.Duration
	manager     *simplenet.Manager
	handler     Handler
}

func New(opts ...Option) *App {
	app := &App{
		manager:     simplenet.NewManger(),
		ReadBufSize: 1024,
		MaxRecvSize: 64 * 1024,
		MaxSendSize: 64 * 1024,
		handler:     &defaultHandler{},
	}
	for _, opt := range opts {
		opt(app)
	}
	return app
}

func (app *App) Dial(network, address string) (*simplenet.Session, error) {
	return simplenet.Dial(network, address, simplenet.ProtocolFunc(app.newClientCodec))
}

func (app *App) Listen(network, address string) (*simplenet.Server, error) {
	listener, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	return app.NewServer(listener), nil
}

func (app *App) NewClient(conn net.Conn) simplenet.ISession {
	codec, _ := app.newClientCodec(conn)
	return app.manager.NewSession(codec)
}

func (app *App) NewServer(listener net.Listener) *simplenet.Server {
	return simplenet.NewServer(
		listener,
		simplenet.ProtocolFunc(app.newServerCodec),
		simplenet.HandlerFunc(func(session simplenet.ISession) {
			app.handleSessoin(session)
		}),
	)
}

func (app *App) handleSessoin(session simplenet.ISession) {
	defer session.Close()

	if app.handler.InitSession(session) != nil {
		return
	}

	for {
		msg, err := session.Receive()
		if err != nil {
			return
		}

		req := msg.(Message)
		app.handler.Transaction(session, req, func() {
			app.services[req.ServiceID()].(Service).HandleRequest(session, req)
		})
	}
}
