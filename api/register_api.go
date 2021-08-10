package api

import (
	"simple/api/role_api"
	"simple/lib/simpleapi"
)

func RegisterApi(app *simpleapi.App) {
	app.Register(1, &role_api.RoleApi{})
}
