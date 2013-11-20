package entities

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/Vladimiroff/vec2d"
)

func TestSunMarshalling(t *testing.T) {
	var uSun Sun

	mSun, err := json.Marshal(sun)
	if err != nil {
		t.Error("Sun marshaling failed:", err)
	}

	err = json.Unmarshal(mSun, &uSun)
	if err != nil {
		t.Error("Sun unmarshaling failed:", err)
	}

	if sun.Key() != uSun.Key() {
		t.Error(
			"Keys of both sun are different!\n",
			sun.Key(),
			"!=",
			uSun.Key(),
		)
	}

	if !reflect.DeepEqual(sun, uSun) {
		t.Error("Suns are different after the marshal->unmarshal step")
	}
}

func TestUpdateSun(t *testing.T) {
	sun := Sun{"gophie", "", 4, vec2d.New(100, 100), vec2d.New(20, 20)}
	sun.update()

	if sun.Position.X != 29.22842712474619 {
		t.Error("Suns's position is wrong: ", sun.Position.X)
	}
}
