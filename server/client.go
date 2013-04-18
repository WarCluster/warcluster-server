package server

import (
	"../db_manager"
	"../entities"
	"fmt"
	"github.com/fzzy/sockjs-go/sockjs"
)

type Client struct {
	session  sockjs.Session
	nickname string
	channel  chan string
	player   *entities.Player
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

func (self *Client) ReadLinesInto(session sockjs.Session, message []byte) {
	for {
		if request, err := UnmarshalRequest(string(message)); err == nil {
			if action, err := ParseRequest(request); err == nil {
				fmt.Println(action)
				// action(ch, self.conn, self.player, request)
			} else {
				fmt.Println(err.Error())
			}
		}
	}
}

func (self *Client) WriteLinesFrom(session sockjs.Session, message []byte) {
	// for msg := range ch {
	// 	if _, err := io.WriteString(self.conn, msg); err != nil {
	// 		return
	// 	}
	// }
}
