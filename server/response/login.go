package response

import "github.com/Vladimiroff/vec2d"

type LoginSuccess struct {
	baseResponse
	Username string
	Position *vec2d.Vector
}

type LoginFailed struct {
	baseResponse
}

func NewLoginSuccess() *LoginSuccess {
	r := new(LoginSuccess)
	r.Command = "login_success"
	return r
}

func NewLoginFailed() *LoginFailed {
	r := new(LoginFailed)
	r.Command = "login_failed"
	return r
}
