package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/fzzy/sockjs-go/sockjs"

	"warcluster/entities"
	"warcluster/server/response"
)

// The information for each person is stored in two seperate structures. Player and Client.
// This is one of them. The purpouse of the Client struct is to hold the server(connection) information.
// 1.Session holds the curent player session socket for comunication.
// 2.Player is a pointer to the player struct for easy access.
type Client struct {
	Session sockjs.Session
	Player  *entities.Player
}

// Authenticate is a function called for every client's new session.
// It manages several important tasks at the start of the session.
// 1.Ask the user for Username and twitter ID.
// 2.Search the DB to find the player if it's not a new one.
// 3.If the player is new there is a subsequence initiated:
// 3.1.Create a new sun with GenerateSun
// 3.2.Choose home planet from the newly created solar sysitem.
// 3.3.Create a reccord of the new player and start comunication.
func authenticate(session sockjs.Session) (*entities.Player, bool, error) {
	var player *entities.Player
	var nickname string
	var twitterId string
	request := new(Request)

	for {
		if message := session.Receive(); message == nil {
			return nil, false, errors.New("No credentials provided")
		} else {
			if err := json.Unmarshal(message, request); err == nil {
				if len(request.Username) > 0 && len(request.TwitterID) > 0 {
					nickname = request.Username
					twitterId = request.TwitterID
					break
				}
			} else {
				log.Print("Error in server.client.authenticate: ", err.Error())
			}
		}
	}

	entity, _ := entities.Get(fmt.Sprintf("player.%s", nickname))
	justRegistered := entity == nil
	if justRegistered {
		allSunsEntities := entities.Find("sun.*")
		allSuns := []entities.Sun{}
		for _, entity := range allSunsEntities {
			allSuns = append(allSuns, *entity.(*entities.Sun))
		}
		sun := entities.GenerateSun(nickname, allSuns, []entities.Sun{})
		planets, homePlanet := entities.GeneratePlanets(nickname, sun)
		player = entities.CreatePlayer(nickname, twitterId, homePlanet)

		//TODO: Remove the bottom three lines when the client is smart enough to invoke
		//      scope of view on all clients in order to osee the generated system
		for _, planet := range planets {
			entities.Save(planet)
			stateChange := response.NewStateChange()
			stateChange.Planets = map[string]entities.Entity{
				planet.Key(): planet,
			}
			response.Send(stateChange, sessions.Broadcast)
		}

		entities.Save(player)
		entities.Save(sun)

		stateChange := response.NewStateChange()
		stateChange.Suns = map[string]entities.Entity{
			sun.Key(): sun,
		}
		response.Send(stateChange, sessions.Broadcast)
	} else {
		player = entity.(*entities.Player)
	}
	return player, justRegistered, nil
}
