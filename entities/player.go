package entities

import (
	"encoding/json"
	"fmt"
	"time"
)

type Player struct {
	username       string
	Color          Color
	TwitterID      string
	HomePlanet     string
	ScreenSize     []int
	ScreenPosition []int
}

func (p *Player) String() string {
	return p.username
}

func (p *Player) GetKey() string {
	return fmt.Sprintf("player.%s", p.username)
}

func (p *Player) StartMission(source *Planet, target *Planet, fleet int, mission_type string) *Mission {
	if fleet > 100 {
		fleet = 100
	} else if fleet <= 0 {
		fleet = 10
	}
	current_time := time.Now().UnixNano() / 1e6
	base_ship_count := source.GetShipCount()
	ship_count := int(base_ship_count * fleet / 100)
	source.SetShipCount(base_ship_count - ship_count)

	mission := Mission{
		Source:      source.GetCoords(),
		Target:      target.GetCoords(),
		Type:        mission_type,
		CurrentTime: current_time,
		StartTime:   current_time,
		TravelTime:  current_time,
		Player:      p.username,
		ShipCount:   ship_count,
	}
	mission.CalculateTravelTime()
	return &mission
}

func (p *Player) Serialize() (string, []byte, error) {
	result, err := json.Marshal(p)
	if err != nil {
		return p.GetKey(), nil, err
	}
	return p.GetKey(), result, nil
}

func CreatePlayer(username, TwitterID string, HomePlanet *Planet) *Player {
	userhash := simplifyHash(usernameHash(username))

	red := []int{151, 218 , 233, 72, 27}
	green := []int{8, 75, 177, 140, 85}
	blue := []int{14, 15, 4, 19, 192}
	hashValue := func(index int) int {
		return int((userhash[0] - 48)/2)
	}

	color := Color{username, red[hashValue(0)], green[hashValue(0)], blue[hashValue(0)]}
	player := Player{username, color, TwitterID, HomePlanet.GetKey(), []int{0, 0}, []int{0, 0}}
	HomePlanet.Owner = username
	HomePlanet.Color = color
	return &player
}
