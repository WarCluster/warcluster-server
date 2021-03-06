package server

import (
	"container/list"
	"testing"

	"github.com/Vladimiroff/vec2d"

	"warcluster/entities"
	"warcluster/server/response"
)

var (
	cp = NewClientPool(16)

	planet = entities.Planet{
		Name:     "GOP6720",
		Position: &vec2d.Vector{2, 2},
	}

	sun = entities.Sun{
		Name:     "GOP672",
		Username: "gophie",
		Position: vec2d.New(20, 20),
	}
)

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

var (
	client1 = *NewFakeClient(&player1)
	client2 = *NewFakeClient(&player1)
	client3 = *NewFakeClient(&player2)
	client4 = *NewFakeClient(&player2)
)

func TestAddClientToClientPool(t *testing.T) {
	cp.pool = make(map[string]*list.List)

	l := len(cp.pool)
	cp.Add(&client1)
	if len(cp.pool) != l+1 {
		t.Fail()
	}
	client1.MoveToAreas([]string{"area:1:1"})
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

func TestSendMessageToSession(t *testing.T) {
	cp.pool = make(map[string]*list.List)
	resp := response.NewSendMissions()
	cp.Send(&player1, resp)

	cp.Add(&client1)
	cp.Add(&client2)
	cp.Add(&client3)

	l1 := len(client1.codec.(*fakeCodec).Messages)
	l2 := len(client2.codec.(*fakeCodec).Messages)
	l3 := len(client3.codec.(*fakeCodec).Messages)
	cp.Send(&player1, resp)

	if len(client1.codec.(*fakeCodec).Messages) != l1+1 {
		t.Errorf("%d", len(client1.codec.(*fakeCodec).Messages))
	}

	if len(client2.codec.(*fakeCodec).Messages) != l2+1 {
		t.Fail()
	}

	if len(client3.codec.(*fakeCodec).Messages) != l3 {
		t.Fail()
	}
}

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
	client1.codec.(*fakeCodec).Messages = make([][]byte, 0)
	cp.Add(&client1)

	client1.pushStateChange(&planet)
	if len(client1.codec.(*fakeCodec).Messages) != 0 {
		t.Error("Client received messages  without a tick", len(client1.codec.(*fakeCodec).Messages))
	}

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

	cp.Broadcast(&sun)
	cp.Broadcast(&planet)
	if client1.stateChange == nil {
		t.Error("Client1 has no stacked planets")
	}

	if len(client1.stateChange.RawPlanets) != 1 {
		t.Errorf("Client1 has %d stacked planets instead of 1", len(client1.stateChange.RawPlanets))
	}

	if client2.stateChange == nil {
		t.Error("Client2 has no stacked planets")
	}

	if len(client2.stateChange.RawPlanets) != 1 {
		t.Errorf("Client2 has %d stacked planets instead of 1", len(client2.stateChange.RawPlanets))
	}

	if client3.stateChange != nil {
		t.Error("Client3 has stacked planets")
	}

	// Just to make sure nothing breaks with ghost users
	client3.MoveToAreas([]string{"area:1:1"})
	delete(cp.pool, client3.Player.Username)
	cp.Broadcast(&sun)
}
