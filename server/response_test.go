package server

import (
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/Vladimiroff/vec2d"

	"warcluster/entities"
	"warcluster/entities/db"
)

var gophie entities.Player = entities.Player{
	Username:       "gophie",
	Color:          entities.Color{22, 22, 22},
	TwitterID:      "gophie92",
	HomePlanet:     "planet.GOP6720",
	ScreenSize:     []uint16{1, 1},
	ScreenPosition: &vec2d.Vector{2, 2},
}

var panda entities.Player = entities.Player{
	Username:       "panda",
	Color:          entities.Color{22, 22, 22},
	TwitterID:      "panda13",
	HomePlanet:     "planet.PAN6720",
	ScreenSize:     []uint16{1, 1},
	ScreenPosition: &vec2d.Vector{2, 2},
}

var client Client = Client{
	Session:        new(testSession),
	Player:         &gophie,
}

var planet1 entities.Planet = entities.Planet{
	Name:           "GOP6720",
	Position:       &vec2d.Vector{2, 2},
	IsHome:         true,
	Owner:          "gophie",
}

var planet2 entities.Planet = entities.Planet{
	Name:           "GOP6724",
	Position:       &vec2d.Vector{4, 4},
	IsHome:         false,
	Owner:          "gophie",
}

var planet3 entities.Planet = entities.Planet{
	Name:           "PAN6720",
	Position:       &vec2d.Vector{10, 10},
	IsHome:         true,
	Owner:          "panda",
}

type ResponseTestSuite struct {
	suite.Suite
	conn    redis.Conn
	request *Request
}

func (suite *ResponseTestSuite) SetupTest() {
	suite.conn = db.Pool.Get()
	suite.conn.Do("FLUSHDB")
	entities.Save(&planet1)
	entities.Save(&planet2)
	entities.Save(&planet3)
	entities.Save(&gophie)
	entities.Save(&panda)

	suite.request = new(Request)
	suite.request.Command = "start_mission"
	suite.request.StartPlanet = "planet.GOP6720"
	suite.request.EndPlanet = "planet.PAN6720"
	suite.request.Position = vec2d.New(2.0, 4.0)
	suite.request.Resolution = []uint16{1920, 1080}
	suite.request.Fleet = 32
	suite.request.Username = "panda"
	suite.request.TwitterID = "panda13"
	suite.request.Fraction = 4
	suite.request.SunTextureId = 0
	suite.request.Client = &client
	suite.request.Type = "spy"
}

func (suite *ResponseTestSuite) TestParseActionWithoutStartPlanet() {
	suite.request.StartPlanet = ""

	err := parseAction(suite.request)

	assert.NotNil(suite.T(), err)

	suite.request.StartPlanet = planet1.Name
}

func (suite *ResponseTestSuite) TestParseActionWithoutEndPlanet() {
	suite.request.EndPlanet = ""

	err := parseAction(suite.request)

	assert.NotNil(suite.T(), err)

	suite.request.EndPlanet = planet3.Name
}

func (suite *ResponseTestSuite) TestParseActionFromForeignPlanet() {
	err := parseAction(suite.request)

	assert.NotNil(suite.T(), err)
}

func TestResponseTestSuite(t *testing.T) {
	suite.Run(t, new(ResponseTestSuite))
}
