package entities

import (
	"encoding/json"
	"fmt"
	"github.com/Vladimiroff/vec2d"
	"math"
	"time"
)

type Planet struct {
	Color               Color
	coords              []int
	Texture             int
	Size                int
	LastShipCountUpdate int64
	ShipCount           int
	MaxShipCount        int
	Owner               string
}

func (p *Planet) String() string {
	return fmt.Sprintf("Planet[%s, %s]", p.coords[0], p.coords[1])
}

func (p *Planet) GetKey() string {
	return fmt.Sprintf("planet.%d_%d", p.coords[0], p.coords[1])
}

func (p *Planet) GetCoords() []int {
	return p.coords
}

func (p *Planet) Serialize() (string, []byte, error) {
	_ = p.GetShipCount()
	result, err := json.Marshal(p)
	if err != nil {
		return p.GetKey(), nil, err
	}
	return p.GetKey(), result, nil
}

func (p *Planet) GetShipCount() int {
	if len(p.Owner) > 0 {
		p.UpdateShipCount()
	}
	return p.ShipCount
}

func (p *Planet) SetShipCount(count int) {
	p.ShipCount = count
	p.LastShipCountUpdate = time.Now().Unix()
}

func (p *Planet) UpdateShipCount() {
	if len(p.Owner) > 0 {
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
func GeneratePlanets(hash string, sun_position *vec2d.Vector) ([]*Planet, *Planet) {
	hashElement := func(index int) float64 {
		return float64(hash[index]) - 48 // The offset of simbol "1" in the ascii table
	}

	result := []*Planet{}
	ring_offset := float64(PLANETS_RING_OFFSET)
	planet_radius := float64(PLANETS_PLANET_RADIUS)

	for ix := 0; ix < PLANETS_PLANET_COUNT; ix++ {
		planet_in_creation := Planet{Color{"Base", 200, 180, 140}, []int{0, 0}, 0, 0, time.Now().Unix(), 10, 0, ""}
		ring_offset += planet_radius + hashElement(4*ix)*5

		planet_in_creation.coords[0] = int(float64(sun_position.X) + ring_offset*math.Cos(
			hashElement(4*ix+1)*40))
		planet_in_creation.coords[1] = int(float64(sun_position.Y) + ring_offset*math.Sin(
			hashElement(4*ix+1)*40))

		planet_in_creation.Texture = int(hashElement(4*ix + 2))
		planet_in_creation.Size = 1 + int(hashElement(4*ix+3))
		//p.LastShipCountUpdate = time.Now()
		result = append(result, &planet_in_creation)
	}
	// + 1 bellow stands for: after all the planet info is read the next element is the user's home planet idx
	return result, result[int(hashElement(PLANETS_PLANET_COUNT*PLANETS_PLANET_HASH_ARGS+1))]
}
