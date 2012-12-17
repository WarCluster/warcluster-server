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

// func (mission *Mission, ) End() {
//     end_point
//     start_point

//     if mission.end_planet.ship_count >= mission.ship_count {
//         mission.end_planet.ship_count -= mission.ship_count
//     } else {
//         mission.ship_count -= mission.end_planet.ship_count
//         mission.end_planet.ship_count = mission.ship_count
//         mission.end_planet.max_ship_count = mission.ship_count
//         mission.end_planet.owner = mission.player
//         //Fuuuuuck
//     }
// }
