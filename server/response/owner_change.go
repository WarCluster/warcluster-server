package response

import "warcluster/entities"

type OwnerChange struct {
	baseResponse
	RawPlanet map[string]*entities.Planet       `json:"-"`
	Planet    map[string]*entities.PlanetPacket `json:",omitempty"`
}

func NewOwnerChange() *OwnerChange {
	r := new(OwnerChange)
	r.Command = "owner_change"
	return r
}

func (o *OwnerChange) Sanitize(player *entities.Player) {
	o.Planet = SanitizePlanets(player, o.RawPlanet)
}
