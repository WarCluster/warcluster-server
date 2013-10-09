package entities

import (
	"testing"
	"time"
)

func TestDatabasePreparations(t *testing.T) {
	player := Player{"gophie", Color{22, 22, 22}, "asdf", "planet.0_0", []int{1, 1}, []int{2, 2}}
	expected_json := "{\"Color\":{\"R\":22,\"G\":22,\"B\":22},\"TwitterID\":\"asdf\",\"HomePlanet\":\"planet.0_0\",\"ScreenSize\":[1,1],\"ScreenPosition\":[2,2]}"
	expected_key := "player.gophie"

	key, byteJson, err := player.Serialize()
	if key != expected_key || string(byteJson) != expected_json {
		t.Error(string(byteJson))
		t.Error("Player JSON formatting gone wrong!")
	}

	if err != nil {
		t.Error("Error during serialization: ", err)
	}
}

func TestDeserializePlayer(t *testing.T) {
	var player *Player
	serialized_player := []byte("{\"Color\":{\"R\":22,\"G\":22,\"B\":22},\"TwitterID\":\"asdf\",\"HomePlanet\":\"planet.3_4\",\"ScreenSize\":[1,1],\"ScreenPosition\":[2,2]}")
	player = Construct("player.gophie", serialized_player).(*Player)
	if player.username != "gophie" {
		t.Error("Player's name is ", player.username)
	}

	if player.TwitterID != "asdf" {
		t.Error("Player's twitter id is ", player.TwitterID)
	}

	if player.HomePlanet != "planet.3_4" {
		t.Error("Player's HomePlanet is ", player.HomePlanet)
	}
	if player.ScreenSize[0] != 1 || player.ScreenSize[1] != 1 {
		t.Error("Player's ScreenSize is ", player.ScreenSize)
	}

	if player.ScreenPosition[0] != 2 || player.ScreenPosition[1] != 2 {
		t.Error("Player's ScreenPosition is ", player.ScreenPosition)
	}
}

func TestCreateMission(t *testing.T) {
	start_time := time.Now()
	planet_start := Planet{Color{22, 22, 22}, []int{271, 203}, true, 3, 1, start_time.Unix(), 100, 1000, "gophie"}
	planet_end := Planet{Color{22, 22, 22}, []int{471, 403}, false, 3, 1, start_time.Unix(), 50, 1000, "gophie"}
	player := Player{"gophie", Color{22, 22, 22}, "asdf", "planet.271_203", []int{1, 1}, []int{2, 2}}

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
