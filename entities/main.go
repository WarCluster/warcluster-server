// Package entities implements the core game logic
package entities

import (
	"bytes"
	"encoding/gob"
	"strings"

	"warcluster/entities/db"
)

const (
	AREA_TEMPLATE       = "area:%d:%d"
	AREA_SIZE           = 10000
	PLANETS_RING_OFFSET = 300
	PLANET_RADIUS       = 300
	PLANET_COUNT        = 10
	PLANET_HASH_ARGS    = 4
	SOLAR_SYSTEM_RADIUS = 9000
	SPY_REPORT_VALIDITY = 300 // in seconds
)

// Entity interface is implemented by all entity types here
type Entity interface {
	AreaSet() string
	Key() string
}

// Simple RGB color struct
type Color struct {
	R float32
	G float32
	B float32
}

// Creates an entity via unmarshaling a json.
// The concrete entity type is given by the user as `key`
func Load(key string, data []byte) Entity {
	var (
		buffer bytes.Buffer
		entity Entity
	)

	buffer.Write(data)
	decoder := gob.NewDecoder(&buffer)
	entityType := strings.Split(key, ".")[0]

	switch entityType {
	case "player":
		entity = new(Player)
	case "planet":
		entity = new(Planet)
	case "mission":
		entity = new(Mission)
	case "sun":
		entity = new(Sun)
	case "ss":
		entity = new(SolarSlot)
	case "spy_report":
		entity = new(SpyReport)
	default:
		return nil
	}
	decoder.Decode(entity)
	return entity
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

	return Load(key, record), nil
}

// Saves an entity to the database. Records' key is entity.Key()
// If there is a record with such key in the database, simply updates
// the record. Otherwise creates a new one.
//
// Failed marshaling of the given entity is pretty much the only
// point of failure in this function... I supose.
func Save(entity Entity) error {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	key := entity.Key()
	err := encoder.Encode(entity)
	if err != nil {
		return err
	}

	setKey := entity.AreaSet()

	conn := db.Pool.Get()
	defer conn.Close()

	return db.Save(conn, key, setKey, buffer.Bytes())
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

		entityList = append(entityList, Load(key, record))
	}

	return entityList
}

// Remove a member from set
func RemoveFromArea(key, from string) error {
	conn := db.Pool.Get()
	defer conn.Close()

	return db.Srem(conn, key, from)
}
