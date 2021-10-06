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
	"simple/internal/module"
	"simple/pkg/simpleapi"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gencode := flag.Bool("gencode", false, "generate code")
	genpath := flag.String("genpath", "", "generate path")
	flag.Parse()

	server := module.NewServer()

	app := simpleapi.New(
		simpleapi.SetReadBufSize(1024),
		simpleapi.SetMaxRecvSize(65536),
		simpleapi.SetMaxSendSize(65536),
		simpleapi.SetHandler(server),
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

	module.InitModule(db)

	server.Start("tcp", "0.0.0.0:5234", 1, db, app)
	defer server.Stop()

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
