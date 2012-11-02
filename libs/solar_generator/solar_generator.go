package solar_generator


type Player string

type Planet struct {
    coords []int
    texture int
    size int
    ship_count int
    max_ship_count int
    owner Player
}
