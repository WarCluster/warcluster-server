package server

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

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

// Thread-safe pool of all clients, with opened sockets.
type ClientPool struct {
	mutex sync.RWMutex
	pool  map[string]*list.List
}

func NewClientPool() *ClientPool {
	cp := new(ClientPool)
	cp.pool = make(map[string]*list.List)
	return cp
}

// Adds the given client to the pool.
func (cp *ClientPool) Add(client *Client) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	element := new(list.Element)
	element.Value = client

	_, ok := cp.pool[client.Player.Username]
	if !ok {
		cp.pool[client.Player.Username] = list.New()
	}
	cp.pool[client.Player.Username].PushBack(element)
}

// Remove the client to the pool.
// It is safe to remove non-existing client.
func (cp *ClientPool) Remove(client *Client) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	element := new(list.Element)
	element.Value = client
	cp.pool[client.Player.Username].Remove(element)

	if cp.pool[client.Player.Username].Len() == 0 {
		delete(cp.pool, client.Player.Username)
	}
}

// Broadcast sends the given message to every session in the pool.
func (cp *ClientPool) Broadcast(m []byte) {
	cp.mutex.RLock()
	defer cp.mutex.RUnlock()

	for username := range cp.pool {
		for element := cp.pool[username].Front(); element != nil; element = element.Next() {
			client := element.Value.(*Client)
			client.Session.Send(m)
		}
	}
}

// This function is called from the message handler to parse the first message for every new connection.
// It check for existing user in the DB and logs him if the password is correct.
// If the user is new he is initiated and a new home planet nad solar system are generated.
func login(session sockjs.Session) (*Client, error) {
	player, justRegistered, err := authenticate(session)
	if err != nil {
		response.Send(response.NewLoginFailed(), session.Send)
		log.Println(err)
		return nil, errors.New("Login failed")
	}

	client := &Client{
		Session: session,
		Player:  player,
	}
	homePlanetEntity, err := entities.Get(player.HomePlanet)
	if err != nil {
		return nil, errors.New("Your home planet is missing!")
	}
	homePlanet := homePlanetEntity.(*entities.Planet)

	loginSuccess := response.NewLoginSuccess(player, homePlanet, justRegistered)
	response.Send(loginSuccess, session.Send)
	return client, nil
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
		allSuns := []*entities.Sun{}
		for _, entity := range allSunsEntities {
			allSuns = append(allSuns, entity.(*entities.Sun))
		}
		sun := entities.GenerateSun(nickname, allSuns, []*entities.Sun{})
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
			response.Send(stateChange, clients.Broadcast)
		}

		entities.Save(player)
		entities.Save(sun)

		stateChange := response.NewStateChange()
		stateChange.Suns = map[string]entities.Entity{
			sun.Key(): sun,
		}
		response.Send(stateChange, clients.Broadcast)
	} else {
		player = entity.(*entities.Player)
	}
	return player, justRegistered, nil
}
