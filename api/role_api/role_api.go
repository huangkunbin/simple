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
		2: {CreateReq{}, CreateRes{}},
	}
}

func (r *RoleApi) Login(session simplenet.ISession, req *LoginReq) *LoginRes {
	userName := module.Role.Login(req.UserName, req.Password)
	return &LoginRes{
		UserName: userName,
	}
}

func (r *RoleApi) Create(session simplenet.ISession, req *CreateReq) *CreateRes {
	userName := module.Role.Create(req.UserName, req.Password)
	return &CreateRes{
		UserName: userName,
	}
}
