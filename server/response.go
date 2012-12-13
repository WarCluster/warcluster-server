package server

import "../db_manager"
import e "../entities"
import "fmt"
import "log"

const BEST_PING = 150
const WORST_PING = 1500
const STEPS = 10

func scopeOfView(position []int, resolution []int, lag int) ([]int, []int) {
    step := int(WORST_PING-BEST_PING/STEPS)
    multiply := 1.1 + float32((lag - BEST_PING) / step) * 0.1
    end_resolution := []int{
        int(float32(resolution[0]) * multiply),
        int(float32(resolution[1]) * multiply),
    }

    top_left := []int{
        position[0] - int((end_resolution[0] - resolution[0]) / 2),
        position[1] - int((end_resolution[1] - resolution[1]) / 2),
    }

    bottom_right := []int{
        position[0] + resolution[0] + int((end_resolution[0] - resolution[0]) / 2),
        position[1] + resolution[1] + int((end_resolution[1] - resolution[1]) / 2),
    }
    return top_left, bottom_right
}

func actionParser(username, start_planet_key, end_planet_key string, count int) error {
    var player e.Player
    var start_planet, end_planet e.Planet

    player = db_manager.GetEntity(fmt.Sprint("player.", username)).(e.Player)
    start_planet = db_manager.GetEntity(fmt.Sprint("planet.", start_planet_key)).(e.Planet)
    end_planet = db_manager.GetEntity(fmt.Sprint("planet.", end_planet_key)).(e.Planet)

    if start_planet.Owner != username {
        log.Fatal("This is not your home!")
    }

    if start_planet.ShipCount < count {
        count = start_planet.ShipCount
    }

    return player.StartMission(start_planet, end_planet, count)
}
