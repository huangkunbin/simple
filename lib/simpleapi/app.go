package simpleapi

import (
	"net"
	"simple/lib/mynet"
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
	manager     *mynet.Manager
}

func New(opts ...Option) *App {
	app := &App{
		manager:     mynet.NewManger(),
		ReadBufSize: 1024,
		MaxRecvSize: 64 * 1024,
		MaxSendSize: 64 * 1024,
	}
	for _, opt := range opts {
		opt(app)
	}
	return app
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
		handler = &defaultHandler{}
	}
	return mynet.NewServer(
		listener,
		mynet.ProtocolFunc(app.newServerCodec),
		mynet.HandlerFunc(func(session *mynet.Session) {
			app.handleSessoin(session, handler)
		}),
	)
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
