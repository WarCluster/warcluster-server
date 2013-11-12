package entities

import (
	"fmt"
	"github.com/Vladimiroff/vec2d"
)

type Mission struct {
	Color      Color
	Source     string
	Target     string
	Type       string
	StartTime  int64
	TravelTime int64
	Player     string
	ShipCount  int32
}

// Database key.
func (m *Mission) Key() string {
	return fmt.Sprintf("mission.%d_%s", m.StartTime, m.Source)
}

// We plan to tweak the missions' speed based on some game logic.
// For now, 10 seems like a fair choice.
func (m *Mission) GetSpeed() int64 {
	return 10
}

// Calculates the travel time in milliseconds between two planets with given speed.
// Traveling is implemented like a simple time.Sleep from our side.
func calculateTravelTime(source, target *vec2d.Vector, speed int64) int64 {
	distance := vec2d.GetDistance(source, target)
	return int64(distance / float64(speed) * 100)
}

// When the missionary is done traveling (a.k.a. sleeping) calls this in order
// to calculate the outcome of the battle/suppliemnt/spying on target planet.
//
// In case of Attack: We have to check if the target planet is owned by the attacker.
// If that's true we simply increment the ship count on that planet. If not we do the
// math and decrease the count ship on the attacked planet. We should check if the attacker
// should own that planet, which comes with all the changing colors and owner stuff.
//
// In case of Supply: We simply increase the ship count and we're done :P
//
// In case of Spy: TODO
func EndMission(target *Planet, missionInfo *Mission) (excessShips int32) {
	switch missionInfo.Type {
	case "Attack":
		if target.Owner == missionInfo.Player {
			target.SetShipCount(target.ShipCount + missionInfo.ShipCount)
		} else {
			if missionInfo.ShipCount < target.ShipCount {
				target.SetShipCount(target.ShipCount - missionInfo.ShipCount)
			} else {
				if target.IsHome {
					target.SetShipCount(0)
					excessShips = missionInfo.ShipCount - target.ShipCount
				} else {
					target.SetShipCount(missionInfo.ShipCount - target.ShipCount)
					target.Owner = missionInfo.Player
					target.Color = missionInfo.Color
				}
			}
		}
	case "Supply":
		target.SetShipCount(target.ShipCount + missionInfo.ShipCount)
	case "Spy":
	}

	return
}
