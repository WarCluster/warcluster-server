package server

import (
	"reflect"
	"testing"

	"github.com/Vladimiroff/vec2d"
)

func TestAllTypesOfMission(t *testing.T) {

	var tableTests = []struct {
		input  string
		output func(*Request) error
	}{
		{"start_mission", parseAction},
		{"scope_of_view", scopeOfView},
		{"something_else", nil},
	}

	request := new(Request)
	request.StartPlanet = []string{"start"}
	request.EndPlanet = "end"
	request.Position = vec2d.New(2.0, 4.0)
	request.Resolution = []uint64{1920, 1080}

	for _, test := range tableTests {
		request.Command = test.input
		result, _ := ParseRequest(request)

		result_value := reflect.ValueOf(result)
		output_value := reflect.ValueOf(test.output)

		if result_value.Pointer() != output_value.Pointer() {
			t.Errorf("Request with %s returnes %#v, expected %#v", test.input, result, test.output)
		}
	}
}

func TestStartMissionWithoutEnoughArguments(t *testing.T) {
	request := new(Request)
	request.Command = "start_mission"
	request.StartPlanet = []string{"start"}
	result, _ := ParseRequest(request)

	if result != nil {
		t.Errorf("Request start_mision without EndPlanet returnes %#v", result)
	}
}

func TestScopeOfViewWithoutEnoughArguments(t *testing.T) {
	request := new(Request)
	request.Command = "scope_of_view"
	request.Position = vec2d.New(2.0, 4.0)
	result, _ := ParseRequest(request)

	if result != nil {
		t.Errorf("Request scope_of_view without EndPlanet returnes %#v", result)
	}
}

func TestStartMissionWithNegativeFleet(t *testing.T) {
	request := new(Request)
	request.Command = "start_mission"
	request.StartPlanet = []string{"start"}
	request.EndPlanet = "end"
	request.Fleet = -10
	ParseRequest(request)

	if request.Fleet != 100 {
		t.Errorf("Request start_mision with negative fleet makes fleet size %d", request.Fleet)
	}
}

func TestStartMissionWithMoreThanHundredFleet(t *testing.T) {
	request := new(Request)
	request.Command = "start_mission"
	request.StartPlanet = []string{"start"}
	request.EndPlanet = "end"
	request.Fleet = 150
	ParseRequest(request)

	if request.Fleet != 100 {
		t.Errorf("Request start_mision with negative fleet makes fleet size %d", request.Fleet)
	}
}
