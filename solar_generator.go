package main

import (
    "fmt"
    "./libs"
)

func main() {
    // Hardocded values
    username := "gophie"
    sun_position := []int{500, 300}

    hash := libs.GenerateHash(username)
    planets, home_planet := libs.GeneratePlanets(hash, sun_position)
    player := libs.CreatePlayer(username, hash, home_planet)
    fmt.Println("Player:", player)
    fmt.Println("Generated Planets:", planets)
}


