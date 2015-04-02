package server

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
	"github.com/pzsz/voronoi"
	"golang.org/x/net/websocket"

	"warcluster/entities"
	"warcluster/leaderboard"
	"warcluster/server/response"
)

// This function is called from the message handler to parse the first message for every new connection.
// It check for existing user in the DB and logs him if the password is correct.
// If the user is new he is initiated and a new home planet nad solar system are generated.
func login(ws *websocket.Conn) (*Client, response.Responser, error) {
	player, err := authenticate(ws)
	if err != nil {
		return nil, response.NewLoginFailed(), err
	}

	client := NewClient(ws, player)
	homePlanetEntity, err := entities.Get(player.HomePlanet)
	if err != nil {
		return nil, nil, errors.New("Player's home planet is missing!")
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
	return client, loginSuccess, nil
}

func FetchSetupData(ws *websocket.Conn) (*entities.SetupData, error) {
	var request Request

	messageStruct := response.NewLoginInformation()
	if err := websocket.JSON.Send(ws, &messageStruct); err != nil {
		return nil, err
	}

	if err := websocket.JSON.Receive(ws, &request); err != nil {
		return nil, err
	}

	accountData := new(entities.SetupData)
	if request.Command != "setup_parameters" {
		return nil, errors.New("Wrong command. Expected 'setup_parameters'")
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
func authenticate(ws *websocket.Conn) (*entities.Player, error) {
	var (
		player    *entities.Player
		nickname  string
		twitterId string
		request   Request
	)

	if err := websocket.JSON.Receive(ws, &request); err != nil {
		return nil, err
	}
	if len(request.Username) <= 0 || len(request.TwitterID) <= 0 {
		return nil, errors.New("Incomplete credentials")
	}

	serverParamsMessage := response.NewServerParams()
	if err := websocket.JSON.Send(ws, &serverParamsMessage); err != nil {
		return nil, err
	}

	nickname = request.Username
	twitterId = request.TwitterID

	entity, _ := entities.Get(fmt.Sprintf("player.%s", nickname))
	if entity == nil {
		setupData, err := FetchSetupData(ws)
		if err != nil {
			return nil, err
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
