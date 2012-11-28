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

func (self Player) String() string {
    return self.username
}

func (self Player) GetKey() string {
    return fmt.Sprintf("player.%s", self.username)
}

func (self *Player) GetMissions() []*Mission {
    return []*Mission{}
}

func (self *Player) GetPlanets() []*Planet {
    return []*Planet{}
}

func (self *Player) StartMission() error {
    // new_mission := CreateMission()
    //db.Do("RPUSH new_mission)
    return nil
}

func (self Player) Serialize() (string, []byte) {
    result, err := json.Marshal(self)
    if err != nil {
        log.Fatal(err)
    }
    return self.GetKey(), result
}

func CreatePlayer(username, Hash string, HomePlanet *Planet) Player {
    player := Player{username, Hash, HomePlanet.GetKey(), []int{0, 0}, []int{0, 0}}
    HomePlanet.Owner = username
    return player
}

