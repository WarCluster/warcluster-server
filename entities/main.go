// Package entities implements the core game logic
package entities

import (
	"bytes"
	"encoding/gob"
	"errors"
	"strings"

	"warcluster/config"
	"warcluster/entities/db"
)

type Race struct {
	ID    uint8
	Name  string
	Color Color
}

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

var (
	Settings config.Entities
	Races    []Race
)

func init() {
	var cfg config.Config
	cfg.Load("../config/config.gcfg")

	Settings = cfg.Entities
	Races = make([]Race, len(cfg.Race), len(cfg.Race))
	for name, params := range cfg.Race {
		Races[params.Id] = Race{params.Id, name, Color{params.Red, params.Green, params.Blue}}
	}
}

//Validate if the color values are in range
func (c *Color) Validate() error {
	if c.R < 0 || c.R > 1 {
		return errors.New("Color component red out of range.")
	}
	if c.G < 0 || c.G > 1 {
		return errors.New("Color component green out of range.")
	}
	if c.B < 0 || c.B > 1 {
		return errors.New("Color component blue out of range.")
	}
	return nil
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

// Move member from one set to another
func moveToArea(key, from, to string) error {
	conn := db.Pool.Get()
	defer conn.Close()

	return db.Smove(conn, from, to, key)
}

// Remove a member from set
func RemoveFromArea(key, set string) error {
	conn := db.Pool.Get()
	defer conn.Close()

	return db.Srem(conn, set, key)
}

// Returns if entity is a member of the set
func InArea(key, set string) bool {
	conn := db.Pool.Get()
	defer conn.Close()

	result, err := db.Sismember(conn, set, key)
	if err != nil {
		return false
	}
	return result
}
