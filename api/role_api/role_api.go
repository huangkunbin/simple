package role_api

import (
	"simple/internal/module"
	"simple/pkg/simpleapi"
	"simple/pkg/simplenet"
)

type RoleApi struct {
}

func (r *RoleApi) APIs() simpleapi.APIs {
	return simpleapi.APIs{
		1: {LoginReq{}, LoginRes{}},
	}
}

func (r *RoleApi) Login(session *simplenet.Session, req *LoginReq) *LoginRes {
	userName := module.Role.Login(req.UserName, req.Password)
	return &LoginRes{
		UserName: userName,
	}
}
