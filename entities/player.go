package entities

import (
	"fmt"
	"sync"
	"time"

	"github.com/Vladimiroff/vec2d"
)

type Player struct {
	Username       string
	Color          Color
	TwitterID      string
	HomePlanet     string
	ScreenSize     []uint16
	ScreenPosition *vec2d.Vector
	SpyReports     []*SpyReport `json:"-" bson:"-"`
	mutex          sync.Mutex
}

// Database key.
func (p *Player) Key() string {
	return fmt.Sprintf("player.%s", p.Username)
}

// Returns the sorted set by X or Y where this entity has to be put in
func (p *Player) AreaSet() string {
	homePlanet, _ := Get(p.HomePlanet)
	return homePlanet.AreaSet()
}

// Starts missions to one of the players planet to some other. Each mission have type
// and the user decides which part of the planet's fleet he would like to send.
func (p *Player) StartMission(source, target *Planet, fleet int32, missionType string) *Mission {
	currentTime := time.Now().UnixNano() / 1e6
	baseShipCount := source.GetShipCount()
	shipCount := int32(baseShipCount * fleet / 100)
	source.SetShipCount(baseShipCount - shipCount)

	mission := Mission{
		Color: p.Color,
		Source: embeddedPlanet{
			Name:     source.Name,
			Owner:    source.Owner,
			Position: source.Position,
		},
		Target: embeddedPlanet{
			Name:     target.Name,
			Owner:    target.Owner,
			Position: target.Position,
		},
		Type:      missionType,
		StartTime: currentTime,
		Player:    p.Username,
		ShipCount: shipCount,
		areaSet:   source.AreaSet(),
	}
	mission.TravelTime = calculateTravelTime(source.Position, target.Position, mission.GetSpeed())
	return &mission
}

func (p *Player) UpdateSpyReports() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	playersReports := Find(fmt.Sprintf("spy_report.%s_*", p.Username))
	spyReports := make([]*SpyReport, 0, len(playersReports))
	for _, reportEntity := range playersReports {
		report := reportEntity.(*SpyReport)
		if report.IsValid() {
			spyReports = append(spyReports, report)
		} else {
			Delete(report.Key())
		}
	}
	p.SpyReports = spyReports
}

// Creates new player after the authentication and generates color based on the unique hash
func CreatePlayer(username, TwitterID string, homePlanet *Planet) *Player {
	userhash := simplifyHash(usernameHash(username))

	red := []float32{0.59215686, 0.85490196, 0.91372549, 0.28235294, 0.96078431, 0.32941176}
	green := []float32{0.031372549, 0.29411765, 0.69411765, 0.54901961, 0.41176471, 0.57254902}
	blue := []float32{0.054901961, 0.058823529, 0.015686275, 0.074509804, 0.56862745, 0.85882353}
	hashValue := func(index uint8) uint8 {
		return uint8((userhash[0] - 48) / 2)
	}

	color := Color{red[hashValue(0)], green[hashValue(0)], blue[hashValue(0)]}
	player := Player{
		Username:       username,
		Color:          color,
		TwitterID:      TwitterID,
		HomePlanet:     homePlanet.Key(),
		ScreenSize:     []uint16{0, 0},
		ScreenPosition: homePlanet.Position,
	}
	homePlanet.Owner = username
	homePlanet.Color = color
	return &player
}
