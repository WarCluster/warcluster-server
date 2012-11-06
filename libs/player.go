package libs


type Player struct {
    username string
    hash string
    home_planet *Planet
    planets []*Planet
    missions []*Mission
}

func (player *Player) String() string {
    return player.hash
}

func CreatePlayer(username, hash string, home_planet *Planet) Player {
    planets := []*Planet{home_planet}
    missions := []*Mission{}
    player := Player{username, hash, home_planet, planets, missions}
    home_planet.owner = &player
    return player
}

