package entities

type Entity interface {
	GetKey() string
	Serialize() (string, []byte, error)
	String() string
}

const (
	PLANETS_RING_OFFSET           = 600
	PLANETS_PLANET_RADIUS         = 160
	PLANETS_PLANET_COUNT          = 10
	PLANETS_PLANET_HASH_ARGS      = 4
	SUNS_RANDOM_SPAWN_ZONE_RADIUS = 50000
	SUNS_SOLAR_SYSTEM_RADIUS      = 4000
)

var Types map[string]Entity
