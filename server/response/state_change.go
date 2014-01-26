package response

import (
	"warcluster/entities"
)

type StateChange struct {
	baseResponse
	Missions map[string]entities.Entity `json:",omitempty"`
	Planets  map[string]entities.Entity `json:",omitempty"`
	Suns     map[string]entities.Entity `json:",omitempty"`
}

func NewStateChange() *StateChange {
	r := new(StateChange)
	r.Command = "state_change"
	return r
}

func (s *StateChange) Send(player *entities.Player, sender func([]byte)) error {
	for name, entity := range s.Planets {
		if planet, ok := entity.(*entities.Planet); ok {
			s.Planets[name] = planet.Sanitize(player)
		}
	}
	return Send(s, sender)
}
