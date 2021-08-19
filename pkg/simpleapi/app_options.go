package simpleapi

import (
	"time"
)

type Option func(options *App)

func SetReadBufSize(sz int) Option {
	return func(options *App) {
		options.ReadBufSize = sz
	}
}

func SetMaxRecvSize(sz int) Option {
	return func(options *App) {
		options.MaxRecvSize = sz
	}
}

func SetMaxSendSize(sz int) Option {
	return func(options *App) {
		options.MaxRecvSize = sz
	}
}

func SetRecvTimeout(t time.Duration) Option {
	return func(options *App) {
		options.RecvTimeout = t
	}
}

func SetSendTimeout(t time.Duration) Option {
	return func(options *App) {
		options.SendTimeout = t
	}
}
