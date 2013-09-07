package entities

import (
	"strings"
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

func TestMissionSerialize(t *testing.T) {
	start_time := time.Date(2013, time.August, 14, 22, 12, 6, 0, time.UTC).UnixNano() / 1e6
	mission := Mission{[]int{32, 64}, []int{2, 2}, "Atack", start_time, start_time, start_time, "gophie", 5}
	expected_json_prefix := "{\"Source\":[32,64],\"Target\":[2,2],\"Type\":\"Atack\",\"CurrentTime\""
	expected_json_suffix := "\"StartTime\":1376518326000,\"TravelTime\":1376518326000,\"Player\":\"gophie\",\"ShipCount\":5}"

	key, json, err := mission.Serialize()

	if key != mission.GetKey() {
		t.Error("You're not using the missions' GetKey()!")
	}

	if !strings.HasPrefix(string(json), expected_json_prefix) || !strings.HasSuffix(string(json), expected_json_suffix) {
		t.Error("Serialized mission is ", string(json))
	}

	if err != nil {
		t.Error("Error during serialization: ", err)
	}
}

func TestMissionDeserialize(t *testing.T) {
	serialized_mission := []byte(strings.Join([]string{"{\"Source\":[32,64],",
		"\"Target\":[2,2],",
		"\"Type\":\"Atack\",",
		"\"CurrentTime\":\"2013-08-14T22:12:06Z\",",
		"\"StartTime\":\"2013-08-14T22:12:06.06Z\",",
		"\"TravelTime\":\"2013-08-14T22:12:06Z\",",
		"\"Player\":\"gophie\",",
		"\"ShipCount\":5}"}, ""))
	mission := Construct("mission.137650752666_32_64", serialized_mission).(*Mission)

	if mission.Player != "gophie" {
		t.Error("Mission's player is ", mission.Player)
	}

	if mission.ShipCount != 5 {
		t.Error("Mission's ShipCount is ", mission.ShipCount)
	}

	if mission.Source[0] != 32 || mission.Source[1] != 64 {
		t.Error("Mission's Source is ", mission.Source)
	}

	if mission.Target[0] != 2 || mission.Target[1] != 2 {
		t.Error("Mission's Target is ", mission.Target)
	}
}

func TestEndMission(t *testing.T) {
	mission := new(Mission)
	secondMission := new(Mission)
	endPlanet := new(Planet)
	player := new(Player)
	start_time := time.Now().UnixNano() * 1e6
	*player = Player{"chochko", Color{"asd1", 22, 22, 22}, "asdf2", "planet.0_0", []int{1, 1}, []int{2, 2}}
	*mission = Mission{[]int{32, 64}, []int{2, 2}, "Atack", start_time, start_time, start_time, "gophie", 15}
	*secondMission = Mission{[]int{32, 64}, []int{2, 2}, "Atack", start_time, start_time, start_time, "chochko", 10}
	*endPlanet = Planet{[]int{2, 2}, 6, 3, start_time, 2, 0, "chochko"}

	endPlanet = EndMission(endPlanet, player, secondMission)
	/* //TODO: Test needs to be revised in order to handle calculation of ship count
	if endPlanet.GetShipCount() != 12 {
		t.Error("End Planet ship count was expected  to be 12 but it is:", endPlanet.GetShipCount())
	}
	*/
	if endPlanet.Owner != "chochko" {
		t.Error("End Planet owner was expected  to be chochko but is:", endPlanet.Owner)
	}

	endPlanet = EndMission(endPlanet, player, mission)
	/* //TODO: Test needs to be revised in order to handle calculation of ship count
	if endPlanet.GetShipCount() != 3 {
		t.Error("End Planet ship count was expected  to be 3 but it is:", endPlanet.GetShipCount())
	}
	*/
	if endPlanet.Owner != "gophie" {
		t.Error("End Planet owner was expected  to be gophie but is:", endPlanet.Owner)
	}
}

func TestEndMissionDenyTakeover(t *testing.T) {
	mission := new(Mission)
	endPlanet := new(Planet)
	player := new(Player)
	start_time := time.Now().UnixNano() * 1e6
	*mission = Mission{[]int{32, 64}, []int{2, 2}, "Atack", start_time, start_time, start_time, "gophie", 15}
	*endPlanet = Planet{[]int{2, 2}, 6, 3, start_time, 2, 0, "chochko"}
	*player = Player{"chochko", Color{"asd1", 22, 22, 22}, "asdf1", "planet.2_2", []int{1, 1}, []int{2, 2}}

	endPlanet = EndMission(endPlanet, player, mission)
	//TODO: Test needs to be revised in order to handle calculation of ship count
	if endPlanet.GetShipCount() != 0 {
		t.Error("End Planet ship count was expected  to be 0 but it is:", endPlanet.GetShipCount())
	}
	//TODO: Test needs to be revised in order to handle feedback mission with exess ships
	if endPlanet.Owner != "chochko" {
		t.Error("End Planet owner was expected  to be chochko but is:", endPlanet.Owner)
	}
}


func TestTravelTime(t *testing.T) {
	mission := new(Mission)
	*mission = Mission{
		Source:      []int{100, 200},
		Target:      []int{800, 150},
		CurrentTime: time.Now().UnixNano() / 1e6,
		StartTime:   time.Now().UnixNano() / 1e6,
		TravelTime:  time.Now().UnixNano() / 1e6,
		Player:      "gophie",
		ShipCount:   50,
	}
	mission.CalculateTravelTime()
	var expectedTravel int64 = 7017

	if mission.TravelTime != expectedTravel {
		t.Error("Wrong arrival time:", mission.TravelTime, "instead of:", expectedTravel)
	}
}