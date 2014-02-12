package server

import (
	"fmt"
	"testing"
	"math/rand"
	"time"

	"github.com/fzzy/sockjs-go/sockjs"
	"github.com/Vladimiroff/vec2d"

	"warcluster/entities"
)


var cp = NewClientPool()

var player1 = entities.Player{
	Username:       "gophie",
	Color:          entities.Color{0.59215686, 0.59215686, 0.59215686},
	TwitterID:      "gop",
	HomePlanet:     "planet.GOP6720",
	ScreenSize:     []uint16{1, 1},
	ScreenPosition: &vec2d.Vector{2, 2},
}

var player2 = entities.Player{
	Username:       "snoopy",
	Color:          entities.Color{0.59215686, 0.59215686, 0.59215686},
	TwitterID:      "snoop",
	HomePlanet:     "planet.SNO6750",
	ScreenSize:     []uint16{1, 1},
	ScreenPosition: &vec2d.Vector{2, 8},
}

type testSession struct {
	session_id string
}

func (s *testSession) Receive() (m []byte) {
	return []byte{}
}

func (s *testSession) Send(m []byte) {
	return
}

func (s *testSession) Close(code int, reason string) {
	return
}

func (s *testSession) End() {
	return
}

func (s *testSession) Info() sockjs.RequestInfo {
	return *new(sockjs.RequestInfo)
}

func (s *testSession) Protocol() sockjs.Protocol {
	return *new(sockjs.Protocol)
}

func (s *testSession) String() string {
	rand.Seed(time.Now().Unix())
	return fmt.Sprintf("%d", rand.Int())
	// return session_id
}

var client1 = Client{
	Session: new(testSession),
	Player:  &player1,
}

var client2 = Client{
	Session: new(testSession),
	Player:  &player1,
}

var client3 = Client{
	Session: new(testSession),
	Player:  &player2,
}

var client4 = Client{
	Session: new(testSession),
	Player:  &player2,
}

func TestAddClientToClientPool(t *testing.T) {
	l := len(cp.pool)
	cp.Add(&client1)
	if len(cp.pool) != l + 1 {
		t.Fail()
	}
}

func TestCloseClientSessionWithMoreThanOneSessions(t *testing.T) {
	cp.Add(&client1)
	cp.Add(&client2)
	l := len(cp.pool)

	cp.Remove(&client1)
	if len(cp.pool) != l {
		t.Errorf("Expected %d received %d", l, len(cp.pool))
	}
}

func TestCloseLastClientSessionAndRemoveIt(t *testing.T) {
	cp.Add(&client3)
	l := len(cp.pool)

	cp.Remove(&client3)
	if len(cp.pool) != l - 1 {
		t.Errorf("Expected %d received %d", l - 1, len(cp.pool))
	}
}