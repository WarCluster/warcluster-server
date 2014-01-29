package entities

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/Vladimiroff/vec2d"
)

func TestCreateMission(t *testing.T) {
	startTime := time.Now()
	planetStart := Planet{
		Name:                "GOP6720",
		Color:               Color{0.59215686, 0.59215686, 0.59215686},
		Position:            vec2d.New(271, 203),
		IsHome:              true,
		Texture:             3,
		Size:                1,
		LastShipCountUpdate: startTime.Unix(),
		ShipCount:           100,
		MaxShipCount:        1000,
		Owner:               "gophie",
	}
	planetEnd := Planet{
		Name:                "GOP6721",
		Color:               Color{0.59215686, 0.59215686, 0.59215686},
		Position:            vec2d.New(471, 403),
		IsHome:              false,
		Texture:             3,
		Size:                1,
		LastShipCountUpdate: startTime.Unix(),
		ShipCount:           50,
		MaxShipCount:        1000,
		Owner:               "gophie",
	}
	player := Player{
		Username:       "gophie",
		Color:          Color{0.59215686, 0.59215686, 0.59215686},
		TwitterID:      "asdf",
		HomePlanet:     "planet.GOP6720",
		ScreenSize:     []uint16{1, 1},
		ScreenPosition: &vec2d.Vector{2, 2},
	}

	validMission := player.StartMission(&planetStart, &planetEnd, 80, "Attack")

	planetStart.ShipCount = 100
	invalidMission := player.StartMission(&planetStart, &planetEnd, 120, "Attack")

	if validMission.Source.Name != "GOP6720" {
		t.Error(validMission.Source)
		t.Error("Planet planet.GOP6720 was expected as start planet!")
	}

	if validMission.Target.Name != "GOP6721" {
		t.Error(validMission.Target)
		t.Error("Planet planet.GOP6720 was expected as end planet!")
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
	uPlayer.Username = player.Username

	if player.Key() != uPlayer.Key() {
		t.Error(
			"Keys of both players are different!\n",
			player.Key(),
			"!=",
			uPlayer.Key(),
		)
	}

	if !reflect.DeepEqual(player, uPlayer) {
		t.Error("Players are different after the marshal->unmarshal step")
	}
}
