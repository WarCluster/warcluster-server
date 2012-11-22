package entities

import (
    "time"
    "encoding/json"
    "fmt"
    "log"
)

type Mission struct {
    start_planet *Planet
    start_time time.Time
    Player string
    ShipCount int
    EndPlanet *Planet
}

// TODO Finish this!
func (self *Mission) PrepareForDB() (string, []byte) {
    key := fmt.Sprintf(
        "mission_%d_%d_%d",
        self.start_planet.coords[0],
        self.start_planet.coords[1],
    )

    result, err := json.Marshal(self)
    if err != nil {
        log.Fatal(err)
    }
    return key, result
}

// func (mission *Mission) End() {
//     if mission.end_planet.ship_count >= mission.ship_count {
//         mission.end_planet.ship_count -= mission.ship_count
//     } else {
//         mission.ship_count -= mission.end_planet.ship_count
//         mission.end_planet.ship_count = mission.ship_count
//         mission.end_planet.max_ship_count = mission.ship_count
//         mission.end_planet.owner = mission.player
//         //Fuuuuuck
//     }
// }
// 
// func StartMission(player *Player, fleet int, start_planet, end_planet *Planet) *Mission {
//     start_time = time.Now()
//     ship_count = int(start_planet.ship_count / 10) * fleet
//     start_planet.ship_count -= ship_count
//     mission := Mission{player, ship_count, start_planet, end_planet, start_time}
//     player.missions = append(player.missions, *mission)
//     return *mission
// }
// 
