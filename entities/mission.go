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

func (self *Mission) GetSpeed() int {
	return 5
}

func (self Mission) Serialize() (string, []byte, error) {
	result, err := json.Marshal(self)
	if err != nil {
		return self.GetKey(), nil, err
	}
	return self.GetKey(), result, nil
}

func (self Mission) GetStartPlanet() string {
	return self.start_planet
}

func (self Mission) GetStartTime() time.Time {
	return self.start_time
}

func EndMission(endPlanet Planet, missionInfo Mission) Planet {
	if endPlanet.Owner == missionInfo.Player {
		endPlanet.SetShipCount(endPlanet.GetShipCount() + missionInfo.ShipCount)
	} else {
		if missionInfo.ShipCount < endPlanet.GetShipCount() {
			endPlanet.SetShipCount(endPlanet.GetShipCount() - missionInfo.ShipCount)
		} else {
			endPlanet.SetShipCount(missionInfo.ShipCount - endPlanet.GetShipCount())
			endPlanet.Owner = missionInfo.Player
		}
	}
	return endPlanet
}
