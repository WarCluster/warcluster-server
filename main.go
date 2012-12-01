package main

import (
    "log"
    "./entities"
    "./db_manager"
)

func main() {
    defer db_manager.Finalize()
    username := "gophie"
    sun_position := []int{500, 300}

    hash := entities.GenerateHash(username)
    _, home_planet := entities.GeneratePlanets(hash, sun_position)
    player := entities.CreatePlayer(username, hash, home_planet)

    log.Println(player)
    log.Println("------------------------------")
    db_manager.SetEntity(player)
    new_player := db_manager.GetEntity(player.GetKey())
    if new_player != nil {
        log.Print(new_player)
    }

    log.Println(db_manager.GetList("planets", "gophie"))

}
