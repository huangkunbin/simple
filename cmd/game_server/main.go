package main

import (
	"context"
	"flag"

	"os"
	"os/signal"
	"simple/api"
	"simple/internal/config"
	"simple/internal/log"
	"simple/internal/mdb"
	"simple/pkg/simpleapi"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gencode := flag.Bool("gencode", false, "generate code")
	genpath := flag.String("genpath", "", "generate path")
	flag.Parse()

	app := simpleapi.New(
		simpleapi.SetReadBufSize(1024),
		simpleapi.SetMaxRecvSize(65536),
		simpleapi.SetMaxSendSize(65536),
	)

	api.RegisterApi(app)

	if *gencode && *genpath != "" {
		simpleapi.GenCode(*genpath, app)
		return
	}

	config.InitConfig()

	log.UpdateLoggers(config.Cfg.LoggerLevel, config.Cfg.LoggerType)

	db := mdb.New(config.Cfg.ServerId)
	db.Start(config.Cfg.DB, config.Cfg.SyncFileDir)

	server, err := app.Listen("tcp", "0.0.0.0:0", nil)
	if err != nil {
		log.Error("setup server failed:", err)
	}
	defer server.Stop()
	go server.Serve()

	waitNotify(ctx)
}

func waitNotify(ctx context.Context) {
	sigTERM := make(chan os.Signal, 1)
	signal.Notify(sigTERM, os.Interrupt, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Info("Done")
	case <-sigTERM:
		log.Info("killed")
	}
}
