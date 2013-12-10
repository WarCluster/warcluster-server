package response

import (
	"testing"

	"github.com/Vladimiroff/vec2d"
)

func TestCalculateCanvasSize(t *testing.T) {
	topLeft, bottomRight := calculateCanvasSize(vec2d.New(20, 50), []uint16{800, 600})
	if *topLeft != *vec2d.New(-5380, 5350) {
		t.Errorf("topLeft is %#v, expected: %#v", *topLeft, *vec2d.New(-5380, 5350))
	}
	if *bottomRight != *vec2d.New(5420, -5250) {
		t.Errorf("bottomRight is %#v, expected: %#v", *bottomRight, *vec2d.New(5420, -5250))
	}
}

func TestListAreas(t *testing.T) {
	topLeft := vec2d.New(-5380, 5350)
	bottomRight := vec2d.New(5420, -5250)

	areas := listAreas(topLeft, bottomRight)
	if len(areas) != 4 {
		t.Errorf("Areas turned out to be %#v", areas)
	}
}
