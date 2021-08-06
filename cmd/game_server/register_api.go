package main

import (
	"simple/api/role_api"
	"simple/lib/simpleapi"
)

func registerApi(app *simpleapi.App) {
	app.Register(1, &role_api.RoleApi{})
}
