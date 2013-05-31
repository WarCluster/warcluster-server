package entities

import (
	"testing"
	"time"
)

func TestMissionGeyKey(t *testing.T) {
	start_time := time.Date(2012, time.November, 10, 23, 0, 0, 0, time.UTC)
	mission := Mission{"planet.32_64", start_time, "gophie", 5, "planet.2_2"}

	if mission.GetKey() != "mission.1352588400_32_64" {
		t.Error("Mission's time is ", mission.GetKey())
	}
}

func TestMissionSerialize(t *testing.T) {
	start_time := time.Date(2012, time.November, 10, 23, 0, 0, 0, time.UTC)
	mission := Mission{"planet.32_64", start_time, "gophie", 5, "planet.2_2"}
	expected_json := "{\"Player\":\"gophie\",\"ShipCount\":5,\"EndPlanet\":\"planet.2_2\"}"

	key, json, err := mission.Serialize()

	if key != mission.GetKey() {
		t.Error("You're not using the missions' GetKey()!")
	}

	if string(json) != expected_json {
		t.Error("Serialized mission is ", json, "but iy should be ", expected_json)
	}

	if err != nil {
		t.Error("Error during serialization: ", err)
	}
}

func TestMissionDeserialize(t *testing.T) {
	serialized_mission := []byte("{\"Player\":\"gophie\",\"ShipCount\":5,\"EndPlanet\":\"planet.2_2\"}")
	mission := Construct("mission.1352588400_32_64", serialized_mission).(Mission)

	if mission.Player != "gophie" {
		t.Error("Mission's player is ", mission.Player)
	}

	if mission.ShipCount != 5 {
		t.Error("Mission's ShipCount is ", mission.ShipCount)
	}

	if mission.EndPlanet != "planet.2_2" {
		t.Error("Mission's EndPlanet is ", mission.EndPlanet)
	}
}

func TestEndMission(t *testing.T) {
	start_time := time.Date(2012, time.November, 10, 23, 0, 0, 0, time.UTC)
	mission := Mission{"planet.32_64", start_time, "gophie", 15, "planet.2_2"}
	seccondMission := Mission{"planet.32_64", start_time, "chochko", 10, "planet.2_2"}
	endPlanet := Planet{[]int{2, 2}, 6, 3, start_time.Unix(), 2, 0, "chochko"}

	endPlanet = EndMission(endPlanet, seccondMission)
	/* //TODO: Test needs to be revised in order to handle calculation of ship count
	if endPlanet.GetShipCount() != 12 {
		t.Error("End Planet ship count was expected  to be 12 but it is:", endPlanet.GetShipCount())
	}
	*/
	if endPlanet.Owner != "chochko" {
		t.Error("End Planet owner was expected  to be chochko but is:", endPlanet.Owner)
	}

	endPlanet = EndMission(endPlanet, mission)
	/* //TODO: Test needs to be revised in order to handle calculation of ship count
	if endPlanet.GetShipCount() != 3 {
		t.Error("End Planet ship count was expected  to be 3 but it is:", endPlanet.GetShipCount())
	}
	*/
	if endPlanet.Owner != "gophie" {
		t.Error("End Planet owner was expected  to be gophie but is:", endPlanet.Owner)
	}
}
