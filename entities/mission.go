package entities

import (
	"fmt"
	"github.com/Vladimiroff/vec2d"
)

type Mission struct {
	Color       Color
	Source      string
	Target      string
	Type        string
	StartTime   int64
	TravelTime  int64
	Player      string
	ShipCount   int
}

func (m *Mission) String() string {
	return m.GetKey()
}

func (m *Mission) GetKey() string {
	return fmt.Sprintf("mission.%d_%s", m.StartTime, m.Source)
}

func (m *Mission) GetSpeed() int64 {
	return 10
}

func calculateTravelTime(source, target *vec2d.Vector, speed int64) int64 {
	distance := vec2d.GetDistance(source, target)
	return int64(distance / float64(speed) * 100)
}

func EndMission(target *Planet, missionInfo *Mission) (excessShips int) {
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
