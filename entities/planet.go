package entities

import (
	"encoding/json"
	"fmt"
	"github.com/Vladimiroff/vec2d"
	"math"
	"time"
)

type Planet struct {
	coords              []int
	Texture             int
	Size                int
	LastShipCountUpdate int64
	ShipCount           int
	MaxShipCount        int
	Owner               string
}

func (self *Planet) String() string {
	return fmt.Sprintf("Planet[%s, %s]", self.coords[0], self.coords[1])
}

func (self *Planet) GetKey() string {
	return fmt.Sprintf("planet.%d_%d", self.coords[0], self.coords[1])
}

func (self *Planet) GetCoords() []int {
	return self.coords
}

func (self *Planet) Serialize() (string, []byte, error) {
	_ = self.GetShipCount()
	result, err := json.Marshal(self)
	if err != nil {
		return self.GetKey(), nil, err
	}
	return self.GetKey(), result, nil
}

func (self *Planet) GetShipCount() int {
	if len(self.Owner) > 0 {
		self.UpdateShipCount()
	}
	return self.ShipCount
}

func (self *Planet) SetShipCount(count int) {
	self.ShipCount = count
	self.LastShipCountUpdate = time.Now().Unix()
}

func (self *Planet) UpdateShipCount() {
	passedTime := time.Now().Unix() - self.LastShipCountUpdate
	timeModifier := int64(self.Size/3) + 1
	//TODO: To be completed for all planet size types
	//if getobject(Owner.getkey).gethomeplanet == self.getkey
	self.ShipCount += int(passedTime / timeModifier)
	self.LastShipCountUpdate = time.Now().Unix()
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
		planet_in_creation := Planet{[]int{0, 0}, 0, 0, time.Now().Unix(), 10, 0, ""}
		ring_offset += planet_radius + hashElement(4*ix)*5

		planet_in_creation.coords[0] = int(float64(sun_position.X) + ring_offset*math.Cos(
			hashElement(4*ix+1)*40))
		planet_in_creation.coords[1] = int(float64(sun_position.Y) + ring_offset*math.Sin(
			hashElement(4*ix+1)*40))

		planet_in_creation.Texture = int(hashElement(4*ix+2))
		planet_in_creation.Size = 1 + int(hashElement(4*ix+3))
		//self.LastShipCountUpdate = time.Now()
		result = append(result, &planet_in_creation)
	}
	// + 1 bellow stands for: after all the planet info is read the next element is the user's home planet idx
	return result, result[int(hashElement(PLANETS_PLANET_COUNT*PLANETS_PLANET_HASH_ARGS+1))]
}
