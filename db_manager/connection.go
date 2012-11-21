package db_manager

import (
    "fmt"
    "log"
    "github.com/garyburd/redigo/redis"
    "encoding/json"
    "../entities"
    // "github.com/garyburd/redigo/redisx"
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
    key, prepared_entity := entity.PrepareForDB()

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

func GetEntity(key string) entity.Entity {
    result, err := redis.Bytes(connection.Do("GET", key))
    if err != nil {
        log.Print(err)
        return nil
    }

    return Construct(key, result)
}
