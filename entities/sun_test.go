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

	testData := [][2]float64{
		{-9000, 0},
		{9000, 0},
		{-4500, 7794},
		{-4500, -7794},
		{4500, 7794},
		{4500, -7794},
	}

	testSlots := testSun.calculateAdjacentSlots()

	for i := 0; i < len(testData); i++ {
		if testSlots[i].Position.X != testData[i][0] || testSlots[i].Position.Y != testData[i][1] {
			t.Errorf("Left solar slot out of place. Coordinates: %s\n", testSlots[0].Position)
		}
	}
}

func TestBasicFindSolarSlotPosition(t *testing.T) {
	friends := []*Sun{
		{
			Username: "MOR001",
			Name:     "MORS01",
			Position: vec2d.New(0, 0),
		},
		{
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
		{
			Username: "MOR001",
			Name:     "MORS01",
			Position: vec2d.New(0, 0),
		},
		{
			Username: "MOR002",
			Name:     "MORS02",
			Position: vec2d.New(-9000, 0),
		},
		{
			Username: "MOR003",
			Name:     "MORS03",
			Position: vec2d.New(-18000, 0),
		},
		{
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
