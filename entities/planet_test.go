package entities

import (
    "testing"
    "strconv"
)

func TestGeneratePlanet(t *testing.T) {
    hash := "5762908447300427353060676895795336101745023746116233389596883"
    sun_position := []int{500, 300}
    expected_planets := []Planet{
        Planet{[]int{375, 247}, 6, 3, 0, 0, "gophie"},
        Planet{[]int{694, 300}, 8, 5, 0, 0, "gophie"},
        Planet{[]int{271, 203}, 3, 1, 0, 0, "gophie"},
        Planet{[]int{209, 365}, 2, 8, 0, 0, "gophie"},
        Planet{[]int{671,  -6}, 3, 1, 0, 0, "gophie"},
        Planet{[]int{907, 300}, 6, 8, 0, 0, "gophie"},
        Planet{[]int{918, 101}, 9, 6, 0, 0, "gophie"},
        Planet{[]int{352, 798}, 5, 4, 0, 0, "gophie"},
        Planet{[]int{686, 841}, 1, 1, 0, 0, "gophie"},
    }
    generated_planets, _ := GeneratePlanets(hash, sun_position)

    if len(generated_planets) != 9 {
        t.Error("Wrong planets count")
    }
    for i:=0; i<9; i++ {
        if generated_planets[i].coords[0] != expected_planets[i].coords[0] {
            t.Error("X coordinate missmatch on Planet[" + strconv.Itoa(i) + "]")
        }
        if generated_planets[i].coords[1] != expected_planets[i].coords[1] {
            t.Error("Y coordinate missmatch on Planet[" + strconv.Itoa(i) + "]")
        }
        if generated_planets[i].Texture != expected_planets[i].Texture {
            t.Error("Ring offset missmatch on Planet[" + strconv.Itoa(i) + "]")
        }
        if generated_planets[i].Size != expected_planets[i].Size {
            t.Error("Ring offset missmatch on Planet[" + strconv.Itoa(i) + "]")
        }
    }
}

func TestDatabasePreparationsWithoutAnOwner(t *testing.T) {
    planet := Planet{[]int{271, 203}, 3, 1, 0, 0, ""}
    expected_json := "{\"Texture\":3,\"Size\":1,\"ShipCount\":0,\"MaxShipCount\":0,\"Owner\":\"\"}"
    expected_key := "planet.271_203"

    key, json := planet.PrepareForDB()
    if key != expected_key || string(json) != expected_json {
        t.Error("Planet JSON formatting gone wrong!")
    }
}

func TestDatabasePreparationsWithAnOwner(t *testing.T) {
    planet := Planet{[]int{271, 203}, 3, 1, 0, 0, "gophie"}
    expected_json := "{\"Texture\":3,\"Size\":1,\"ShipCount\":0,\"MaxShipCount\":0,\"Owner\":\"gophie\"}"
    expected_key := "planet.271_203"

    key, json := planet.PrepareForDB()
    if key != expected_key || string(json) != expected_json {
        t.Error(json)
        t.Error("Planet JSON formatting gone wrong!")
    }
}

func TestDeserializePlanet(t *testing.T) {
    var planet Planet
    serialized_planet := []byte("{\"Texture\":3,\"Size\":1,\"ShipCount\":10,\"MaxShipCount\":15,\"Owner\":\"gophie\"}")
    planet = Construct("planet.10_12", serialized_planet).(Planet)

    if planet.Texture != 3 {
        t.Error("Planet's texture is ", planet.Texture)
    }

    if planet.Size != 1 {
        t.Error("Planet's tize is ", planet.Size)
    }

    if planet.ShipCount != 10 {
        t.Error("Planet's ship count is ", planet.ShipCount)
    }

    if planet.MaxShipCount != 15 {
        t.Error("Planet's max ship count is ", planet.MaxShipCount)
    }

    if planet.Owner != "gophie" {
        t.Error("Planet's owner is ", planet.Owner)
    }

    if planet.coords[0] != 10  && planet.coords[1] != 12 {
        t.Error("Planet's coords are ", planet.coords)
    }
}
