package entities

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
)

type Planet struct {
	coords       []int
	Texture      int
	Size         int
	ShipCount    int
	MaxShipCount int
	Owner        string
}

func (self Planet) String() string {
	return fmt.Sprintf("Planet[%s, %s]", self.coords[0], self.coords[1])
}

func (self Planet) GetKey() string {
	return fmt.Sprintf("planet.%d_%d", self.coords[0], self.coords[1])
}

func (self Planet) GetCoords() []int {
	return self.coords
}

func (self Planet) Serialize() (string, []byte) {
	result, err := json.Marshal(self)
	if err != nil {
		log.Fatal(err)
	}
	return self.GetKey(), result
}

func GeneratePlanets(hash string, sun_position []int) ([]Planet, *Planet) {

	hashElement := func(index int) float64 {
		return float64(hash[index]) - 48
	}

	result := []Planet{}
	ring_offset := float64(80)
	planet_radius := float64(50)

	for ix := 0; ix < 9; ix++ {
		planet_in_creation := Planet{[]int{0, 0}, 0, 0, 0, 0, ""}
		ring_offset += planet_radius + hashElement(4*ix)

		planet_in_creation.coords[0] = int(float64(sun_position[0]) + ring_offset*math.Cos(
			hashElement(4*ix+1)*40))
		planet_in_creation.coords[1] = int(float64(sun_position[1]) + ring_offset*math.Sin(
			hashElement(4*ix+1)*40))

		planet_in_creation.Texture = int(hashElement(4*ix + 2))
		planet_in_creation.Size = 1 + int(hashElement(4*ix+3))
		result = append(result, planet_in_creation)
	}
	return result, &result[int(hashElement(37))-1]
}
