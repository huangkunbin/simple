package role_api

import (
	"simple/lib/mynet"
	"simple/lib/simpleapi"
)

type RoleApi struct {
}

func (r *RoleApi) APIs() simpleapi.APIs {
	return simpleapi.APIs{
		1: {LoginReq{}, LoginRes{}},
	}
}

func (r *RoleApi) Login(session *mynet.Session, req *LoginReq) *LoginRes {
	return &LoginRes{
		UserName: req.UserName,
	}
}
