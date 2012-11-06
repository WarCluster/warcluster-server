package main

import (
    "fmt"
    "./libs"
)

func main() {
    player := libs.CreatePlayer("kiril")
    fmt.Println(player.String())
    // fmt.Println(player.generatePlanets([]int{1,2}))
}


