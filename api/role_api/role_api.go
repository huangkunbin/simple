package role_api

import (
	"simple/biz"
	"simple/lib/simpleapi"
	"simple/lib/simplenet"
)

type RoleApi struct {
}

func (r *RoleApi) APIs() simpleapi.APIs {
	return simpleapi.APIs{
		1: {LoginReq{}, LoginRes{}},
	}
}

func (r *RoleApi) Login(session *simplenet.Session, req *LoginReq) *LoginRes {
	userName := biz.Role.Login(req.UserName, req.Password)
	return &LoginRes{
		UserName: userName,
	}
}
