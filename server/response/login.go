package response

import (
	"github.com/Vladimiroff/vec2d"

	"warcluster/entities"
)

type Fraction struct {
	Id    uint16
	Color entities.Color
	Name  string
}

type LoginSuccess struct {
	baseResponse
	Username   string
	Position   *vec2d.Vector
	Fraction   Fraction
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
	r.Fraction = Fraction{player.Race.ID, player.Race.Color(), player.Race.Name()}
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
