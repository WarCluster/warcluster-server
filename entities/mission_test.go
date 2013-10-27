package entities

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestMissionGetKey(t *testing.T) {
	start_time := time.Date(2012, time.November, 10, 23, 0, 0, 0, time.UTC).UnixNano() / 1e6
	mission := new(Mission)
	mission.Source = []int{32, 64}
	mission.StartTime = start_time

	if mission.GetKey() != "mission.1352588400000_32_64" {
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

	uMission.CurrentTime = timeStamp

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
	secondMission := new(Mission)
	endPlanet := new(Planet)
	start_time := time.Now().UnixNano() * 1e6
	*secondMission = Mission{Color{22, 22, 22}, []int{32, 64}, []int{2, 2}, "Attack", start_time, start_time, start_time, "chochko", 10}
	*endPlanet = Planet{Color{22, 22, 22}, []int{2, 2}, false, 6, 3, start_time, 2, 0, "chochko"}

	t.Skip()
	endPlanet = EndMission(endPlanet, secondMission)
	if endPlanet.GetShipCount() != 12 {
		t.Error("End Planet ship count was expected  to be 12 but it is:", endPlanet.GetShipCount())
	}

	if endPlanet.Owner != "chochko" {
		t.Error("End Planet owner was expected  to be chochko but is:", endPlanet.Owner)
	}

	endPlanet = EndMission(endPlanet, &mission)
	if endPlanet.GetShipCount() != 3 {
		t.Error("End Planet ship count was expected  to be 3 but it is:", endPlanet.GetShipCount())
	}

	if endPlanet.Owner != "gophie" {
		t.Error("End Planet owner was expected  to be gophie but is:", endPlanet.Owner)
	}
}

//TODO: Test needs to be revised in order to handle calculation of ship count
//TODO: Test needs to be revised in order to handle feedback mission with exess ships
func TestEndMissionDenyTakeover(t *testing.T) {
	endPlanet := new(Planet)
	*endPlanet = Planet{Color{22, 22, 22}, []int{2, 2}, true, 6, 3, timeStamp, 2, 0, "chochko"}

	endPlanet = EndMission(endPlanet, &mission)
	if endPlanet.GetShipCount() != 0 {
		t.Error("End Planet ship count was expected  to be 0 but it is:", endPlanet.GetShipCount())
	}

	if endPlanet.Owner != "chochko" {
		t.Error("End Planet owner was expected  to be chochko but is:", endPlanet.Owner)
	}
}

func TestTravelTime(t *testing.T) {
	mission.CalculateTravelTime()
	var expectedTravel int64 = 7017

	if mission.TravelTime != expectedTravel {
		t.Error("Wrong arrival time:", mission.TravelTime, "instead of:", expectedTravel)
	}
}
