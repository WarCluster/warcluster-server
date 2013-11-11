package entities

import (
	"github.com/Vladimiroff/vec2d"
	"time"
)

var (
	timeStamp int64   = time.Date(2012, time.November, 10, 23, 0, 0, 0, time.UTC).UnixNano() / 1e6
	mission   Mission = Mission{
		Color:      Color{22, 22, 22},
		Source:     "GOP6720",
		Target:     "GOP6721",
		Type:       "Attack",
		StartTime:  timeStamp,
		TravelTime: timeStamp,
		Player:     "gophie",
		ShipCount:  5,
	}
	planet Planet = Planet{
		Name:                "GOP6720",
		Color:               Color{22, 22, 22},
		Position:            vec2d.New(271, 203),
		IsHome:              false,
		Texture:             3,
		Size:                1,
		LastShipCountUpdate: timeStamp,
		ShipCount:           0,
		MaxShipCount:        0,
		Owner:               "gophie",
	}
	player Player = Player{
		Username:       "gophie",
		Color:          Color{22, 22, 22},
		TwitterID:      "asdf",
		HomePlanet:     "planet.271_203",
		ScreenSize:     []uint16{1, 1},
		ScreenPosition: []int64{2, 2},
	}
	sun Sun = Sun{
		Name:     "GOP672",
		Username: "gophie",
		Position: vec2d.New(20, 20),
	}
)
