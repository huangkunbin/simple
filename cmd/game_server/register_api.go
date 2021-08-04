package main

import (
	"simple/api/role_api"
	"simple/lib/simpleapi"
)

var app = simpleapi.New()

func init() {
	app.Register(1, &role_api.RoleApi{})
}
