package api

import (
	"simple/api/role_api"
	"simple/pkg/simpleapi"
)

func RegisterApi(app *simpleapi.App) {
	app.Register(1, &role_api.RoleApi{})
}
