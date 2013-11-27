package entities

import (
	"encoding/json"
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

// This type is used as a proxy type while marshaling Planet.
type planetMarshalHook Planet

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
		"entities:%d:%d",
		int64(p.Position.X/ENTITIES_RANGE_SIZE),
		int64(p.Position.Y/ENTITIES_RANGE_SIZE),
	)
}

// We need to define the MarshalJSON in order to automatically
// update the ship count right before sending this entity to
// the client or to the database.
func (p *Planet) MarshalJSON() ([]byte, error) {
	p.UpdateShipCount()
	return json.Marshal((*planetMarshalHook)(p))
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

// Updates the ship count based on last time this count has
// been updated and of course the planet size.
// NOTE: If the planet is somebody's home we set a static increasion rate.
func (p *Planet) UpdateShipCount() {
	var timeModifier int64
	if p.HasOwner() {
		passedTime := time.Now().Unix() - p.LastShipCountUpdate
		if p.IsHome {
			timeModifier = 2
		} else {
			timeModifier = 6 - int64(p.Size/3)
		}
		p.ShipCount += int32(passedTime / (timeModifier * 10))
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
	planetRadius := float64(PLANETS_PLANET_RADIUS)

	for ix := 0; ix < PLANETS_PLANET_COUNT; ix++ {
		planet := Planet{
			Color:        Color{200, 180, 140},
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
	homePlanetIdx := int8(hashElement(PLANETS_PLANET_COUNT*PLANETS_PLANET_HASH_ARGS + 1))
	result[homePlanetIdx].IsHome = true
	result[homePlanetIdx].ShipCount = 80
	return result, result[homePlanetIdx]
}
