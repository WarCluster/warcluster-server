package entities

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Player struct {
	username       string
	Hash           string
	HomePlanet     string
	ScreenSize     []int
	ScreenPosition []int
}

func (self Player) String() string {
	return self.username
}

func (self Player) GetKey() string {
	return fmt.Sprintf("player.%s", self.username)
}

func (self *Player) StartMission(start_planet, end_planet Planet, fleet int) Mission {
	if fleet > 10 {
		fleet = 10
	} else if fleet <= 0 {
		fleet = 1
	}
	start_time := time.Now()
	ship_count := int(start_planet.ShipCount/10) * fleet
	start_planet.ShipCount -= ship_count
	mission := Mission{start_planet.GetKey(), start_time, self.GetKey(), ship_count, end_planet.GetKey()}
	// da se vpishe missiona v bazata i da se appendne link kum neq v player missions
	return mission
}

func (self Player) Serialize() (string, []byte) {
	result, err := json.Marshal(self)
	if err != nil {
		log.Fatal(err)
	}
	return self.GetKey(), result
}

func CreatePlayer(username, Hash string, HomePlanet *Planet) Player {
	player := Player{username, Hash, HomePlanet.GetKey(), []int{0, 0}, []int{0, 0}}
	HomePlanet.Owner = username
	return player
}
