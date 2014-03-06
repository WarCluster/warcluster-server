package server

import (
	"testing"

	"github.com/garyburd/redigo/redis"

	"warcluster/entities/db"
)

const (
	user        = "{\"Command\": \"login\", \"Username\": \"JohnDoe\", \"TwitterId\": \"some twitter ID\"}"
	setupParams = "{\"Command\": \"SetupParameters\", \"Fraction\": 0, \"SunTextureId\": 0}"
)

func TestRegisterNewUser(t *testing.T) {
	conn := db.Pool.Get()
	defer conn.Close()
	conn.Do("FLUSHALL")

	session := new(testSession)
	session.Send([]byte(user))
	session.Send([]byte(setupParams))

	players_before, err := redis.Strings(conn.Do("KEYS", "player.*"))
	before := len(players_before)

	if _, err := authenticate(session); err != nil {
		t.Errorf("authenticate() failed with %s", err.Error())
	}

	players_after, err := redis.Strings(conn.Do("KEYS", "player.*"))
	after := len(players_after)

	if err != nil {
		t.Error(err)
	}

	if after != before+1 {
		t.Errorf("%#s\n", players_after)
		t.Fail()
	}
}

func TestAuthenticateExcistingUser(t *testing.T) {
	conn := db.Pool.Get()
	defer conn.Close()
	conn.Do("FLUSHALL")

	session := new(testSession)
	session.Send([]byte(user))
	session.Send([]byte(setupParams))
	session.Send([]byte(user))

	players_before, err := redis.Strings(conn.Do("KEYS", "player.*"))
	before := len(players_before)

	authenticate(session)
	authenticate(session)

	players_after, err := redis.Strings(conn.Do("KEYS", "player.*"))
	after := len(players_after)

	if err != nil {
		t.Error(err)
	}

	if after != before+1 {
		t.Fail()
	}
}

func TestAuthenticateUserWithIncompleteData(t *testing.T) {
	conn := db.Pool.Get()
	defer conn.Close()
	conn.Do("FLUSHALL")

	var session testSession
	session.Send([]byte("{\"Command\": \"login\", \"TwitterId\": \"some twitter ID\"}"))

	players_before, err := redis.Strings(conn.Do("KEYS", "player.*"))
	before := len(players_before)

	authenticate(&session)

	players_after, err := redis.Strings(conn.Do("KEYS", "player.*"))
	after := len(players_after)

	if err != nil {
		t.Error(err)
	}

	if before != after {
		t.Fail()
	}
}

func TestAuthenticateUserWithNilData(t *testing.T) {
	conn := db.Pool.Get()
	defer conn.Close()
	conn.Do("FLUSHALL")

	session := new(testSession)
	session.Send(nil)
	_, err := authenticate(session)

	if err == nil {
		t.Fail()
	}
}

func TestAuthenticateUserWithInvalidJSONData(t *testing.T) {
	conn := db.Pool.Get()
	defer conn.Close()
	conn.Do("FLUSHALL")

	session := new(testSession)
	session.Send([]byte("panda"))
	_, err := authenticate(session)

	if err == nil {
		t.Fail()
	}
}

func TestAuthenticateUserWithNilSetupData(t *testing.T) {
	conn := db.Pool.Get()
	defer conn.Close()
	conn.Do("FLUSHALL")

	session := new(testSession)
	session.Send([]byte(user))
	session.Send(nil)
	_, err := authenticate(session)

	if err == nil {
		t.Fail()
	}
}
