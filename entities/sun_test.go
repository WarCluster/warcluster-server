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

func TestAdjacentSlotsCalculation(t *testing.T) {
	testSun := Sun{
		Username: "MOR001",
		Name:     "MORS01",
		Position: vec2d.New(0, 0),
	}

	testSlots := testSun.calculateAdjacentSlots()

	if testSlots[0].Position.X != -9000 || testSlots[0].Position.Y != 0 {
		t.Error("Left solar slot out of place. Coordinates: ")
		t.Error(testSlots[0].Position)
	}

	if testSlots[1].Position.X != 9000 || testSlots[1].Position.Y != 0 {
		t.Error("Right solar slot out of place. Coordinates: ")
		t.Error(testSlots[1].Position)
	}

	if testSlots[2].Position.X != -4500 || testSlots[2].Position.Y != 7794 {
		t.Error("Left Top solar slot out of place. Coordinates: ")
		t.Error(testSlots[2].Position)
	}

	if testSlots[3].Position.X != -4500 || testSlots[3].Position.Y != -7794 {
		t.Error("Left Bottom solar slot out of place. Coordinates: ")
		t.Error(testSlots[3].Position)
	}

	if testSlots[4].Position.X != 4500 || testSlots[4].Position.Y != 7794 {
		t.Error("Right Top solar slot out of place. Coordinates: ")
		t.Error(testSlots[4].Position)
	}

	if testSlots[5].Position.X != 4500 || testSlots[5].Position.Y != -7794 {
		t.Error("Right Bottom solar slot out of place. Coordinates: ")
		t.Error(testSlots[5].Position)
	}
}

func TestBasicFindSolarSlotPosition(t *testing.T) {
	friends := []*Sun{
		&Sun{
			Username: "MOR001",
			Name:     "MORS01",
			Position: vec2d.New(0, 0),
		},
		&Sun{
			Username: "MOR002",
			Name:     "MORS02",
			Position: vec2d.New(200, 200),
		},
	}

	targetSlot := getStartSolarSlotPosition(friends)

	if targetSlot.Position.X != 0 || targetSlot.Position.Y != 0 {
		t.Error("Target solar slot out of place. Coordinates: ")
		t.Error(targetSlot.Position)
	}
}

func TestFindSolarSlotPosition(t *testing.T) {
	friends := []*Sun{
		&Sun{
			Username: "MOR001",
			Name:     "MORS01",
			Position: vec2d.New(0, 0),
		},
		&Sun{
			Username: "MOR002",
			Name:     "MORS02",
			Position: vec2d.New(-9000, 0),
		},
		&Sun{
			Username: "MOR003",
			Name:     "MORS03",
			Position: vec2d.New(-18000, 0),
		},
		&Sun{
			Username: "MOR004",
			Name:     "MORS04",
			Position: vec2d.New(-4500, -7794),
		},
	}

	targetSlot := getStartSolarSlotPosition(friends)

	if targetSlot.Position.X != -9000 || targetSlot.Position.Y != 0 {
		t.Error("Target solar slot out of place. Coordinates: ")
		t.Error(targetSlot.Position)
	}
}
