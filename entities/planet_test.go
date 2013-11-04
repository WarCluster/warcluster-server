package entities

import (
	"encoding/json"
	"github.com/Vladimiroff/vec2d"
	"reflect"
	"strconv"
	"testing"
)

func TestGeneratePlanets(t *testing.T) {
	hash := "5762908447300427353060676895795336101745023746116233389596883"
	sunPosition := vec2d.New(500, 300)
	expectedPlanets := []Planet{
		Planet{Color{22, 22, 22}, []int{-76, 57}, false, 6, 3, timeStamp, 10, 0, "gophie"},
		Planet{Color{22, 22, 22}, []int{1470, 300}, false, 8, 5, timeStamp, 10, 0, "gophie"},
		Planet{Color{22, 22, 22}, []int{-689, -200}, false, 3, 1, timeStamp, 10, 0, "gophie"},
		Planet{Color{22, 22, 22}, []int{-1051, 648}, false, 2, 8, timeStamp, 10, 0, "gophie"},
		Planet{Color{22, 22, 22}, []int{1428, -1363}, false, 3, 1, timeStamp, 10, 0, "gophie"},
		Planet{Color{22, 22, 22}, []int{2735, 300}, false, 6, 8, timeStamp, 10, 0, "gophie"},
		Planet{Color{22, 22, 22}, []int{2818, -798}, false, 9, 6, timeStamp, 10, 0, "gophie"},
		Planet{Color{22, 22, 22}, []int{-322, 3080}, false, 5, 4, timeStamp, 10, 0, "gophie"},
		Planet{Color{22, 22, 22}, []int{1547, 3339}, false, 1, 1, timeStamp, 10, 0, "gophie"},
		Planet{Color{22, 22, 22}, []int{-2744, -1065}, false, 4, 6, timeStamp, 10, 0, "gophie"},
	}
	generatedPlanets, _ := GeneratePlanets(hash, sunPosition)

	if len(generatedPlanets) != 10 {
		t.Error("Wrong planets count")
	}
	for i := 0; i < 10; i++ {
		if generatedPlanets[i].Coords[0] != expectedPlanets[i].Coords[0] {
			t.Error("X coordinate missmatch on Planet[", strconv.Itoa(i), "] Expected", expectedPlanets[i].Coords[0], "Actual ", generatedPlanets[i].Coords[0])
		}
		if generatedPlanets[i].Coords[1] != expectedPlanets[i].Coords[1] {
			t.Error("Y coordinate missmatch on Planet[", strconv.Itoa(i), "] Expected", expectedPlanets[i].Coords[1], "Actual ", generatedPlanets[i].Coords[1])
		}
		if generatedPlanets[i].Texture != expectedPlanets[i].Texture {
			t.Error("Ring offset missmatch on Planet[", strconv.Itoa(i), "] Expected", expectedPlanets[i].Texture, "Actual ", generatedPlanets[i].Texture)
		}
		if generatedPlanets[i].Size != expectedPlanets[i].Size {
			t.Error("Size missmatch on Planet[", strconv.Itoa(i), "] Expected", expectedPlanets[i].Size, "Actual ", generatedPlanets[i].Size)
		}
	}
}

func TestPlanetMarshalling(t *testing.T) {
	var uPlanet Planet

	mPlanet, err := json.Marshal(planet)
	if err != nil {
		t.Error("Planet marshaling failed:", err)
	}

	err = json.Unmarshal(mPlanet, &uPlanet)
	if err != nil {
		t.Error("Planet unmarshaling failed:", err)
	}

	if planet.GetKey() != uPlanet.GetKey() {
		t.Error(
			"Keys of both planets are different!\n",
			planet.GetKey(),
			"!=",
			uPlanet.GetKey(),
		)
	}

	if !reflect.DeepEqual(planet, uPlanet) {
		t.Error("Planets are different after the marshal->unmarshal step")
	}
}

func TestPlanetHasOwner(t *testing.T) {
	if !planet.HasOwner() {
		t.Fail()
	}

	planet.Owner = ""
	if planet.HasOwner() {
		t.Fail()
	}
}
