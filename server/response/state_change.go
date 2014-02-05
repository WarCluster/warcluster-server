package response

import (
	"warcluster/entities"
)

type StateChange struct {
	baseResponse
	Missions   map[string]*entities.Mission      `json:",omitempty"`
	RawPlanets map[string]*entities.Planet       `json:"-"`
	Planets    map[string]*entities.PlanetPacket `json:",omitempty"`
	Suns       map[string]*entities.Sun          `json:",omitempty"`
}

func NewStateChange() *StateChange {
	r := new(StateChange)
	r.Command = "state_change"
	return r
}

func (s *StateChange) Sanitize(player *entities.Player) {
	s.Planets = SanitizePlanets(player, s.RawPlanets)
}
