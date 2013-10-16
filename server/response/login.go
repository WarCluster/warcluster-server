package response

type LoginSuccess struct {
	BaseResponse
	Username string
	Position []int
}

type LoginFail struct {
	BaseResponse
}

func NewLoginSuccess() *LoginSuccess {
	r := new(LoginSuccess)
	r.Command = "login_success"
	return r
}

func NewLoginFail() *LoginFail {
	r := new(LoginFail)
	r.Command = "login_fail"
	return r
}
