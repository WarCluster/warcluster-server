package response

import (
	"fmt"

	"warcluster/entities"
)

type ServerParams struct {
	baseResponse
	HomeSPM            float64 //ships per minute
	PlanetsSPM         map[string]float64
	ShipsDeathModifier float64
	Races              map[string]entities.Race
}

func NewServerParams() *ServerParams {
	var planetSizeIdx int8

	r := new(ServerParams)
	r.Races = make(map[string]entities.Race)
	r.PlanetsSPM = make(map[string]float64)

	r.Command = "server_params"

	for _, race := range entities.Races {
		r.Races[race.Name] = race
	}
	r.HomeSPM = 60 / float64(entities.ShipCountTimeMod(1, true))
	for planetSizeIdx = 1; planetSizeIdx <= 10; planetSizeIdx++ {
		planetSPM := float64(entities.ShipCountTimeMod(planetSizeIdx, false))
		r.PlanetsSPM[fmt.Sprintf("%v", planetSizeIdx)] = 60 / planetSPM
	}
	r.ShipsDeathModifier = entities.Settings.ShipsDeathModifier
	return r
}

func (_ *ServerParams) Sanitize(*entities.Player) {}
