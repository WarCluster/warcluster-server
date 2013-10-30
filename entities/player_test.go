package entities

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestCreateMission(t *testing.T) {
	start_time := time.Now()
	planet_start := Planet{
		Color:               Color{22, 22, 22},
		Coords:              []int{271, 203},
		IsHome:              true,
		Texture:             3,
		Size:                1,
		LastShipCountUpdate: start_time.Unix(),
		ShipCount:           100,
		MaxShipCount:        1000,
		Owner:               "gophie",
	}
	planet_end := Planet{
		Color:               Color{22, 22, 22},
		Coords:              []int{471, 403},
		IsHome:              false,
		Texture:             3,
		Size:                1,
		LastShipCountUpdate: start_time.Unix(),
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

	valid_mission := player.StartMission(&planet_start, &planet_end, 80, "Attack")

	planet_start.ShipCount = 100
	invalid_mission := player.StartMission(&planet_start, &planet_end, 120, "Attack")

	if valid_mission.Source[0] != 271 || valid_mission.Source[1] != 203 {
		t.Error(valid_mission.Source)
		t.Error("Planet planet.271_203 was expected as start planet!")
	}

	if valid_mission.Target[0] != 471 || valid_mission.Target[1] != 403 {
		t.Error(valid_mission.Target)
		t.Error("Planet planet.471_403 was expected as end planet!")
	}

	if valid_mission.ShipCount != 80 {
		t.Error(valid_mission.ShipCount)
		t.Error("Mission ShipCount was expected to be 80!")
	}

	if invalid_mission.ShipCount != 100 {
		t.Error(invalid_mission.ShipCount)
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
