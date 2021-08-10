package main

import (
	"log"
	"simple/api"
	"simple/api/role_api"
	"simple/lib/simpleapi"
)

func main() {
	app := simpleapi.New(
		simpleapi.SetReadBufSize(1024),
		simpleapi.SetMaxRecvSize(65536),
		simpleapi.SetMaxSendSize(65536),
	)

	api.RegisterApi(app)

	server, err := app.Listen("tcp", "0.0.0.0:0", nil)
	if err != nil {
		log.Fatal("setup server failed:", err)
	}
	defer server.Stop()

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

}
