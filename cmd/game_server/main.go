package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"simple/api"
	"simple/lib/simpleapi"
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

	server, err := app.Listen("tcp", "0.0.0.0:0", nil)
	if err != nil {
		log.Fatal("setup server failed:", err)
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
		log.Print("Done")
	case <-sigTERM:
		log.Print("killed")
	}
}
