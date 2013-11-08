package entities

import (
	"encoding/json"
	"fmt"
	"github.com/Vladimiroff/vec2d"
	"math"
	"time"
)

type Planet struct {
	Name                string
	Color               Color
	Position            *vec2d.Vector
	IsHome              bool
	Texture             int
	Size                int
	LastShipCountUpdate int64
	ShipCount           int
	MaxShipCount        int
	Owner               string
}

// This type is used as a proxy type while marshaling Planet
type planetMarshalHook Planet

func (p *Planet) Key() string {
	return fmt.Sprintf("planet.%s", p.Name)
}

func (p *Planet) HasOwner() bool {
	return len(p.Owner) > 0
}

func (p *Planet) MarshalJSON() ([]byte, error) {
	p.UpdateShipCount()
	return json.Marshal((*planetMarshalHook)(p))
}

func (p *Planet) GetShipCount() int {
	p.UpdateShipCount()
	return p.ShipCount
}

func (p *Planet) SetShipCount(count int) {
	p.UpdateShipCount()
	p.ShipCount = count
	p.LastShipCountUpdate = time.Now().Unix()
}

func (p *Planet) UpdateShipCount() {
	if p.HasOwner() {
		passedTime := time.Now().Unix() - p.LastShipCountUpdate
		timeModifier := int64(p.Size/3) + 1
		//TODO: To be completed for all planet size types
		//if getobject(Owner.getkey).gethomeplanet == p.getkey
		p.ShipCount += int(passedTime / (timeModifier * 10))
		p.LastShipCountUpdate = time.Now().Unix()
	}
}

/*
TODO: We need to add ship count on new planet creation
TODO: Put all funny numbers in a constans in our config file
NOTE: 5 in ring_offset is the distance between planets
*/
func GeneratePlanets(nickname string, sun *Sun) ([]*Planet, *Planet) {
	hash := generateHash(nickname)
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
		ringOffset += planetRadius + hashElement(4*ix)*5

		planet.Name = fmt.Sprintf("%s%v", sun.Name, ix)
		planet.Position.X = math.Floor(sun.Position.X + ringOffset*math.Cos(hashElement(4*ix+1)*40))
		planet.Position.Y = math.Floor(sun.Position.Y + ringOffset*math.Sin(hashElement(4*ix+1)*40))
		planet.Texture = int(hashElement(4*ix + 2))
		planet.Size = 1 + int(hashElement(4*ix+3))
		planet.LastShipCountUpdate = time.Now().Unix()
		result = append(result, &planet)
	}
	// + 1 bellow stands for: after all the planet info is read the next element is the user's home planet idx
	homePlanetIdx := int(hashElement(PLANETS_PLANET_COUNT*PLANETS_PLANET_HASH_ARGS + 1))
	result[homePlanetIdx].IsHome = true
	return result, result[homePlanetIdx]
}
