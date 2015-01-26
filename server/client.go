package server

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sync"

	"github.com/ChimeraCoder/anaconda"
	"github.com/fzzy/sockjs-go/sockjs"
	"github.com/pzsz/voronoi"

	"warcluster/entities"
	"warcluster/entities/db"
	"warcluster/leaderboard"
	"warcluster/server/response"
)

// The information for each person is stored in two seperate structures. Player and Client.
// This is one of them. The purpouse of the Client struct is to hold the server(connection) information.
// 1.Session holds the curent player session socket for comunication.
// 2.Player is a pointer to the player struct for easy access.
type Client struct {
	Session     sockjs.Session
	Player      *entities.Player
	areas       map[string]struct{}
	poolElement *list.Element
	stateChange *response.StateChange
	mutex       sync.Mutex
}

func NewClient(session sockjs.Session, player *entities.Player) *Client {
	return &Client{
		Session: session,
		Player:  player,
		areas:   make(map[string]struct{}),
	}
}

// Send response directly to the client
func (c *Client) Send(response response.Responser) {
	response.Sanitize(c.Player)
	message, _ := json.Marshal(response)
	c.Session.Send(message)
}

// Send all changes to the client and flush them
func (c *Client) sendStateChange() {
	if c.stateChange != nil {
		c.Send(c.stateChange)
		c.stateChange = nil
	}
}

// Add a change to the stateChange
func (c *Client) pushStateChange(entity entities.Entity) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.stateChange == nil {
		c.stateChange = response.NewStateChange()
	}

	switch e := entity.(type) {
	case *entities.Mission:
		c.stateChange.Missions[e.Key()] = e
	case *entities.Planet:
		c.stateChange.RawPlanets[e.Key()] = e
	case *entities.Sun:
		c.stateChange.Suns[e.Key()] = e
	}
}

// Moves the client to another area
func (c *Client) MoveToAreas(areaSlice []string) {
	conn := db.Pool.Get()
	defer conn.Close()
	c.mutex.Lock()
	defer c.mutex.Unlock()

	player := c.Player.Key()
	// Create map of the new areas in order
	// to search in them more easily
	areas := make(map[string]struct{})
	for _, area := range areaSlice {
		areas[area] = empty
	}

	// Remove left areas
	for area, _ := range c.areas {
		if _, in := areas[area]; !in {
			delete(c.areas, area)
			db.Srem(conn, area, player)
		}
	}

	// Add newly occupied areas
	for area, _ := range areas {
		if _, in := c.areas[area]; !in {
			c.areas[area] = empty
			db.Sadd(conn, area, player)
		}
	}
}

// This function is called from the message handler to parse the first message for every new connection.
// It check for existing user in the DB and logs him if the password is correct.
// If the user is new he is initiated and a new home planet nad solar system are generated.
func login(session sockjs.Session) (*Client, response.Responser, error) {
	player, err := authenticate(session)
	if err != nil {
		return nil, response.NewLoginFailed(), errors.New("Login failed")
	}

	client := NewClient(session, player)
	homePlanetEntity, err := entities.Get(player.HomePlanet)
	if err != nil {
		return nil, nil, errors.New("Your home planet is missing!")
	}
	homePlanet := homePlanetEntity.(*entities.Planet)

	loginSuccess := response.NewLoginSuccess(player, homePlanet)
    
    planetEntities := entities.Find("planet.*")
	planets := make([]*entities.Planet, 0, len(planetEntities))
	sites := make([]voronoi.Vertex, 0, len(planetEntities))
	x0, xn, y0, yn := 0.0, 0.0, 0.0, 0.0
	for i, planetEntity := range planetEntities {
		planets = append(planets, planetEntity.(*entities.Planet))
		if x0 > planets[i].Position.X {
			x0 = planets[i].Position.X
		}
		if xn < planets[i].Position.X {
			xn = planets[i].Position.X
		}
		if y0 > planets[i].Position.Y {
			y0 = planets[i].Position.Y
		}
		if yn < planets[i].Position.Y {
			yn = planets[i].Position.Y
		}
		sites = append(sites, voronoi.Vertex{planets[i].Position.X, planets[i].Position.Y})
	}

	bbox := voronoi.NewBBox(x0, xn, y0, yn)

    response.Diagram = voronoi.ComputeDiagram(sites, bbox, true)
	
    return client, loginSuccess, err
}

func FetchSetupData(session sockjs.Session) (*entities.SetupData, error) {
	messageStruct := response.NewLoginInformation()
	marshalledMessage, err := json.Marshal(messageStruct)
	if err != nil {
		return nil, err
	}
	session.Send(marshalledMessage)

	request := new(Request)
	message := session.Receive()
	if message == nil {
		return nil, errors.New("No credentials provided in setup data")
	}

	if err := json.Unmarshal(message, request); err != nil {
		return nil, err
	}

	accountData := new(entities.SetupData)
	if request.Command != "setup_parameters" {
		return nil, errors.New("Wrong command")
	}

	accountData.Race = request.Race
	accountData.SunTextureId = request.SunTextureId

	if err := accountData.Validate(); err != nil {
		return nil, err
	}
	return accountData, nil
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

	message := session.Receive()
	if message == nil {
		return nil, errors.New("No credentials provided")
	}

	if err := json.Unmarshal(message, request); err != nil {
		return nil, err
	}

	if len(request.Username) <= 0 || len(request.TwitterID) <= 0 {
		return nil, errors.New("Incomplete credentials")
	}

	serverParamsMessage := response.NewServerParams()
	marshalledMessage, err := json.Marshal(serverParamsMessage)
	if err != nil {
		return nil, errors.New("Failed to provide server params.")
	}
	session.Send(marshalledMessage)

	nickname = request.Username
	twitterId = request.TwitterID

	entity, _ := entities.Get(fmt.Sprintf("player.%s", nickname))
	if entity == nil {
		setupData, err := FetchSetupData(session)
		if err != nil {
			return nil, errors.New("Reading client data failed.")
		}
		player = register(setupData, nickname, twitterId)
	} else {
		player = entity.(*entities.Player)
	}
	return player, nil
}

// Registration process is quite simple:
//
// - Gather all twitter friends.
// - Create a new sun with GenerateSun.
// - Choose home planet from the newly created solar sysitem.
// - Create a reccord of the new player and start comunication.
func register(setupData *entities.SetupData, nickname, twitterId string) *entities.Player {
	friendsSuns := fetchFriendsSuns(nickname)
	sun := entities.GenerateSun(nickname, friendsSuns, setupData)
	planets, homePlanet := entities.GeneratePlanets(nickname, sun)
	player := entities.CreatePlayer(nickname, twitterId, homePlanet, setupData)

	for _, planet := range planets {
		entities.Save(planet)
		clients.Broadcast(planet)
	}

	entities.Save(player)
	entities.Save(sun)

	clients.Broadcast(sun)
	leaderBoard.Add(&leaderboard.Player{
		Username:   player.Username,
		RaceId:     player.RaceID,
		HomePlanet: homePlanet.Name,
		Planets:    1,
	})
	return player
}

// Returns a slice with twitter ids of the given user's friends
func fetchTwitterFriends(screenName string) ([]string, error) {
	anaconda.SetConsumerKey(cfg.Twitter.ConsumerKey)
	anaconda.SetConsumerSecret(cfg.Twitter.ConsumerSecret)
	api := anaconda.NewTwitterApi(cfg.Twitter.AccessToken, cfg.Twitter.AccessTokenSecret)

	v := url.Values{}
	v.Set("count", "100")
	v.Set("cursor", "-1")
	v.Set("screen_name", screenName)

	friendsIds, err := api.GetFriendsIds(v)
	if err != nil {
		return nil, err
	}

	friends, err := api.GetUsersLookupByIds(friendsIds.Ids, url.Values{})
	if err != nil {
		return nil, err
	}

	var friendsNames []string
	for _, friend := range friends {
		friendsNames = append(friendsNames, friend.ScreenName)
	}

	return friendsNames, nil
}

// Returns a slice of friend's suns
func fetchFriendsSuns(twitterName string) (suns []*entities.Sun) {
	friendsNames, twitterErr := fetchTwitterFriends(twitterName)
	if twitterErr != nil {
		return
	}

	for _, name := range friendsNames {
		playerEntity, err := entities.Get(fmt.Sprintf("player.%s", name))
		if playerEntity == nil || err != nil {
			continue
		}

		player := playerEntity.(*entities.Player)
		if err != nil {
			continue
		}

		if friendSun, sunErr := entities.Get(player.Sun()); sunErr == nil {
			suns = append(suns, friendSun.(*entities.Sun))
		}
	}
	return
}
