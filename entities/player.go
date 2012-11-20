package entities

import (
    "fmt"
    "encoding/json"
    "log"
)

type Player struct {
    username string
    Hash string
    HomePlanet string
    ScreenSize []int
    ScreenPosition []int

}

func (player *Player) String() string {
    return player.username
}

func (player *Player) GetKey() string {
    return fmt.Sprintf("player.%s", player.username)
}

func (player *Player) GetMissions() []*Mission {
    return []*Mission{}
}

func (player *Player) StartMission() error {
    // new_mission := CreateMission()
    //db.Do("RPUSH new_mission)
    return nil
}

func (player Player) PrepareForDB() (string, []byte) {
    key := player.GetKey()
    result, err := json.Marshal(player)
    if err != nil {
        log.Fatal(err)
    }
    return key, result
}

func CreatePlayer(username, Hash string, HomePlanet *Planet) Player {
    player := Player{username, Hash, HomePlanet.GetKey(), []int{0, 0}, []int{0, 0}}
    HomePlanet.Owner = username
    return player
}

