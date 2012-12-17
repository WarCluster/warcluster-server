package db_manager

import (
	"../entities"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"strings"
)

func SetEntity(entity entities.Entity) bool {
	key, prepared_entity := entity.Serialize()

	send_err := connection.Send("SET", key, prepared_entity)
	if send_err != nil {
		log.Print(send_err)
		return false
	}

	flush_err := connection.Flush()
	if flush_err != nil {
		log.Print(flush_err)
		return false
	}
	return true
}

func GetEntity(key string) entities.Entity {
	result, err := redis.Bytes(connection.Do("GET", key))
	if err != nil {
		log.Print(err)
		return nil
	}

	return entities.Construct(key, result)
}

func GetList(group_type string, username string) []entities.Entity {
	var entity_list []entities.Entity
	var coord string

	result, err := redis.String(connection.Do("GET", fmt.Sprintf("%v.%v", group_type, username)))
	if err != nil {
		log.Print(err)
		return nil
	}

	for _, coord = range strings.Split(result, ",") {
		key := fmt.Sprintf("%v.%v", group_type[:len(group_type)-1], coord)
		entity_list = append(entity_list, GetEntity(key))
	}
	return entity_list
}
