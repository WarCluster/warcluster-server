package main

import (
    "./db_manager"
    "./entities"
    "log"
)

func main() {
    defer db_manager.Finalize()
    username := "gophie"
    sun_position := []int{500, 300}

    hash := entities.GenerateHash(username)
    _, home_planet := entities.GeneratePlanets(hash, sun_position)
    player := entities.CreatePlayer(username, hash, home_planet)

    log.Println("Created player:", player)
    log.Println("------------------------------")
    db_manager.SetEntity(player)
    if new_player := db_manager.GetEntity(player.GetKey()); new_player != nil {
        log.Println("Fetched player from the db:", new_player)
    }
}
