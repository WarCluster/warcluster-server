package db_manager

import (
    "fmt"
    "log"
    "github.com/garyburd/redigo/redis"
    "../entities"
    "strings"
)

var connection redis.Conn

const (
    HOSTNAME = "localhost"
    PORT = 6379
    NETWORK = "tcp"
)


func init() {
    var err error
    log.Print("Initializing database connection... ")
    connection, err = connect()
    if err != nil {
        log.Fatal(err)
    }
}

func connect() (redis.Conn, error) {
    return redis.Dial("tcp", fmt.Sprintf("%v:%v", HOSTNAME, PORT))
}

func Finalize() {
    log.Print("Closing database connection... ")
    err := connection.Close()
    if err != nil {
        log.Fatal(err)
    }
}

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
    result, err := redis.String(connection.Do("GET", fmt.Sprintf("%v.%v", group_type, username)))
    if err != nil {
        log.Print(err)
        return nil
    }
    var entity_list []entities.Entity
    var coord string
    for _, coord = range strings.Split(result, ",") {
        key := fmt.Sprintf("%v.%v", group_type[:len(group_type)-1], coord)
        entity_list = append(entity_list, GetEntity(key))
    }
    return entity_list
}


// func GetPlayerMissions(username string) []*entities.Mission {
// 
// }
// 
// func GetPlayerPlanets(username string) []*entities.Planet {
// 
// }
