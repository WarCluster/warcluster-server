package response

import (
	"fmt"
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

//SPM ShipsPerMinute
type ServerParams struct {
	baseResponse
	HomeSPM    float64
	PlanetsSPM map[string]float64
	Teams      map[string]entities.Color
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

func NewServerParams() *ServerParams {
	var planetSizeIdx int8

	r := new(ServerParams)
	r.Teams = make(map[string]entities.Color)
	r.PlanetsSPM = make(map[string]float64)

	r.Command = "server_params"
	for raceIdx := 0; raceIdx < entities.RACE_VARIATION_CNT; raceIdx++ {
		race, _ := entities.AssignRace(uint16(raceIdx))
		r.Teams[race.Name()] = race.Color()
	}
	r.HomeSPM = 60 / float64(entities.ShipCountTimeMod(1, true))
	for planetSizeIdx = 1; planetSizeIdx <= 10; planetSizeIdx++ {
		r.PlanetsSPM[fmt.Sprintf("%v", planetSizeIdx)] = 60 / float64(entities.ShipCountTimeMod(planetSizeIdx, false))
	}
	return r
}

func (l *LoginSuccess) Sanitize(*entities.Player)     {}
func (l *LoginFailed) Sanitize(*entities.Player)      {}
func (l *LoginInformation) Sanitize(*entities.Player) {}
func (l *ServerParams) Sanitize(*entities.Player)     {}
