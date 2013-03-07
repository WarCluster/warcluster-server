package db_manager

import (
	"../entities"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"strings"
)

func SetEntity(entity entities.Entity) bool {
	key, prepared_entity, err := entity.Serialize()
	if err != nil {
		return false
	}

	mutex.Lock()
	send_err := connection.Send("SET", key, prepared_entity)
	mutex.Unlock()
	if send_err != nil {
		log.Print(send_err)
	}

	mutex.Lock()
	flush_err := connection.Flush()
	mutex.Unlock()
	if flush_err != nil {
		log.Print(flush_err)
	}
	return true
}

func GetEntity(key string) (entities.Entity, error) {
	mutex.Lock()
	result, err := redis.Bytes(connection.Do("GET", key))
	mutex.Unlock()
	if err != nil {
		return nil, err
	}
	return entities.Construct(key, result), nil
}

func GetList(group_type string, username string) []entities.Entity {
	var entity_list []entities.Entity
	var coord string

	mutex.Lock()
	result, err := redis.String(connection.Do("GET", fmt.Sprintf("%v.%v", group_type, username)))
	mutex.Unlock()
	if err != nil {
		log.Print(err)
		return nil
	}

	for _, coord = range strings.Split(result, ",") {
		key := fmt.Sprintf("%v.%v", group_type[:len(group_type)-1], coord)
		entity, _ := GetEntity(key)
		entity_list = append(entity_list, entity)
	}
	return entity_list
}

func GetEntities(pattern string) []entities.Entity {
	mutex.Lock()
	result, err := redis.Values(connection.Do("KEYS", pattern))
	mutex.Unlock()
	if err != nil {
		log.Print(err)
		return nil
	}

	results := fmt.Sprintf("%s", result)
	var entity_list []entities.Entity
	for _, key := range strings.Split(results, " ") {
		if entity, err := GetEntity(key); err == nil {
			entity_list = append(entity_list, entity)
		}
	}
	return entity_list
}

func DeleteEntity(key string) error {
	mutex.Lock()
	_, err := redis.Bytes(connection.Do("DEL", key))
	mutex.Unlock()
	return err
}
