package entities

import (
	"encoding/json"
	"reflect"
	"testing"
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
