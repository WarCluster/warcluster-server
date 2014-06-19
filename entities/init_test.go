package entities

import (
	"time"

	"github.com/Vladimiroff/vec2d"
)

var (
	timeStamp int64   = time.Date(2012, time.November, 10, 23, 0, 0, 0, time.UTC).UnixNano() / 1e6
	now       int64   = time.Now().UnixNano() * 1e6
	mission   Mission = Mission{
		Color: Color{22, 22, 22},
		Source: embeddedPlanet{
			Name:     "GOP6720",
			Position: vec2d.New(271, 203),
		},
		Target: embeddedPlanet{
			Name:     "GOP6721",
			Position: vec2d.New(2, 2),
		},
		Type:       "Attack",
		StartTime:  timeStamp,
		TravelTime: time.Duration(timeStamp),
		Player:     "gophie",
		ShipCount:  5,
	}
	secondMission = Mission{
		Color: Color{22, 22, 22},
		Source: embeddedPlanet{
			Name:     "GOP6720",
			Owner:    "gophie",
			Position: vec2d.New(271, 203),
		},
		Target: embeddedPlanet{
			Name:     "GOP6721",
			Owner:    "chochko",
			Position: vec2d.New(2, 2),
		},
		Type:       "Attack",
		StartTime:  now,
		TravelTime: time.Duration(now),
		Player:     "chochko",
		ShipCount:  10,
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
	endPlanet = Planet{
		Name:                "GOP6721",
		Color:               Color{22, 22, 22},
		Position:            vec2d.New(2, 2),
		IsHome:              false,
		Texture:             6,
		Size:                3,
		LastShipCountUpdate: now,
		ShipCount:           2,
		MaxShipCount:        0,
		Owner:               "chochko",
	}
	player Player = Player{
		Username:       "gophie",
		RaceID:         1,
		TwitterID:      "asdf",
		HomePlanet:     "planet.GOP6720",
		ScreenSize:     []uint16{1, 1},
		ScreenPosition: &vec2d.Vector{2, 2},
	}
	sun Sun = Sun{
		Name:     "GOP672",
		Username: "gophie",
		Position: vec2d.New(20, 20),
	}
)
