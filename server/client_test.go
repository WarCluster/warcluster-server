package server

import (
	"testing"

	"github.com/garyburd/redigo/redis"

	"warcluster/entities/db"
)

func TestLoginUser(t *testing.T) {
	var session testSession
	session.Send([]byte("{\"Command\": \"login\", \"Username\": \"JohnDoe\", \"TwitterId\": \"some twitter ID\"}"))
	session.Send([]byte("{\"Command\": \"SetupParameters\", \"Fraction\": 0, \"SunTextureId\": 0}"))
	conn := db.Pool.Get()
	defer conn.Close()

	players_before, err := redis.Strings(conn.Do("KEYS", "player.*"))
	before := len(players_before)

	authenticate(&session)

	players_after, err := redis.Strings(conn.Do("KEYS", "player.*"))
	after := len(players_after)

	if err != nil {
		t.Error(err)
	}

	if after == before + 1 {
		t.Fail()
	}
}

func TestLoginUserWithIncompleteData(t *testing.T) {
	var session testSession
	session.Send([]byte("{\"Command\": \"login\", \"TwitterId\": \"some twitter ID\"}"))
	conn := db.Pool.Get()
	defer conn.Close()

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
