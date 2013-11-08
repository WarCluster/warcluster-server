package entities

import (
	"encoding/json"
	"github.com/Vladimiroff/vec2d"
	"reflect"
	"testing"
	"time"
)

func TestMissionGetKey(t *testing.T) {
	startTime := time.Date(2012, time.November, 10, 23, 0, 0, 0, time.UTC).UnixNano() / 1e6
	mission := new(Mission)
	mission.Source = "GOP5610"
	mission.StartTime = startTime

	if mission.GetKey() != "mission.1352588400000_GOP5610" {
		t.Error("Mission's key is ", mission.GetKey())
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

	if mission.GetKey() != uMission.GetKey() {
		t.Error(
			"Keys of both missions are different!\n",
			mission.GetKey(),
			"!=",
			uMission.GetKey(),
		)
	}

	if !reflect.DeepEqual(mission, uMission) {
		t.Error("Missions are different after the marshal->unmarshal step")
	}
}

//TODO: Test needs to be revised in order to handle calculation of ship count
func TestEndMission(t *testing.T) {
	var excessShips int
	secondMission := new(Mission)
	endPlanet := new(Planet)
	startTime := time.Now().UnixNano() * 1e6
	*secondMission = Mission{Color{22, 22, 22}, "GOP5610", "GOP5611", "Attack", startTime, startTime, "chochko", 10}
	*endPlanet = Planet{"", Color{22, 22, 22}, vec2d.New(2, 2), false, 6, 3, startTime, 2, 0, "chochko"}

	t.Skip()
	excessShips = EndMission(endPlanet, secondMission)
	if endPlanet.GetShipCount() != 12 {
		t.Error("End Planet ship count was expected  to be 12 but it is:", endPlanet.GetShipCount())
	}

	if endPlanet.Owner != "chochko" {
		t.Error("End Planet owner was expected  to be chochko but is:", endPlanet.Owner)
	}

	excessShips = EndMission(endPlanet, &mission)
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

//TODO: Test needs to be revised in order to handle calculation of ship count
//TODO: Test needs to be revised in order to handle feedback mission with excess ships
func TestEndMissionDenyTakeover(t *testing.T) {
	var excessShips int
	endPlanet := new(Planet)
	*endPlanet = Planet{"", Color{22, 22, 22}, vec2d.New(2, 2), true, 6, 3, timeStamp, 2, 0, "chochko"}

	excessShips = EndMission(endPlanet, &mission)
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
	expectedTime := int64(7017)
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
