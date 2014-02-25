package server

import (
	"container/list"
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
	Session     sockjs.Session
	Player      *entities.Player
	poolElement *list.Element
}

// This function is called from the message handler to parse the first message for every new connection.
// It check for existing user in the DB and logs him if the password is correct.
// If the user is new he is initiated and a new home planet nad solar system are generated.
func login(session sockjs.Session) (*Client, response.Responser, error) {
	player, err := authenticate(session)
	if err != nil {
		log.Println(err)
		return nil, response.NewLoginFailed(), errors.New("Login failed")
	}

	client := &Client{
		Session: session,
		Player:  player,
	}
	homePlanetEntity, err := entities.Get(player.HomePlanet)
	if err != nil {
		return nil, nil, errors.New("Your home planet is missing!")
	}
	homePlanet := homePlanetEntity.(*entities.Planet)

	loginSuccess := response.NewLoginSuccess(player, homePlanet)
	return client, loginSuccess, err
}

func FetchSetupData(session sockjs.Session) (*entities.SetupData, error) {
	var accountData *entities.SetupData
	accountData = new(entities.SetupData)

	messageStruct := response.NewLoginInformation()
	message, err := json.Marshal(messageStruct)
	if err != nil {
		log.Println(err.Error())
	}
	session.Send(message)

	request := new(Request)

	for {
		if message := session.Receive(); message == nil {
			return nil, errors.New("No credentials provided")
		} else {
			if err := json.Unmarshal(message, request); err == nil {
				if len(request.Command) > 0 && request.Command == "SetupParameters" {
					accountData.Fraction = request.Fraction
					accountData.SunTextureId = request.SunTextureId
					break
				}
			} else {
				log.Print("Error in server.client.authenticate: ", err.Error())
			}
		}
	}

	err = accountData.Validate()

	return accountData, err
}

// Authenticate is a function called for every client's new session.
// It manages several important tasks at the start of the session.
// 1.Ask the user for Username and twitter ID.
// 2.Search the DB to find the player if it's not a new one.
// 3.If the player is new there is a subsequence initiated:
// 3.1.Create a new sun with GenerateSun
// 3.2.Choose home planet from the newly created solar sysitem.
// 3.3.Create a reccord of the new player and start comunication.
func authenticate(session sockjs.Session) (*entities.Player, error) {
	var player *entities.Player
	var nickname string
	var twitterId string
	request := new(Request)

	for {
		if message := session.Receive(); message == nil {
			return nil, errors.New("No credentials provided")
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

		setupInfo, err := FetchSetupData(session)
		if err != nil {
			log.Println(err)
			return nil, errors.New("Reading client data failed.")
		}

		allSunsEntities := entities.Find("sun.*")
		allSuns := []*entities.Sun{}
		for _, sunEntity := range allSunsEntities {
			allSuns = append(allSuns, sunEntity.(*entities.Sun))
		}
		sun := entities.GenerateSun(nickname, allSuns, []*entities.Sun{}, setupInfo)
		planets, homePlanet := entities.GeneratePlanets(nickname, sun)
		player = entities.CreatePlayer(nickname, twitterId, homePlanet, setupInfo)

		//TODO: Remove the bottom three lines when the client is smart enough to invoke
		//      scope of view on all clients in order to osee the generated system
		for _, planet := range planets {
			entities.Save(planet)
			stateChange := response.NewStateChange()
			stateChange.RawPlanets = map[string]*entities.Planet{
				planet.Key(): planet,
			}
			clients.BroadcastToAll(stateChange)
		}

		entities.Save(player)
		entities.Save(sun)

		stateChange := response.NewStateChange()
		stateChange.Suns = map[string]*entities.Sun{
			sun.Key(): sun,
		}
		clients.BroadcastToAll(stateChange)
	} else {
		player = entity.(*entities.Player)
	}
	return player, nil
}
