package entities

import (
	"fmt"
	"math"

	"github.com/Vladimiroff/vec2d"
)

type SolarSlot struct {
	Data     string
	Position *vec2d.Vector
}

func newSolarSlot(x, y float64) *SolarSlot {
	ss := new(SolarSlot)
	ss.Position = vec2d.New(x, y)
	ss.Data = ""
	return ss
}

func (ss *SolarSlot) Key() string {
	return fmt.Sprintf("ss.%v_%v", ss.Position.X, ss.Position.Y)
}

// Returns the set by X or Y where this entity has to be put in
func (ss *SolarSlot) AreaSet() string {
	return fmt.Sprintf(
		Settings.AreaTemplate,
		RoundCoordinateTo(ss.Position.X),
		RoundCoordinateTo(ss.Position.Y),
	)
}

func (ss *SolarSlot) fetchSolarSlotsLayer(zuLevel uint32) (results []string) {
	zLevel := float64(zuLevel)

	angeledOffsetStepX := math.Floor((Settings.SolarSystemRadius / 2) + 0.5)
	verticalOffset := math.Floor((Settings.SolarSystemRadius * math.Sqrt(3) / 2) + 0.5)
	angeledOffsetX := angeledOffsetStepX * zLevel

	horizontalOffsetStepX := Settings.SolarSystemRadius
	horizontalOffsetX := horizontalOffsetStepX * zLevel

	results = append(results, newSolarSlot(ss.Position.X-horizontalOffsetX, ss.Position.Y).Key())
	results = append(results, newSolarSlot(ss.Position.X+horizontalOffsetX, ss.Position.Y).Key())
	results = append(results, newSolarSlot(ss.Position.X-angeledOffsetX, ss.Position.Y+verticalOffset*zLevel).Key())
	results = append(results, newSolarSlot(ss.Position.X-angeledOffsetX, ss.Position.Y-verticalOffset*zLevel).Key())
	results = append(results, newSolarSlot(ss.Position.X+angeledOffsetX, ss.Position.Y+verticalOffset*zLevel).Key())
	results = append(results, newSolarSlot(ss.Position.X+angeledOffsetX, ss.Position.Y-verticalOffset*zLevel).Key())

	for i := 1; i < int(zLevel); i++ {
		results = append(results, newSolarSlot(ss.Position.X-horizontalOffsetX+angeledOffsetStepX*float64(i), ss.Position.Y+verticalOffset*float64(i)).Key())
		results = append(results, newSolarSlot(ss.Position.X-horizontalOffsetX+angeledOffsetStepX*float64(i), ss.Position.Y-verticalOffset*float64(i)).Key())
		results = append(results, newSolarSlot(ss.Position.X+horizontalOffsetX-angeledOffsetStepX*float64(i), ss.Position.Y+verticalOffset*float64(i)).Key())
		results = append(results, newSolarSlot(ss.Position.X+horizontalOffsetX-angeledOffsetStepX*float64(i), ss.Position.Y-verticalOffset*float64(i)).Key())
		results = append(results, newSolarSlot(ss.Position.X-angeledOffsetX+horizontalOffsetStepX*float64(i), ss.Position.Y+verticalOffset*zLevel).Key())
		results = append(results, newSolarSlot(ss.Position.X-angeledOffsetX+horizontalOffsetStepX*float64(i), ss.Position.Y-verticalOffset*zLevel).Key())
	}
	return results
}
