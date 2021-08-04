package main

import (
	"flag"
	"log"
	"simple/api/role_api"
	"simple/lib/simpleapi"
)

func main() {
	gencode := flag.Bool("gencode", false, "generate code")
	genpath := flag.String("genpath", "", "generate path")
	flag.Parse()

	if *gencode && *genpath != "" {
		simpleapi.GenCode(*genpath, app)
		return
	}

	server, err := app.Listen("tcp", "0.0.0.0:0", nil)
	if err != nil {
		log.Fatal("setup server failed:", err)
	}
	go server.Serve()

	client, err := app.Dial("tcp", server.Listener().Addr().String())
	if err != nil {
		log.Fatal("setup client failed:", err)
	}

	for i := 0; i < 10; i++ {
		err := client.Send(&role_api.LoginReq{
			UserName: "hkb",
			Password: "123456",
		})
		if err != nil {
			log.Fatal("send failed:", err)
		}

		rsp, err := client.Receive()
		if err != nil {
			log.Fatal("recv failed:", err)
		}

		log.Printf("AddRsp: %s", rsp.(*role_api.LoginRes).String())
	}

	server.Stop()
}
