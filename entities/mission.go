package entities

import (
	"encoding/json"
	"fmt"
	"log"
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

func (self Mission) Serialize() (string, []byte) {
	result, err := json.Marshal(self)
	if err != nil {
		log.Fatal(err)
	}
	return self.GetKey(), result
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
