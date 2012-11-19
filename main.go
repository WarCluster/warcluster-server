package main

import (
    "log"
    "github.com/garyburd/redigo/redis"
    // "./server"
    "./libs"
    // "encoding/json"
)

func main() {
    db, err := redis.Dial("tcp", ":6379")
    defer db.Close()
    if err != nil {
        log.Fatal(err)
    }

    username := "gophie"
    sun_position := []int{500, 300}

    hash := libs.GenerateHash(username)
    _, home_planet := libs.GeneratePlanets(hash, sun_position)
    player := libs.CreatePlayer(username, hash, home_planet)
    key, prepared_player := player.PrepareForDB()

    log.Println(player)
    log.Println(key, string(prepared_player))


    db.Send("SET", key, prepared_player)
    db.Flush()

    result, err := redis.String(db.Do("GET", key))
    if err != nil {
        log.Fatal(err)
    }
    log.Println(key, result)
}
