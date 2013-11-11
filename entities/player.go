package entities

import (
	"fmt"
	"time"
)

type Player struct {
	Username       string
	Color          Color
	TwitterID      string
	HomePlanet     string
	ScreenSize     []int
	ScreenPosition []int
}

// Database key.
func (p *Player) Key() string {
	return fmt.Sprintf("player.%s", p.Username)
}

// Starts missions to one of the players planet to some other. Each mission have type
// and the user decides which part of the planet's fleet he would like to send.
func (p *Player) StartMission(source, target *Planet, fleet int, missionType string) *Mission {
	if fleet > 100 {
		fleet = 100
	} else if fleet <= 0 {
		fleet = 10
	}
	currentTime := time.Now().UnixNano() / 1e6
	baseShipCount := source.GetShipCount()
	shipCount := int(baseShipCount * fleet / 100)
	source.SetShipCount(baseShipCount - shipCount)

	mission := Mission{
		Color:     p.Color,
		Source:    source.Name,
		Target:    target.Name,
		Type:      missionType,
		StartTime: currentTime,
		Player:    p.Username,
		ShipCount: shipCount,
	}
	mission.TravelTime = calculateTravelTime(source.Position, target.Position, mission.GetSpeed())
	return &mission
}

// Creates new player after the authentication and generates color based on the unique hash
func CreatePlayer(username, TwitterID string, HomePlanet *Planet) *Player {
	userhash := simplifyHash(usernameHash(username))

	red := []int{151, 218, 233, 72, 245, 84}
	green := []int{8, 75, 177, 140, 105, 146}
	blue := []int{14, 15, 4, 19, 145, 219}
	hashValue := func(index int) int {
		return int((userhash[0] - 48) / 2)
	}

	color := Color{red[hashValue(0)], green[hashValue(0)], blue[hashValue(0)]}
	player := Player{username, color, TwitterID, HomePlanet.Key(), []int{0, 0}, []int{0, 0}}
	HomePlanet.Owner = username
	HomePlanet.Color = color
	return &player
}
