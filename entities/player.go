package entities

import (
	"encoding/json"
	"fmt"
	"time"
)

type Player struct {
	username       string
	TwitterID      string
	HomePlanet     string
	ScreenSize     []int
	ScreenPosition []int
}

func (self *Player) String() string {
	return self.username
}

func (self *Player) GetKey() string {
	return fmt.Sprintf("player.%s", self.username)
}

func (self *Player) StartMission(start_planet, end_planet *Planet, fleet int) *Mission {
	if fleet > 10 {
		fleet = 10
	} else if fleet <= 0 {
		fleet = 1
	}
	start_time := time.Now()
	ship_count := int(start_planet.ShipCount/10) * fleet
	start_planet.ShipCount -= ship_count
	mission := Mission{start_planet.GetKey(), start_time, self.GetKey(), ship_count, end_planet.GetKey()}
	return &mission
}

func (self *Player) Serialize() (string, []byte, error) {
	result, err := json.Marshal(self)
	if err != nil {
		return self.GetKey(), nil, err
	}
	return self.GetKey(), result, nil
}

func CreatePlayer(username, TwitterID string, HomePlanet *Planet) *Player {
	player := Player{username, TwitterID, HomePlanet.GetKey(), []int{0, 0}, []int{0, 0}}
	HomePlanet.Owner = username
	return &player
}
