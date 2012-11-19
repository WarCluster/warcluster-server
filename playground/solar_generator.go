package main

import (
    "log"
    "./libs"
)

func main() {
    username := "gophie"
    sun_position := []int{500, 300}

    hash := libs.GenerateHash(username)
    planets, home_planet := libs.GeneratePlanets(hash, sun_position)
    player := libs.CreatePlayer(username, hash, home_planet)
    log.Println("Player:", player)
    log.Println("Generated Planets:", planets)
}

