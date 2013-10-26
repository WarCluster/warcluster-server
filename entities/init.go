package entities

const (
	PLANETS_RING_OFFSET           = 300
	PLANETS_PLANET_RADIUS         = 300
	PLANETS_PLANET_COUNT          = 10
	PLANETS_PLANET_HASH_ARGS      = 4
	SUNS_RANDOM_SPAWN_ZONE_RADIUS = 50000
	SUNS_SOLAR_SYSTEM_RADIUS      = 9000
)

type Entity interface {
	GetKey() string
	String() string
}

type missionMarshalHook Mission
type planetMarshalHook Planet
