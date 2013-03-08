package server

import (
	"../db_manager"
	"../entities"
	"bufio"
	"fmt"
	"io"
	"net"
)

type Client struct {
	conn     net.Conn
	nickname string
	channel  chan string
	player   *entities.Player
}

func authenticate(c net.Conn, bufc *bufio.Reader) (string, *entities.Player) {
	var player entities.Player

	io.WriteString(c, "Twitter Authenticating:\n")
	io.WriteString(c, "Username: ")
	nick, _, _ := bufc.ReadLine()
	nickname := string(nick)

	io.WriteString(c, "TwitterID: ")
	twitter, _, _ := bufc.ReadLine()
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

func (self *Client) ReadLinesInto(ch chan<- string) {
	bufc := bufio.NewReader(self.conn)
	for {
		line, err := bufc.ReadString('\n')
		if err != nil {
			break
		}

		if request, err := UnmarshalRequest(line); err == nil {
			if action, err := ParseRequest(request); err == nil {
				action(ch, self.conn, self.player, request)
			} else {
				fmt.Println(err.Error())
			}
		}
	}
}

func (self *Client) WriteLinesFrom(ch <-chan string) {
	for msg := range ch {
		if _, err := io.WriteString(self.conn, msg); err != nil {
			return
		}
	}
}
