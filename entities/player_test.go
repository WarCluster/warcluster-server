package entities

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestCreateMission(t *testing.T) {
	startTime := time.Now()
	planetStart := Planet{
		Color:               Color{22, 22, 22},
		Coords:              []int{271, 203},
		IsHome:              true,
		Texture:             3,
		Size:                1,
		LastShipCountUpdate: startTime.Unix(),
		ShipCount:           100,
		MaxShipCount:        1000,
		Owner:               "gophie",
	}
	planetEnd := Planet{
		Color:               Color{22, 22, 22},
		Coords:              []int{471, 403},
		IsHome:              false,
		Texture:             3,
		Size:                1,
		LastShipCountUpdate: startTime.Unix(),
		ShipCount:           50,
		MaxShipCount:        1000,
		Owner:               "gophie",
	}
	player := Player{
		username:       "gophie",
		Color:          Color{22, 22, 22},
		TwitterID:      "asdf",
		HomePlanet:     "planet.271_203",
		ScreenSize:     []int{1, 1},
		ScreenPosition: []int{2, 2},
	}

	validMission := player.StartMission(&planetStart, &planetEnd, 80, "Attack")

	planetStart.ShipCount = 100
	invalidMission := player.StartMission(&planetStart, &planetEnd, 120, "Attack")

	if validMission.Source[0] != 271 || validMission.Source[1] != 203 {
		t.Error(validMission.Source)
		t.Error("Planet planet.271_203 was expected as start planet!")
	}

	if validMission.Target[0] != 471 || validMission.Target[1] != 403 {
		t.Error(validMission.Target)
		t.Error("Planet planet.471_403 was expected as end planet!")
	}

	if validMission.ShipCount != 80 {
		t.Error(validMission.ShipCount)
		t.Error("Mission ShipCount was expected to be 80!")
	}

	if invalidMission.ShipCount != 100 {
		t.Error(invalidMission.ShipCount)
		t.Error("Mission ShipCount was expected to be 100!")
	}
}

func TestPlayerMarshalling(t *testing.T) {
	var uPlayer Player

	mPlayer, err := json.Marshal(player)
	if err != nil {
		t.Error("Player marshaling failed:", err)
	}

	err = json.Unmarshal(mPlayer, &uPlayer)
	if err != nil {
		t.Error("Player unmarshaling failed:", err)
	}
	uPlayer.username = player.username

	if player.GetKey() != uPlayer.GetKey() {
		t.Error(
			"Keys of both players are different!\n",
			player.GetKey(),
			"!=",
			uPlayer.GetKey(),
		)
	}

	if !reflect.DeepEqual(player, uPlayer) {
		t.Error("Players are different after the marshal->unmarshal step")
	}
}
