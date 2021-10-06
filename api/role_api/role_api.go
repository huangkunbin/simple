package role_api

import (
	"simple/internal/module"
	"simple/pkg/simpleapi"
	"simple/pkg/simplenet"
)

type RoleApi struct {
}

func (api *RoleApi) APIs() simpleapi.APIs {
	return simpleapi.APIs{
		1: {LoginReq{}, LoginRes{}},
	}
}

func (api *RoleApi) Login(session simplenet.ISession, req *LoginReq) *LoginRes {
	userName := module.Role.Login(session, req.UserName, req.Password)
	return &LoginRes{
		UserName: userName,
	}
}
