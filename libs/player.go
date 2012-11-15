package libs


type Player struct {
    username string
    Hash string
    HomePlanet *Planet
    ScreenSize []int
    ScreenPosition []int

}

func (player *Player) String() string {
    return player.username
}

func CreatePlayer(username, Hash string, HomePlanet *Planet) Player {
    player := Player{username, Hash, HomePlanet, []int{0, 0}, []int{0, 0}}
    HomePlanet.Owner = username
    return player
}

