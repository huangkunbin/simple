package main

import (
	"log"
	"simple/api"
	"simple/api/role_api"
	"simple/pkg/simpleapi"
)

func main() {
	app := simpleapi.New(
		simpleapi.SetReadBufSize(1024),
		simpleapi.SetMaxRecvSize(65536),
		simpleapi.SetMaxSendSize(65536),
	)

	api.RegisterApi(app)

	client, err := app.Dial("tcp", "127.0.0.1:5234")
	if err != nil {
		log.Fatal("setup client failed:", err)
	}

	err = client.Send(&role_api.LoginReq{
		UserName: "hkb1",
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
