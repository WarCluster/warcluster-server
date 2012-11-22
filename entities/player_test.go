package entities

import (
    "testing"
)

func TestDatabasePreparations(t *testing.T) {
    player := Player{"gophie", "asdf", "planet.0_0", []int{1, 1}, []int{2, 2}}
    expected_json := "{\"Hash\":\"asdf\",\"HomePlanet\":\"planet.0_0\",\"ScreenSize\":[1,1],\"ScreenPosition\":[2,2]}"
    expected_key := "player.gophie"

    key, json := player.Serialize()
    if key != expected_key || string(json) != expected_json {
        t.Error(string(json))
        t.Error("Player JSON formatting gone wrong!")
    }
}

func TestDeserialize(t *testing.T) {
    var player Player
    serialized_player := []byte("{\"Hash\":\"asdf\",\"HomePlanet\":\"planet.3_4\",\"ScreenSize\":[1,1],\"ScreenPosition\":[2,2]}")
    player = Construct("player.gophie", serialized_player).(Player)

    if player.username != "gophie" {
        t.Error("Player's name is ", player.username)
    }

    if player.Hash != "asdf" {
        t.Error("Player's hash is ", player.Hash)
    }

    if player.HomePlanet != "planet.3_4" {
        t.Error("Player's HomePlanet is ", player.HomePlanet)
    }

    if player.ScreenSize[0] != 1  && player.ScreenSize[1] != 1 {
        t.Error("Player's ScreenSize is ", player.ScreenSize)
    }

    if player.ScreenPosition[0] != 2  && player.ScreenPosition[1] != 2 {
        t.Error("Player's ScreenPosition is ", player.ScreenPosition)
    }
}
