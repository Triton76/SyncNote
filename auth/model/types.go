package model

import "context"

type AuthStore interface {
	Login(ctx context.Context, loginId string, password string) (*LoginResp, error)
	Register(ctx context.Context, registerreq *RegisterReq) (*RegisterResp, error)
}

type LoginResp struct {
	Code   int64
	UserId string
	Token  string
	Expire int64
}

type RegisterReq struct {
	Email    string
	Captcha  string
	Password string
	Username string
}

type RegisterResp struct {
	Code   int64
	UserId string
	Token  string
	Expire int64
}
