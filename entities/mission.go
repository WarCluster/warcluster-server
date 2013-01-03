package entities

import (
	"encoding/json"
	"fmt"
	"time"
)

type Mission struct {
	start_planet string
	start_time   time.Time
	Player       string
	ShipCount    int
	EndPlanet    string
}

func (self Mission) String() string {
	return self.GetKey()
}

func (self Mission) GetKey() string {
	start_planet_coords := ExtractPlanetCoords(self.start_planet)
	return fmt.Sprintf(
		"mission.%d_%d_%d",
		self.start_time.Unix(),
		start_planet_coords[0],
		start_planet_coords[1],
	)
}

func (self Mission) Serialize() (string, []byte, error) {
	result, err := json.Marshal(self)
	if err != nil {
		return self.GetKey(), nil, err
	}
	return self.GetKey(), result, nil
}

func EndMission(endPlanet Planet, missionInfo Mission) Planet {
	if endPlanet.Owner == missionInfo.Player {
		endPlanet.ShipCount += missionInfo.ShipCount
	} else {
		if missionInfo.ShipCount < endPlanet.ShipCount {
			endPlanet.ShipCount -= missionInfo.ShipCount
		} else {
			endPlanet.ShipCount = missionInfo.ShipCount - endPlanet.ShipCount
			endPlanet.Owner = missionInfo.Player
		}
	}
	return endPlanet
}
