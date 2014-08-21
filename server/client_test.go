package server

import (
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"warcluster/entities/db"
	"warcluster/leaderboard"
)

const (
	user        = "{\"Command\": \"login\", \"Username\": \"JohnDoe\", \"TwitterId\": \"some twitter ID\"}"
	setupParams = "{\"Command\": \"setup_parameters\", \"Race\": 0, \"SunTextureId\": 0}"
)

type ClientTestSuite struct {
	suite.Suite
	conn    redis.Conn
	session *testSession
}

func (suite *ClientTestSuite) SetupTest() {
	suite.conn = db.Pool.Get()
	suite.conn.Do("FLUSHDB")
	suite.session = new(testSession)

	cfg.Load()
	InitLeaderboard(leaderboard.New())
}

func (suite *ClientTestSuite) TearDownTest() {
	suite.conn.Close()
}

func (suite *ClientTestSuite) TestRegisterNewUser() {
	suite.session.Send([]byte(user))
	suite.session.Send([]byte(setupParams))

	players_before, err := redis.Strings(suite.conn.Do("KEYS", "player.*"))
	before := len(players_before)

	_, err = authenticate(suite.session)

	assert.Nil(suite.T(), err)

	players_after, err := redis.Strings(suite.conn.Do("KEYS", "player.*"))
	after := len(players_after)

	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), before+1, after)
}

func (suite *ClientTestSuite) TestAuthenticateExcistingUser() {
	suite.session.Send([]byte(user))
	suite.session.Send([]byte(setupParams))
	suite.session.Send([]byte(user))

	players_before, err := redis.Strings(suite.conn.Do("KEYS", "player.*"))
	before := len(players_before)

	authenticate(suite.session)
	authenticate(suite.session)

	players_after, err := redis.Strings(suite.conn.Do("KEYS", "player.*"))
	after := len(players_after)

	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), before+1, after)
}

func (suite *ClientTestSuite) TestAuthenticateUserWithIncompleteData() {
	suite.session.Send([]byte("{\"Command\": \"login\", \"TwitterId\": \"some twitter ID\"}"))

	players_before, err := redis.Strings(suite.conn.Do("KEYS", "player.*"))
	before := len(players_before)

	authenticate(suite.session)

	players_after, err := redis.Strings(suite.conn.Do("KEYS", "player.*"))
	after := len(players_after)

	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), before, after)
}

func (suite *ClientTestSuite) TestUnableToRegisterNewUserWithWrongCommand() {
	setup := "{\"Command\": \"setup\", \"Race\": 0, \"SunTextureId\": 0}"

	suite.session.Send([]byte(user))
	suite.session.Send([]byte(setup))

	players_before, err := redis.Strings(suite.conn.Do("KEYS", "player.*"))
	before := len(players_before)

	_, err = authenticate(suite.session)

	assert.NotNil(suite.T(), err)

	players_after, err := redis.Strings(suite.conn.Do("KEYS", "player.*"))
	after := len(players_after)

	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), before, after)
}

func (suite *ClientTestSuite) TestAuthenticateUserWithNilData() {
	suite.session.Send(nil)
	_, err := authenticate(suite.session)

	assert.NotNil(suite.T(), err)
}

func (suite *ClientTestSuite) TestAuthenticateUserWithInvalidJSONData() {
	suite.session.Send([]byte("panda"))
	_, err := authenticate(suite.session)

	assert.NotNil(suite.T(), err)
}

func (suite *ClientTestSuite) TestAuthenticateUserWithNilSetupData() {
	suite.session.Send([]byte(user))
	suite.session.Send(nil)
	_, err := authenticate(suite.session)

	assert.NotNil(suite.T(), err)
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
