package entities

import (
	"github.com/Vladimiroff/vec2d"
	"testing"
)

func TestBasePreparations(t *testing.T) {
	sun := Sun{"gophie", 1, vec2d.New(0, 0), vec2d.New(20, 20)}
	expected_json := "{\"Username\":\"gophie\"}"
	expected_key := "sun.20_20"

	key, json, err := sun.Serialize()
	if key != expected_key || string(json) != expected_json {
		t.Error(string(json))
		t.Error("Sun JSON formatting gone wrong!")
	}

	if err != nil {
		t.Error("Error during serialization: ", err)
	}
}

func TestDeserializeSun(t *testing.T) {
	var sun Sun
	serialized_Sun := []byte("{\"Username\":\"gophie\"}")
	sun = Construct("sun.20_20", serialized_Sun).(Sun)

	if sun.Username != "gophie" {
		t.Error("Player's name is ", sun.Username)
	}

	if sun.position.Y != 20 || sun.position.X != 20 {
		t.Error("Kiro da napravi serialize na vektori is ", sun.position.Y, sun.position.X)
	}
}

func TestUpdateSun(t *testing.T) {
	sun := Sun{"gophie", 4, vec2d.New(100, 100), vec2d.New(20, 20)}
	sun.Update()

	if sun.position.X != 29.22842712474619 {
		t.Error("Suns's position is wrong: ", sun.position.X)
	}
}
