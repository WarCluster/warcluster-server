package entities

import (
    "math"
    "encoding/json"
    "fmt"
    "log"
)

type Planet struct {
    coords []int
    Texture int
    Size int
    ShipCount int
    MaxShipCount int
    Owner string
}

func (planet *Planet) GetKey() string {
    return fmt.Sprintf("planet.%d_%d", planet.coords[0], planet.coords[1])
}

func (planet Planet) PrepareForDB() (string, []byte) {
    key := planet.GetKey()
    result, err := json.Marshal(planet)
    if err != nil {
        log.Fatal(err)
    }
    return key, result
}

func GeneratePlanets(hash string, sun_position []int) ([]Planet, *Planet) {

    hashElement := func(index int) float64 {
        return float64(hash[index]) - 48
    }

    result := []Planet{}
    ring_offset := float64(80)
    planet_radius := float64(50)

    for ix:=0; ix<9; ix++ {
        planet_in_creation := Planet{[]int{0,0}, 0, 0, 0, 0, ""}
        ring_offset += planet_radius + hashElement(4 * ix)

        planet_in_creation.coords[0] = int(float64(sun_position[0]) + ring_offset * math.Cos(
            hashElement(4 * ix + 1) * 40))
        planet_in_creation.coords[1] = int(float64(sun_position[1]) + ring_offset * math.Sin(
            hashElement(4 * ix + 1) * 40))

        planet_in_creation.Texture = int(hashElement(4 * ix + 2))
        planet_in_creation.Size = 1 + int(hashElement(4 * ix + 3))
        result = append(result, planet_in_creation)
    }
    return result, &result[int(hashElement(37)) - 1]
}

