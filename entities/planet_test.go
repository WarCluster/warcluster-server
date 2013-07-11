package entities

import (
	"fmt"
	"github.com/Vladimiroff/vec2d"
	"strconv"
	"testing"
	"time"
)

func TestGeneratePlanets(t *testing.T) {
	start_time := time.Date(2012, time.November, 10, 23, 0, 0, 0, time.UTC)
	hash := "5762908447300427353060676895795336101745023746116233389596883"
	sun_position := vec2d.New(500, 300)
	expected_planets := []Planet{
		Planet{[]int{-223, -4}, 6, 3, start_time.Unix(), 0, 0, "gophie"},
		Planet{[]int{1490, 300}, 8, 5, start_time.Unix(), 00, 0, "gophie"},
		Planet{[]int{-578, -153}, 3, 1, start_time.Unix(), 00, 0, "gophie"},
		Planet{[]int{-797, 591}, 2, 8, start_time.Unix(), 00, 0, "gophie"},
		Planet{[]int{1233, -1014}, 3, 1, start_time.Unix(), 00, 0, "gophie"},
		Planet{[]int{2195, 300}, 6, 8, start_time.Unix(), 00, 0, "gophie"},
		Planet{[]int{2203, -507}, 9, 6, start_time.Unix(), 00, 0, "gophie"},
		Planet{[]int{-90 , 2294}, 5, 4, start_time.Unix(), 00, 0, "gophie"},
		Planet{[]int{1234, 2431}, 1, 1, start_time.Unix(), 00, 0, "gophie"},
		Planet{[]int{-1730, -638}, 4, 6, start_time.Unix(), 00, 0, "gophie"},
	}
	generated_planets, _ := GeneratePlanets(hash, sun_position)

	if len(generated_planets) != 10 {
		t.Error("Wrong planets count")
	}
	for i := 0; i < 10; i++ {
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
			t.Error("Size missmatch on Planet[" + strconv.Itoa(i) + "]")
		}
	}
}

func TestDatabasePreparationsWithoutAnOwner(t *testing.T) {
	start_time := time.Now()
	planet := Planet{[]int{271, 203}, 3, 1, start_time.Unix(), 0, 0, ""}
	json_base := "{\"Texture\":3,\"Size\":1,\"LastShipCountUpdate\":%v,\"ShipCount\":0,\"MaxShipCount\":0,\"Owner\":\"\"}"
	expected_json := fmt.Sprintf(json_base, start_time.Unix())
	expected_key := "planet.271_203"

	key, json, err := planet.Serialize()
	if key != expected_key || string(json) != expected_json {
		t.Error(string(json))
		t.Error("Planet JSON formatting gone wrong!")
	}

	if err != nil {
		t.Error("Error during serialization: ", err)
	}
}

func TestDatabasePreparationsWithAnOwner(t *testing.T) {
	start_time := time.Now()
	planet := Planet{[]int{271, 203}, 3, 1, start_time.Unix(), 0, 0, "gophie"}
	json_base := "{\"Texture\":3,\"Size\":1,\"LastShipCountUpdate\":%v,\"ShipCount\":0,\"MaxShipCount\":0,\"Owner\":\"gophie\"}"
	expected_json := fmt.Sprintf(json_base, start_time.Unix())
	expected_key := "planet.271_203"

	key, json, err := planet.Serialize()
	if key != expected_key || string(json) != expected_json {
		t.Error(string(json))
		t.Error(string(expected_json))
		t.Error("Planet JSON formatting gone wrong!")
	}

	if err != nil {
		t.Error("Error during serialization: ", err)
	}
}

func TestDeserializePlanet(t *testing.T) {
	var planet *Planet
	serialized_planet := []byte("{\"Texture\":3,\"Size\":1,\"LastShipCountUpdate\":1352588400,\"ShipCount\":10,\"MaxShipCount\":15,\"Owner\":\"gophie\"}")
	planet = Construct("planet.10_12", serialized_planet).(*Planet)

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

	if planet.coords[0] != 10 && planet.coords[1] != 12 {
		t.Error("Planet's coords are ", planet.coords)
	}
}
