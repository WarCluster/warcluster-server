package server

import (
	"log"
	"testing"
	"time"

	"code.google.com/p/go.net/websocket"
	"github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"warcluster/entities/db"
	"warcluster/leaderboard"
)

var (
	incompleteUser = Request{Command: "login", TwitterID: "some twitter ID"}
	user           = Request{Command: "login", Username: "JohnDoe", TwitterID: "some twitter ID"}
	setupParams    = Request{Command: "setup_parameters", Race: 0, SunTextureId: 0}
	setup          = Request{Command: "setup", Race: 0, SunTextureId: 0}
)

type ClientTestSuite struct {
	suite.Suite
	conn    redis.Conn
	ws      *websocket.Conn
	message map[string]interface{}
}

func (s *ClientTestSuite) Dial() (*websocket.Conn, error) {
	origin := "http://localhost/"
	url := "ws://localhost:7013/websocket"
	return websocket.Dial(url, "", origin)
}

func (suite *ClientTestSuite) SetupTest() {
	var err error

	suite.message = make(map[string]interface{})
	suite.conn = db.Pool.Get()
	suite.conn.Do("FLUSHDB")
	suite.ws, err = suite.Dial()
	if err != nil {
		log.Fatal(err)
	}

	cfg.Load()
	InitLeaderboard(leaderboard.New())
}

func (suite *ClientTestSuite) TearDownTest() {
	suite.ws.Close()
	suite.conn.Close()
}

func (s *ClientTestSuite) assertReceive(command string) {
	s.message = make(map[string]interface{})

	receive := func() <-chan struct{} {
		ch := make(chan struct{})
		go func() {
			websocket.JSON.Receive(s.ws, &s.message)
			ch <- struct{}{}
		}()
		return ch
	}

	select {
	case <-time.After(5 * time.Second):
		s.T().Fatalf("Did not receive %s after 5 seconds", command)
	case <-receive():
		assert.Equal(s.T(), s.message["Command"], command)
	}
}

func (s *ClientTestSuite) assertSend(request *Request) {
	send := func() <-chan struct{} {
		ch := make(chan struct{})
		go func() {
			websocket.JSON.Send(s.ws, request)
			ch <- struct{}{}
		}()
		return ch
	}

	select {
	case <-time.After(5 * time.Second):
		s.T().Fatalf("Did not send %s after 5 seconds", request.Command)
	case <-send():
	}
}

func (s *ClientTestSuite) TestRegisterNewUser() {
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

func (s *ClientTestSuite) TestAuthenticateExcistingUser() {
	s.assertSend(&user)
	s.assertReceive("server_params")
	s.assertReceive("request_setup_params")

	s.assertSend(&setupParams)
	s.assertReceive("login_success")

	s.ws.Close()
	s.Dial()

	players_before, err := redis.Strings(s.conn.Do("KEYS", "player.*"))
	before := len(players_before)
	assert.Nil(s.T(), err)

	s.assertSend(&user)
	s.assertReceive("login_success")

	players_after, err := redis.Strings(s.conn.Do("KEYS", "player.*"))
	after := len(players_after)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), before, after)
}

func (s *ClientTestSuite) TestAuthenticateUserWithIncompleteData() {
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

func (s *ClientTestSuite) TestUnableToRegisterNewUserWithWrongCommand() {

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

func (s *ClientTestSuite) TestAuthenticateUserWithNilData() {
	s.assertSend(nil)
	s.assertReceive("login_failed")
}

func (s *ClientTestSuite) TestAuthenticateUserWithInvalidJSONData() {
	websocket.Message.Send(s.ws, "panda")
	s.assertReceive("login_failed")
}

func (s *ClientTestSuite) TestAuthenticateUserWithNilSetupData() {
	s.assertSend(&user)
	s.assertReceive("server_params")
	s.assertReceive("request_setup_params")

	s.assertSend(nil)
	s.assertReceive("login_failed")
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
