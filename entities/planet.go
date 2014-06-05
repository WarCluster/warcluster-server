package entities

import (
	"fmt"
	"math"
	"time"

	"github.com/Vladimiroff/vec2d"
)

type Planet struct {
	Name                string
	Color               Color
	Position            *vec2d.Vector
	IsHome              bool
	Texture             int8
	Size                int8
	LastShipCountUpdate int64
	ShipCount           int32
	MaxShipCount        int32
	Owner               string
}

// Used only when a planet is being marshalled
type PlanetPacket struct {
	Planet
	IsSpied bool `json:",omitempty"`
}

// Database key.
func (p *Planet) Key() string {
	return fmt.Sprintf("planet.%s", p.Name)
}

// Checks if the planet has an owner or not.
func (p *Planet) HasOwner() bool {
	return len(p.Owner) > 0
}

// Returns the set by X or Y where this entity has to be put in
func (p *Planet) AreaSet() string {
	return fmt.Sprintf(
		AREA_TEMPLATE,
		RoundCoordinateTo(p.Position.X),
		RoundCoordinateTo(p.Position.Y),
	)
}

// Checks what the player could see and strips it if not
// Also updates the ship count right before marshaling
func (p *Planet) Sanitize(player *Player) *PlanetPacket {
	p.UpdateShipCount()
	packet := PlanetPacket{Planet: *p}

	if p.Owner != player.Username {
		packet.ShipCount = -1
		for _, spyReport := range player.SpyReports {
			if spyReport.Name == p.Name && spyReport.IsValid() {
				packet.ShipCount = spyReport.ShipCount
				packet.IsSpied = true
			}
		}
	}
	return &packet
}

// Returns the ship count right after the ship count update.
func (p *Planet) GetShipCount() int32 {
	p.UpdateShipCount()
	return p.ShipCount
}

// Changes the ship count right after the ship count update.
// NOTE: I'm still not sure if we need a mutex here...
func (p *Planet) SetShipCount(count int32) {
	p.UpdateShipCount()
	p.ShipCount = count
	p.LastShipCountUpdate = time.Now().Unix()
}

func ShipCountTimeMod(size int8, isHome bool) int64 {
	var timeModifier int64
	if isHome {
		timeModifier = 2
	} else {
		timeModifier = 6 - int64(size/3)
	}
	return timeModifier * 10
}

// Updates the ship count based on last time this count has
// been updated and of course the planet size.
// NOTE: If the planet is somebody's home we set a static increasion rate.
func (p *Planet) UpdateShipCount() {
	if p.HasOwner() {
		passedTime := time.Now().Unix() - p.LastShipCountUpdate
		p.ShipCount += int32(passedTime / ShipCountTimeMod(p.Size, p.IsHome))
		p.LastShipCountUpdate = time.Now().Unix()
	}
}

// Generates all planets in a solar system, based on the user's hash.
func GeneratePlanets(nickname string, sun *Sun) ([]*Planet, *Planet) {
	hash := GenerateHash(nickname)
	hashElement := func(index int) float64 {
		return float64(hash[index]) - 48 // The offset of simbol "1" in the ascii table
	}

	result := []*Planet{}
	ringOffset := float64(PLANETS_RING_OFFSET)
	planetRadius := float64(PLANET_RADIUS)

	for ix := 0; ix < PLANET_COUNT; ix++ {
		planet := Planet{
			Color:        Color{0.78431373, 0.70588235, 0.54901961},
			Position:     new(vec2d.Vector),
			IsHome:       false,
			ShipCount:    10,
			MaxShipCount: 0,
			Owner:        "",
		}
		// NOTE: 5 is the distance between planets
		ringOffset += planetRadius + hashElement(4*ix)*5

		planet.Name = fmt.Sprintf("%s%v", sun.Name, ix)
		planet.Position.X = math.Floor(sun.Position.X + ringOffset*math.Cos(hashElement(4*ix+1)*40))
		planet.Position.Y = math.Floor(sun.Position.Y + ringOffset*math.Sin(hashElement(4*ix+1)*40))
		planet.Texture = int8(hashElement(4*ix + 2))
		planet.Size = 1 + int8(hashElement(4*ix+3))
		planet.LastShipCountUpdate = time.Now().Unix()
		result = append(result, &planet)
	}
	// + 1 bellow stands for: after all the planet info is read the next element is the user's home planet idx
	homePlanetIdx := int8(hashElement(PLANET_COUNT*PLANET_HASH_ARGS + 1))
	result[homePlanetIdx].IsHome = true
	result[homePlanetIdx].ShipCount = 80
	return result, result[homePlanetIdx]
}
