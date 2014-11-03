package server

import (
	"container/list"
	"testing"

	"code.google.com/p/go.net/websocket"
	"github.com/Vladimiroff/vec2d"

	"warcluster/entities"
)

var cp = NewClientPool(16)

var planet = entities.Planet{
	Name:     "GOP6720",
	Position: &vec2d.Vector{2, 2},
}

var player1 = entities.Player{
	Username:       "gophie",
	RaceID:         1,
	TwitterID:      "gop",
	HomePlanet:     "planet.GOP6720",
	ScreenSize:     []uint64{1, 1},
	ScreenPosition: &vec2d.Vector{2, 2},
}

var player2 = entities.Player{
	Username:       "snoopy",
	RaceID:         1,
	TwitterID:      "snoop",
	HomePlanet:     "planet.SNO6750",
	ScreenSize:     []uint64{1, 1},
	ScreenPosition: &vec2d.Vector{2, 8},
}

var client1 = *NewClient(new(websocket.Conn), &player1)

var client2 = *NewClient(new(websocket.Conn), &player1)

var client3 = *NewClient(new(websocket.Conn), &player2)

var client4 = *NewClient(new(websocket.Conn), &player2)

func TestAddClientToClientPool(t *testing.T) {
	cp.pool = make(map[string]*list.List)

	l := len(cp.pool)
	cp.Add(&client1)
	if len(cp.pool) != l+1 {
		t.Fail()
	}
	cp.Remove(&client1)
}

func TestCloseClientSessionWithMoreThanOneSessions(t *testing.T) {
	cp.pool = make(map[string]*list.List)

	cp.Add(&client1)
	cp.Add(&client2)
	l := len(cp.pool)

	cp.Remove(&client1)
	if len(cp.pool) != l {
		t.Errorf("Expected %d received %d", l, len(cp.pool))
	}
	cp.Remove(&client2)
}

func TestCloseLastClientSessionAndRemoveIt(t *testing.T) {
	cp.pool = make(map[string]*list.List)

	cp.Add(&client3)
	l := len(cp.pool)

	cp.Remove(&client3)
	if len(cp.pool) != l-1 {
		t.Errorf("Expected %d received %d", l-1, len(cp.pool))
	}
}

func TestRemoveUnexistingClient(t *testing.T) {
	cp.pool = make(map[string]*list.List)

	cp.Remove(&client1)
	if len(cp.pool) != 0 {
		t.Fail()
	}
}

// TODO
/* func TestSendMessageToSession(t *testing.T) { */
/* 	cp.pool = make(map[string]*list.List) */
/* 	resp := response.NewSendMissions() */
/* 	cp.Send(&player1, resp) */

/* 	cp.Add(&client1) */
/* 	cp.Add(&client2) */
/* 	cp.Add(&client3) */

/* 	l1 := len(client1.Session.(*testSession).Messages) */
/* 	l2 := len(client2.Session.(*testSession).Messages) */
/* 	l3 := len(client3.Session.(*testSession).Messages) */
/* 	cp.Send(&player1, resp) */

/* 	if len(client1.Session.(*testSession).Messages) != l1+1 { */
/* 		t.Errorf("%d", len(client1.Session.(*testSession).Messages)) */
/* 	} */

/* 	if len(client2.Session.(*testSession).Messages) != l2+1 { */
/* 		t.Fail() */
/* 	} */

/* 	if len(client3.Session.(*testSession).Messages) != l3 { */
/* 		t.Fail() */
/* 	} */

/* } */

func TestPlayer(t *testing.T) {
	cp.pool = make(map[string]*list.List)
	cp.Add(&client1)

	if player, err := cp.Player("gophie"); player == nil || err != nil {
		t.Error(err)
	}

	if player, err := cp.Player("snoopy"); player != nil || err == nil {
		t.Errorf("Received %v as player, nil expected", player.Username)
	}
}

func TestStackingStateChanges(t *testing.T) {
	cp := NewClientPool(1)
	cp.ticker.Stop()
	cp.Add(&client1)

	client1.pushStateChange(&planet)
	// TODO
	/* if len(client1.Session.Receive()) != 0 { */
	/* 	t.Error("Client received messages  without a tick") */
	/* } */

	if client1.stateChange == nil {
		t.Error("Client has no stacked planets")
	}

	if len(client1.stateChange.RawPlanets) != 1 {
		t.Errorf("Client has %d stacked planets instead of 1")
	}
}

func TestBroadcast(t *testing.T) {
	cp := NewClientPool(3)
	cp.ticker.Stop()
	cp.Add(&client1)
	cp.Add(&client2)
	cp.Add(&client3)

	client1.MoveToAreas([]string{"area:1:1"})
	client2.MoveToAreas([]string{"area:1:1"})
	client3.MoveToAreas([]string{"area:2:1"})

	cp.Broadcast(&planet)
	if client1.stateChange == nil {
		t.Error("Client1 has no stacked planets")
	}

	if len(client1.stateChange.RawPlanets) != 1 {
		t.Errorf("Client1 has %d stacked planets instead of 1")
	}

	if client2.stateChange == nil {
		t.Error("Client2 has no stacked planets")
	}

	if len(client2.stateChange.RawPlanets) != 1 {
		t.Errorf("Client2 has %d stacked planets instead of 1")
	}

	if client3.stateChange != nil {
		t.Error("Client3 has stacked planets")
	}
}
