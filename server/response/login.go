package response

type LoginSuccess struct {
	baseResponse
	Username string
	Position []int
}

type LoginFail struct {
	baseResponse
}

func NewLoginSuccess() *LoginSuccess {
	l := new(LoginSuccess)
	l.Command = "login_success"
	return l
}

func NewLoginFail() *LoginFail {
	l := new(LoginFail)
	l.Command = "login_fail"
	return l
}
