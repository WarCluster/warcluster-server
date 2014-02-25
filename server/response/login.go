package response

import (
	"github.com/Vladimiroff/vec2d"

	"warcluster/entities"
)

type LoginSuccess struct {
	baseResponse
	Username   string
	Position   *vec2d.Vector
	HomePlanet struct {
		Name     string
		Position *vec2d.Vector
	}
}

type LoginFailed struct {
	baseResponse
}

type LoginInformation struct {
	baseResponse
}

func NewLoginSuccess(player *entities.Player, homePlanet *entities.Planet) *LoginSuccess {
	r := new(LoginSuccess)
	r.Command = "login_success"
	r.Username = player.Username
	r.Position = player.ScreenPosition
	r.HomePlanet.Name = homePlanet.Name
	r.HomePlanet.Position = homePlanet.Position
	return r
}

func NewLoginFailed() *LoginFailed {
	r := new(LoginFailed)
	r.Command = "login_failed"
	return r
}

func NewLoginInformation() *LoginInformation {
	r := new(LoginInformation)
	r.Command = "request_setup_params"
	return r
}

func (l *LoginSuccess) Sanitize(*entities.Player)     {}
func (l *LoginFailed) Sanitize(*entities.Player)      {}
func (l *LoginInformation) Sanitize(*entities.Player) {}
