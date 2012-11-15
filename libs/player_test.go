package libs

import (
    "testing"
)

func TestDatabasePreparations(t *testing.T) {
    player := Player{"gophie", "asdf", "planet_0_0", []int{1, 1}, []int{2, 2}}
    expected_json := "{\"Hash\":\"asdf\",\"HomePlanet\":\"planet_0_0\",\"ScreenSize\":[1,1],\"ScreenPosition\":[2,2]}"
    expected_key := "player_gophie"

    key, json := player.PrepareForDB()
    if key != expected_key || string(json) != expected_json {
        t.Error(string(key))
        t.Error("Player JSON formatting gone wrong!")
    }
}
