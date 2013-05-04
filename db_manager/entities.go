package db_manager

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"strings"
	"warcluster/entities"
)

/*
SetEntity takes an Entity (struct used as template for all data containers to ease the managing of the DB)
and generates an unique key in order to add the record to the DB.
*/
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

/*
GetEntity is used to pull information from the DB in order to be used by the server.
GetEntity operates as read only function and does not modify the data in the DB.
*/
func GetEntity(key string) (entities.Entity, error) {
	mutex.Lock()
	result, err := redis.Bytes(connection.Do("GET", key))
	mutex.Unlock()
	if err != nil {
		return nil, err
	}
	return entities.Construct(key, result), nil
}

/*
GetList is a special function needed to parse a list of keys stored in the DB for quick acsess.
For instance provide a userna in order to acsess the list of planets owned by this player.
The list will be iterated upon in order to call GetEntity for every key.
*/
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

/*
GetEntities operates as GetEntity but instead of an unique key it takes a patern in order to return
a lyst of Entitys that reflect the entered patern.
*/
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

/*
I think DeleteEntity speaks for itself but still. This function is used to remove entrys from the DB.
*/
func DeleteEntity(key string) error {
	mutex.Lock()
	_, err := redis.Bytes(connection.Do("DEL", key))
	mutex.Unlock()
	return err
}
