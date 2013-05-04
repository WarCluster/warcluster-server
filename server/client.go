package server

import (
	"warcluster/db_manager"
	"warcluster/entities"
	"fmt"
	"github.com/fzzy/sockjs-go/sockjs"
)

type Client struct {
	Session  sockjs.Session
	Nickname string
	Player   *entities.Player
}

func authenticate(session sockjs.Session) (string, *entities.Player) {
	var player entities.Player

	session.Send([]byte("Twitter Authenticating:\n"))
	session.Send([]byte("Username: "))
	nick := session.Receive()
	nickname := string(nick)

	session.Send([]byte("TwitterID: "))
	twitter := session.Receive()
	twitter_id := string(twitter)

	entity, _ := db_manager.GetEntity(fmt.Sprintf("player.%s", nick))
	if entity == nil {
		all_suns_entities := db_manager.GetEntities("sun.*")
		all_suns := []entities.Sun{}
		for _, entity := range all_suns_entities {
			all_suns = append(all_suns, entity.(entities.Sun))
		}
		sun := entities.GenerateSun(nickname, all_suns, []entities.Sun{})
		hash := entities.GenerateHash(nickname)
		planets, home_planet := entities.GeneratePlanets(hash, sun.GetPosition())
		player = entities.CreatePlayer(nickname, twitter_id, home_planet)
		db_manager.SetEntity(player)
		db_manager.SetEntity(sun)
		for i := 0; i < len(planets); i++ {
			db_manager.SetEntity(planets[i])
		}
	} else {
		player = entity.(entities.Player)
	}
	return nickname, &player
}
