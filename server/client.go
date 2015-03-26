package server

import (
	"container/list"
	"sync"

	"golang.org/x/net/websocket"
	"warcluster/entities"
	"warcluster/entities/db"
	"warcluster/server/response"
)

// Codec is implemented by objects that send and receive via websocket
type Codec interface {
	Receive(ws *websocket.Conn, v interface{}) (err error)
	Send(ws *websocket.Conn, v interface{}) (err error)
}

// The information for each person is stored in two seperate structures. Player and Client.
// This is one of them. The purpouse of the Client struct is to hold the server(connection) information.
// 1.Session holds the curent player session socket for comunication.
// 2.Player is a pointer to the player struct for easy access.
type Client struct {
	Conn        *websocket.Conn
	Player      *entities.Player
	areas       map[string]struct{}
	poolElement *list.Element
	stateChange *response.StateChange
	mutex       sync.Mutex
	codec       Codec
	twitter     *anaconda.TwitterApi
}

func NewClient(ws *websocket.Conn, player *entities.Player, twitter *anaconda.TwitterApi) *Client {
	return &Client{
		Conn:    ws,
		Player:  player,
		areas:   make(map[string]struct{}),
		codec:   websocket.JSON,
		twitter: twitter,
	}
}

// Send response directly to the client
func (c *Client) Send(response response.Responser) {
	response.Sanitize(c.Player)
	c.codec.Send(c.Conn, &response)
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
