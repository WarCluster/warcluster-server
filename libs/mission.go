package libs

import "time"

type Mission struct {
    start_planet *Planet
    start_time time.Time
    Player string
    ShipCount int
    EndPlanet *Planet
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
