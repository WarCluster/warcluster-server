package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fzzy/sockjs-go/sockjs"
	"log"
	"warcluster/db_manager"
	"warcluster/entities"
)

// The information for each person is stored in two seperate structures. Player and Client.
// This is one of them. The purpouse of the Client struct is to hold the server(connection) information.
// 1.Session holds the curent player session socket for comunication.
// 2.Nickname is used as a key because its the common data between the two structures and the database.
// 3.This is a pointer to the player struct for easy acsess.
type Client struct {
	Session  sockjs.Session
	Nickname string
	Player   *entities.Player
}

// Authenticate is a function called for every client's new session.
// It manages several important tasks at the start of the session.
// 1.Ask the user for Username and twitter ID.
// 2.Search the DB to find the player if it's not a new one.
// 3.If the player is new there is a subsequence initiated:
// 3.1.Create a new sun with GenerateSun
// 3.2.Choose home planet from the newly created solar sysitem.
// 3.3.Create a reccord of the new player and start comunication.
func authenticate(session sockjs.Session) (string, *entities.Player, error) {
	var player *entities.Player
	var nickname string
	var twitter_id string
	request := new(Request)

	for {
		if message := session.Receive(); message == nil {
			return "", nil, errors.New("No credentials provided")
		} else {
			if err := json.Unmarshal(message, request); err == nil {
				if len(request.Username) > 0 && len(request.TwitterID) > 0 {
					nickname = request.Username
					twitter_id = request.TwitterID
					break
				}
			} else {
				log.Print("Error in server.client.authenticate: ", err.Error())
			}
		}
	}

	entity, _ := db_manager.GetEntity(fmt.Sprintf("player.%s", []byte(nickname)))
	if entity == nil {
		all_suns_entities := db_manager.GetEntities("sun.*")
		all_suns := []entities.Sun{}
		for _, entity := range all_suns_entities {
			all_suns = append(all_suns, *entity.(*entities.Sun))
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
		player = entity.(*entities.Player)
	}
	return nickname, player, nil
}
