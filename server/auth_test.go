package server

import (
	"testing"

	"github.com/Vladimiroff/vec2d"
	"github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/websocket"

	"warcluster/entities"
)

type AuthTest struct {
	WebSocketTestSuite
}

func (s *AuthTest) TestRegisterNewUser() {
	players_before, err := redis.Strings(s.conn.Do("KEYS", "player.*"))
	assert.Nil(s.T(), err)
	before := len(players_before)

	s.assertSend(&user)
	s.assertReceive("server_params")
	s.assertReceive("request_setup_params")

	s.assertSend(&setupParams)
	s.assertReceive("login_success")

	players_after, err := redis.Strings(s.conn.Do("KEYS", "player.*"))
	assert.Nil(s.T(), err)
	after := len(players_after)

	assert.Equal(s.T(), before+1, after)
}

func (s *AuthTest) TestAuthenticateExcistingUser() {
	entities.Save(&entities.Planet{
		Name:     "GOP6720",
		Position: &vec2d.Vector{2, 2},
	})
	entities.Save(&entities.Player{
		Username:       "gophie",
		RaceID:         1,
		TwitterID:      "gop",
		HomePlanet:     "planet.GOP6720",
		ScreenSize:     []uint64{1, 1},
		ScreenPosition: &vec2d.Vector{2, 2},
	})

	players_before, err := redis.Strings(s.conn.Do("KEYS", "player.*"))
	before := len(players_before)
	assert.Nil(s.T(), err)

	s.assertSend(&user)
	s.assertReceive("server_params")
	s.assertReceive("login_success")

	players_after, err := redis.Strings(s.conn.Do("KEYS", "player.*"))
	after := len(players_after)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), before, after)
}

func (s *AuthTest) TestAuthenticateUserWithIncompleteData() {
	players_before, err := redis.Strings(s.conn.Do("KEYS", "player.*"))
	before := len(players_before)
	assert.Nil(s.T(), err)

	s.assertSend(&incompleteUser)
	s.assertReceive("login_failed")

	players_after, err := redis.Strings(s.conn.Do("KEYS", "player.*"))
	after := len(players_after)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), before, after)
}

func (s *AuthTest) TestUnableToRegisterNewUserWithWrongCommand() {
	players_before, err := redis.Strings(s.conn.Do("KEYS", "player.*"))
	before := len(players_before)
	assert.Nil(s.T(), err)

	s.assertSend(&user)
	s.assertReceive("server_params")
	s.assertReceive("request_setup_params")

	s.assertSend(&setup)
	s.assertReceive("login_failed")

	players_after, err := redis.Strings(s.conn.Do("KEYS", "player.*"))
	after := len(players_after)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), before, after)
}

func (s *AuthTest) TestAuthenticateUserWithNilData() {
	s.assertSend(nil)
	s.assertReceive("login_failed")
}

func (s *AuthTest) TestAuthenticateUserWithInvalidJSONData() {
	websocket.Message.Send(s.ws, "panda")
	s.assertReceive("login_failed")
}

func (s *AuthTest) TestAuthenticateUserWithNilSetupData() {
	s.assertSend(&user)
	s.assertReceive("server_params")
	s.assertReceive("request_setup_params")

	s.assertSend(nil)
	s.assertReceive("login_failed")
}

func TestAuthTest(t *testing.T) {
	suite.Run(t, new(AuthTest))
}
