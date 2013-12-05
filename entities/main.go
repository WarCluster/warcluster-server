// Package entities implements the core game logic
package entities

import (
	"encoding/json"

	"warcluster/entities/db"
)

const (
	ENTITIES_AREA_TEMPLATE   = "area:%d:%d"
	ENTITIES_AREA_SIZE       = 10000
	PLANETS_RING_OFFSET      = 300
	PLANETS_PLANET_RADIUS    = 300
	PLANETS_PLANET_COUNT     = 10
	PLANETS_PLANET_HASH_ARGS = 4
	SUNS_SOLAR_SYSTEM_RADIUS = 9000
)

// Entity interface is implemented by all entity types here
type Entity interface {
	AreaSet() string
	Key() string
}

// Simple RGB color struct
type Color struct {
	R uint8
	G uint8
	B uint8
}

// Finds records in the database, by given key
// All Redis wildcards are allowed.
func Find(query string) []Entity {
	var entityList []Entity

	if records, err := GetList(query); err == nil {
		for _, key := range records {
			if entity, err := Get(key); err == nil {
				entityList = append(entityList, entity)
			}
		}
	}

	return entityList
}

// Returns keys of entities from the database
func GetList(pattern string) ([]string, error) {
	conn := db.Pool.Get()
	defer conn.Close()

	return db.GetList(conn, pattern)
}

// Fetches a single record in the database, by given concrete key.
// If there is no entity with such key, returns error.
func Get(key string) (Entity, error) {
	conn := db.Pool.Get()
	defer conn.Close()

	record, err := db.Get(conn, key)
	if err != nil {
		return nil, err
	}

	return Construct(key, record), nil
}

// Saves an entity to the database. Records' key is entity.Key()
// If there is a record with such key in the database, simply updates
// the record. Otherwise creates a new one.
//
// Failed marshaling of the given entity is pretty much the only
// point of failure in this function... I supose.
func Save(entity Entity) error {
	conn := db.Pool.Get()
	defer conn.Close()

	key := entity.Key()
	value, err := json.Marshal(entity)
	if err != nil {
		return err
	}

	setKey := entity.AreaSet()
	err = db.Save(conn, key, setKey, value)
	return err
}

// Deletes a record by the given key
func Delete(key string) error {
	conn := db.Pool.Get()
	defer conn.Close()

	return db.Delete(conn, key)
}

// Get and serialize all members of a set
func GetAreasMembers(areas []string) []Entity {
	conn := db.Pool.Get()
	defer conn.Close()

	keys := []string{}
	entityList := []Entity{}

	for _, area := range areas {
		result, err := db.Smembers(conn, area)
		if err != nil {
			continue
		}
		keys = append(keys, result...)
	}

	for _, key := range keys {
		record, err := db.Get(conn, key)
		if err != nil {
			continue
		}

		entityList = append(entityList, Construct(key, record))
	}

	return entityList
}
