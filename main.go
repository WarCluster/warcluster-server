package main

import (
    "log"
    "github.com/garyburd/redigo/redis"
    // "./server"
    "./libs"
    "encoding/json"
)

func main() {
    db, err := redis.Dial("tcp", ":6379")
    if err != nil {
        log.Fatal(err)
    }
    // server.Run()

    username := "gophie"
    value1 := &libs.Player{}

    sun_position := []int{499, 300}

    hash := libs.GenerateHash(username)
    _, home_planet := libs.GeneratePlanets(hash, sun_position)
    player := libs.CreatePlayer(username, hash, home_planet)
    log.Println("Player:", js)


    db.Send("SET", "foo", player)
    reply, err := redis.Values(db.Do("MGET", "foo"))
    if err != nil {
        log.Fatal(err)
    }
    if err := redis.ScanStruct(reply, &value1); err != nil {
        log.Fatal(err)
    }
    log.Print(value1)
}
