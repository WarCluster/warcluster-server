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

func (p *Player) String() string {
	return p.username
}

func (p *Player) GetKey() string {
	return fmt.Sprintf("player.%s", p.username)
}

func (p *Player) StartMission(source, target *Planet, fleet int) *Mission {
	if fleet > 100 {
		fleet = 100
	} else if fleet <= 0 {
		fleet = 10
	}
	current_time := time.Now().UnixNano() / 1e6
	ship_count := int(source.ShipCount/100) * fleet
	source.ShipCount -= ship_count
	mission := Mission{
		Source: source.GetCoords(),
		Target: target.GetCoords(),
		CurrentTime: current_time,
		StartTime: current_time,
		ArrivalTime: current_time,
		Player: p.username,
		ShipCount: ship_count,
	}
	mission.CalculateArrivalTime()
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
	player := Player{username, TwitterID, HomePlanet.GetKey(), []int{0, 0}, []int{0, 0}}
	HomePlanet.Owner = username
	return &player
}
