package libs

import "math"

type Planet struct {
    coords []int
    texture int
    size int
    ship_count int
    max_ship_count int
    owner *Player
}

func GeneratePlanets(hash string, sun_position []int) ([]Planet, *Planet) {

    hashElement := func(index int) float64 {
        return float64(hash[index]) - 48
    }

    result := []Planet{}
    ring_offset := float64(80)
    planet_radius := float64(50)

    for ix:=0; ix<9; ix++ {
        planet_in_creation := Planet{[]int{0,0}, 0, 0, 0, 0, nil}
        ring_offset += planet_radius + hashElement(4 * ix)

        planet_in_creation.coords[0] = int(float64(sun_position[0]) + ring_offset * math.Cos(
            hashElement(4 * ix + 1) * 40))
        planet_in_creation.coords[1] = int(float64(sun_position[1]) + ring_offset * math.Sin(
            hashElement(4 * ix + 1) * 40))

        planet_in_creation.texture = int(hashElement(4 * ix + 2))
        planet_in_creation.size = 1 + int(hashElement(4 * ix + 3))
        result = append(result, planet_in_creation)
    }
    return result, &result[int(hashElement(37)) - 1]
}
