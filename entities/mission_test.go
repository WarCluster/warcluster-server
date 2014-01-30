package entities

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/Vladimiroff/vec2d"
)

func TestMissionKey(t *testing.T) {
	startTime := time.Date(2012, time.November, 10, 23, 0, 0, 0, time.UTC).UnixNano() / 1e6
	mission := new(Mission)
	mission.Source = embeddedPlanet{planet.Name, planet.Owner, planet.Position}
	mission.StartTime = startTime

	if mission.Key() != "mission.1352588400000_GOP6720" {
		t.Error("Mission's key is ", mission.Key())
	}
}

func TestMissionMarshalling(t *testing.T) {
	var uMission Mission

	mMission, err := json.Marshal(mission)
	if err != nil {
		t.Error("Mission marshaling failed:", err)
	}

	err = json.Unmarshal(mMission, &uMission)
	if err != nil {
		t.Error("Mission unmarshaling failed:", err)
	}

	if mission.Key() != uMission.Key() {
		t.Error(
			"Keys of both missions are different!\n",
			mission.Key(),
			"!=",
			uMission.Key(),
		)
	}

	if !reflect.DeepEqual(mission, uMission) {
		t.Error("Missions are different after the marshal->unmarshal step")
	}
}

func TestEndAttackMission(t *testing.T) {
	var excessShips int32
	excessShips = secondMission.EndAttackMission(&endPlanet)
	if endPlanet.GetShipCount() != 12 {
		t.Error("End Planet ship count was expected  to be 12 but it is:", endPlanet.GetShipCount())
	}

	if endPlanet.Owner != "chochko" {
		t.Error("End Planet owner was expected  to be chochko but is:", endPlanet.Owner)
	}

	mission.ShipCount = 15
	excessShips = mission.EndAttackMission(&endPlanet)
	if endPlanet.GetShipCount() != 3 {
		t.Error("End Planet ship count was expected  to be 3 but it is:", endPlanet.GetShipCount())
	}

	if endPlanet.Owner != "gophie" {
		t.Error("End Planet owner was expected  to be gophie but is:", endPlanet.Owner)
	}

	if excessShips != 0 {
		t.Error("There shouldn't be any excess ships, but the value is", excessShips)
	}
}

func TestEndAttackMissionDenyTakeover(t *testing.T) {
	var excessShips int32
	endPlanet := new(Planet)
	*endPlanet = Planet{"", Color{0.59215686, 0.59215686, 0.59215686}, vec2d.New(2, 2), true, 6, 3, timeStamp, 2, 0, "chochko"}

	mission.ShipCount = 5
	excessShips = mission.EndAttackMission(endPlanet)
	if endPlanet.GetShipCount() != 0 {
		t.Error("End Planet ship count was expected to be 0 but it is:", endPlanet.GetShipCount())
	}

	if endPlanet.Owner != "chochko" {
		t.Error("End Planet owner was expected to be chochko but is:", endPlanet.Owner)
	}

	if excessShips != 5 {
		t.Error("There should be 5 excess ships, but the value is", excessShips)
	}
}

func TestTravelTime(t *testing.T) {
	source := vec2d.New(100, 200)
	target := vec2d.New(800, 150)
	expectedTime := time.Duration(7017)
	time := calculateTravelTime(source, target, 10)

	if time != expectedTime {
		t.Errorf(
			"CalculateTravelTime(%v, %v, 10) = %d instead of %d",
			source,
			target,
			time,
			expectedTime,
		)
	}
}
