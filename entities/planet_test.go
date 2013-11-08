package entities

import (
	"encoding/json"
	"github.com/Vladimiroff/vec2d"
	"reflect"
	"strconv"
	"testing"
)

func TestGeneratePlanets(t *testing.T) {
	expectedPlanets := []Planet{
		Planet{"ABC1234", Color{22, 22, 22}, vec2d.New(-77, 57), false, 6, 3, timeStamp, 10, 0, "gophie"},
		Planet{"ABC1234", Color{22, 22, 22}, vec2d.New(1470, 300), false, 8, 5, timeStamp, 10, 0, "gophie"},
		Planet{"ABC1234", Color{22, 22, 22}, vec2d.New(-690, -201), false, 3, 1, timeStamp, 10, 0, "gophie"},
		Planet{"ABC1234", Color{22, 22, 22}, vec2d.New(-1052, 648), false, 2, 8, timeStamp, 10, 0, "gophie"},
		Planet{"ABC1234", Color{22, 22, 22}, vec2d.New(1428, -1364), false, 3, 1, timeStamp, 10, 0, "gophie"},
		Planet{"ABC1234", Color{22, 22, 22}, vec2d.New(2735, 300), false, 6, 8, timeStamp, 10, 0, "gophie"},
		Planet{"ABC1234", Color{22, 22, 22}, vec2d.New(2818, -799), false, 9, 6, timeStamp, 10, 0, "gophie"},
		Planet{"ABC1234", Color{22, 22, 22}, vec2d.New(-323, 3080), false, 5, 4, timeStamp, 10, 0, "gophie"},
		Planet{"ABC1234", Color{22, 22, 22}, vec2d.New(1547, 3339), false, 1, 1, timeStamp, 10, 0, "gophie"},
		Planet{"ABC1234", Color{22, 22, 22}, vec2d.New(-2745, -1066), false, 4, 6, timeStamp, 10, 0, "gophie"},
	}
	sun.Position = vec2d.New(500, 300)
	generatedPlanets, _ := GeneratePlanets("gophie", &sun)

	if len(generatedPlanets) != 10 {
		t.Error("Wrong planets count")
	}
	for i := 0; i < 10; i++ {
		if generatedPlanets[i].Position.X != expectedPlanets[i].Position.X {
			t.Error("X coordinate missmatch on Planet[", strconv.Itoa(i), "] Expected", expectedPlanets[i].Position.X, "Actual ", generatedPlanets[i].Position.X)
		}
		if generatedPlanets[i].Position.Y != expectedPlanets[i].Position.Y {
			t.Error("Y coordinate missmatch on Planet[", strconv.Itoa(i), "] Expected", expectedPlanets[i].Position.Y, "Actual ", generatedPlanets[i].Position.Y)
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

	if planet.Key() != uPlanet.Key() {
		t.Error(
			"Keys of both planets are different!\n",
			planet.Key(),
			"!=",
			uPlanet.Key(),
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
