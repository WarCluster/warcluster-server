package entities

import (
	"fmt"
	"sync"
	"time"

	"github.com/Vladimiroff/vec2d"
)

type Player struct {
	Username       string
	RaceID         uint8
	TwitterID      string
	HomePlanet     string
	ScreenSize     []uint64
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

// Returns the name of sun in player's solar system
func (p *Player) Sun() string {
	return p.HomePlanet[:len(p.HomePlanet)-1]
}

// Starts missions to one of the players planet to some other. Each mission have type
// and the user decides which part of the planet's fleet he would like to send.
func (p *Player) StartMission(source, target *Planet, path []*vec2d.Vector, fleet int32, missionType string) *Mission {
	currentTime := time.Now().UnixNano() / 1e6
	baseShipCount := source.GetShipCount()
	shipCount := int32(baseShipCount * fleet / 100)
	source.SetShipCount(baseShipCount - shipCount)

	mission := Mission{
		Color: Races[p.RaceID].Color,
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
		Path:      path,
		Type:      missionType,
		StartTime: currentTime,
		Player:    p.Username,
		ShipCount: shipCount,
		areaSet:   source.AreaSet(),
	}
	mission.TravelTime = calculateMissionTravelTime(source.Position, target.Position, path, Settings.MissionSpeed)
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
func CreatePlayer(username, TwitterID string, homePlanet *Planet, setupData *SetupData) *Player {
	player := Player{
		Username:       username,
		TwitterID:      TwitterID,
		HomePlanet:     homePlanet.Key(),
		ScreenSize:     []uint64{0, 0},
		ScreenPosition: homePlanet.Position,
	}

	player.RaceID = setupData.Race

	homePlanet.Owner = username
	homePlanet.Color = Races[player.RaceID].Color
	return &player
}
