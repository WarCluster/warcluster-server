package entities

import (
    "testing"
    "time"
)

func TestMissionGeyKey(t *testing.T) {
    start_time := time.Date(2012, time.November, 10, 23, 0, 0, 0, time.UTC)
    mission := Mission{"planet.32_64", start_time, "gophie", 5, "planet.2_2"}

    if mission.GetKey() != "mission_1352588400_32_64" {
        t.Error("Mission's time is ", mission.GetKey())
    }
}

func TestMissionSerialize(t *testing.T) {
    start_time := time.Date(2012, time.November, 10, 23, 0, 0, 0, time.UTC)
    mission := Mission{"planet.32_64", start_time, "gophie", 5, "planet.2_2"}
    expected_json := "{\"Player\":\"gophie\",\"ShipCount\":5,\"EndPlanet\":\"planet.2_2\"}"

    key, json := mission.Serialize()
    if key != mission.GetKey() {
        t.Error("You're not using the missions' GetKey()!")
    }

    if string(json) != expected_json {
        t.Error("Serialized mission is ", json, "but iy should be ", expected_json)
    }

}
